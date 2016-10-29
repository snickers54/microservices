package handlers

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/snickers54/microservices/gateway/context"
)

func handlerLogic(router *mux.Router, route string, handlers ...context.AppHandler) *mux.Route {
	return router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		context := context.AppContext{Writer: w, Request: r, HandlerInterruption: false}
		context.Set("route_time", time.Now())
		for _, handler := range handlers {
			if context.HandlerInterruption == false {
				handler(&context)
			}
		}
	})
}

func GET(router *mux.Router, route string, handlers ...context.AppHandler) {
	handlerLogic(router, route, handlers...).Methods("GET")
}

func POST(router *mux.Router, route string, handlers ...context.AppHandler) {
	handlerLogic(router, route, handlers...).Methods("POST")
}

func PUT(router *mux.Router, route string, handlers ...context.AppHandler) {
	handlerLogic(router, route, handlers...).Methods("PUT")
}

func DELETE(router *mux.Router, route string, handlers ...context.AppHandler) {
	handlerLogic(router, route, handlers...).Methods("DELETE")
}
