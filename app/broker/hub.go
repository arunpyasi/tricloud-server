package broker

import (
	"context"

	"github.com/indrenicloud/tricloud-agent/wire"
	"github.com/indrenicloud/tricloud-server/app/logg"
	"github.com/indrenicloud/tricloud-server/app/noti"
)

type Hub struct {
	userName string // hub is run per user basis

	AllUserConns  map[wire.UID]*NodeConn // TODO lock if ..
	AllAgentConns map[wire.UID]*NodeConn

	ListOfAgents map[string]wire.UID // agentkey (not deploy key) with one most current connection

	AddConnection    chan *NodeConn
	RemoveConnection chan *NodeConn

	PacketChan      chan *packet
	queryAgentsChan chan *agentsQuery
	removeagentChan chan string

	// broadcast new agent info to all user conn or decrease
	broadcastAgentsToUser chan struct{}

	Ctx       context.Context
	CtxCancel context.CancelFunc

	IDGenerator *generator

	event *noti.EventManager

	BroadCastEvent chan []byte
}

func NewHub(ctx context.Context, e *noti.EventManager, user string) *Hub {

	ctx1, ctxcancel := context.WithCancel(ctx)

	return &Hub{
		userName:      user,
		AllUserConns:  make(map[wire.UID]*NodeConn),
		AllAgentConns: make(map[wire.UID]*NodeConn),

		ListOfAgents: make(map[string]wire.UID),

		AddConnection:         make(chan *NodeConn),
		RemoveConnection:      make(chan *NodeConn),
		PacketChan:            make(chan *packet),
		queryAgentsChan:       make(chan *agentsQuery),
		removeagentChan:       make(chan string),
		broadcastAgentsToUser: make(chan struct{}),
		BroadCastEvent:        make(chan []byte),
		Ctx:                   ctx1,
		CtxCancel:             ctxcancel,
		IDGenerator:           newGenerator(),
		event:                 e,
	}
}

func (h *Hub) Run() {
	go h.debugg()

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

				go node.Reader()
				go node.Writer()

				logg.Info("starting system staticstcs service")
				s := &wire.SysStatCmd{
					Interval: 5,
					Timeout:  0,
				}
				b, err := wire.Encode(
					node.Connectionid,
					wire.CMD_SYSTEM_STAT,
					wire.BroadcastUsers,
					s,
				)
				if err != nil {
					logg.Warn("Encoding sysemstat cmd error")
				}
				node.send <- b
				h.signalToUpdade()

			case UserType:
				h.AllUserConns[node.Connectionid] = node

				go node.Reader()
				go node.Writer()
				ags := &wire.AgentsCountMsg{Agents: make(map[string]wire.UID)}

				for s, id := range h.ListOfAgents {
					ags.Agents[s] = id
				}
				byt, err := wire.Encode(node.Connectionid, wire.CMD_AGENTS_NO, wire.BroadcastUsers, ags)

				if err != nil {
					logg.Debug("Encoading Mistake 🧩🧩 ")
					logg.Debug(err)
					return
				}

				if err == nil {
					node.send <- byt
				}

			}

		case nconn := <-h.RemoveConnection:
			logg.Debug("Removing magic balls ⚽️🏀⚽️🏀 ")
			if nconn.Type == AgentType {
				delete(h.ListOfAgents, nconn.Identifier)
				delete(h.AllAgentConns, nconn.Connectionid)
				h.signalToUpdade()
			} else if nconn.Type == UserType {
				delete(h.AllUserConns, nconn.Connectionid)
			}
			nconn.close()

		case receivedPacket := <-h.PacketChan:
			logg.Info("packet received")

			h.processPacket(receivedPacket)
		case q := <-h.queryAgentsChan:
			activeagents := make(map[string]wire.UID)
			for key, val := range h.ListOfAgents {
				activeagents[key] = val
			}
			q.responseChan <- activeagents
		case agentid := <-h.removeagentChan:
			// this is used when agent is deleted form ui
			// it has to forcefully removed
			agentconid, ok := h.ListOfAgents[agentid]
			if !ok {
				break
			}

			conn, ok := h.AllAgentConns[agentconid]
			if ok {
				e := &wire.Exit{}
				b, _ := wire.Encode(conn.Connectionid, wire.CMD_EXIT, wire.DefaultFlow, e)
				conn.send <- b
				logg.Warn("removing agent from hub")
				go func() {
					h.RemoveConnection <- conn
				}()

			}
		case <-h.broadcastAgentsToUser:
			h.broadcastAgentsInfo()
		case bt := <-h.BroadCastEvent:
			logg.Debug("They are here alert everyone 👮‍👮‍‍👮‍‍👮‍‍👮‍")
			h.broadcastEvent(bt)
		}

	}
	logg.Info("hub exitting")
}

func (h *Hub) signalToUpdade() {
	go func() {
		h.broadcastAgentsToUser <- struct{}{}
	}()
}

func (h *Hub) broadcastAgentsInfo() {
	logg.Debug("Broadcasting👏 ")
	ags := &wire.AgentsCountMsg{Agents: make(map[string]wire.UID)}
	for s, id := range h.ListOfAgents {
		ags.Agents[s] = id

	}

	byt, err := wire.Encode(wire.UID(0), wire.CMD_AGENTS_NO, wire.BroadcastUsers, ags)
	if err != nil {
		logg.Debug("Encoading Mistake 🧩🧩 ")
		logg.Debug(err)
		return
	}
	for _, conn := range h.AllUserConns {
		logg.Debug("👏 ")
		select {
		case conn.send <- byt:
		default:
		}
	}

}

func (h *Hub) processPacket(p *packet) {

	if p.head.CmdType == wire.CMD_GCM_TOKEN {
		t := wire.TokenMessage{}
		wire.Decode(p.rawdata, &t)
		h.event.SaveToken(p.conn.Identifier, t.Token)
		return
	}

	switch p.head.Flow {
	case wire.UserToAgent:
		h.handleUserPacket(p)
	case wire.AgentToUser:
		h.handleAgentPacket(p)
	case wire.BroadcastUsers:
		if p.head.CmdType == wire.CMD_SYSTEM_STAT {
			go h.consumePacket(p)
		}
		h.broadcastUsers(p)
		return
	default:
		logg.Info("Not Implemented")
	}
}

func (h *Hub) broadcastUsers(pak *packet) {

	pak.head.Connid = pak.conn.Connectionid
	updatedbytes := wire.UpdateHeader(pak.head, pak.rawdata)

	for _, conn := range h.AllUserConns {

		select {
		case conn.send <- updatedbytes:
		default:
		}
	}

}

func (h *Hub) broadcastEvent(bt []byte) {

	for _, conn := range h.AllUserConns {

		select {
		case conn.send <- bt:
		default:
		}
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
