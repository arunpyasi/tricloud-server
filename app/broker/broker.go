package broker

import (
	"net/http"
	"sync"

	"github.com/indrenicloud/tricloud-agent/wire"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/app/auth"
	"github.com/indrenicloud/tricloud-server/app/database"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

type Broker struct {
	BLock sync.Mutex
	Hubs  map[string]*Hub
}

func NewBroker() *Broker {
	return &Broker{
		BLock: sync.Mutex{},
		Hubs:  make(map[string]*Hub),
	}
}

func (b *Broker) GetHub(user string) *Hub {

	b.BLock.Lock()
	defer b.BLock.Unlock()
	hub, ok := b.Hubs[user]
	if ok {
		return hub
	}

	hub = NewHub()
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

	hub := b.GetHub(agent.Owner)
	node := NewNodeConn(key, AgentType, conn, hub)

	hub.AddConnection <- node

}

// ServeUserWebsocket serves users websocket conn
func (b *Broker) ServeUserWebsocket(w http.ResponseWriter, r *http.Request) {
	logg.Info("user websocket connn recived")

	token := auth.ParseAPIKey(r.Header.Get("Api-key"))
	claims, ok := token.Claims.(auth.MyClaims)

	if !ok || !token.Valid {
		http.Error(w, "not authorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logg.Warn(err)
		return
	}

	hub := b.GetHub(claims.User)
	node := NewNodeConn(claims.User, UserType, conn, hub)

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
