package handlers

import (
	"errors"
	"net/http"

	"github.com/snickers54/microservices/gateway/network"
)

func clusterDescribe(context *AppContext) {
	context.WriteJSON(network.GetCluster())
	context.Done()
}

func clusterRegister(context *AppContext) {
	node := network.Node{}
	context.BindJSON(&node)
	node.Status = network.STATUS_ACTIVE
	if ok := network.GetCluster().Nodes.Add(node); ok == false {
		context.Error(errors.New("This node already exists."), http.StatusConflict)
		return
	}
	context.WriteJSON(node)
	context.Done()
}
