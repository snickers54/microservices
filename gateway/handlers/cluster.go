package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/snickers54/microservices/gateway/context"
	"github.com/snickers54/microservices/gateway/network"
)

func clusterDescribe(c *context.AppContext) {
	c.WriteJSON(network.GetCluster())
}

func clusterRegister(c *context.AppContext) {
	node := network.Node{}
	c.BindJSON(&node)
	node.Status = network.STATUS_ACTIVE
	node.IP = strings.Split(c.Request.RemoteAddr, ":")[0]
	if ok := network.GetCluster().Nodes.Add(node); ok == false {
		c.Error(errors.New("This node already exists."), http.StatusConflict)
		return
	}
	c.WriteJSON(node)
}
