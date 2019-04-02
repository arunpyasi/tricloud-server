package database

import (
	"encoding/json"
	"errors"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type User struct {
	ID       string `json:"id"`
	UserName string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	FullName string `json:"fullname,omitempty"`
	Email    string `json:"email,omitempty"`
	Active   string `json:"active,omitempty"`
	Agent    *Agent `json:"agents,omitempty"`
}

type Agent struct {
	ID         string `json:"id"`
	OS         string `json:"os,omitempty"`
	Key        string `json:"key,omitempty"`
	LastLogin  string `json:"lastlogin,omitempty"`
	FirstAdded string `json:"firstadded,omitempty"`
}

func CreateUser(user_data User) error {
	err := Conn.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return fmt.Errorf("Failed to create bucket: %v", err)
		}
		enc, _ := json.Marshal(user_data)
		var dec []byte
		json.Unmarshal(enc, &dec)
		fmt.Print(string(enc))
		if err := bk.Put([]byte(user_data.ID), enc); err != nil {
			return fmt.Errorf("Failed to insert '%s'", user_data.ID)
		}
		return err
	})
	return err
}
func GetUser(id string) ([]byte, error) {
	var user_details []byte
	err := Conn.View(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}
		x := bk
		user_details = x.Get([]byte(id))
		if user_details == nil {
			err = errors.New("No user with ID " + id + " found")
		}
		return err
	})
	m := make(map[string]interface{})
	m["data"] = user_details
	json_data, err := json.Marshal(m)
	return json_data, err
}

func GetAllUsers() ([]byte, error) {
	var users []User
	Conn.View(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return fmt.Errorf("Failed to create bucket: %v", err)
		}
		x := bk
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

func DeleteUser(user_id string) error {
	err := Conn.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte("users"))
		err := bk.Delete([]byte(user_id))
		return err
	})
	return err
}
