package network

import "time"

const STATUS_ACTIVE = "ACTIVE"
const STATUS_INACTIVE = "INACTIVE"
const STATUS_UNSTABLE = "UNSTABLE"

type routeInstance struct {
	Service *Service `json:"service"`
	Status  string   `json:"status"`
	Retries uint     `json:"retries"`
}

type routeInstances []routeInstance

type endpoint struct {
	RouteInstances routeInstances `json:"services"`
	Stats          endpointStats  `json:"stats"`
}

type endpointStats struct {
	Calls           uint      `json:"number_calls"`
	SuccessfulCalls uint      `json:"successful_calls"`
	FailedCalls     uint      `json:"failed_calls"`
	AverageTime     time.Time `json:"average_latency"`
}

func (self *routeInstances) Exists(service *Service) bool {
	for _, instanceService := range *self {
		if instanceService.Service.String() == service.String() {
			return true
		}
	}
	return false
}
