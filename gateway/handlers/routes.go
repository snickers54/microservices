package handlers

import (
	"errors"
	"net/http"

	"github.com/snickers54/microservices/gateway/context"
	"github.com/snickers54/microservices/gateway/network"
	"github.com/snickers54/microservices/library/models"
)

func routeRegister(c *context.AppContext) {
	tempStruct := struct {
		Route   string         `json:"route"`
		Service models.Service `json:"service"`
	}{}
	if ok := c.BindJSON(&tempStruct); !ok {
		return
	}
	if ok := network.InsertRoute(tempStruct.Route, &tempStruct.Service); ok == false {
		c.Error(errors.New("Inserting route failed."), http.StatusConflict)
		return
	}
	c.WriteJSON(tempStruct)
}
