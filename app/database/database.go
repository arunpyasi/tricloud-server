package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	bolt "go.etcd.io/bbolt"
)

var DB = &Boltdb{}

func Start() {

}

func Close() {
	DB.Close()
}

func init() {
	// TODO get this from config or ENV
	dev := true
	path := "mybolt.db"

	var runmigration bool
	if _, err := os.Stat(path); os.IsNotExist(err) {
		runmigration = true
	}

	err := DB.Open(path)
	if err != nil {
		panic(err)
	}

	if !runmigration {
		return
	}

	if dev {
		devMigration()
	} else {
		normalMigration()
	}
}

// if it is devenvironment create some fake users and agents for testing
func devMigration() {
	usr, err := NewUser(map[string]string{
		"id":       "batman47",
		"password": "hard123",
		"fullname": "Batman Kickass",
		"email":    "batman47@gentelmanclub.com",
	}, true)

	CreateUser(usr)

	if err != nil {
		log.Println(err)
	}

	AddapiKey("batman47", "agent")

	agentsbyte, err := DB.ReadAll(AgentBucketName)

	if err != nil {
		log.Println(err)
		return
	}

	// logging the all agents
	for _, agentbyte := range agentsbyte {
		agent := &Agent{}
		json.Unmarshal(agentbyte, agent)
		log.Println(agent.ID, "::", agent.Owner)
	}
}

// else just make sure essential buckets are created
func normalMigration() {
	err := DB.conn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(UserBucketName)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(AgentBucketName)
		if err != nil {
			return err
		}
		return nil

	})

	if err != nil {
		log.Println(err)
	}
}

func checkFields(info map[string]string, fields []string) error {

	if len(info) > len(fields) {
		return errors.New("More than required fields")
	}

	if len(info) > len(fields) {
		return errors.New("Less than required fields")
	}

	for _, field := range fields {

		value, ok := info[field]
		if !ok {
			return fmt.Errorf("field %s not found", value)
		}
	}
	return nil
}
