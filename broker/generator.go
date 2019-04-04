package broker

import (
	"sync"

	"github.com/indrenicloud/tricloud-server/core"
)

// uid generator
// is using int32 atomic and reset in max int16 size, return as int16 (casting) better??

type generator struct {
	l     sync.Mutex
	used  map[core.UID]struct{}
	value core.UID
}

func newGenerator() *generator {
	return &generator{
		l:     sync.Mutex{},
		used:  make(map[core.UID]struct{}),
		value: 1,
	}

}

func (g *generator) generate() core.UID {
	g.l.Lock()
	defer g.l.Unlock()
	cuid := g.value

	for { // TODO check infinite loop using a counter and panic maybe
		_, ok := g.used[cuid]
		g.value++
		if !ok {
			return g.value
		}
	}

}

func (g *generator) free(u core.UID) {
	g.l.Lock()
	defer g.l.Unlock()

	delete(g.used, u)

}
