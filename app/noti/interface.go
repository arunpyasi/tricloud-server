package noti

import (
	"context"
)

// Provider is a notification provider interface
type Provider interface {
	Init(CredentialStore)
	PushNotification(context.Context, interface{}) error
}

// CredentialStore provides storage of keys
type CredentialStore interface {
	Set(string, string)
	Get(string) []string
}
