package main

import (
	"gopkg.in/mgo.v2/bson"
)

type ProcessInfo struct {
	ObjectId       bson.ObjectId `bson:"_id"`
	Description    string        `json:"description" bson:"Process"`
	ExecutablePath string        `json:"path" bson:"ExecutablePath"`
}
