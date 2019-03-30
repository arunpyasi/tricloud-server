package broker

import (
	"context"

	"github.com/gorilla/websocket"
)

// NodeConn is used to represent both userconnection and agent connection
type NodeConn struct {
	Connectionid uid
	Identifier   string   // key for agent userid for user
	Type         NodeType // UserType or AgentType
	readerCtx    *context.Context
	CloseReader  context.CancelFunc
	writerCtx    *context.Context
	CloseWriter  context.CancelFunc
	MyHub        *Hub
	Running      bool

	conn *websocket.Conn
	send chan []byte
}

func NewNodeConn(identifier string, t NodeType, conn *websocket.Conn, h *Hub) *NodeConn {
	return nil
}

func (n *NodeConn) Reader() {

}

func (n *NodeConn) Writer() {

}
