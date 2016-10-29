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
		currentRoute := v.(*Route)
		for key, instance := range currentRoute.Endpoints {
			if instance.Service.String() == service.String() {
				currentRoute.Endpoints[key] = currentRoute.Endpoints[len(currentRoute.Endpoints)-1]
				currentRoute.Endpoints = currentRoute.Endpoints[:len(currentRoute.Endpoints)-1]
				break
			}
		}
		return false
	})
}

func GetRoute(path string) *Route {
	rawRoute, ok := _map.Get(path)
	if ok == false {
		return nil
	}
	return rawRoute.(*Route)
}

func InsertEndpoint(path string, service *Service) bool {
	_, ok := _map.Get(path)
	if ok == false {
		_map.Insert(path, new(Route))
	}
	rawRoute, ok := _map.Get(path)
	route := rawRoute.(*Route)
	if route.Endpoints.Exists(service) == false {
		route.Endpoints = append(route.Endpoints, &Endpoint{
			Service: service,
			Status:  STATUS_ACTIVE,
			Retries: 0,
		})
		log.Println("Registered route : ", route, " from ", service.String())
		return true
	}
	log.Println("Route : ", route, " from ", service.String(), " has already been registered.")
	return false
}
