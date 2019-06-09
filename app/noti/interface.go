package noti

import (
	"context"
)

// Provider is a notification provider interface
type Provider interface {
	Init()
	GetName() string
	PushNotification(context.Context, string, map[string]string) error
}

// CredentialStore provides storage of keys
type CredentialStore interface {
	GetAPIFile(string) string
	GetAPIKey(string) string
	SetToken(string, string, string)
	GetToken(string, string) []string
	SetOption(string, string, string)
	GetOption(string, string) string
}
