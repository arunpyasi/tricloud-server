package restapi

import (
	"net/http"

	"github.com/indrenicloud/tricloud-server/app/auth"
)

func rootRoute(h http.ResponseWriter, r *http.Request) {

	api := r.Header.Get("Api-key")
	if api == "" {
		http.ServeFile(h, r, "./public/login.html")
		return
	}
	token := auth.ParseAPIKey(api)
	_, ok := token.Claims.(*auth.MyClaims)
	if !ok || !token.Valid {
		http.ServeFile(h, r, "./public/login.html")
		return
	}
	http.ServeFile(h, r, "./public/dashboard.html")
}
