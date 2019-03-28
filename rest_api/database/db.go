package database

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

type User struct {
	ID       string
	Username string
}

func AddUsers(db *bolt.DB, id string, username string) error {
	err := db.Update(func(tx *bolt.Tx) error {
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

func GetAllUsers(db *bolt.DB) ([]User, error) {
	var users []User
	err := db.View(func(tx *bolt.Tx) error {
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
	return users, err
}
