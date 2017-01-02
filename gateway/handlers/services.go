package handlers

import (
	"errors"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

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
	service.Status = models.STATUS_ACTIVE
	service.ID = bson.NewObjectIdWithTime(time.Now()).String()
	if ok := network.GetCluster().Services.Add(service); ok == false {
		c.Error(errors.New("This service is already registered."), http.StatusConflict)
		return
	}
	c.WriteJSON(service)
}
