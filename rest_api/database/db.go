package database

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type User struct {
	ID       string `json:"id"`
	UserName string `json:"username,omitempty"`
	FullName string `json:"fullname,omitempty"`
	Email    string `json:"email,omitempty"`
	Agent    *Agent `json:"agent,omitempty"`
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

func GetAllUsers() ([]byte, error) {
	var users []User
	err := Conn.View(func(tx *bolt.Tx) error {
		x := tx.Bucket([]byte("users"))
		c := x.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var data User
			json.Unmarshal(v, &data)
			users = append(users, data)
		}
		return nil
	})
	user_json, _ := json.Marshal(users)

	return user_json, err
}
