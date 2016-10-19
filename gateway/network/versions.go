package network

import "fmt"

type Version struct {
	Name  string `json:"name"`
	Major uint   `json:"major"`
	Minor uint   `json:"minor"`
	Patch uint   `json:"patch"`
}

func (self *Version) String() string {
	return fmt.Sprintf("%s (%d.%d.%d)", self.Name, self.Major, self.Minor, self.Patch)
}
