package handlers

import (
	"errors"
	"net/http"

	"github.com/snickers54/microservices/gateway/context"
	"github.com/snickers54/microservices/gateway/network"
)

func servicesDescribe(c *context.AppContext) {
	c.WriteJSON(network.GetCluster().Services)
}

func servicesRegister(c *context.AppContext) {
	service := network.Service{}
	if ok := c.BindJSON(&service); !ok {
		return
	}
	service.Status = network.STATUS_ACTIVE
	if ok := network.GetCluster().Services.Add(service); ok == false {
		c.Error(errors.New("This node already exists."), http.StatusConflict)
		return
	}
	c.WriteJSON(service)
}
