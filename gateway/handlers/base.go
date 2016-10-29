package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/snickers54/microservices/gateway/context"
	"github.com/snickers54/microservices/gateway/middlewares"
)

func statsHandlers(router *mux.Router) {
	subRouter := router.PathPrefix("/stats").Subrouter()
	GET(subRouter, "/", statsSummarize)
}

func gatewaysHandlers(router *mux.Router) {
	subRouter := router.PathPrefix("/cluster").Subrouter()
	GET(subRouter, "/", clusterDescribe)
	POST(subRouter, "/nodes", clusterRegister, middlewares.Sync, middlewares.CloseBody)
}

func servicesHandlers(router *mux.Router) {
	subRouter := router.PathPrefix("/services").Subrouter()
	GET(subRouter, "/", servicesDescribe)
	POST(subRouter, "/", servicesRegister, middlewares.Sync, middlewares.CloseBody)
}

func dispatchHandlers(router *mux.Router) {
	subRouter := router.PathPrefix("/api").Subrouter()
	patternDispatch := "/{path:.*}"
	replayMiddlewares := []context.AppHandler{
		dispatchToService,
		middlewares.WriteReplay,
		middlewares.CloseBody,
		middlewares.StatsReplay,
	}
	GET(subRouter, patternDispatch, replayMiddlewares...)
	POST(subRouter, patternDispatch, replayMiddlewares...)
	PUT(subRouter, patternDispatch, replayMiddlewares...)
	DELETE(subRouter, patternDispatch, replayMiddlewares...)
}

func Start() {
	router := mux.NewRouter()
	statsHandlers(router)
	gatewaysHandlers(router)
	servicesHandlers(router)
	dispatchHandlers(router)
	log.Fatal(http.ListenAndServe(":8000", router))
}
