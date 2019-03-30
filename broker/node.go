package broker

import (
	"context"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// NodeConn is used to represent both userconnection and agent connection
type NodeConn struct {
	Connectionid uuid.UUID
	Identifier   string   // key for agent userid for user
	Agents       []string // only applicable for Usertype
	Parent       string   //only applicable for Parent
	Type         NodeType // UserType or AgentType
	readerCtx    *context.Context
	CloseReader  context.CancelFunc
	writerCtx    *context.Context
	CloseWriter  context.CancelFunc

	conn *websocket.Conn
	send chan []byte
}
