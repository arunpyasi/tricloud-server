package database

import (
	"encoding/json"
	"log"
	"os"

	bolt "go.etcd.io/bbolt"
)

var DB = &Boltdb{}

func init() {
	path := "mybolt.db"
	runmigration := false
	if _, err := os.Stat(path); os.IsNotExist(err) {
		runmigration = true
	}

	DB.Open(path)

	if runmigration {
		RunMigration()
	}
}

func RunMigration() {

	// TODO get this from config or ENV
	dev := true

	// if it is devenvironment create some fake users and agents for testing
	// else just make sure essential buckets are created

	if dev {
		devMigration()
	} else {
		normalMigration()
	}

}

func devMigration() {
	err := CreateUser(map[string]string{
		"id":       "batman47",
		"password": "hard123",
		"fullname": "Batman Kickass",
		"email":    "batman47@gentelmanclub.com",
	})

	if err != nil {
		log.Println(err)
	}

	err = CreateAgent("batman47")
	if err != nil {
		log.Println(err)
	}

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

func normalMigration() {
	err := DB.conn.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists([]byte("agents"))
		if err != nil {
			return err
		}
		return nil

	})

	if err != nil {
		log.Println(err)
	}
}
