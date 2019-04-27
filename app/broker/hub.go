package broker

import (
	"context"
	"time"

	"github.com/indrenicloud/tricloud-agent/wire"
	"github.com/indrenicloud/tricloud-server/app/database/statstore"
	"github.com/indrenicloud/tricloud-server/app/logg"
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
			logg.Info("adding connection to hub")
			switch node.Type {
			case AgentType:
				oldconn, ok := h.AllAgentConns[node.Connectionid]
				if ok {
					oldconn.close()
				}
				h.AllAgentConns[node.Connectionid] = node
				h.ListOfAgents[node.Identifier] = node.Connectionid
			case UserType:
				h.AllUserConns[node.Connectionid] = node
			}
			go node.Reader()
			go node.Writer()

			s := &wire.SysStatCmd{
				Interval: 5,
				Timeout:  20,
			}

			b, err := wire.Encode(
				node.Connectionid,
				wire.CMD_SYSTEM_STAT,
				wire.BroadcastUsers,
				s,
			)
			if err == nil {
				node.send <- b
			} else {
				logg.Warn("Encoding sysemstat cmd error")
			}

		case nconn := <-h.RemoveConnection:
			if nconn.Type == AgentType {
				delete(h.ListOfAgents, nconn.Identifier)
				delete(h.AllAgentConns, nconn.Connectionid)
			} else if nconn.Type == UserType {
				delete(h.AllUserConns, nconn.Connectionid)
			}

			nconn.close()
		case receivedPacket := <-h.PacketChan:
			logg.Info("packet received")

			h.processPacket(receivedPacket)

		}
	}
}

func (h *Hub) processPacket(p *packet) {

	switch p.head.Flow {
	case wire.UserToAgent:
		h.handleUserPacket(p)
	case wire.AgentToUser:
		h.handleAgentPacket(p)
	case wire.BroadcastUsers:
		h.consumePacket(p)
		return
	default:
		logg.Info("Not Implemented")
	}
}

func (h *Hub) consumePacket(pak *packet) {
	switch pak.head.CmdType {
	case wire.CMD_SYSTEM_STAT:
		go func() {
			statstore.StoreStat(pak.conn.Identifier, time.Now().UnixNano(), pak.body)
		}()
	}
	h.broadcastUsers(pak)

}

func (h *Hub) broadcastUsers(pak *packet) {

	for _, conn := range h.AllUserConns {
		pak.head.Connid = pak.conn.Connectionid
		conn.send <- wire.UpdateHeader(pak.head, pak.rawdata)
	}

}

func (h *Hub) handleUserPacket(pak *packet) {

	if pak.head.Connid == 0 {
		logg.Warn("Don't know where to send packet")
		return
	}
	conn, ok := h.AllAgentConns[pak.head.Connid]
	if !ok {
		logg.Warn("Agent connection not found")
		return
	}
	pak.head.Connid = pak.conn.Connectionid
	conn.send <- wire.UpdateHeader(pak.head, pak.rawdata)
}

func (h *Hub) handleAgentPacket(pak *packet) {
	if pak.head.Connid == 0 {
		logg.Warn("msg donot have recevier conn id")
		return
	}

	conn, ok := h.AllUserConns[pak.head.Connid]
	if !ok {
		logg.Warn("couldnot find connection with id")
		return
	}
	pak.head.Connid = pak.conn.Connectionid
	conn.send <- wire.UpdateHeader(pak.head, pak.rawdata)
}
