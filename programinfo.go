package main

import (
	"gopkg.in/mgo.v2/bson"
)

type ProgramInfo struct {
	ObjectId      bson.ObjectId `bson:"_id"`
	Name          string        `json:"name" bson:"DisplayName"`
	Version       string        `json:"version" bson:"DisplayVersion"`
	LatestVersion string        `json:"latestversion" bson:"LatestVersion"`
	NeedsUpdate   bool          `json:"needsupdate" bson:"NeedsUpdate"`
}
