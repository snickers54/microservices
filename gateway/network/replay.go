package network

import (
	"net/http"
	"net/url"
	"strings"
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
	newHost := ""
	switch node.(type) {
	case *Node:
		newHost = "http://" + node.(*Node).IP + ":" + node.(*Node).Port
		newHost += c.Request.RequestURI
	case *models.Service:
		newHost = "http://" + node.(*models.Service).IP + ":" + node.(*models.Service).Port
		newHost += strings.Replace(c.Request.RequestURI, "/api", "", 1)
	}
	c.Request.RequestURI = ""
	URL, errURL := url.Parse(newHost)
	if errURL != nil {
		log.WithError(errURL).Debug("Couldn't parse the newHost variable")
	}
	c.Request.URL = URL
	log.WithField("request", c.Request).Debug("Replaying request !")
	response, err := client.Do(c.Request)
	log.WithField("response", response).Info("Response from service")
	if err != nil {
		log.WithError(err).WithField("Host", c.Request.URL.Host).Debug("Failed to replay to service.")
		return
	}
	c.Responses = append(c.Responses, response)
}
