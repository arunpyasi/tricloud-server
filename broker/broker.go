package broker

import (
	"log"
	"net/http"
	"sync"

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

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("could not upgrade to wesocket:", err)
		return
	}

	key, ok := r.Form["key"] // todo or cookies
	if !ok {
		log.Println("does not have key")
		return
	}

	parent, err := getParent(key[0]) // TODO auth

	if err != nil {
		log.Println(err)
		return
	}

	hub := b.GetHub(parent)
	node := NewNodeConn(key[0], AgentType, conn, hub)

	hub.AddConnection <- node

}

// ServeUserWebsocket serves users websocket conn
func (b *Broker) ServeUserWebsocket(w http.ResponseWriter, r *http.Request) {

	//return conn id (uuid)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("could not upgrade to wesocket:", err)
		return
	}

	c, err := r.Cookie("wscookie")
	if err != nil {
		log.Println("nocookie! but i am hungry", err)
		return
	}
	user := getUserFromCookie(c.Value)

	hub := b.GetHub(user)
	node := NewNodeConn(user, UserType, conn, hub)

	hub.AddConnection <- node

}

func (b *Broker) GetAgents(user string) []string {
	return nil
}
