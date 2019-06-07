package noti

import (
	"context"
)

// Provider is a notification provider interface
type Provider interface {
	Init(string)
	PushNotification(context.Context, interface{}) error
}

// CredentialStore provides storage of keys
type CredentialStore interface {
	GetAPIFile() string
	GetAPI() string
	Set(string, string)
	Get(string) []string
}
