package broker

import (
	"time"

	"github.com/indrenicloud/tricloud-server/app/database/statstore"
)

// this consumes packet that is, saves log to database
// unmarsel them to check agent vitals
// and emits Events/notification if needed
func (h *Hub) consumePacket(pak *packet) {

	// we are in seperate coroutine than hub
	// so should not modify hub 's state that might
	// result race condition
	t := time.Now()
	statstore.StoreStat(pak.conn.Identifier, t.UnixNano(), pak.body)
	h.event.ProcessLog(pak.conn.Identifier, h.userName, pak.rawdata, h.BroadCastEvent)
}
