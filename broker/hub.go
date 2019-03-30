package broker

import (
	"context"

	"github.com/google/uuid"
)

type Hub struct {
	UserConns  map[uuid.UUID]*NodeConn
	AgentConns map[uuid.UUID]*NodeConn

	UserDevices map[string]uuid.UUID
	UserAgents  map[string]uuid.UUID

	AddConnection    chan *NodeConn
	RemoveConnection chan *NodeConn

	PacketChan chan *packet

	Ctx       context.Context
	CtxCancel context.CancelFunc
}

func NewHub() *Hub {

	ctx, ctxcancel := context.WithCancel(context.Background())

	return &Hub{
		UserConns:  make(map[uuid.UUID]*NodeConn),
		AgentConns: make(map[uuid.UUID]*NodeConn),

		UserDevices: make(map[string]uuid.UUID),
		UserAgents:  make(map[string]uuid.UUID),

		AddConnection:    make(chan *NodeConn),
		RemoveConnection: make(chan *NodeConn),
		PacketChan:       make(chan *packet),
		Ctx:              ctx,
		CtxCancel:        ctxcancel,
	}
}

func (h *Hub) Run() {

	for {
		select {
		case _ = <-h.Ctx.Done():
			// cleanup & exit here
		case node := <-h.AddConnection:
			//pass
		case node := <-h.RemoveConnection:
			//pass
		case receivedPacket := <-h.PacketChan:

			go func() {
				// packet routing
			}

		}
	}
}
