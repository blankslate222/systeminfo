package main

import (
	"fmt"
	"sync"

	"gopkg.in/mgo.v2"
)

type DBConnection struct {
	sync.Mutex
	session *mgo.Session
	err     error
}

var db *DBConnection = nil

func GetDB() *DBConnection {
	if db == nil {
		fmt.Println("instantiating db")
		db = new(DBConnection)
	}
	return db
}
func GetConnection(db *DBConnection) *mgo.Session {
	if db.session == nil {
		db.session, db.err = mgo.Dial(`mongodb://golang:golang@ds047085.mongolab.com:47085/sysinfo`)
		if db.err != nil {
			panic(db.err)
		}
	}
	return db.session
}
