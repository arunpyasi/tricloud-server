package database

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/indrenicloud/tricloud-server/restapi/auth"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string   `json:"id"`
	Password string   `json:"password,omitempty"`
	FullName string   `json:"fullname,omitempty"`
	Email    string   `json:"email,omitempty"`
	APIKeys  []string `json:"apikey"`
}

var (
	UserBucketName = []byte("users")
)

func NewUser(userInfo map[string]string) (*User, error) {
	fields := []string{"id", "password", "fullname", "email"}

	for _, field := range fields {

		value, ok := userInfo[field]
		if !ok {
			return nil, fmt.Errorf("field %s not found", value)
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInfo["password"]), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("couldnot create password hash")
	}

	return &User{
		ID:       userInfo["id"],
		Password: string(hashedPassword),
		FullName: userInfo["fullname"],
		Email:    userInfo["email"],
	}, nil

}

func CreateUser(userInfo map[string]string) error {
	user, err := NewUser(userInfo)
	if err != nil {
		return err
	}
	userbyte, err := json.Marshal(user)
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

	err = json.Unmarshal(userbyte, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetapiKeys(id string) ([]string, error) {

	userbyte, err := DB.Read([]byte(id), UserBucketName)
	if err != nil {
		return nil, err
	}

	u := &User{}
	err = json.Unmarshal(userbyte, u)
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
	err = json.Unmarshal(userbyte, u)
	if err != nil {
		return err
	}

	newkey := auth.NewAPIKey(keytype, id)
	if newkey == "" {
		return errors.New("could not create api")
	}

	u.APIKeys = append(u.APIKeys, newkey)
	userbyte, err = json.Marshal(u)

	return nil
}

func RemoveapiKey(id, key string) error {
	userbyte, err := DB.Read([]byte(id), UserBucketName)
	if err != nil {
		return err
	}

	u := &User{}
	err = json.Unmarshal(userbyte, u)
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
	return nil
}

/*
func GetAllUsers() ([]byte, error) {
	var users []User
	Conn.View(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("users"))
		x := tx.Bucket([]byte("users"))
		c := x.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var data User
			json.Unmarshal(v, &data)
			users = append(users, data)
		}
		return nil
	})
	m := make(map[string]interface{})
	m["data"] = users
	json_data, err := json.Marshal(m)
	return json_data, err
}

func UpdateUser(id string, user User) error {
	err := Conn.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return fmt.Errorf("Failed to create bucket: %v", err)
		}
		enc, _ := json.Marshal(user)
		var dec []byte
		json.Unmarshal(enc, &dec)
		if err := bk.Put([]byte(id), enc); err != nil {
			return fmt.Errorf("Failed to update '%s'", user)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Failed to update : %v", err)
	}
	return nil
}

func DeleteUser(id string) error {
	err := Conn.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte("users"))
		err := bk.Delete([]byte(id))
		return err
	})
	return err
}
*/
