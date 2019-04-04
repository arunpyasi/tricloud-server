package database

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
)

type Agent struct {
	ID         string    `json:"id"`
	OS         string    `json:"os,omitempty"`
	LastLogin  string    `json:"lastlogin,omitempty"`
	FirstAdded string    `json:"firstadded,omitempty"`
	Active     bool      `json:"active,omitempty"`
	HostInfo   *HostInfo `json:"hostinfo"`
}
type HostInfo struct {
	Hostname             string `json:"hostname"`
	Uptime               string `json:"uptime"`
	Procs                string `json:"procs"`
	OS                   string `json:"os"`
	Platform             string `json:"platform"`
	PlatformFamily       string `json:"platformfamily"`
	PlatformVersion      string `json:"platformversion"`
	VirtualizationSystem string `json:"virtualizationsystem"`
}

func CreateAgent(agent Agent) error {
	err := Conn.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte("agents"))
		if err != nil {
			return fmt.Errorf("Failed to create bucket: %v", err)
		}
		id, _ := uuid.NewRandom()
		agent.ID = id.String()
		enc, _ := json.Marshal(agent)
		if err := bk.Put([]byte(agent.ID), enc); err != nil {
			return fmt.Errorf("Failed to insert '%s'", agent.ID)
		}
		return nil
	})
	return err
}

func GetAllAgents() ([]byte, error) {
	var agents []Agent
	Conn.View(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("agents"))
		x := tx.Bucket([]byte("agents"))
		c := x.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var data Agent
			json.Unmarshal(v, &data)
			agents = append(agents, data)
		}
		return nil
	})
	m := make(map[string]interface{})
	m["data"] = agents
	json_data, err := json.Marshal(m)
	return json_data, err
}

func GetAgent(id string) ([]byte, error) {
	var agent_details Agent
	err := Conn.View(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("agents"))
		x := tx.Bucket([]byte("agents"))
		agent := x.Get([]byte(id))
		if agent == nil {
			return errors.New("No user with ID " + id + " found")
		}
		json.Unmarshal(agent, &agent_details)
		return nil
	})
	m := make(map[string]interface{})
	m["data"] = agent_details
	json_data, err := json.Marshal(m)
	return json_data, err
}
func UpdateAgent(id string, agent Agent) error {
	err := Conn.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte("agents"))
		if err != nil {
			return fmt.Errorf("Failed to create bucket: %v", err)
		}

		enc, _ := json.Marshal(agent)
		var dec []byte
		json.Unmarshal(enc, &dec)
		if err := bk.Put([]byte(id), enc); err != nil {
			return fmt.Errorf("Failed to update '%v'", agent)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Failed to update : %v", err)
	}
	return nil
}
func DeleteAgent(id string) error {
	err := Conn.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("agents"))
		bk := tx.Bucket([]byte("agents"))
		err := bk.Delete([]byte(id))
		return err
	})
	return err
}
