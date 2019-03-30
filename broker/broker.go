package broker

import (
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

	return hub
}

func (b *Broker) ServeAgentWebsocket(w http.ResponseWriter, r *http.Request) {
	//check deploykey/per agent key if its valid
	// find the user it agent belongs to
	// get hub from userid
	// hub.addconnection
	//return conn id (uuid)
}

func (b *Broker) ServeUserWebsocket(w http.ResponseWriter, r *http.Request) {

	//return conn id (uuid)
}

func registerAgent() {

}

func registerNewConnection() {
	// key / userid
	// authorization
	// create node type
	// if its agent/key if its userid and get hub
	//
}
