package database

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/indrenicloud/tricloud-server/app/config"
	"github.com/indrenicloud/tricloud-server/app/logg"

	bolt "go.etcd.io/bbolt"
)

var DB = &Boltdb{}

func Start() {

}

func Close() {
	DB.Close()
}

func init() {
	dev := config.GetConfig().Dev

	firstrun := false
	if _, err := os.Stat(config.GetConfig().DBpath); os.IsNotExist(err) {
		firstrun = true
	}

	err := DB.Open(config.GetConfig().DBpath)
	if err != nil {
		panic(err)
	}

	if dev {
		devMigration(firstrun)
	} else {
		normalMigration(firstrun)
	}

}

// if it is devenvironment create some fake users and agents for testing
func devMigration(firstrun bool) {

	userinfo := map[string]string{
		"id":       "batman47",
		"password": "hard123",
		"fullname": "Batman Kickass",
		"email":    "batman47@gentelmanclub.com",
	}

	if firstrun {
		initilizeBuckets([][]byte{UserBucketName, AgentBucketName})
		usr, err := NewUser(userinfo, true)
		logg.Info("Creating a demo user")
		logg.Info(userinfo)

		CreateUser(usr)
		if err != nil {
			log.Println(err)
		}

		AddapiKey(usr.ID, "agent")
	}

	user, err := GetUser(userinfo["id"])
	if err != nil {
		logg.Error("didnot create user :(")
		panic(err)
	}

	logg.Info("Demo api key")
	logg.Info(user.APIKeys)

	agentsbyte, _ := DB.ReadAll(AgentBucketName)

	// logging the all agents
	if len(agentsbyte) > 0 {
		logg.Info("All agents in db:")
	}
	for _, agentbyte := range agentsbyte {
		agent := &Agent{}
		Decode(agentbyte, agent)
		logg.Info(agent)
	}
}

// else just make sure essential buckets are created
func normalMigration(firstrun bool) {
	initilizeBuckets([][]byte{UserBucketName, AgentBucketName})
}

func initilizeBuckets(buckets [][]byte) {
	for _, b := range buckets {
		err := DB.conn.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists(b)
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			logg.Error("err in initilizing buckets")
		}
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
