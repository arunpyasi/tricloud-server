package noti

import (
	"context"
	_ "errors"
	"log"

	"github.com/indrenicloud/tricloud-server/app/logg"
	"google.golang.org/api/option"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

// Firebase bashed notification provider
type Firebase struct {
	firebaseConfigFile string

	opt option.ClientOption

	app *firebase.App

	mClient *messaging.Client
}

// NewFirebase is a constuctor
func NewFirebase(confFile string) *Firebase {
	f := new(Firebase)
	f.firebaseConfigFile = confFile

	return f
}

// Init should initilize
func (f *Firebase) Init() {

	f.opt = option.WithCredentialsFile(f.firebaseConfigFile)
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, f.opt)
	if err != nil {
		logg.Info("error initializing app:")
		logg.Info(err)
	}
	f.app = app

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}
	f.mClient = client
}

func (f *Firebase) GetName() string {
	return "firebase"
}

// PushNotification pushes the notification
func (f *Firebase) PushNotification(ctx context.Context, _token string, _data string) error {
	logg.Info("pushing NOTIFICATION HYPEE!!! 🤩")

	data := &messaging.Notification{Title: "TriCloud Notification", Body: _data}

	_message := &messaging.Message{
		Notification: data,
		Token:        _token,
	}

	response, err := f.mClient.Send(ctx, _message)
	if err != nil {
		logg.Debug("Firebase push failed 😕😕😕😕😕😕")
		logg.Debug(err)
	}
	// Response is a message ID string.
	logg.Info("Successfully sent message 📟📠📟: ")
	logg.Info(response)

	return nil

}
