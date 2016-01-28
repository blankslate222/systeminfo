package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func Programinfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("program info called")
	db := GetDB()
	db.Lock()
	defer db.Unlock()
	sess := GetConnection(db)
	collection := sess.DB("sysinfo").C("programs")
	length, _ := collection.Count()
	result := make([]ProgramInfo, length)
	qry := collection.Find(bson.M{})
	// fmt.Println(qry)
	err := qry.All(&result)
	if err != nil {
		fmt.Println("error in mongo find query")
	}
	//fmt.Println(result)
	j, _ := json.Marshal(result)
	//fmt.Fprint(w, j)
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func Processinfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("process info called")
	db := GetDB()
	db.Lock()
	defer db.Unlock()
	sess := GetConnection(db)
	collection := sess.DB("sysinfo").C("processes")
	length, _ := collection.Count()
	result := make([]ProcessInfo, length)
	qry := collection.Find(bson.M{})
	// fmt.Println(qry)
	err := qry.All(&result)
	if err != nil {
		fmt.Println("error in mongo find query")
	}
	//fmt.Println(result)
	j, _ := json.Marshal(result)
	//fmt.Fprint(w, j)
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func HostMachineinfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("host info called")
	db := GetDB()
	db.Lock()
	defer db.Unlock()
	sess := GetConnection(db)
	collection := sess.DB("sysinfo").C("machineinfo")
	length, _ := collection.Count()
	result := make([]interface{}, length)
	qry := collection.Find(bson.M{})
	// fmt.Println(qry)
	err := qry.All(&result)
	if err != nil {
		fmt.Println("error in mongo find query")
	}
	//fmt.Println(result)
	j, _ := json.Marshal(result)
	//fmt.Fprint(w, j)
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func Generateinfo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("generate info called")
	progs := DetectInstalledPrograms()

	if progs == true {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
