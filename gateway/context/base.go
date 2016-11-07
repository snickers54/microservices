package context

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type Payload map[string]interface{}

type AppHandler func(*AppContext)
type Responses []*http.Response
type AppContext struct {
	Writer              http.ResponseWriter
	Request             *http.Request
	HandlerInterruption bool
	Keys                map[string]interface{}
	Responses           Responses
	Mutex               sync.Mutex
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
	self.Writer.WriteHeader(http.StatusOK)
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

func (c *AppContext) MustGet(key string) interface{} {
	if value, exists := c.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

func ExtractIpPort(value string) (ip string, port string, err error) {
	tuple := strings.Split(value, ":")
	if len(tuple) != 2 {
		err = errors.New("Addr doesn't seems valid, can't extract IP and Port")
	}
	return tuple[0], tuple[1], err
}
