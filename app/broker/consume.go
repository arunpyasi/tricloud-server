package broker

import (
	"fmt"
	"time"

	"github.com/indrenicloud/tricloud-agent/wire"
	"github.com/indrenicloud/tricloud-server/app/database/statstore"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

const (
	MemHightThreshold = 20
	CPUHightThreshold = 90
	EventTimeout      = 300 * time.Second
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
	defer h.lEvent.Unlock()
	oldt, ok := h.eventTimestampLog[pak.conn.Identifier]
	fmt.Println(oldt)
	if ok {
		fmt.Println("📠")
		t2 := oldt.Add(EventTimeout)
		if t.Before(t2) {
			return
		}
	}

	sd := wire.SysStatData{}
	wire.Decode(pak.rawdata, &sd)

	Events := make(map[string]string)

	memper := ((sd.TotalMem - sd.AvailableMem) * 100) / (sd.TotalMem)

	if memper > MemHightThreshold {
		Events["memory"] = "high"
		Events["mem_per"] = fmt.Sprintf("%d", memper)
		logg.Info("Nearly out of memory")

	}

	cpuper := float64(0)

	for _, i := range sd.CPUPercent {
		cpuper = cpuper + i
	}
	cpuper = cpuper / float64(len(sd.CPUPercent))
	//time.Microsecond

	if cpuper > CPUHightThreshold {
		Events["cpu"] = "high"
		Events["cpu_per"] = fmt.Sprintf("%.6f", cpuper)
	}

	if len(Events) > 0 {
		h.eventTimestampLog[pak.conn.Identifier] = t
		Events["agent"] = pak.conn.Identifier
		Events["user"] = h.userName
		Events["type"] = "resurce_spike"
		fmt.Printf("%+v", h.eventTimestampLog)

		h.event.SendEvent(h.userName, Events)
	}

}
