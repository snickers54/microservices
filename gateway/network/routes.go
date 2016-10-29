package network

import "time"

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
