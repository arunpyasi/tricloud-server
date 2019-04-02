package database

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

var Conn *bolt.DB

func init() {
	Conn = openDB()
}

func openDB() *bolt.DB {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	fmt.Print("opened db")
	return db
}

func Close() {
	log.Println(Conn.Close())
}
