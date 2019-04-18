package app

import (
	"net/http"

	"github.com/indrenicloud/tricloud-server/app/broker"
	"github.com/indrenicloud/tricloud-server/app/database"
	"github.com/indrenicloud/tricloud-server/app/logg"
	"github.com/indrenicloud/tricloud-server/app/restapi"
)

func Run() {

	b := broker.NewBroker()
	database.Start()
	defer database.Close()

	go listenAgentsConnection(b)

	logg.Warn(http.ListenAndServe(":8080", restapi.GetMainRouter(b)))

}

func listenAgentsConnection(b *broker.Broker) {
	logg.Warn(http.ListenAndServe(":8081", restapi.GetAgentRouter(b)))
}
