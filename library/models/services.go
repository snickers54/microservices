package models

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
)

type Service struct {
	Name    string  `json:"name"`
	IP      string  `json:"ip"`
	Port    string  `json:"port"`
	Version Version `json:"version"`
	Status  string  `json:"status"`
}
type Services []*Service

func (self Services) Len() int           { return len(self) }
func (self Services) Swap(i, j int)      { self[i], self[j] = self[j], self[i] }
func (self Services) Less(i, j int) bool { return self[i].Name < self[j].Name }

func (self *Service) String() string {
	return fmt.Sprintf("%s - %s | %s:%s", self.Name, self.Version.String(), self.IP, self.Port)
}

func (self *Services) Add(service Service) bool {
	for _, currentService := range *self {
		if currentService.IP == service.IP && currentService.Port == service.Port {
			return false
		}
	}
	log.WithField("service", service.String()).Debug("Add service to list of services.")
	*self = append(*self, &service)
	return true
}

func (self *Service) Register(URL string) []error {
	_, _, errs := gorequest.New().Post(URL).Send(*self).End()
	return errs
}