package database

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	bolt "go.etcd.io/bbolt"
)

type User struct {
	ID       string   `json:"id"`
	Password string   `json:"password,omitempty"`
	FullName string   `json:"fullname,omitempty"`
	Email    string   `json:"email,omitempty"`
	Active   string   `json:"active,omitempty"`
	Agents   []string `json:"agents"`
}

type Agent struct {
	ID         string `json:"id"`
	OS         string `json:"os,omitempty"`
	LastLogin  string `json:"lastlogin,omitempty"`
	FirstAdded string `json:"firstadded,omitempty"`
	Active     bool   `json:"active,omitempty"`
}

type SysLog struct {
	AgentID string `json:"agentid"`
	MemInfo string `json:"meminfo"`
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
		return nil
	})
	return err
}
func GetUser(id string) ([]byte, error) {
	var user_details User
	err := Conn.View(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("users"))
		x := tx.Bucket([]byte("users"))
		user := x.Get([]byte(id))
		if user == nil {
			return errors.New("No user with ID " + id + " found")
		}
		json.Unmarshal(user, &user_details)
		return nil
	})
	m := make(map[string]interface{})
	m["data"] = user_details
	json_data, err := json.Marshal(m)
	return json_data, err
}

func GetAllUsers() ([]byte, error) {
	var users []User
	Conn.View(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("users"))
		x := tx.Bucket([]byte("users"))
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

func UpdateUser(id string, user User) error {
	err := Conn.Update(func(tx *bolt.Tx) error {
		bk, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return fmt.Errorf("Failed to create bucket: %v", err)
		}
		enc, _ := json.Marshal(user)
		var dec []byte
		json.Unmarshal(enc, &dec)
		if err := bk.Put([]byte(id), enc); err != nil {
			return fmt.Errorf("Failed to update '%s'", user)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Failed to update : %v", err)
	}
	return nil
}

func DeleteUser(id string) error {
	err := Conn.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte("users"))
		err := bk.Delete([]byte(id))
		return err
	})
	return err
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
			return fmt.Errorf("Failed to update '%s'", agent)
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
