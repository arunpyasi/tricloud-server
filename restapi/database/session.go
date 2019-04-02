package database

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

//var lock sync.Mutex
var store = sessions.NewCookieStore([]byte("DevPassword11111"))

func GetUserFromSession(r *http.Request) (string, error) {

	session, err := store.Get(r, "session")
	if err != nil {
		return "", err
	}

	user, ok := session.Values["username"]

	if !ok {
		return "", fmt.Errorf("user not set")
	}

	value, ok2 := user.(string)

	if !ok2 || user == "" {
		return "", fmt.Errorf("invalid type set or blank")
	}

	return value, nil
}

func SetUserSession(user string, w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, "session")
	if err != nil {
		return err
	}

	session.Values["username"] = user

	err = session.Save(r, w)

	if err != nil {
		return err
	}
	return nil
}
