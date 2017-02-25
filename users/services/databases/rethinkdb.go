package databases

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	r "gopkg.in/gorethink/gorethink.v3"
)

var rethinkSession *r.Session

func rethinkConnect() {
	session, err := r.Connect(r.ConnectOpts{
		Address: viper.GetString("database.host") + ":" + viper.GetString("database.port"),
	})
	if err != nil {
		log.WithError(err).Fatal("Couldn't connect to rethinkDB")
	}
	rethinkSession = session
}

func GetRethinkDBInstance() *r.Session {
	if rethinkSession == nil {
		rethinkConnect()
	}
	return rethinkSession
}
