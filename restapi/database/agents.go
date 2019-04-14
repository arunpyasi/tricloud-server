package database

import (
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

func CreateAgent(owner string) (string, error) {

	user, err := GetUser(owner)

	if err != nil {
		return "", errors.New("cannot create agent, user doesnot exist")
	}

	agent := NewAgent(owner)
	if agent == nil {
		return "", errors.New("could not create agent")
	}

	user.Agents = append(user.Agents, agent.ID)

	agentencoded, err := Encode(agent)
	if err != nil {
		return "", err
	}

	userencoded, err := Encode(user)
	if err != nil {
		return "", err
	}

	DB.Update([]byte(user.ID), userencoded, UserBucketName)

	err = DB.Create([]byte(agent.ID), agentencoded, AgentBucketName)
	if err != nil {
		return "", err
	}
	return agent.ID, nil
}

func GetAgent(id string) (*Agent, error) {

	agent := &Agent{}
	agentbyte, err := DB.Read([]byte(id), AgentBucketName)

	if err != nil {
		return nil, err
	}

	err = Decode(agentbyte, agent)
	if err != nil {
		return nil, err
	}

	return agent, nil
}
func GetAllUserAgents(id string) ([]*Agent, error) {

	var agents []*Agent

	user, err := GetUser(id)
	if err != nil {
		return nil, err
	}

	for _, val := range user.Agents {
		agent, err := GetAgent(val)
		if err != nil {
			continue
		}
		agents = append(agents, agent)
	}
	return agents, nil
}

func GetAllAgents() ([]*Agent, error) {

	var agents []*Agent

	agentsbyte, err := DB.ReadAll(AgentBucketName)
	if err != nil {
		return nil, err
	}

	for _, val := range agentsbyte {
		var agent Agent
		err = Decode(val, &agent)
		if err != nil {
			return nil, err
		}
		agents = append(agents, &agent)
	}
	return agents, nil
}

func UpdateAgent(agent *Agent) error {

	agentbyte, err := Encode(agent)
	if err != nil {
		return err
	}

	return DB.Update([]byte(agent.ID), agentbyte, AgentBucketName)
}

func UpdateSystemInfo(id string, systeminfo map[string]string) error {

	agent, err := GetAgent(id)
	if err != nil {
		return err
	}

	err = checkFields(systeminfo, []string{
		"hostname", "uptime", "procs", "os", "platform", "platformfamily", "platformversion",
		"virtualizationsystem", "cpu", "vendorid", "family", "model", "physicalid", "coreid",
		"cores", "modelname", "mhz", "cachesize", "flags", "microcode",
	})
	if err != nil {
		return err
	}

	agent.SystemInfo = systeminfo
	return UpdateAgent(agent)
}

func GetSysteminfo(id string) (map[string]string, error) {

	agent, err := GetAgent(id)
	if err != nil {
		return nil, err
	}

	return agent.SystemInfo, nil
}

func DeleteAgent(id string) error {
	return DB.Delete([]byte(id), AgentBucketName)
}
