package noti

import (
	"sync"

	"github.com/indrenicloud/tricloud-server/app/config"
)

type CredStore struct {
	lmtx sync.RWMutex
}

func (c *CredStore) GetAPIFile() string {

	conf := config.GetConfig()
	return conf.AppFirebaseKeyFile
}

func (c *CredStore) GetAPI() string {
	panic("Not IMPLEMENTED!!")
}

func (c *CredStore) Set(user string, authkey string) {
	c.lmtx.Lock()
	defer c.lmtx.Unlock()

	conf := config.GetConfig()
	keys, ok := conf.FirebaseKeys[user]

	if !ok {
		keys = []string{authkey}
		//keys = authkey
	} else {
		keys = append(keys, authkey)
		//keys = authkey
	}

	conf.FirebaseKeys[user] = keys
	conf.Update()
}

func (c *CredStore) Get(user string) []string {
	c.lmtx.RLock()
	defer c.lmtx.RUnlock()

	conf := config.GetConfig()
	keys, ok := conf.FirebaseKeys[user]
	if ok {
		return keys
	}
	return nil
}
