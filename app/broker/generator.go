package broker

import (
	"sync"

	"github.com/indrenicloud/tricloud-agent/wire"
)

// uid generator
// is using int32 atomic and reset in max int16 size, return as int16 (casting) better??

type generator struct {
	l     sync.Mutex
	used  map[wire.UID]struct{}
	value wire.UID
}

func newGenerator() *generator {
	return &generator{
		l:     sync.Mutex{},
		used:  make(map[wire.UID]struct{}),
		value: 1,
	}

}

func (g *generator) generate() wire.UID {
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

func (g *generator) free(u wire.UID) {
	g.l.Lock()
	defer g.l.Unlock()

	delete(g.used, u)

}
