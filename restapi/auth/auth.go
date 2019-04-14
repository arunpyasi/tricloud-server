package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hashedPassword)
}

func IsPasswordCorrect(pass, expected string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(expected))
	if err == nil {
		return true
	}
	return false
}
