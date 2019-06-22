package broker

import (
	"github.com/indrenicloud/tricloud-agent/wire"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

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
	logg.Debug("They are here alert everyone 👮‍👮‍‍👮‍‍👮‍‍👮‍")

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
