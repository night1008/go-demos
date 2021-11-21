package auth

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	AuthnUserIDKey           = "authn-user-id"
	AuthnUserLastActiveAtKey = "authn-user-last-active-at"
)

type Auther interface {
	GetDB() *gorm.DB
	Authenticate(username, password string) (bool, interface{}, error)
}

func HashPassword(password string, cost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
