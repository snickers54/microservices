package handlers

import "github.com/snickers54/microservices/gateway/context"

func nodePing(c *context.AppContext) {
	c.WriteJSON(context.Payload{})
}
