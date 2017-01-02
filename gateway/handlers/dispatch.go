package handlers

import (
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/snickers54/microservices/gateway/context"
	"github.com/snickers54/microservices/gateway/network"
	"github.com/snickers54/microservices/library/models"
)

func dispatchToService(c *context.AppContext) {
	headerVersion, exists := c.Get("μS-Version-Matching")
	route := network.GetRoute(strings.Replace(c.Request.URL.String(), "/api", "", -1))
	if route == nil {
		c.Error(errors.New("This route doesn't exists, is your μservice registered to the gateway ?"), http.StatusNotFound)
		c.Done()
		return
	}
	// we get from the endpoint all ACTIVE services
	endpoints := route.GetValidEndpoints()
	// We filter by version if the header exists
	if exists == true {
		semver := models.Version{}
		semver.Parse(headerVersion.(string))
		endpoints = endpoints.FindByVersion(semver)
	}

	endpoint := network.RoundRobin(endpoints)
	c.Set("route", route)
	var wg sync.WaitGroup
	wg.Add(1)
	go network.ReplayHTTP(c, endpoint.Service, &wg)
	wg.Wait()
}
