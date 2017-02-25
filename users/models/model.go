package models

import (
	log "github.com/Sirupsen/logrus"
	"github.com/snickers54/microservices/users/services/databases"
	r "gopkg.in/gorethink/gorethink.v3"
)

const USER_DB = "users"

func Start() {
	session := databases.GetRethinkDBInstance()
	_, err := r.DBList().Contains(USER_DB).Do(func(dbExists r.Term) r.Term {
		return r.Branch(
			dbExists,
			nil,
			r.DBCreate(USER_DB),
		)
	}).Run(session)
	if err != nil {
		log.WithError(err).Fatal("Couldn't create ")
	}
}

type DriverRethinkDB struct{}

func (d *DriverRethinkDB) session() *r.Session {
	return databases.GetRethinkDBInstance()
}
