package models

import (
	"fmt"
	"time"

	"git.raoulh.pw/raoulh/my_greenhouse_backend/models/myfood"
)

//GetUserByEmailAndID get a user entry from db. if full is true, loads all measurements from db too
func GetUserByEmailAndID(email, deviceID string, full bool) (u *User, err error) {
	if full {
		err = db.Preload("UnitMeasurements").Where("mf_username = ? and device_id = ?", email, deviceID).Last(u).Error
	} else {
		err = db.Where("mf_username = ? and device_id = ?", email, deviceID).Last(u).Error
	}
	return
}

//GetUserByTokenAndID get a user entry from db. if full is true, loads all measurements from db too
func GetUserByTokenAndID(token, deviceID string, full bool) (u *User, err error) {
	if full {
		err = db.Preload("UnitMeasurements").Where("mf_token = ? and device_id = ?", token, deviceID).Last(u).Error
	} else {
		err = db.Where("mf_token = ? and device_id = ?", token, deviceID).Last(u).Error
	}
	return
}

func Login(username, pass, deviceID string) (u *User, err error) {
	mf := myfood.NewMyFoodApi(myfood.MyFoodApiHost)

	//Check if username exists in db
	u, err = GetUserByEmailAndID(username, deviceID, false)
	usersLock.Lock(u)
	defer usersLock.Unlock(u)

	if err == nil {
		//user on this device found in db, use this entry and update the token from MF
		t, err := mf.GetToken(username, pass)
		if err != nil {
			logging.Warnf("MF.GetToken failed: %v", err)
			return nil, err
		}

		u.MF_Token = t.Token
		u.MF_RefreshToken = t.RefreshToken
		u.MF_RefreshTokenFailureCout = 0

		//As of now, the validity from MF api is broken.
		//Use something that is almost less than 3 months. After it expires, the MF token will still be valid for
		//some time, but next time the app will validate it's token, it will get a renew with a new one.
		//This way if the app has not been launched in 2 weeks, the token is invalidated and user must login again
		u.MF_TokenValidity = time.Now().AddDate(0, 0, 75)
		u.MF_TokenValid = true

		u.LastLogin = time.Now()

		go RefreshUserData(u.ID)

		//save the new token to db
		return u, db.Save(&u).Error
	}

	//User does not exist in db, login to MF api and create a new user with token and deviceID
	t, err := mf.GetToken(username, pass)
	if err != nil {
		logging.Warnf("MF.GetToken failed: %v", err)
		return
	}

	u.LastLogin = time.Now()
	u.DeviceID = deviceID
	u.MF_Username = username
	u.MF_Token = t.Token
	u.MF_RefreshToken = t.RefreshToken
	u.MF_TokenValidity = time.Now().AddDate(0, 0, 75) //see comment ^
	u.MF_TokenValid = true
	u.MF_RefreshTokenFailureCout = 0

	err = db.Create(u).Error

	go RefreshUserData(u.ID)

	return u, err
}

func CheckToken(u *User) (err error) {
	if !u.MF_TokenValid {
		return fmt.Errorf("token is not valid anymore")
	}
	return
}
