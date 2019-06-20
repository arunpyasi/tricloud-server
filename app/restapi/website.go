package restapi

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/app/database"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

func GetWebsites(w http.ResponseWriter, r *http.Request) {

	websites, err := database.GetAllWebsites()
	if err != nil {
		errorResp(w, err)
		return
	}
	generateResp(w, websites, err)

}

func CreateWebsite(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		errorResp(w, err)
		return
	}
	defer r.Body.Close()
	var websiteinfo map[string]string
	err = deJson(body, &websiteinfo)
	if err != nil {
		errorResp(w, err)
		return
	}
	website, err := database.NewWebsite(websiteinfo, false)
	if err != nil {
		errorResp(w, err)
		return
	}
	database.CreateWebsite(website)
	go sMonitor.AddWebsite(website)
	updatedwesbsites, err := database.GetAllWebsites()
	generateResp(w, updatedwesbsites, err)

}

func DeleteWebsite(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["name"]
	logg.Warn("Deleting website url")
	err := database.DeleteWebsite(name)
	if err != nil {
		errorResp(w, err)
	}
	go sMonitor.RemoveWebsite(name)

	updatedwebsites, err := database.GetAllWebsites()
	if err != nil {
		errorResp(w, err)
		return
	}

	generateResp(w, updatedwebsites, err)

}

func GetWebsite(w http.ResponseWriter, r *http.Request) {
	// only if user owns agent or superuser
	vars := mux.Vars(r)
	name := vars["name"]
	website, err := database.GetWebsite(name)

	if err != nil {
		errorResp(w, err)
		return
	}

	generateResp(w, website, err)
}
