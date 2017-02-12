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
	Major int64  `json:"major"`
	Minor int64  `json:"minor"`
	Patch int64  `json:"patch"`
}

func (self *Version) String() string {
	return fmt.Sprintf("%s (%d.%d.%d)", self.Name, self.Major, self.Minor, self.Patch)
}

func (self *Version) Parse(str string) {
	re := regexp.MustCompile("(?P<Major>[0-9]+).(?P<Minor>[0-9|x]+).(?P<Patch>[0-9|x]+)")
	pureString := re.FindString(str)
	list := strings.Split(pureString, ".")
	log.Info(list)
	if len(list) >= 3 {
		self.Major, _ = strconv.ParseInt(list[0], 10, 64)
		self.Minor, _ = strconv.ParseInt(list[1], 10, 64)
		if list[1] == "x" {
			self.Minor = -1
		}
		self.Patch, _ = strconv.ParseInt(list[2], 10, 64)
		if list[2] == "x" {
			self.Patch = -1
		}
	}
	self.Name = strings.Replace(str, " ("+pureString+")", "", -1)
	log.WithField("string", str).Debug("Result of parsing semantic versionning : ", self)
}

func (self *Version) Match(semver Version) bool {
	minorTruth := self.Minor == semver.Minor || semver.Minor < 0 || self.Minor < 0
	patchTruth := self.Patch == semver.Patch || semver.Patch < 0 || self.Patch < 0
	return self.Major == semver.Major && minorTruth && patchTruth
}
