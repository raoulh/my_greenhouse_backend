package models

import (
	"crypto/rand"
	"fmt"
	"time"
)

type AuthToken struct {
	ID         uint      `gorm:"primarykey" json:"-"`
	ValidUntil time.Time `json:"valid_until"`
	Token      string    `json:"token"`
	Name       string    `json:"name"`
	FullName   string    `json:"fullname"`
	LoginName  string    `json:"loginname"`
	Email      string    `json:"email"`
}

const (
	tokenDuration = time.Hour * 10 //10 hours token validity
)

func IsTokenValid(t string) bool {
	defer cleanupExpired()

	var auth AuthToken
	if err := db.Where("token = ? AND valid_until > ?", t, time.Now()).Last(&auth).Error; err != nil {
		return false
	}
	return true
}

func cleanupExpired() {
	db.Where("valid_until < ?", time.Now()).Delete(&AuthToken{})
}

func Login(username, pass string) (at *AuthToken, err error) {
	defer cleanupExpired()

	//create a new token for this user
	at = &AuthToken{
		ValidUntil: time.Now().Add(tokenDuration),
		Token:      generateToken(),
		Name:       username,
		LoginName:  username,
	}

	err = db.Create(at).Error

	return
}

func Logout(at *AuthToken) (err error) {
	defer cleanupExpired()

	return db.Where("token = ?", at.Token).Delete(&AuthToken{}).Error
}

func generateToken() string {
	b := make([]byte, 24)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
