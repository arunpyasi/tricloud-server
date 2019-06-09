package noti

import (
	"sync"

	"github.com/indrenicloud/tricloud-server/app/config"
)

type CredStore struct {
	config *config.Config
	lmtx   sync.RWMutex
}

func NewCredStore() *CredStore {
	cs := new(CredStore)
	cs.config = config.GetConfig()
	return cs
}

func (c *CredStore) GetAPIFile(provderType string) string {
	c.lmtx.RLock()
	defer c.lmtx.RUnlock()

	pc, ok := c.config.EventProviders[provderType]
	if ok {
		return pc.ConfigFile
	}
	return ""
}

func (c *CredStore) GetAPIKey(provderType string) string {
	c.lmtx.RLock()
	defer c.lmtx.RUnlock()

	pc, ok := c.config.EventProviders[provderType]
	if ok {
		return pc.Apikey
	}
	return ""
}

func (c *CredStore) SetToken(provderType string, user string, authkey string) {
	c.lmtx.Lock()
	defer c.lmtx.Unlock()

	pc, ok := c.config.EventProviders[provderType]
	keys, ok := pc.TokenPerUser[user]

	if !ok {
		keys = []string{authkey}
		//keys = authkey
	} else {
		keys = append(keys, authkey)
		//keys = authkey
	}

	pc.TokenPerUser[user] = keys
	c.config.Update()
}

func (c *CredStore) GetToken(provderType string, user string) []string {
	c.lmtx.RLock()
	defer c.lmtx.RUnlock()

	pc, ok := c.config.EventProviders[provderType]
	ks, ok := pc.TokenPerUser[user]
	if ok {
		return ks
	}
	return nil
}

func (c *CredStore) SetOption(provderType string, user string, option string) {
	c.lmtx.Lock()
	defer c.lmtx.Unlock()

	pc, ok := c.config.EventProviders[provderType]
	if ok {
		pc.Options[user] = option
		c.config.Update()
	}
}

func (c *CredStore) GetOption(provderType string, user string) string {
	c.lmtx.RLock()
	defer c.lmtx.RUnlock()

	pc, ok := c.config.EventProviders[provderType]
	if ok {
		ks, ok := pc.Options[user]
		if ok {
			return ks
		}
	}

	return ""
}
