package broker

import (
	"net/http"

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

type DirectPacket struct {
	Name  string
	Body  []byte
	Ntype NodeType
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type agentsQuery struct {
	responseChan chan map[string]wire.UID
}
