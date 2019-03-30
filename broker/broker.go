package broker

import (
	"log"
	"net/http"
	"sync"
)

type Broker struct {
	bLock sync.Mutex
	Hubs  map[string]*Hub
}

func (b *Broker) GetHub(user string) *Hub {

	b.bLock.Lock()
	defer b.bLock.Unlock()
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
		log.Println(err)
		w.Write([]byte("error"))
		return
	}

	key, ok := r.Form["key"]
	if !ok {
		log.Println("does not have key")
		return
	}

	parent, err := getParent(key[0]) // TODO auth

	if err != nil {
		return
		log.Println(err)
	}

	hub := b.GetHub(parent)
	node := NewNodeConn(key[0], AgentType, conn, hub)

	hub.AddConnection <- node

}

// ServeUserWebsocket serves users websocket conn
func (b *Broker) ServeUserWebsocket(w http.ResponseWriter, r *http.Request) {

	//return conn id (uuid)
}
