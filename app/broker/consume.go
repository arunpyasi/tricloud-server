package broker

import (
	"time"

	"github.com/indrenicloud/tricloud-agent/wire"
	"github.com/indrenicloud/tricloud-server/app/database/statstore"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

const (
	MemHightThreshold = 20
	CPUHightThreshold = 90
	EventTimeout      = 30
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

	h.lEvent.Lock()
	oldt, ok := h.eventTimestampLog[pak.conn.Identifier]
	if ok {
		t2 := oldt.Add(time.Duration(EventTimeout) * time.Second)
		if t2.Before(t) {
			h.lEvent.Unlock()
			return
		}
	}
	h.eventTimestampLog[pak.conn.Identifier] = t
	h.lEvent.Unlock()

	sd := wire.SysStatData{}
	wire.Decode(pak.rawdata, &sd)

	Events := make(map[string]string)

	memper := ((sd.TotalMem - sd.AvailableMem) * 100) / (sd.TotalMem)

	if memper > MemHightThreshold {
		Events["Memory"] = "LOW"
		logg.Info("Nearly out of memory")

	}

	cpuper := float64(0)

	for _, i := range sd.CPUPercent {
		cpuper = cpuper + i
	}
	cpuper = cpuper / float64(len(sd.CPUPercent))

	if cpuper > CPUHightThreshold {
		Events["CPU"] = "HIGH"
	}

	if len(Events) > 0 {

		h.event.SendEvent(h.userName, Events)
	}

}
