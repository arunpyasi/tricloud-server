package noti

import (
	"context"

	"firebase.google.com/go/messaging"
)

type EventManager struct {
	cs CredentialStore //its should be threadsafe

	eventProvider  Provider
	chanLogConsume chan []byte
}

func NewEventManager() *EventManager {
	em := new(EventManager)
	em.cs = new(CredStore)
	em.eventProvider = NewFirebase()
	em.eventProvider.Init(em.cs.GetAPIFile())
	return em
}

func (e *EventManager) SendEvent(user string, data map[string]string) {

	tokens := e.cs.Get(user)
	for _, token := range tokens {
		_message := &messaging.Message{
			Data:  data,
			Token: token,
		}
		e.eventProvider.PushNotification(context.Background(), _message)
	}

}

func (e *EventManager) SaveToken(user string, token string) {
	e.cs.Set(user, token)

}
