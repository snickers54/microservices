package handlers

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/snickers54/microservices/gateway/context"
)

type GSubRouter struct {
	router *mux.Router
}

func NewGSubRouter(router *mux.Router) *GSubRouter {
	return &GSubRouter{
		router: router,
	}
}

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

func (self *GSubRouter) GET(route string, handlers ...context.AppHandler) {
	handlerLogic(self.router, route, handlers...).Methods("GET")
}

func (self *GSubRouter) POST(route string, handlers ...context.AppHandler) {
	handlerLogic(self.router, route, handlers...).Methods("POST")
}

func (self *GSubRouter) PUT(route string, handlers ...context.AppHandler) {
	handlerLogic(self.router, route, handlers...).Methods("PUT")
}

func (self *GSubRouter) DELETE(route string, handlers ...context.AppHandler) {
	handlerLogic(self.router, route, handlers...).Methods("DELETE")
}
