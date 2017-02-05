package middlewares

import (
	"io/ioutil"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/snickers54/microservices/gateway/context"
	"github.com/snickers54/microservices/library/models"
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
	route := rawRoute.(*models.Route)
	route.Stats.Calls++
	route.Stats.AverageTime = travelling_average(route.Stats.AverageTime, startingTime, route.Stats.SuccessfulCalls)
	if len(c.Responses) <= 0 {
		route.Stats.FailedCalls++
	} else {
		route.Stats.SuccessfulCalls++
	}
}

func WriteReplay(c *context.AppContext) {
	// for now we consider only one ..
	log.WithField("responses", c.Responses).Info("array of responses")
	if len(c.Responses) <= 0 {
		log.WithField("responses.length", 0).Warn("No responses from replay ?!")
		return
	}
	response := c.Responses[0]
	data := []byte{}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		response.StatusCode = 500
	}
	c.Writer.WriteHeader(response.StatusCode)
	c.Writer.Write(data)
}

func CloseBody(c *context.AppContext) {
	for _, response := range c.Responses {
		response.Body.Close()
	}
}
