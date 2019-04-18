package restapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/indrenicloud/tricloud-server/app/auth"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

// when using refelection to find type  using custom type avoids collision in contex.value
type contextkey int

const ContextUser contextkey = iota

var (
	ErrorNotAuthorized = errors.New("Not authorized")
)

func MiddlewareJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}

// MiddlewareSession checks the session for request and tags username to request context
func MiddlewareSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := auth.ParseAPIKey(r.Header.Get("Api-key"))
		if !token.Valid {
			logg.Warn("Token invalid")
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*auth.MyClaims)
		logg.Info(claims)

		if !ok {
			logg.Warn("Cannot cast clam type")
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}

		ctx := context.WithValue(r.Context(), ContextUser, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logg.Info(r.URL.Path)
		h.ServeHTTP(w, r)
	})
}
