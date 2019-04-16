package restapi

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/restapi/auth"
	"github.com/indrenicloud/tricloud-server/restapi/database"
)

var (
	ErrorPasswordIncorrect = errors.New("Incorrect password")
	ErrorAPI               = errors.New("could not generate api")
)

func RegisterAuthHandlers(r *mux.Router) {

	r.HandleFunc("/signin", SignIn).Methods("POST")
}

func SignIn(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	var logininfo map[string]string
	err = json.Unmarshal(body, &logininfo) //todo check fields
	if err != nil {
		log.Println(err)
	}

	log.Print(logininfo)

	user, err := database.GetUser(logininfo["userid"])

	if err != nil {
		w.Write(GenerateResponse(nil, err))
		return
	}

	if !auth.IsPasswordCorrect(logininfo["password"], user.Password) {
		w.Write(GenerateResponse(nil, ErrorPasswordIncorrect))
		return
	}
	var apierr error
	api := auth.NewAPIKey("session", user.ID, user.SuperUser)

	if api == "" {
		apierr = ErrorAPI
	}

	resp := GenerateResponse(map[string]string{"api": api}, apierr)
	w.Write(resp)
}
