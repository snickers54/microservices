package network

import (
	"net/http"
	"sync"
	"time"

	"github.com/snickers54/microservices/gateway/context"
)

func ReplayHTTP(c *context.AppContext, node interface{}, wg sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	switch node.(type) {
	case Node:
		c.Request.URL.Host = node.(Node).IP + ":" + node.(Node).Port
	case Service:
		c.Request.URL.Host = node.(Service).IP + ":" + node.(Service).Port
	}
	response, err := client.Do(c.Request)
	if err != nil {
		return
	}
	c.Responses = append(c.Responses, response)
}
