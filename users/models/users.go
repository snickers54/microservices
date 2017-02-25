package models

import (
	log "github.com/Sirupsen/logrus"
	r "gopkg.in/gorethink/gorethink.v3"
)

const USER_TABLE = "members"

type User struct {
	ID       string `gorethink:"id,omitempty"`
	Username string `gorethink:"username"`
	DriverRethinkDB
}

type Users []User

func (self *User) GetUsers() Users {
	rows, err := r.DB(USER_DB).Table(USER_TABLE).Run(self.session())
	if err != nil {
		log.WithError(err).Error("GetUsers failed to execute query.")
		return Users{}
	}
	var users Users
	err2 := rows.All(&users)
	if err2 != nil {
		log.WithError(err).Error("GetUsers failed to retrieve data.")
		return Users{}
	}
	return users
}
