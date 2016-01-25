package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"Programinfo",
		"GET",
		"/programs",
		Programinfo,
	},
	Route{
		"Processinfo",
		"GET",
		"/processes",
		Processinfo,
	},
	Route{
		"HostMachineinfo",
		"GET",
		"/hostmachineinfo",
		HostMachineinfo,
	},
	Route{
		"Generateinfo",
		"POST",
		"/systeminfo",
		Generateinfo,
	},
}
