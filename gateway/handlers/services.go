package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/snickers54/microservices/gateway/context"
	"github.com/snickers54/microservices/gateway/network"
	"github.com/snickers54/microservices/library/models"
)

func servicesDescribe(c *context.AppContext) {
	c.WriteJSON(network.GetCluster().Services)
}

func servicesRegister(c *context.AppContext) {
	service := models.Service{}
	if ok := c.BindJSON(&service); !ok {
		return
	}
	service.IP = strings.Split(c.Request.RemoteAddr, ":")[0]
	service.Status = models.STATUS_ACTIVE
	if ok := network.GetCluster().Services.Add(service); ok == false {
		c.Error(errors.New("This service is already registered."), http.StatusConflict)
		return
	}
	c.WriteJSON(service)
}
