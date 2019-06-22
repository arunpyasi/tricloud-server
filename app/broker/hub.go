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

	CDirectSend chan *DirectPacket

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
		CDirectSend:           make(chan *DirectPacket),
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
			h.addConnection(node)
		case nconn := <-h.RemoveConnection:
			h.removeConnection(nconn)

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
		case dp := <-h.CDirectSend:
			h.directSend(dp)
		case bt := <-h.BroadCastEvent:
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
