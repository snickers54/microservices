package middlewares

import (
	"strings"
	"sync"

	"github.com/snickers54/microservices/gateway/context"
	"github.com/snickers54/microservices/gateway/network"
)

func Sync(c *context.AppContext) {
	nodes := network.GetCluster().Nodes
	gatewaysToReplay := network.Nodes{}
	notified := c.Request.Header.Get("gateway-notified")
	if len(notified) == 0 {
		gatewaysToReplay = nodes
	}
	gatewaysToReplay = gatewaysNotFoundInNotified(notified, nodes)
	gatewaysToReplay = excludeMyself(gatewaysToReplay)
	for _, node := range gatewaysToReplay {
		notified = notified + "," + node.IP + ":" + node.Port
	}
	c.Request.Header.Set("gateway-notified", notified)
	var wg sync.WaitGroup
	wg.Add(len(gatewaysToReplay))
	for _, node := range gatewaysToReplay {
		go network.ReplayHTTP(c, node, wg)
	}
	wg.Wait()
	c.Next()
}

func excludeMyself(gateways network.Nodes) network.Nodes {
	gatewaysToReplay := network.Nodes{}
	for _, gateway := range gateways {
		if gateway.Myself == false {
			gatewaysToReplay = append(gatewaysToReplay, gateway)
		}
	}
	return gatewaysToReplay
}

func gatewaysNotFoundInNotified(notified string, nodes network.Nodes) network.Nodes {
	gatewaysToReplay := network.Nodes{}
	for _, gateway := range strings.Split(notified, ",") {
		tempNode := network.Node{}
		tempNode.UpdateFromAddr(gateway)
		if nodes.Exists(tempNode) == false {
			gatewaysToReplay = append(gatewaysToReplay, tempNode)
		}
	}
	return gatewaysToReplay
}
