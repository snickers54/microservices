package middlewares

import (
	"net/http"
	"strings"
	"sync"

	"github.com/snickers54/microservices/gateway/context"
	"github.com/snickers54/microservices/gateway/network"
)

func Sync(c *context.AppContext) {
	// if it's coming from a gateway we have to identify gateway that were not notified..
	nodes := network.GetCluster().Nodes
	gatewaysToReplay := network.Nodes{}
	notified := c.Request.Header.Get("gateway-notified")
	if isGateway(nodes, c.Request) {
		for _, gateway := range strings.Split(notified, ",") {
			tempNode := network.Node{}
			tempNode.UpdateFromAddr(gateway)
			if nodes.Exists(tempNode) == false {
				gatewaysToReplay = append(gatewaysToReplay, tempNode)
			}
		}
	} else {
		gatewaysToReplay = network.GetCluster().Nodes
	}

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

func isGateway(nodes network.Nodes, r *http.Request) bool {
	node := network.Node{}
	node.UpdateFromAddr(r.RemoteAddr)
	return nodes.Exists(node)
}
