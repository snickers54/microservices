package handlers

import (
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

func (self *GSubRouter) GET(route string, handlers ...context.AppHandler) {
	context.HandlerLogic(self.router, route, handlers...).Methods("GET")
}

func (self *GSubRouter) POST(route string, handlers ...context.AppHandler) {
	context.HandlerLogic(self.router, route, handlers...).Methods("POST")
}

func (self *GSubRouter) PUT(route string, handlers ...context.AppHandler) {
	context.HandlerLogic(self.router, route, handlers...).Methods("PUT")
}

func (self *GSubRouter) DELETE(route string, handlers ...context.AppHandler) {
	context.HandlerLogic(self.router, route, handlers...).Methods("DELETE")
}
