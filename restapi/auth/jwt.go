package auth

import (
	"errors"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var signkey = []byte("DEVpASsWord_flyingpinkelephant")

type MyClaims struct {
	Apitype  string `json:"t,omitempty"`
	User     string `json:"u,omitempty"`
	Super    bool   `json:"s,omitempty"`
	IssuedAt int64  `json:"i,omitempty"`
	Expiry   int64  `json:"e,omitempty"`
	//jwt.StandardClaims
}

func (c MyClaims) Valid() error {
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

func NewAPIKey(apitype, username string, superuser bool) string {

	expiryDate := time.Time{}

	if apitype == "session" {
		expiryDate = time.Now().Add(24 * 30 * time.Hour)
	}

	claims := &MyClaims{
		Apitype:  apitype,
		IssuedAt: time.Now().UnixNano(),
		Super:    superuser,
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

func ParseAPIKey(tokenstr string) *jwt.Token {

	token, err := jwt.ParseWithClaims(tokenstr, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signkey, nil
	})
	if err != nil {
		log.Println(err)
		return nil
	}

	return token
}
