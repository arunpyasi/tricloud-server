package broker

import (
	"github.com/indrenicloud/tricloud-agent/wire"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

func (h *Hub) addConnection(_node *NodeConn) {
	logg.Info("adding connection to hub")
	switch _node.Type {
	case AgentType:
		oldconn, ok := h.AllAgentConns[_node.Connectionid]
		if ok {
			oldconn.close()
		}
		h.AllAgentConns[_node.Connectionid] = _node
		h.ListOfAgents[_node.Identifier] = _node.Connectionid

		go _node.Reader()
		go _node.Writer()

		logg.Info("starting system staticstcs service")
		s := &wire.SysStatCmd{
			Interval: 5,
			Timeout:  0,
		}
		b, err := wire.Encode(
			_node.Connectionid,
			wire.CMD_SYSTEM_STAT,
			wire.BroadcastUsers,
			s,
		)
		if err != nil {
			logg.Warn("Encoding sysemstat cmd error")
		}
		_node.send <- b
		h.signalToUpdade()

	case UserType:
		h.AllUserConns[_node.Connectionid] = _node

		go _node.Reader()
		go _node.Writer()
		ags := &wire.AgentsCountMsg{Agents: make(map[string]wire.UID)}

		for s, id := range h.ListOfAgents {
			ags.Agents[s] = id
		}
		byt, err := wire.Encode(_node.Connectionid, wire.CMD_AGENTS_NO, wire.BroadcastUsers, ags)

		if err != nil {
			logg.Debug("Encoading Mistake 🧩🧩 ")
			logg.Debug(err)
			return
		}

		if err == nil {
			_node.send <- byt
		}

	}
}

func (h *Hub) removeConnection(_node *NodeConn) {
	logg.Debug("Removing magic balls ⚽️🏀⚽️🏀 ")
	if _node.Type == AgentType {
		delete(h.ListOfAgents, _node.Identifier)
		delete(h.AllAgentConns, _node.Connectionid)
		h.signalToUpdade()
	} else if _node.Type == UserType {
		delete(h.AllUserConns, _node.Connectionid)
	}
	h.IDGenerator.free(_node.Connectionid)
	_node.close()
}

func (h *Hub) directSend(pak *DirectPacket) {
	logg.Debug("runnning ")
	switch pak.Ntype {
	case UserType:

	case AgentType:
		//logg.Debug("1")
		logg.Debug(pak.Name)
		id, ok := h.ListOfAgents[pak.Name]
		if !ok {
			return
		}
		//logg.Debug("2")
		conn, ok := h.AllAgentConns[id]
		if !ok {
			return
		}
		logg.Debug("3")
		conn.send <- pak.Body
	}

}
