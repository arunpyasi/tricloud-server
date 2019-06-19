package noti

import (
	"encoding/json"
)

type EventType string

const (
	AgentDown          EventType = "Agent down"
	AgentUp                      = "Agent up"
	MemorySpike                  = "Memory spike"
	MemorySpikeOver              = "Memory spike over"
	CPUSpike                     = "cpu spike"
	CPUSpikeOver                 = "cpu spike over "
	BandwithUsageLimit           = "Bandwidth usage limit exceed"
	NetworkSpike                 = "Network spike"
	NetworkSpikeOver             = "Network spike over"
	ServiceCrashed               = "Service Crashed"
	ServiveRunning               = "Service running"
	ProcessCrashed               = "process crashed"
	ProcessRunning               = "process running"
)

type EventContainer struct {
	Agentid string
	Events  []Event

	Timestamp int64
	Timestr   string
}

type Event struct {
	Type    EventType
	Message string
}

func newEventContainer(agentid string, timestamp int64, timestr string) *EventContainer {
	return &EventContainer{
		Agentid:   agentid,
		Timestamp: timestamp,
		Timestr:   timestr,
	}
}

func Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Decode(raw []byte, out interface{}) error {
	return json.Unmarshal(raw, out)
}
