package database

import (
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type User struct {
	ID       string
	Username string
}

func AddUsers(id string, username string) error {
	err := Conn.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return fmt.Errorf("Failed to create bucket: %v", err)
		}

		if err := bk.Put([]byte(id), []byte(username)); err != nil {
			return fmt.Errorf("Failed to insert '%s': %v", id, username)
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
			users = append(users, User{
				ID:       string(k),
				Username: string(v),
			})
		}
		return nil
	})
	user_json, _ := json.Marshal(users)

	return user_json, err
}
