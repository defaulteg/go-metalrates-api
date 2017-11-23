package listener

import (
	"net/http"
	"fmt"
	"net/url"
	"github.com/defaulteg/api/modules/core"
	"strings"
	"github.com/gorilla/mux"
	"github.com/gorilla/context"
	"regexp"
)


func Index(w http.ResponseWriter, r *http.Request) {
	pathToIndex := r.URL.Path[1:]
	http.ServeFile(w, r, pathToIndex)
}

func BaseRatesLatestHandler(w http.ResponseWriter, r *http.Request) {
	request, _ := url.ParseQuery(r.URL.RawQuery)

	// If no parameters specified
	if len(request) == 0 {

		// Get rates from database
		getJsonData(w, core.QueryGetAll, core.QueryGetLatest, core.JsonOneObject, nil, "", "", "EUR")

	} else {
		var base string
		if base = request.Get("base"); base == "" {
			base = "EUR"
		}

		// If currencies parameter specified
		if symbols := request.Get("symbols"); symbols != "" {
			currencies := strings.Split(symbols, ",")
			getJsonData(w, core.QueryGetParticular, core.QueryGetLatest, core.JsonOneObject, currencies, "", "", base)
		} else {
			getJsonData(w, core.QueryGetAll, core.QueryGetLatest, core.JsonOneObject, nil, "", "", base)
		}
	}
}

func BaseRatesHistoricalHandler(w http.ResponseWriter, r *http.Request) {
	request, _ := url.ParseQuery(r.URL.RawQuery)

	// If no parameters specified
	if len(request) == 0 {
		returnError(w, http.StatusBadRequest, core.ERR_SPECIFY_DATE)
		return
	} else {
		if !checkRequestDateParam(w, r) {
			return
		}


		if symbols := request.Get("symbols"); symbols != "" {
			// Return Particular
			currencies := strings.Split(symbols, ",")
			getJsonData(w, core.QueryGetParticular, core.QueryGetHistorical, core.JsonSliceObject, currencies, request.Get("from"), request.Get("to"), "EUR")
		} else {
			// Return ALL
			getJsonData(w, core.QueryGetAll, core.QueryGetHistorical,core.JsonSliceObject, nil, request.Get("from"), request.Get("to"), "EUR")
		}
	}
}

func LocalRatesLatestHandler(w http.ResponseWriter, r *http.Request) {
	request, _ := url.ParseQuery(r.URL.RawQuery)

	if len(request) == 0 {

		GetLocalJsonData(w, "all", "", "", false)

	} else {
		if rate := request.Get("rate"); rate != "" {
			r := regexp.MustCompile(`^(\bsell\b|\bbuy\b|\ball\b)$`)

			if !(r.MatchString(rate)) {
				returnError(w, http.StatusBadRequest, core.ERR_INVALID_DATE_FORMAT)
				return
			}

			GetLocalJsonData(w, rate, "", "", false)
		}
	}
}

func LocalRatesLatestHistoricalHandler(w http.ResponseWriter, r *http.Request) {

	request, _ := url.ParseQuery(r.URL.RawQuery)

	// If no parameters specified
	if len(request) == 0 {
		returnError(w, http.StatusBadRequest, core.ERR_SPECIFY_DATE)
		return
	} else {
		if !checkRequestDateParam(w, r) {
			return
		}

		if rate := request.Get("rate"); rate != "" {
			r := regexp.MustCompile(`^(\bsell\b|\bbuy\b|\ball\b)$`)

			if !(r.MatchString(rate)) {
				returnError(w, http.StatusBadRequest, core.ERR_INVALID_RATE_FORMAT)
				return
			}

			GetLocalJsonData(w, rate, request.Get("from"), request.Get("to"), true)
		} else {
			GetLocalJsonData(w, "all", request.Get("from"), request.Get("to"), true)
		}
	}
}

func MetalRatesLatestHandler(w http.ResponseWriter, r *http.Request) {
	request, _ := url.ParseQuery(r.URL.RawQuery)

	if len(request) == 0 {
		GetMetalJsonData(w, nil, "", "", "", false)
	} else {
		weight := request.Get("weight")
		var metals []string

 		if symbols := request.Get("symbols"); symbols != "" {
			metals = strings.Split(symbols, ",")
		}

		GetMetalJsonData(w, metals, weight, "", "", false)
	}
}

func MetalRatesHistoricalHandler(w http.ResponseWriter, r *http.Request) {

}

// TODO: metal rates!
// LATEST: symbols, weight | HISTORICAL: symbols

// TODO: cryptocurrency rates
// LATEST: symbols, market, pair | HISTORICAL: symbols, market

// TODO: Migration script

//http://localhost:8080/localrates/all?key1=val1,val5&key2=val2,val3
func LocalRatesAllHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	fmt.Println(mux.Vars(r))
	fmt.Println(r.URL.Query())
	fmt.Println(r.URL.Query().Get("key1"))
	fmt.Fprintf(w, "all: %v and %v", variables, r.URL.Query())
}

func LocalRatesAllBaseHandler(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	v := context.GetAll(r)
	fmt.Printf("%v", len(v))
	fmt.Fprintf(w, "all 1: %v and %v", variables, r.URL.Query())
}

func LocalRatesAllBaseHandler2(w http.ResponseWriter, r *http.Request) {
	variables := mux.Vars(r)
	v := context.GetAll(r)
	fmt.Println(v)
	fmt.Println(r.URL.Query())
	fmt.Fprintf(w, "all 2: %v and %v", variables,r.URL.Query())
}

func parseRequest(values map[string][] string) {
	for value := range values {

		fmt.Println("key: " + value + " value: ")
	}
}
