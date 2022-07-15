package models

import (
	"time"

	"git.raoulh.pw/raoulh/my_greenhouse_backend/models/myfood"
)

func GetFullUser(userID uint) (u *User, err error) {
	err = db.Preload("UnitMeasurements").Where("id = ?", userID).Last(u).Error
	return
}

//RefreshUserData will call MyFood api to retrieve latest data and measurements
func RefreshUserData(userID uint) {
	u, err := GetFullUser(userID)
	if err != nil {
		logging.Warnf("unable to get user with id: %d: %v", userID, err)
		return
	}

	if !u.MF_TokenValid {
		return
	}

	usersLock.Lock(u)
	defer usersLock.Unlock(u)

	mf := myfood.NewMyFoodApi(myfood.MyFoodApiHost)

	//check token validity
	if u.MF_TokenValidity.After(time.Now()) {
		//get a new token
		logging.Info("token is going to expire, refresh token")

		refreshToken := &myfood.TokenData{
			Token:        u.MF_Token,
			RefreshToken: u.MF_RefreshToken,
		}

		//try to get a new token
		t, err := mf.RefreshToken(refreshToken)
		if err != nil {
			logging.Warnf("MF.RefreshToken failed: %v", err)

			u.MF_RefreshTokenFailureCout++

			if u.MF_RefreshTokenFailureCout >= 5 {
				//Can't get a new token. Mark it as invalid and stop updating this user...
				logging.Warnf("RefreshToken failed after 5 attempts. Aborting")

				u.MF_TokenValid = false
			}

			db.Save(u)
			return
		}

		u.MF_Token = t.Token
		u.MF_RefreshToken = t.RefreshToken
		u.MF_TokenValidity = time.Now().AddDate(0, 0, 75)
		u.MF_TokenValid = true
	}

	units, err := mf.GetAllProductionUnitIdsForCurrentUser(u.MF_Token)
	if err != nil {
		logging.Warn(err.Error())
		return
	}

	for i := range units.Data.ProductUnits {
		idx := -1
		for ii := range u.Meas {
			if u.Meas[ii].ProdUnitID == units.Data.ProductUnits[i].ID {
				idx = ii
			}
		}

		if idx >= 0 { //this prodId exists already
			u.Meas[idx].ProdUnitID = units.Data.ProductUnits[i].ID
		} else {
			//create a new entry
			umeas := UnitMeasurements{
				ProdUnitID: units.Data.ProductUnits[i].ID,
			}
			u.Meas = append(u.Meas, umeas)
			idx = len(u.Meas) - 1
		}

		//refresh measurements
		prodDetail, err := mf.GetProductionUnitDetailForUser(u.MF_Token, u.Meas[idx].ProdUnitID, 0)
		if err != nil {
			logging.Warn(err.Error())
			return
		}

		u.Meas[idx].ProdUnitType = prodDetail.Data.ProductionUnitType

		u.Meas[idx].PH.CurrentValue = prodDetail.Data.CurrentPhValue
		u.Meas[idx].PH.DayAverage = prodDetail.Data.AverageDayPhValue
		u.Meas[idx].PH.HourAverage = prodDetail.Data.AverageHourPhValue
		u.Meas[idx].PH.CurrentTime = prodDetail.Data.CurrentPhCaptureTime
		u.Meas[idx].PH.LastDayTime = prodDetail.Data.LastDayPhCaptureTime

		u.Meas[idx].Water.CurrentValue = prodDetail.Data.CurrentWaterTempValue
		u.Meas[idx].Water.DayAverage = prodDetail.Data.AverageDayWaterTempValue
		u.Meas[idx].Water.HourAverage = prodDetail.Data.AverageHourWaterTempValue
		u.Meas[idx].Water.CurrentTime = prodDetail.Data.CurrentWaterTempCaptureTime
		u.Meas[idx].Water.LastDayTime = prodDetail.Data.LastDayWaterTempCaptureTime

		u.Meas[idx].Air.CurrentValue = prodDetail.Data.CurrentAirTempValue
		u.Meas[idx].Air.DayAverage = prodDetail.Data.AverageDayAirTempValue
		u.Meas[idx].Air.HourAverage = prodDetail.Data.AverageHourAirTempValue
		u.Meas[idx].Air.CurrentTime = prodDetail.Data.CurrentAirTempCaptureTime
		u.Meas[idx].Air.LastDayTime = prodDetail.Data.LastDayAirTempCaptureTime

		u.Meas[idx].Humidity.CurrentValue = prodDetail.Data.CurrentHumidityValue
		u.Meas[idx].Humidity.DayAverage = prodDetail.Data.AverageDayHumidityValue
		u.Meas[idx].Humidity.HourAverage = prodDetail.Data.AverageHourHumidityValue
		u.Meas[idx].Humidity.CurrentTime = prodDetail.Data.CurrentHumidityCaptureTime
		u.Meas[idx].Humidity.LastDayTime = prodDetail.Data.LastDayHumidityCaptureTime
	}

	if db.Save(u).Error != nil {
		logging.Warnf("Failed to save user data into db: %v", err)
	}
}