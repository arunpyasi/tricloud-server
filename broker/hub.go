package broker

import (
	"context"
	"log"

	"github.com/indrenicloud/tricloud-agent/wire"
)

type Hub struct {
	AllUserConns  map[wire.UID]*NodeConn // TODO lock if ..
	AllAgentConns map[wire.UID]*NodeConn

	ListOfAgents map[string]wire.UID // agentkey (not deploy key) with one most current connection

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
		AllUserConns:  make(map[wire.UID]*NodeConn),
		AllAgentConns: make(map[wire.UID]*NodeConn),

		ListOfAgents: make(map[string]wire.UID),

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
			log.Println("adding connection to hub")
			switch node.Type {
			case AgentType:
				//todo connection of this identifier may be present
				// do we remove/close that
				h.AllAgentConns[node.Connectionid] = node
				h.ListOfAgents[node.Identifier] = node.Connectionid
			case UserType:
				h.AllUserConns[node.Connectionid] = node
			}
			go node.Reader()
			go node.Writer()

		case _ = <-h.RemoveConnection:
			//pass
		case receivedPacket := <-h.PacketChan:
			log.Println("packet received")

			h.processPacket(receivedPacket)

		}
	}
}

func (h *Hub) processPacket(p *packet) {

	header, _ := wire.GetHeader(p.Data)

	switch header.Flow {
	case wire.AgentToServer, wire.UserToServer:
		h.consumePacket(p, header)
		return
	case wire.UserToAgent:
		h.handleUserPacket(p, header)
	case wire.AgentToUser:
		h.handleAgentPacket(p, header)
	case wire.BroadcastUsers:
		//pass
	default:
		log.Println("Not Implemented")
	}
}

func (h *Hub) consumePacket(pak *packet, header *wire.Header) {

}

func (h *Hub) broadcastUsers(pak *packet, header *wire.Header) {

	for _, conn := range h.AllUserConns {
		conn.send <- pak.Data
	}

}

func (h *Hub) handleUserPacket(pak *packet, header *wire.Header) {

	if header.Connid == 0 {
		log.Println("Don't know where to send packet")
		return
	}
	conn, ok := h.AllAgentConns[header.Connid]
	if !ok {
		log.Println("Agent connection not found")
		return
	}
	conn.send <- pak.Data
}

func (h *Hub) handleAgentPacket(pak *packet, header *wire.Header) {
	if header.Connid == 0 {
		log.Println("msg donot have recevier conn id")
		return
	}

	conn, ok := h.AllUserConns[header.Connid]
	if !ok {
		log.Println("couldnot find connection with id")
		return
	}
	conn.send <- pak.Data

}
