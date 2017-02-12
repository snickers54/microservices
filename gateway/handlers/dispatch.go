package handlers

import (
	"errors"
	"net/http"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/snickers54/microservices/gateway/context"
	"github.com/snickers54/microservices/gateway/network"
	"github.com/snickers54/microservices/library/models"
)

func dispatchToService(c *context.AppContext) {
	headerVersion := c.Request.Header.Get("Ms-Version-Matching")
	log.WithFields(log.Fields{"version": headerVersion}).Info("Header found ?")
	route := network.GetRoute(strings.Replace(c.Request.URL.String(), "/api", "", -1))
	if route == nil {
		c.Error(errors.New("This route doesn't exists, is your μservice registered to the gateway ?"), http.StatusNotFound)
		return
	}
	// we get from the endpoint all ACTIVE services
	endpoints := route.GetValidEndpoints()
	// We filter by version if the header exists
	if headerVersion != "" {
		semver := models.Version{}
		semver.Parse(headerVersion)
		endpoints = endpoints.FindByVersion(semver)
		log.WithFields(log.Fields{"endpoints": endpoints, "version": semver}).Info("endpoints found for asked semantic versioning.")
	}
	if len(endpoints) == 0 {
		c.Error(errors.New("No endpoint available for this configuration."), http.StatusNotFound)
		return
	}
	endpoint := network.RoundRobin(endpoints)
	c.Set("route", route)
	var wg sync.WaitGroup
	wg.Add(1)
	go network.ReplayHTTP(c, endpoint.Service, &wg)
	wg.Wait()
}
