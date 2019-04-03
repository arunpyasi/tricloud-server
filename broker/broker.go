package broker

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/restapi/database"
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

	if key != "456456" {
		log.Println("invalid key")
	}
	owner := "root"

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("could not upgrade to wesocket:", err)
		return
	}

	//parent, err := getParent(key)

	hub := b.GetHub(owner)
	node := NewNodeConn(key, AgentType, conn, hub)

	hub.AddConnection <- node

}

// ServeUserWebsocket serves users websocket conn
func (b *Broker) ServeUserWebsocket(w http.ResponseWriter, r *http.Request) {
	log.Println("user websocket connn recived")

	user, err := database.GetUserFromSession(r)
	if err != nil {
		log.Println("session not set:", err)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("could not upgrade to wesocket:", err)
		return
	}

	hub := b.GetHub(user)
	node := NewNodeConn(user, UserType, conn, hub)

	hub.AddConnection <- node

}

func (b *Broker) GetAgents(user string) []string {
	return nil
}
