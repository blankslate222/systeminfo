package main

import (
	"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

func GetConnection() *mgo.Session {
	session, err := mgo.Dial(`mongodb://golang:golang@ds047085.mongolab.com:47085/sysinfo`)
	if err != nil {
		panic(err)
	}
	return session
}
