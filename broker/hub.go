package broker

import (
	"context"
	"encoding/json"

	t "github.com/indrenicloud/tricloud-server/core"
)

type Hub struct {
	UserConns  map[uid]*NodeConn // TODO lock if ..
	AgentConns map[uid]*NodeConn

	//UserDevices map[string]uid
	UserAgents map[string]uid // agentkey (not deploy key) with one most current key

	AddConnection    chan *NodeConn
	RemoveConnection chan *NodeConn

	PacketChan chan *packet

	Ctx       context.Context
	CtxCancel context.CancelFunc

	IDGenerator *generator
}

func NewHub() *Hub {

	ctx, ctxcancel := context.WithCancel(context.Background())

	return &Hub{
		UserConns:  make(map[uid]*NodeConn),
		AgentConns: make(map[uid]*NodeConn),

		//UserDevices: make(map[string]uid),
		UserAgents: make(map[string]uid),

		AddConnection:    make(chan *NodeConn),
		RemoveConnection: make(chan *NodeConn),
		PacketChan:       make(chan *packet),
		Ctx:              ctx,
		CtxCancel:        ctxcancel,
		IDGenerator:      newGenerator(),
	}
}

func (h *Hub) Run() {

	for {
		select {
		case _ = <-h.Ctx.Done():
			// cleanup may be
			h.CtxCancel()
			return
		case node := <-h.AddConnection:
			switch node.Type {
			case AgentType:
				//todo connection of this identifier may be present
				// do we remove/close that
				h.AgentConns[node.Connectionid] = node
				h.UserAgents[node.Identifier] = node.Connectionid
			case UserType:
				h.UserConns[node.Connectionid] = node
			}
			go node.Reader()
			go node.Writer()

		case _ = <-h.RemoveConnection:
			//pass
		case receivedPacket := <-h.PacketChan:
			//process
			var msg t.MessageFormat
			json.Unmarshal(receivedPacket.Data, &msg)
			conn, ok := h.UserConns[uid(msg.Receiver)]
			if ok {
				conn.send <- receivedPacket.Data
			}

		}
	}
}
