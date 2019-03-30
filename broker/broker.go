package broker

import (
	"sync"

	"github.com/google/uuid"
)

type Broker struct {
	BLock          sync.RWMutex
	connectionPool map[uuid.UUID]*NodeConn
	Nodes          map[string][]uuid.UUID
}

func (b *Broker) RegisterAgent(conn *NodeConn) error {
	return nil
}

func (b *Broker) RegisterUser(conn *NodeConn) error {
	return nil
}

func (b *Broker) getConnection(connId uuid.UUID) (*NodeConn, error) {
	return nil, nil
}

/*
func (b *Broker) UnregisterAgent() error {

}

func (b *Broker) UnregisterUser() error {

}
*/
