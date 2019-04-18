package broker

import "github.com/gorilla/websocket"

type NodeType byte

const (
	UserType NodeType = iota
	AgentType
)

type packet struct {
	Conn *NodeConn
	Data []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
