package broker

import (
	"context"
	"net/http"
	"sync"

	"github.com/indrenicloud/tricloud-agent/wire"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/app/auth"
	"github.com/indrenicloud/tricloud-server/app/database"
	"github.com/indrenicloud/tricloud-server/app/logg"
	"github.com/indrenicloud/tricloud-server/app/noti"
)

type Broker struct {
	BLock sync.Mutex
	Hubs  map[string]*Hub
	event *noti.EventManager
}

func NewBroker() *Broker {
	e := noti.NewEventManager()
	return &Broker{
		BLock: sync.Mutex{},
		event: e,
		Hubs:  make(map[string]*Hub),
	}
}

func (b *Broker) getHub(user string) *Hub {

	b.BLock.Lock()
	defer b.BLock.Unlock()
	hub, ok := b.Hubs[user]
	if ok {
		return hub
	}

	hub = NewHub(context.Background(), b.event, user)
	b.Hubs[user] = hub
	go hub.Run()

	return hub
}

// ServeAgentWebsocket serves agents
func (b *Broker) ServeAgentWebsocket(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	key := vars["key"]

	agent, err := database.GetAgent(key)
	if err != nil {
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logg.Warn(err)
		return
	}

	hub := b.getHub(agent.Owner)
	node := NewNodeConn(key, AgentType, conn, hub)
	hub.AddConnection <- node

}

// ServeUserWebsocket serves users websocket conn
func (b *Broker) ServeUserWebsocket(w http.ResponseWriter, r *http.Request) {
	logg.Info("user websocket connn recived")

	vars := mux.Vars(r)
	apikey, ok := vars["apikey"]

	if !ok {
		logg.Warn("Not athorized 1")
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}
	logg.Info(apikey)

	token := auth.ParseAPIKey(apikey)
	claims, ok := token.Claims.(*auth.MyClaims)
	if !ok || !token.Valid {
		logg.Warn("Not athorized 2")
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logg.Warn("could not upgrade")
		logg.Warn(err)
		return
	}
	logg.Warn("upgraded")

	hub := b.getHub(claims.User)
	node := NewNodeConn(claims.User, UserType, conn, hub)

	logg.Warn("adding conn")
	hub.AddConnection <- node

}

func (b *Broker) GetActiveAgents(user string) map[string]wire.UID {
	b.BLock.Lock()
	defer b.BLock.Unlock()

	hub, ok := b.Hubs[user]
	if !ok {
		return nil
	}
	req := &agentsQuery{
		responseChan: make(chan map[string]wire.UID),
	}
	hub.queryAgentsChan <- req
	am := <-req.responseChan
	logg.Info("LOCKOFF recived from channel")
	return am
}
func (b *Broker) RemoveAgent(agentid, user string) {
	b.BLock.Lock()
	defer b.BLock.Unlock()

	hub, ok := b.Hubs[user]
	if ok {
		logg.Info("removing agent")
		hub.removeagentChan <- agentid
	}
	logg.Info("LOCKOFF removing agent  done")
}

func (b *Broker) GetEventManager() *noti.EventManager {
	return b.event
}

func (b *Broker) GetHub(user string) *Hub {
	b.BLock.Lock()
	defer b.BLock.Unlock()
	h, ok := b.Hubs[user]
	if ok {
		return h
	}
	return nil
}
