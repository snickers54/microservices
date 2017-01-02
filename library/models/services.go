package models

import (
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
)

type Service struct {
	ID      string  `json:"id"`
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
	return fmt.Sprintf("#%s %s - %s | %s:%s", self.ID, self.Name, self.Version.String(), self.IP, self.Port)
}

func (self *Services) Add(service Service) bool {
	for _, currentService := range *self {
		if currentService.IP == service.IP && currentService.Port == service.Port {
			return false
		}
	}
	service.ID = bson.NewObjectIdWithTime(time.Now()).String()
	log.WithField("service", service.String()).Debug("Add service to list of services.")
	*self = append(*self, &service)
	return true
}

func (self *Service) Register(URL string) []error {
	_, body, errs := gorequest.New().Post(URL).Send(*self).End()
	err := json.Unmarshal([]byte(body), self)
	if err != nil {
		errs = append(errs, err)
	}
	return errs
}
