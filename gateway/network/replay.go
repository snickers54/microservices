package network

import (
	"net/http"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/snickers54/microservices/gateway/context"
	"github.com/snickers54/microservices/library/models"
)

func ReplayHTTP(c *context.AppContext, node interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}
	switch node.(type) {
	case Node:
		c.Request.URL.Host = node.(Node).IP + ":" + node.(Node).Port
	case models.Service:
		c.Request.URL.Host = node.(models.Service).IP + ":" + node.(models.Service).Port
	}
	log.WithField("request", c.Request).Debug("Replaying request !")
	response, err := client.Do(c.Request)
	if err != nil {
		return
	}
	c.Responses = append(c.Responses, response)
}
