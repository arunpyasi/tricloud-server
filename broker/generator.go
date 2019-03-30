package broker

import "sync"

// uid generator

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
	for {
		_, ok := g.used[cuid]
		g.value++
		if !ok {
			return g.value
		}
	}

	return g.value
}

func (g *generator) free(u uid) {
	g.l.Lock()
	defer g.l.Unlock()

	delete(g.used, u)

}
