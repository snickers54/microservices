package network

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	radix "github.com/armon/go-radix"
	"github.com/snickers54/microservices/library/models"
)

var _map *radix.Tree

type radixTree struct{ *radix.Tree }

func init()                                          { _map = radix.New() }
func (self *radixTree) MarshalJSON() ([]byte, error) { return json.Marshal(self.ToMap()) }
func Len() int                                       { return _map.Len() }
func ToMap() map[string]interface{}                  { return _map.ToMap() }

func RemoveInstance(service *models.Service) {
	_map.Walk(func(s string, v interface{}) bool {
		currentRoute := v.(*models.Route)
		for key, instance := range currentRoute.Endpoints {
			if instance.Service.String() == service.String() {
				log.WithField("service", service.String()).Debug("Removing a service from radix tree.")
				currentRoute.Endpoints[key] = currentRoute.Endpoints[len(currentRoute.Endpoints)-1]
				currentRoute.Endpoints = currentRoute.Endpoints[:len(currentRoute.Endpoints)-1]
				break
			}
		}
		return false
	})
}

func GetRoute(path string) *models.Route {
	rawRoute, ok := _map.Get(path)
	if ok == false {
		return nil
	}
	return rawRoute.(*models.Route)
}

func InsertRoute(path string, service *models.Service) bool {
	_, ok := _map.Get(path)
	if ok == false {
		_map.Insert(path, new(models.Route))
	}
	rawRoute, ok := _map.Get(path)
	route := rawRoute.(*models.Route)
	contextLogger := log.WithFields(log.Fields{"route": route, "from": service.String()})
	if route.Endpoints.Exists(service) == false {
		route.Endpoints = append(route.Endpoints, &models.Endpoint{
			Service: service,
			Status:  models.STATUS_ACTIVE,
			Retries: 0,
		})
		contextLogger.Info("Route registered.")
		return true
	}
	contextLogger.Info("Route already exists.")
	return false
}
