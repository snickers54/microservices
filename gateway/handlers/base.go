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
	subRouter.GET("", statsSummarize)
}

func clusterHandlers(router *mux.Router) {
	subRouter := NewGSubRouter(router.PathPrefix("/cluster").Subrouter())
	subRouter.GET("", clusterDescribe)
	subRouter.POST("/nodes", middlewares.Sync, middlewares.CloseBody, clusterRegister)
}

func servicesHandlers(router *mux.Router) {
	subRouter := NewGSubRouter(router.PathPrefix("/services").Subrouter())
	subRouter.GET("", servicesDescribe)
	subRouter.POST("", servicesRegister, middlewares.Sync, middlewares.CloseBody)
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

func baseHandlers(router *mux.Router) {
	subRouter := NewGSubRouter(router.PathPrefix("/").Subrouter())
	subRouter.GET("/healthcheck", nodePing)
}

func Start() {
	router := mux.NewRouter()
	baseHandlers(router)
	statsHandlers(router)
	clusterHandlers(router)
	servicesHandlers(router)
	dispatchHandlers(router)
	log.Fatal(http.ListenAndServe(":"+viper.GetString("node.port"), router))
}
