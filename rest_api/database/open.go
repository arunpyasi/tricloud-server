package database

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

func OpenRead() *bolt.DB {
	return openDB(true)
}

func OpenWrite() *bolt.DB {
	return openDB(false)
}

func openDB(readOnly bool) *bolt.DB {
	db, err := bolt.Open("my.db", 0600, &bolt.Options{ReadOnly: readOnly})
	if err != nil {
		log.Fatalf("error: %s", err)

	}
	fmt.Print("opened db")
	return db
}
