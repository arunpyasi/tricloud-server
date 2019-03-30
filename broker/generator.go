package broker

import "sync"

// uid generator
// is using int32 atomic and reset in max int16 size, return as int16 (casting) better??

type uid uint16

type generator struct {
	l     sync.Mutex
	used  map[uid]struct{}
	value uid
}

func newGenerator() *generator {
	return &generator{
		l:     sync.Mutex{},
		used:  make(map[uid]struct{}),
		value: 0,
	}

}

func (g *generator) generate() uid {
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

func (g *generator) free(u uid) {
	g.l.Lock()
	defer g.l.Unlock()

	delete(g.used, u)

}
