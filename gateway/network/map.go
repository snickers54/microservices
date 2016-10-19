package network

import (
	"encoding/json"
	"log"

	"github.com/armon/go-radix"
)

var _map *radix.Tree

type radixTree struct{ *radix.Tree }

func init()                                          { _map = radix.New() }
func (self *radixTree) MarshalJSON() ([]byte, error) { return json.Marshal(self.ToMap()) }
func Len() int                                       { return _map.Len() }
func ToMap() map[string]interface{}                  { return _map.ToMap() }

func RemoveInstance(service *Service) {
	_map.Walk(func(s string, v interface{}) bool {
		currentEndpoint := v.(*endpoint)
		for key, instance := range currentEndpoint.RouteInstances {
			if instance.Service.String() == service.String() {
				currentEndpoint.RouteInstances[key] = currentEndpoint.RouteInstances[len(currentEndpoint.RouteInstances)-1]
				currentEndpoint.RouteInstances = currentEndpoint.RouteInstances[:len(currentEndpoint.RouteInstances)-1]
				break
			}
		}
		return false
	})
}

func InsertRoute(route string, service *Service) bool {
	_, ok := _map.Get(route)
	if ok == false {
		_map.Insert(route, new(endpoint))
	}
	rawEndpoint, ok := _map.Get(route)
	endpoint := rawEndpoint.(*endpoint)
	if endpoint.RouteInstances.Exists(service) == false {
		endpoint.RouteInstances = append(endpoint.RouteInstances, routeInstance{
			Service: service,
			Status:  STATUS_ACTIVE,
		})
		log.Println("Registered route : ", route, " from ", service.String())
		return true
	}
	log.Println("Route : ", route, " from ", service.String(), " has already been registered.")
	return false
}
