package network

import "fmt"

type Service struct {
	Name    string  `json:"name"`
	DNS     string  `json:"dns_name"`
	IP      string  `json:"ip"`
	Port    string  `json:"port"`
	Version Version `json:"version"`
}
type Services []Service

func (self *Service) String() string {
	return fmt.Sprintf("%s - %s | %s(%s:%s)", self.Name, self.Version.String(), self.DNS, self.IP, self.Port)
}
