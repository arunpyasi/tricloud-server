package script

import (
	"time"

	"github.com/indrenicloud/tricloud-agent/wire"
	"github.com/indrenicloud/tricloud-server/app/broker"
	"github.com/indrenicloud/tricloud-server/app/database"
)

type ScriptManager struct {
	cAddScript    chan *database.Script
	cRemoveScript chan string
	cRunScript    chan *database.Script
	//cUpdateScript chan *database.Script
	scripts map[string]*database.Script
	broker  *broker.Broker
}

func New(b *broker.Broker) *ScriptManager {
	return &ScriptManager{
		broker:        b,
		cAddScript:    make(chan *database.Script),
		cRemoveScript: make(chan string),
		cRunScript:    make(chan *database.Script),
		//cUpdateScript: make(chan *database.Script),
		scripts: make(map[string]*database.Script),
	}

}
func (s *ScriptManager) AddScript(sc *database.Script) {
	go func() {
		s.cAddScript <- sc
	}()

}

func (s *ScriptManager) RunScript(sc *database.Script) {
	//logg.Debug("___RUN")
	go func() {

		s.cRunScript <- sc
	}()

}

func (s *ScriptManager) RemoveScript(sc string) {
	go func() {
		s.cRemoveScript <- sc
	}()

}

func (s *ScriptManager) Run() {
	//logg.Debug("lallala")

	addScript := func(sc *database.Script) {
		s.scripts[sc.Name] = sc
	}

	removeScript := func(scname string) {
		delete(s.scripts, scname)
	}

	runScript := func(sc *database.Script) {

		hub := s.broker.GetHub(sc.User)
		if hub == nil {
			return
		}

		req := wire.ScriptReq{
			Code: sc.Code,
		}
		bytes, err := wire.Encode(0, wire.CMD_SCRIPT, wire.DefaultFlow, req)
		if err != nil {
			return
		}

		pak := broker.DirectPacket{
			Name:  sc.Agent,
			Ntype: broker.AgentType,
			Body:  bytes,
		}

		hub.CDirectSend <- &pak
	}

	checkRunnableScript := func(tm time.Time) {

	}

	T := time.NewTicker(5 * time.Second)
	for {
		select {
		case t := <-T.C:
			checkRunnableScript(t)
		case sc := <-s.cAddScript:
			addScript(sc)
		case sc := <-s.cRemoveScript:
			removeScript(sc)
		case sc := <-s.cRunScript:
			go runScript(sc)
			//case sc := <-s.cUpdateScript:
			//	updateScript(sc)
		}
	}

}

func worker() {

}
