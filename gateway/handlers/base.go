package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func statsHandlers(router *mux.Router) {
	subRouter := router.PathPrefix("/stats").Subrouter()
	GET(subRouter, "/", statsSummarize)
}

func gatewaysHandlers(router *mux.Router) {
	subRouter := router.PathPrefix("/cluster").Subrouter()
	GET(subRouter, "/", clusterDescribe)
	POST(subRouter, "/nodes", clusterRegister)
}

func servicesHandlers(router *mux.Router) {
	subRouter := router.PathPrefix("/services").Subrouter()
	GET(subRouter, "/", servicesDescribe)
}

func Start() {
	router := mux.NewRouter()
	statsHandlers(router)
	gatewaysHandlers(router)
	log.Fatal(http.ListenAndServe(":8000", router))
}
