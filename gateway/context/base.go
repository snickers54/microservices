package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Payload map[string]interface{}

type AppHandler func(*AppContext)
type Responses []*http.Response
type AppContext struct {
	Writer       http.ResponseWriter
	Request      *http.Request
	Keys         map[string]interface{}
	Responses    Responses
	Mutex        sync.Mutex
	IndexHandler int
	Handlers     []AppHandler
}

func (self *AppContext) Next() {
	self.IndexHandler += 1
	for {
		if self.IndexHandler < 0 || self.IndexHandler >= len(self.Handlers) {
			break
		}
		self.Handlers[self.IndexHandler](self)
		self.IndexHandler += 1
	}
}
func (self *AppContext) Done() { self.IndexHandler = -100 }

func HandlerLogic(router *mux.Router, route string, handlers ...AppHandler) *mux.Route {
	return router.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		ctx := AppContext{Writer: w, Request: r, IndexHandler: -1, Handlers: handlers}
		ctx.Set("route_time", time.Now())
		ctx.Next()
	})
}

func (self *AppContext) WriteJSON(value interface{}) {
	result, err := json.Marshal(value)
	if err != nil {
		self.Error(err, http.StatusInternalServerError)
		return
	}
	self.Writer.Header().Set("Content-Type", "application/json")
	self.Writer.WriteHeader(http.StatusOK)
	self.Writer.Write(result)
}

func (self *AppContext) BindJSON(value interface{}) bool {
	if self.Request.Body == nil {
		self.Error(errors.New("Please send a request body"), http.StatusBadRequest)
		return false
	}
	err := json.NewDecoder(self.Request.Body).Decode(&value)
	if err != nil {
		self.Error(err, http.StatusBadRequest)
		return false
	}
	return true
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

func (c *AppContext) Set(key string, value interface{}) {
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}
	c.Keys[key] = value
}

func (c *AppContext) Get(key string) (value interface{}, exists bool) {
	if c.Keys != nil {
		value, exists = c.Keys[key]
	}
	return
}

func ExtractIpPort(value string) (ip string, port string, err error) {
	tuple := strings.Split(value, ":")
	if len(tuple) != 2 {
		err = errors.New("Addr doesn't seems valid, can't extract IP and Port")
		return "", "", err
	}
	return tuple[0], tuple[1], err
}
