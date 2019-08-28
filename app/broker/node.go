package broker

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/indrenicloud/tricloud-agent/wire"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

const (
	pongWait   = 30 * time.Second
	pingPeriod = 10 * time.Second
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
	conn         *websocket.Conn
	send         chan []byte
	lastOnline   time.Time
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

	n.lastOnline = time.Now()

	n.conn.SetPongHandler(func(data string) error {
		n.lastOnline = time.Now()
		return nil
	})

	for {
		_, data, err := n.conn.ReadMessage() // todo byte[:read]
		if err != nil {
			logg.Info(err)
			n.MyHub.RemoveConnection <- n
			// todo check the type of error then continue/return depending on it
			return
		}

		head, body := wire.GetHeader(data)

		sendPacket := &packet{
			conn:    n,
			head:    head,
			body:    body,
			rawdata: data,
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

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	defer n.conn.Close()
	for {
		select {
		case _ = <-n.writerCtx.Done():
			//n.conn.WriteControl(websocket.CloseGoingAway, []byte("goodbye"), time.Now())
			return
		case out := <-n.send:
			err := n.conn.WriteMessage(websocket.BinaryMessage, out)

			if err != nil {
				n.MyHub.RemoveConnection <- n
				return
			}
		case t := <-ticker.C:
			logg.Debug("tick.")
			if n.lastOnline.Add(pongWait).After(t) {
				err := n.conn.WriteMessage(websocket.PingMessage, []byte("👍"))
				logg.Debug("ping")
				if err != nil {
					n.MyHub.RemoveConnection <- n
					return
				}
				continue
			}
			logg.Debug("Missed your appoitment bro")
			n.MyHub.RemoveConnection <- n
			return

		}
	}

}

func (n *NodeConn) close() {
	if n.readerCtx.Err() == nil {
		n.CloseReader()
	}
	if n.writerCtx.Err() == nil {
		n.CloseWriter()
	}
	logg.Debug("Closing conn")
	n.conn.Close()

}
