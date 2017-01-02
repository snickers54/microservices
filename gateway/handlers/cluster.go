package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/snickers54/microservices/gateway/context"
	"github.com/snickers54/microservices/gateway/network"
	"github.com/snickers54/microservices/library/models"
)

func clusterDescribe(c *context.AppContext) {
	c.WriteJSON(network.GetCluster())
}

func clusterRegister(c *context.AppContext) {
	node := network.Node{}
	if ok := c.BindJSON(&node); !ok {
		return
	}
	node.Status = models.STATUS_ACTIVE
	node.IP = strings.Split(c.Request.RemoteAddr, ":")[0]
	if ok := network.GetCluster().Nodes.Add(node); ok == false {
		c.Error(errors.New("This node already exists."), http.StatusConflict)
		return
	}
	c.WriteJSON(node)
}
