package broker

import (
	"context"
	"encoding/json"
	"log"

	"github.com/indrenicloud/tricloud-server/core"
)

type Hub struct {
	AllUserConns  map[core.UID]*NodeConn // TODO lock if ..
	AllAgentConns map[core.UID]*NodeConn

	ListOfAgents map[string]core.UID // agentkey (not deploy key) with one most current connection

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
		AllUserConns:  make(map[core.UID]*NodeConn),
		AllAgentConns: make(map[core.UID]*NodeConn),

		ListOfAgents: make(map[string]core.UID),

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

			h.ProcessPacket(receivedPacket)

		}
	}
}

func (h *Hub) ProcessPacket(p *packet) {

	var msg core.MessageFormat
	json.Unmarshal(p.Data, &msg)

	if msg.CmdType > core.CMD_SERVER_MAX {

		switch p.Conn.Type {
		case AgentType:
			h.ProcessAgentPacket(p, &msg)
		case UserType:
			h.ProcessUserPacket(p, &msg)
		}
		return
	}

	//server packet process here

	switch msg.CmdType {
	case core.CMD_LIST_AGENTS:
		agents := make([]string, 10)

		for k := range h.ListOfAgents {
			agents = append(agents, k)
		}
		reciver := p.Conn.Connectionid
		response := core.MessageFormat{

			CmdType: core.CMD_LIST_AGENTS,
			Results: agents,
		}

		out := response.GetBytes()
		conn, ok := h.AllUserConns[reciver]

		if !ok {
			log.Println("sending conn not found")
		}
		log.Println("sending to user")
		conn.send <- out

	}

}

func (h *Hub) ProcessUserPacket(p *packet, msg *core.MessageFormat) {

	if msg.ReceiverConnid == 0 {
		if msg.ReceiverIdentity == "" {
			log.Println("Don't know where to send packet")
			return
		}
		//msg.ReceiverConnid
		id, ok := h.ListOfAgents[msg.ReceiverIdentity]
		if !ok {
			log.Println("Don't have agent of that identity")
			return
		}
		conn, ok := h.AllAgentConns[id]
		if !ok {
			log.Println("Agent connection not found")
			return
		}
		msg.ReceiverConnid = p.Conn.Connectionid
		conn.send <- msg.GetBytes()

	}

}

func (h *Hub) ProcessAgentPacket(p *packet, msg *core.MessageFormat) {

	if msg.ReceiverConnid == 0 {
		log.Println("msg donot have recevier conn id")
		return
	}
	conn, ok := h.AllUserConns[msg.ReceiverConnid]
	if !ok {
		log.Println("couldnot find connection with id")
		return
	}
	conn.send <- p.Data

}
