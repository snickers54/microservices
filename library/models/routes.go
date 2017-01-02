package models

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

type Route struct {
	Endpoints Endpoints  `json:"services"`
	Stats     RouteStats `json:"stats"`
}

type RouteStats struct {
	Calls           uint          `json:"number_calls"`
	SuccessfulCalls uint          `json:"successful_calls"`
	FailedCalls     uint          `json:"failed_calls"`
	AverageTime     time.Duration `json:"average_latency"`
}

func (self *Route) GetValidEndpoints() Endpoints {
	endpoints := Endpoints{}
	for _, endpoint := range self.Endpoints {
		if endpoint.Status == STATUS_ACTIVE {
			endpoints = append(endpoints, endpoint)
		}
	}
	return endpoints
}

func (self *Route) Register(route, URL string, service *Service) []error {
	_, _, errs := gorequest.New().Post(URL).Send(gin.H{
		"route":   route,
		"service": service,
	}).End()
	return errs
}
