package database

import (
	"errors"
	"fmt"

	"github.com/indrenicloud/tricloud-server/app/auth"
)

type User struct {
	ID        string   `json:"id"`
	Password  string   `json:"password,omitempty"`
	FullName  string   `json:"fullname,omitempty"`
	Email     string   `json:"email,omitempty"`
	SuperUser bool     `json:"superuser,omitempty"`
	APIKeys   []string `json:"apikey"`
	Agents    []string `json:"agents"`
}

var (
	UserBucketName = []byte("users")
)

func NewUser(userInfo map[string]string, superuser bool) (*User, error) {
	fields := []string{"id", "password", "fullname", "email"}

	for _, field := range fields {
		value, ok := userInfo[field]
		if !ok {
			return nil, fmt.Errorf("field %s not found", value)
		}
	}

	return &User{
		ID:        userInfo["id"],
		Password:  auth.GeneratePassword(userInfo["password"]),
		FullName:  userInfo["fullname"],
		SuperUser: superuser,
		Email:     userInfo["email"],
	}, nil

}

func CreateUser(user *User) error {

	userbyte, err := Encode(user)
	if err != nil {
		return err
	}

	err = DB.Create([]byte(user.ID), userbyte, UserBucketName)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(id string) (*User, error) {
	user := &User{}
	userbyte, err := DB.Read([]byte(id), UserBucketName)

	if err != nil {
		return nil, err
	}

	err = Decode(userbyte, user)
	if err != nil {
		return nil, err
	}
	user.Password = ""

	return user, nil
}

func GetapiKeys(id string) ([]string, error) {

	userbyte, err := DB.Read([]byte(id), UserBucketName)
	if err != nil {
		return nil, err
	}

	u := &User{}
	err = Decode(userbyte, u)
	if err != nil {
		return nil, err
	}

	return u.APIKeys, nil
}

func AddapiKey(id, keytype string) error {

	userbyte, err := DB.Read([]byte(id), UserBucketName)
	if err != nil {
		return err
	}

	u := &User{}
	err = Decode(userbyte, u)
	if err != nil {
		return err
	}

	newkey := auth.NewAPIKey(keytype, id, false)
	if newkey == "" {
		return errors.New("could not create api")
	}
	u.APIKeys = append(u.APIKeys, newkey)
	userbyte, err = Encode(u)
	if err != nil {
		return err
	}
	return DB.Update([]byte(u.ID), userbyte, UserBucketName)

}

func RemoveapiKey(id, key string) error {
	userbyte, err := DB.Read([]byte(id), UserBucketName)
	if err != nil {
		return err
	}

	u := &User{}
	err = Decode(userbyte, u)
	if err != nil {
		return err
	}

	var newkeys []string
	for _, value := range u.APIKeys {
		if !(value == key) {
			newkeys = append(newkeys, value)
		}
	}
	u.APIKeys = newkeys

	userbyte, err = Encode(u)
	if err != nil {
		return err
	}
	return DB.Update([]byte(u.ID), userbyte, UserBucketName)
}

func GetAllUsers() ([]*User, error) {
	users := []*User{}

	usersbyte, err := DB.ReadAll(UserBucketName)
	if err != nil {
		return nil, err
	}

	for index, val := range usersbyte {
		user := &User{}
		err = Decode(val, user)
		user.Password = ""
		users[index] = user
		if err != nil {
			return nil, err
		}

	}
	return users, nil
}

func UpdateUser(userinfo map[string]string) error {
	fields := []string{"id", "password", "fullname", "email"}

	olduserbyte, err := DB.Read([]byte(userinfo["id"]), UserBucketName)
	if err != nil {
		return err
	}

	olduser := make(map[string]string)

	err = Decode(olduserbyte, olduser)
	if err != nil {
		return err
	}

	pass, ok := userinfo["password"]
	if ok {
		userinfo["password"] = auth.GeneratePassword(pass)
	}

	for _, val := range fields {
		_, ok := userinfo[val]
		if ok {
			olduser[val] = userinfo[val]
		}
	}

	userbyte, err := Encode(olduser)
	if err != nil {
		return err
	}

	return DB.Update([]byte(olduser["id"]), userbyte, UserBucketName)
}

func DeleteUser(id string) error {
	return DB.Delete([]byte(id), UserBucketName)
}
