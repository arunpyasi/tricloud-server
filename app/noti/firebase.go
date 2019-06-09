package noti

import (
	"context"
	"errors"
	_ "errors"
	"fmt"
	"log"

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
func NewFirebase() *Firebase {
	f := new(Firebase)

	return f
}

// Init should initilize
func (f *Firebase) Init(confFile string) {
	f.firebaseConfigFile = confFile

	f.opt = option.WithCredentialsFile(f.firebaseConfigFile)
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, f.opt)
	if err != nil {
		fmt.Println("error initializing app:", err)
	}
	f.app = app

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}
	f.mClient = client
}

// PushNotification pushes the notification
func (f *Firebase) PushNotification(ctx context.Context, message interface{}) error {
	fmt.Println("pushing NOTIFICATION HYPEE!!! 🤩")

	_message, ok := message.(*messaging.Message)
	if !ok {
		fmt.Println("Invalid Message type")
		return errors.New("Invalid type of message")
	}
	response, err := f.mClient.Send(ctx, _message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message 📟📠📟: ", response)

	return nil

}
