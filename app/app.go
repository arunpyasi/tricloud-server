package app

import (
	"net/http"

	"github.com/indrenicloud/tricloud-server/app/broker"
	"github.com/indrenicloud/tricloud-server/app/database"
	"github.com/indrenicloud/tricloud-server/app/logg"
	"github.com/indrenicloud/tricloud-server/app/restapi"
	"github.com/rs/cors"
)

func Run() {

	b := broker.NewBroker()
	database.Start()
	defer database.Close()

	go listenAgentsConnection(b)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions, http.MethodPut},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Api-key", "Content-Type", "Accept"},
	})
	logg.Warn(http.ListenAndServe(":8080", c.Handler(restapi.GetMainRouter(b))))

}

func listenAgentsConnection(b *broker.Broker) {
	logg.Warn(http.ListenAndServe(":8081", restapi.GetAgentRouter(b)))
}
