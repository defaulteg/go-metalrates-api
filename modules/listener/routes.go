package listener

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route


var routes = Routes {

	// Index Page
	Route {
		"Index",		//name
		"GET",			//method
		"/",			//pattern
		Index,			//func
	},

	// Base currency rates ALL, PARTICULAR
	Route {
		"Base Rates: ALL/PARTICULAR; Latest",
		"GET",
		"/baserates/latest",
		BaseRatesLatestHandler,
	},

	// Base currency rates HISTORICAL
	Route {
		"Base Rates: ALL/PARTICULAR; Historical",
		"GET",
		"/baserates/historical",
		BaseRatesHistoricalHandler,
	},

	// Local currency rates ALL
	Route {
		"LocalRates",
		"GET",
		"/localrates/latest",
		LocalRatesLatestHandler,
	},

	Route {
		"LocalRates",
		"GET",
		"/localrates/historical",
		LocalRatesLatestHistoricalHandler,
	},
}
