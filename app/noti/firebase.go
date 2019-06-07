package noti

import (
	"context"
	"log"
	"fmt"
	_"errors"
	"google.golang.org/api/option"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

// Firebase bashed notification provider
type Firebase struct {
	store  CredentialStore
	
	firebaseConfigFile string

	opt option.ClientOption
	
	app *firebase.App

	mClient *messaging.Client


}

// NewFirebase is a constuctor
func NewFirebase() *Firebase {
	f := new(Firebase)
	f.firebaseConfigFile = "./tcloud-42ebf-firebase-adminsdk-ma9t8-d5a2581857.json"
	return f
}

// Init should initilize
func (f *Firebase) Init(store CredentialStore) {
	f.store = store
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
func (f *Firebase) PushNotification(ctx context.Context, message interface{} ) error {
	
	//_message, ok := message.(*messaging.Message)
	//if !ok {
	//	return errors.New("Invalid message type")
	//}

	// This registration token comes from the client FCM SDKs.
	registrationToken := "fegnEF0AXtY:APA91bG4f6R6S0I1vtAkf7ngd0z6Vo3aaUiMnCMpy7pmgDZF0aplQ41tt4F4ww0FRhK1BEkZFnEk1nEa79D0hFeGk5ydYldwjSX67P17a71sbCT9iwiJ5JLmXizEOz9xVGzA9i8Ux3M9"

	// See documentation on defining a message payload.
	
	_message := &messaging.Message{
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Token: registrationToken,
	}
	

	response, err := f.mClient.Send(ctx, _message)
	if err != nil {
		log.Fatalln(err)
	}
	// Response is a message ID string.
	fmt.Println("Successfully sent message:", response)

	
	
	return nil

}
