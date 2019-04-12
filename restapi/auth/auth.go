package auth

import (
	"errors"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var signkey = []byte("DEVpASsWord_flyingpinkelephant")

type myClaims struct {
	Apitype  string `json:"type,omitempty"`
	User     string `json:"user,omitempty"`
	IssuedAt int64  `json:"issued,omitempty"`
	Expiry   int64  `json:"expiry,omitempty"`
}

func (c *myClaims) Valid() error {
	if !(c.Apitype == "agent" || c.Apitype == "user" || c.Apitype == "session") {
		return errors.New("invalid type")
	}

	if c.User == "" {
		return errors.New("User not set")
	}

	if c.IssuedAt == 0 {
		return errors.New("issued date not valid")
	}

	if c.Apitype == "session" {
		if c.Expiry == 0 {
			return errors.New("expirydate needed")
			// TODO check expiry date to current date
		}
	}

	return nil
}

func NewAPIKey(apitype, username string) string {

	expiryDate := time.Time{}

	if apitype == "session" {
		expiryDate = time.Now().Add(24 * 30 * time.Hour)
	}

	claims := &myClaims{
		Apitype:  apitype,
		IssuedAt: time.Now().UnixNano(),
		User:     username,
		Expiry:   expiryDate.UnixNano(),
	}

	signer := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	api, err := signer.SignedString(signkey)
	if err != nil {
		log.Println("could not generate api:", err)
		return ""
	}

	return api
}

func CheckAPIKey(tokenstr string) *jwt.Token {

	token, err := jwt.Parse(tokenstr, func(token *jwt.Token) (interface{}, error) {
		return signkey, nil
	})
	if err != nil {
		return nil
	}

	return token
}
