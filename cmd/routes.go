package cmd

import "net/http"

// Route represents information for url routing
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes represents a slice of Route
type Routes []Route

var routes = Routes{
	Route{
		"CaseShow",
		"GET",
		`/cases/{accNum:[a-zA-Z0-9]{5,11}}`,
		CaseShow,
	},
}
