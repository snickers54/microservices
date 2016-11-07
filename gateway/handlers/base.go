package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/snickers54/microservices/gateway/context"
	"github.com/snickers54/microservices/gateway/middlewares"
	"github.com/spf13/viper"
)

func statsHandlers(router *mux.Router) {
	subRouter := NewGSubRouter(router.PathPrefix("/stats").Subrouter())
	subRouter.GET("/", statsSummarize)
}

func gatewaysHandlers(router *mux.Router) {
	subRouter := NewGSubRouter(router.PathPrefix("/cluster").Subrouter())
	subRouter.GET("/", clusterDescribe)
	subRouter.POST("/nodes", clusterRegister, middlewares.Sync, middlewares.CloseBody)
}

func servicesHandlers(router *mux.Router) {
	subRouter := NewGSubRouter(router.PathPrefix("/services").Subrouter())
	subRouter.GET("/", servicesDescribe)
	subRouter.POST("/", middlewares.Sync, middlewares.CloseBody, servicesRegister)
}

func dispatchHandlers(router *mux.Router) {
	subRouter := NewGSubRouter(router.PathPrefix("/api").Subrouter())
	patternDispatch := "/{path:.*}"
	replayMiddlewares := []context.AppHandler{
		dispatchToService,
		middlewares.WriteReplay,
		middlewares.CloseBody,
		middlewares.StatsReplay,
	}
	subRouter.GET(patternDispatch, replayMiddlewares...)
	subRouter.POST(patternDispatch, replayMiddlewares...)
	subRouter.PUT(patternDispatch, replayMiddlewares...)
	subRouter.DELETE(patternDispatch, replayMiddlewares...)
}

func Start() {
	router := mux.NewRouter()
	statsHandlers(router)
	gatewaysHandlers(router)
	servicesHandlers(router)
	dispatchHandlers(router)
	log.Fatal(http.ListenAndServe(":"+viper.GetString("cluster.port"), router))
}
