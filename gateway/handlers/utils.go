package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Payload map[string]interface{}

type AppHandler func(*AppContext)

type AppContext struct {
	Writer              http.ResponseWriter
	Request             *http.Request
	HandlerInterruption bool
}

func (self *AppContext) Next() { self.HandlerInterruption = false }
func (self *AppContext) Done() { self.HandlerInterruption = true }

func (self *AppContext) WriteJSON(value interface{}) {
	result, err := json.Marshal(value)
	if err != nil {
		self.Error(err, http.StatusInternalServerError)
		return
	}
	self.Writer.Header().Set("Content-Type", "application/json")
	self.Writer.Write(result)
}

func (self *AppContext) BindJSON(value interface{}) {
	if self.Request.Body == nil {
		self.Error(errors.New("Please send a request body"), http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(self.Request.Body).Decode(&value)
	if err != nil {
		self.Error(err, http.StatusBadRequest)
		return
	}
}

func (self *AppContext) Error(err error, code int) {
	defer self.Done()
	bytes, err := json.Marshal(
		Payload{
			"message": err.Error(),
			"code":    fmt.Sprintf("%d", code),
		})
	if err != nil {
		http.Error(self.Writer, "{\"message\":\"Internal server error\", \"code\":\"500\"}", http.StatusInternalServerError)
		return
	}
	http.Error(self.Writer, string(bytes), code)
}

func handlerLogic(router *mux.Router, route string, handlers ...AppHandler) *mux.Route {
	return router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		context := AppContext{Writer: w, Request: r, HandlerInterruption: false}
		for _, handler := range handlers {
			if context.HandlerInterruption == false {
				handler(&context)
			}
		}
	})
}

func GET(router *mux.Router, route string, handlers ...AppHandler) {
	handlerLogic(router, route, handlers...).Methods("GET")
}

func POST(router *mux.Router, route string, handlers ...AppHandler) {
	handlerLogic(router, route, handlers...).Methods("POST")
}

func PUT(router *mux.Router, route string, handlers ...AppHandler) {
	handlerLogic(router, route, handlers...).Methods("PUT")
}

func DELETE(router *mux.Router, route string, handlers ...AppHandler) {
	handlerLogic(router, route, handlers...).Methods("DELETE")
}
