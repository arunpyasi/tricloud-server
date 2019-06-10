package noti

import (
	"context"
	"sync"

	"github.com/indrenicloud/tricloud-server/app/logg"
)

type EventManager struct {
	cs             CredentialStore //its should be threadsafe
	eventProviders []Provider
}

func NewEventManager() *EventManager {
	em := new(EventManager)
	em.cs = NewCredStore()
	em.eventProviders = notificationProviders(em.cs)
	for _, ep := range em.eventProviders {
		ep.Init()
	}

	return em
}

func (e *EventManager) SendEvent(user string, data map[string]string) {

	var wg sync.WaitGroup

	for _, ee := range e.eventProviders {
		tokens := e.cs.GetToken(ee.GetName(), user)
		logg.Info("Notification Loop")
		logg.Info(ee.GetName())
		logg.Info(tokens)
		for _, token := range tokens {
			go func(_ee Provider, t string) {
				wg.Add(1)
				logg.Info("Inside Loop 👽👽👽")
				logg.Info(_ee.GetName())
				_ee.PushNotification(context.Background(), t, data)
				wg.Done()
			}(ee, token)

		}
	}
	wg.Wait()

}

func (e *EventManager) SaveToken(user string, token string) {
	// TODO
	e.cs.SetToken("firebase", user, token)
}
