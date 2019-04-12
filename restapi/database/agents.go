package database

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Agent struct {
	ID         string            `json:"id"`
	Owner      string            `json:"owner"`
	LastLogin  time.Time         `json:"lastlogin,omitempty"`
	FirstAdded time.Time         `json:"firstadded,omitempty"`
	SystemInfo map[string]string `json:"systeminfo,omitempty"`
}

var (
	AgentBucketName = []byte("agents")
)

func NewAgent(owner string) *Agent {

	id, err := uuid.NewRandom()

	if err != nil {
		return nil
	}

	return &Agent{
		ID:         id.String(),
		Owner:      owner,
		FirstAdded: time.Now(),
	}
}

func CreateAgent(owner string) error {

	agent := NewAgent(owner)
	if agent == nil {
		return errors.New("could not create agent")
	}

	agentencoded, err := json.Marshal(agent)
	if err != nil {
		return err
	}

	err = DB.Create([]byte(agent.ID), agentencoded, AgentBucketName)
	if err != nil {
		return err
	}
	return nil
}

func GetAgent(id string) (*Agent, error) {

	agent := &Agent{}
	agentbyte, err := DB.Read([]byte(id), AgentBucketName)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(agentbyte, agent)
	if err != nil {
		return nil, err
	}

	return agent, nil
}

/*
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

func CreateHostInfo(id string, hostinfo *HostInfo) error {
	err := Conn.Update(func(tx *bolt.Tx) error {
		var agent_details Agent
		bk, err := tx.CreateBucketIfNotExists([]byte("agents"))
		if err != nil {
			return fmt.Errorf("Failed to create bucket: %v", err)
		}
		x := tx.Bucket([]byte("agents"))
		agent := x.Get([]byte(id))
		if agent == nil {
			return errors.New("No user with ID " + id + " found")
		}
		json.Unmarshal(agent, &agent_details)
		agent_details.HostInfo = hostinfo
		enc, _ := json.Marshal(agent_details)
		if err := bk.Put([]byte(agent_details.ID), enc); err != nil {
			return fmt.Errorf("Failed to insert '%s'", agent_details.ID)
		}
		return nil
	})
	return err
}

func GetHostInfo(id string) ([]byte, error) {
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
	m["data"] = agent_details.HostInfo
	json_data, err := json.Marshal(m)
	return json_data, err
}

func CreateCPUInfo(id string, cpuinfo *CPUInfo) error {
	err := Conn.Update(func(tx *bolt.Tx) error {
		var agent_details Agent
		bk, err := tx.CreateBucketIfNotExists([]byte("agents"))
		if err != nil {
			return fmt.Errorf("Failed to create bucket: %v", err)
		}
		x := tx.Bucket([]byte("agents"))
		agent := x.Get([]byte(id))
		if agent == nil {
			return errors.New("No user with ID " + id + " found")
		}
		json.Unmarshal(agent, &agent_details)
		agent_details.CPUInfo = cpuinfo
		enc, _ := json.Marshal(agent_details)
		if err := bk.Put([]byte(agent_details.ID), enc); err != nil {
			return fmt.Errorf("Failed to insert '%s'", agent_details.ID)
		}
		return nil
	})
	return err
}

func GetCPUInfo(id string) ([]byte, error) {
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
	m["data"] = agent_details.CPUInfo
	json_data, err := json.Marshal(m)
	return json_data, err
}
*/
