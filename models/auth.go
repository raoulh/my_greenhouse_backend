package models

import (
	"time"
)

func cleanupExpired() {
	db.Where("valid_until < ?", time.Now()).Delete(&User{})
}

func Login(username, pass string) (at *User, err error) {
	defer cleanupExpired()

	//create a new token for this user
	at = &User{}

	err = db.Create(at).Error

	return
}

func Logout(at *User) (err error) {
	defer cleanupExpired()

	return db.Where("token = ?", at.ID).Delete(&User{}).Error
}
