package network

const STATUS_ACTIVE = "ACTIVE"
const STATUS_INACTIVE = "INACTIVE"
const NB_RETRIES uint = 3

type Endpoint struct {
	Service *Service `json:"service"`
	Status  string   `json:"status"`
	Retries uint     `json:"retries"`
}

type Endpoints []*Endpoint

func (self *Endpoints) Exists(service *Service) bool {
	for _, instanceService := range *self {
		if instanceService.Service.String() == service.String() {
			return true
		}
	}
	return false
}
func (self *Endpoints) FindByVersion(semver Version) Endpoints {
	list := Endpoints{}

	for _, route := range *self {
		if semver.Match(route.Service.Version) {
			list = append(list, route)
		}
	}
	return list
}
