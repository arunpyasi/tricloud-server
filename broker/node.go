package broker

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
	"github.com/indrenicloud/tricloud-agent/wire"
)

// NodeConn is used to represent both userconnection and agent connection
type NodeConn struct {
	Connectionid wire.UID
	Identifier   string   // key for agent userid for user
	Type         NodeType // UserType or AgentType
	readerCtx    context.Context
	CloseReader  context.CancelFunc
	writerCtx    context.Context
	CloseWriter  context.CancelFunc
	MyHub        *Hub
	Running      bool

	conn *websocket.Conn
	send chan []byte
}

func NewNodeConn(identifier string, t NodeType, conn *websocket.Conn, h *Hub) *NodeConn {

	rctx, rctxCancelFunc := context.WithCancel(h.Ctx)
	wctx, wctxCancelFunc := context.WithCancel(h.Ctx)

	return &NodeConn{
		Connectionid: h.IDGenerator.generate(),
		Identifier:   identifier,
		Type:         t,
		readerCtx:    rctx,
		CloseReader:  rctxCancelFunc,
		writerCtx:    wctx,
		CloseWriter:  wctxCancelFunc,
		MyHub:        h,
		Running:      true,
		conn:         conn,
		send:         make(chan []byte),
	}
}

func (n *NodeConn) Reader() {

	for {
		_, data, err := n.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			// todo check the type of error then continue/return depending on it
			return
		}
		sendPacket := &packet{
			Conn: n,
			Data: data,
		}

		n.MyHub.PacketChan <- sendPacket

		select {
		case _ = <-n.readerCtx.Done():
			return
		default:
		}

	}

}

func (n *NodeConn) Writer() {
	defer n.conn.Close()
	for {
		select {
		case _ = <-n.writerCtx.Done():
			return
		case out := <-n.send:
			err := n.conn.WriteMessage(websocket.TextMessage, out)

			if err != nil {
				return
			}
		}
	}

}
