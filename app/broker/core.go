package broker

import (
	"github.com/gorilla/websocket"
	"github.com/indrenicloud/tricloud-agent/wire"
)

type NodeType byte

const (
	UserType NodeType = iota
	AgentType
)

type packet struct {
	conn    *NodeConn
	head    *wire.Header
	body    []byte
	rawdata []byte //all bytes including header bytes
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type agentsQuery struct {
	responseChan chan map[string]wire.UID
}
