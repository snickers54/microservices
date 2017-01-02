package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type Version struct {
	Name  string `json:"name"`
	Major uint64 `json:"major"`
	Minor uint64 `json:"minor"`
	Patch uint64 `json:"patch"`
}

func (self *Version) String() string {
	return fmt.Sprintf("%s (%d.%d.%d)", self.Name, self.Major, self.Minor, self.Patch)
}

func (self *Version) Parse(str string) {
	re := regexp.MustCompile("(?P<Major>[0-9]+).(?P<Minor>[0-9]+).(?P<Patch>[0-9]+)")
	pureString := re.FindString(str)
	list := strings.Split(pureString, ".")
	if len(list) >= 3 {
		self.Major, _ = strconv.ParseUint(list[0], 10, 64)
		self.Minor, _ = strconv.ParseUint(list[1], 10, 64)
		self.Patch, _ = strconv.ParseUint(list[2], 10, 64)
	}
	self.Name = strings.Replace(str, " ("+pureString+")", "", -1)
	log.WithField("string", str).Debug("Result of parsing semantic versionning : ", self)
}

func (self *Version) Match(semver Version) bool {
	return (self.Major == semver.Major && self.Minor == semver.Minor && self.Patch == semver.Patch)
}
