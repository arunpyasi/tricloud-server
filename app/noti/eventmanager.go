package noti

import (
	"context"
	"fmt"
	"sync"
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
		fmt.Println("Notification Loop")
		fmt.Println(ee.GetName())
		fmt.Println(tokens)
		for _, token := range tokens {
			go func(_ee Provider, t string) {
				wg.Add(1)
				fmt.Println("Inside Loop 👽👽👽")
				fmt.Println(_ee.GetName())
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
