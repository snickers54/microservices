package middlewares

import (
	"time"

	"github.com/snickers54/microservices/gateway/context"
	"github.com/snickers54/microservices/gateway/network"
)

func travelling_average(average time.Duration, startingTime time.Time, count uint) time.Duration {
	newCount := int64(count)
	value := (average.Nanoseconds() * newCount) + time.Now().Sub(startingTime).Nanoseconds()
	value = value / (newCount + 1)
	return time.Duration(value)
}

func StatsReplay(c *context.AppContext) {
	rawTime, _ := c.Get("route_time")
	startingTime := rawTime.(time.Time)
	rawRoute, ok := c.Get("route")
	if ok == false {
		return
	}
	route := rawRoute.(*network.Route)
	route.Stats.Calls += 1
	route.Stats.AverageTime = travelling_average(route.Stats.AverageTime, startingTime, route.Stats.SuccessfulCalls)
	if len(c.Responses) <= 0 {
		route.Stats.FailedCalls += 1
	} else {
		route.Stats.SuccessfulCalls += 1
	}
}

func WriteReplay(c *context.AppContext) {
	// for now we consider only one ..
	if len(c.Responses) <= 0 {
		return
	}
	response := c.Responses[0]
	data := []byte{}
	c.Writer.WriteHeader(response.StatusCode)
	response.Body.Read(data)
	c.Writer.Write(data)
}

func CloseBody(c *context.AppContext) {
	for _, response := range c.Responses {
		response.Body.Close()
	}
	c.Next()
}
