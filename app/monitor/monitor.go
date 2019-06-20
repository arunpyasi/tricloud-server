package monitor

import (
	"context"
	"net/http"
	"time"

	"github.com/indrenicloud/tricloud-server/app/database"
	"github.com/indrenicloud/tricloud-server/app/logg"
	"github.com/indrenicloud/tricloud-server/app/noti"
)

type Monitor struct {
	cSiteState    chan *database.Website
	cNewSite      chan *database.Website
	cRemoveSite   chan string
	oldSiteStates map[string]*database.Website
	siteWorkers   map[string]context.CancelFunc
	eventm        *noti.EventManager
}

func NewMonitor(eventm *noti.EventManager) *Monitor {
	return &Monitor{
		eventm:        eventm,
		cSiteState:    make(chan *database.Website),
		cNewSite:      make(chan *database.Website),
		cRemoveSite:   make(chan string),
		oldSiteStates: make(map[string]*database.Website),
		siteWorkers:   make(map[string]context.CancelFunc),
	}
}
func (m *Monitor) AddWebsite(ws *database.Website) {
	m.cNewSite <- ws
}

func (m *Monitor) RemoveWebsite(name string) {
	m.cRemoveSite <- name
}

func (m *Monitor) Run() {
	sites, err := database.GetAllWebsites()
	if err != nil {
		logg.Debug("GET SITES ERROR")
	}

	ctx := context.Background()
	newSite := func(_site *database.Website) {
		logg.Debug("adding suspect to watchlist👮‍♀️👮‍♀️👮‍♀️")
		ctxw, ctxwFuc := context.WithCancel(ctx)
		go m.worker(_site.Name, _site.Url, ctxw)
		m.oldSiteStates[_site.Name] = _site
		m.siteWorkers[_site.Name] = ctxwFuc
	}

	for _, site := range sites {
		newSite(site)
	}

	for {
		select {
		case s := <-m.cSiteState:
			//logg.Debug("adding to db")
			oldstate, ok := m.oldSiteStates[s.Name]
			if !ok {
				logg.Debug("Map empty :(")
				continue
			}
			s.Subscriber = oldstate.Subscriber
			if oldstate.Active != s.Active {
				m.emitAlert(s)
			}

			m.oldSiteStates[s.Name] = s
			go database.UpdateWebsite(s)
		case site := <-m.cNewSite:
			newSite(site)

		case site := <-m.cRemoveSite:
			cf := m.siteWorkers[site]
			cf()
			delete(m.siteWorkers, site)
			delete(m.oldSiteStates, site)
			logg.Debug("removing " + site)
		}
	}

}

func (m *Monitor) emitAlert(site *database.Website) {

	bytes, err := database.Encode(site)
	if err != nil {
		return
	}
	m.eventm.SendEvent(site.Subscriber, string(bytes))

}

func (m *Monitor) worker(name, url string, ctx context.Context) {

	for {
		var site = database.Website{
			Name: name,
			Url:  url,
		}

		resp, err := http.Get(url)
		if err != nil {
			print(err.Error())
			site.Active = false
		} else {
			print(string(resp.StatusCode) + resp.Status)
			site.Active = true
		}
		site.Timestamp = time.Now().Unix()
		m.cSiteState <- &site

		time.Sleep(5 * time.Second)

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}
