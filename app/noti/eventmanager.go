package noti

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/indrenicloud/tricloud-agent/wire"
	alertdb "github.com/indrenicloud/tricloud-server/app/database/statstore"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

const (
	MemHightThreshold = 20
	CPUHightThreshold = 50
	EventTimeout      = 300 * time.Second
)

type EventManager struct {
	cs             CredentialStore //its should be threadsafe
	eventProviders []Provider

	eTimestamp map[string]map[EventType]time.Time
	lEvent     sync.Mutex
}

func NewEventManager() *EventManager {
	em := new(EventManager)
	em.cs = NewCredStore()
	em.eventProviders = notificationProviders(em.cs)
	em.eTimestamp = make(map[string]map[EventType]time.Time)

	for _, ep := range em.eventProviders {
		ep.Init()
	}

	return em
}

func (e *EventManager) ProcessLog(agentName, username string,
	rawlog []byte,
	ch chan<- []byte) {

	logg.Debug("Lets check if this log needs to emit events")
	t := time.Now()

	e.lEvent.Lock()
	defer e.lEvent.Unlock()

	sd := wire.SysStatData{}
	wire.Decode(rawlog, &sd)

	events := newEventContainer(agentName, t.UnixNano(), t.Format("Mon Jan _2 2006 15:04:05 "))

	oldrecord, isthereOldRecord := e.eTimestamp[agentName]

	checkTimeout := func(et EventType) bool {
		if !isthereOldRecord {
			return false
		}
		oldt, ok := oldrecord[et]
		if ok {
			t2 := oldt.Add(EventTimeout)
			if t.Before(t2) {
				return true
			}
		}
		return false
	}

	setTime := func(et EventType) {
		if oldrecord == nil {
			oldrecord = make(map[EventType]time.Time)
		}
		oldrecord[et] = t
		e.eTimestamp[agentName] = oldrecord
	}

	// mem usage
	memper := ((sd.TotalMem - sd.AvailableMem) * 100) / (sd.TotalMem)
	if memper > MemHightThreshold {

		if checkTimeout(MemorySpike) {
			goto AFTERMEMCHECK
		}

		events.Events = append(events.Events,
			Event{Type: MemorySpike,
				Message: fmt.Sprintf("%d", memper),
			})

		setTime(MemorySpike)

	} else {
		if checkTimeout(MemorySpike) {
			events.Events = append(events.Events,
				Event{Type: MemorySpikeOver,
					Message: fmt.Sprintf("%4d", memper)})

			delete(oldrecord, MemorySpike)
		}

	}

AFTERMEMCHECK:
	//cpu usage
	cpuper := float64(0)
	for _, i := range sd.CPUPercent {
		cpuper = cpuper + i
	}
	cpuper = cpuper / float64(len(sd.CPUPercent))

	if cpuper > CPUHightThreshold {
		if checkTimeout(CPUSpike) {
			goto AFTERCPUCHECK
		}

		events.Events = append(events.Events,
			Event{Type: CPUSpike,
				Message: fmt.Sprintf("%.6f", cpuper),
			})
		setTime(CPUSpike)

	} else {
		if checkTimeout(CPUSpike) {
			events.Events = append(events.Events,
				Event{Type: CPUSpikeOver,
					Message: fmt.Sprintf("%.6f", cpuper),
				})
		}
		delete(oldrecord, CPUSpike)

	}
AFTERCPUCHECK:

	if len(events.Events) == 0 {
		return
	}
	e.sendEvent(username, events, ch)

}

func (e *EventManager) sendEvent(user string, ec *EventContainer, ch chan<- []byte) {
	eventsbyte, err := Encode(ec)
	if err != nil {
		return
	}
	go func() {
		h := wire.NewHeader(wire.UID(0), wire.CMD_EVENTS, wire.BroadcastUsers)
		ch <- wire.AttachHeader(h, eventsbyte)
	}()

	go alertdb.StoreAlert(eventsbyte, []byte(user), ec.Timestamp)

	e.SendEvent(user, string(eventsbyte))

}

func (e *EventManager) SendEvent(user string, _data string) {

	var wg sync.WaitGroup

	for _, ee := range e.eventProviders {
		tokens := e.cs.GetToken(ee.GetName(), user)
		logg.Info("Notification Loop")
		logg.Info(ee.GetName())
		logg.Info(tokens)
		for _, token := range tokens {
			go func(_ee Provider, t string) {
				wg.Add(1)
				logg.Info("Inside Loop 👽👽👽")
				logg.Info(_ee.GetName())
				_ee.PushNotification(context.Background(), t, _data)
				wg.Done()
			}(ee, token)

		}
	}
	wg.Wait()

}

func (e *EventManager) SaveToken(user string, token string) {
	// TODO
	e.cs.SetToken("firebase", user, token)
}
