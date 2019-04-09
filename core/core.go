package core

import (
	"encoding/json"
	"log"
)

// CommandType could be byte since we won't have more than 255 commands
type CommandType int
type UID uint16

const (
	CMD_SERVER_HELLO CommandType = iota
	CMD_LIST_AGENTS
	CMD_SERVER_MAX
	CMD_SYSTEM_INFO //= CMD_SERVER_MAX + 1
	CMD_TERMINAL
	CMD_BUILTIN_MAX
)

// CMD_ALL_MAX future use :wink
var CMD_ALL_MAX CommandType = CMD_BUILTIN_MAX

// MessageFormat is core message format
type MessageFormat struct {
	ReceiverConnid   UID               `json:"rcid,omitempty"`
	ReceiverIdentity string            `json:"rident,omitempty"`
	CmdType          CommandType       `json:"cmdtype,omitempty"`
	Arguments        map[string]string `json:"args,omitempty"`
	Results          []string          `json:"results,omitempty"`
}

func (m *MessageFormat) GetBytes() []byte {

	outByte, err := json.Marshal(m)
	if err != nil {
		log.Fatal("could not marsal msg")
		return nil
	}

	return outByte
}
