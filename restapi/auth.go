package restapi

import (
	"log"
	"net/http"

	"github.com/indrenicloud/tricloud-server/restapi/database"

	"github.com/gorilla/mux"
)

func RegisterAuthHandlers(r *mux.Router) {

	r.HandleFunc("/signin", SignIn)
	r.HandleFunc("/signup", SignUp)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "root" && password == "123" {
		err := database.SetUserSession(username, w, r)
		if err != nil {
			log.Println("couldnot set session:", err)
		}
		log.Println("best password of year")
	}
}

func SignUp(w http.ResponseWriter, r *http.Request) {

}
