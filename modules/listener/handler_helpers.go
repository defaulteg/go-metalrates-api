package listener

import (
	"net/http"
	"encoding/json"
	"log"
	"github.com/defaulteg/api/modules/core"
	"github.com/defaulteg/api/utils"
	"regexp"
	"net/url"
)

func getJsonData(w http.ResponseWriter, queryCountType, queryType, sliceOrOneObject int, currencies []string, dateFrom, dateTo, base string) {
	// Get rates from database
	if rates, err := core.GetJsonRates(queryCountType, queryType, currencies, dateFrom, dateTo); err != nil {
		log.Fatal(err)
	} else {
		// Set base currency
		for i, _ := range rates {
			rates[i].Base = base

			rates[i], err = utils.SetBaseCurrency(rates[i], base)
			if err != nil {
				returnError(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

			// If returned data is slice with one array of data
		if sliceOrOneObject == core.JsonOneObject {
			if err := json.NewEncoder(w).Encode(rates[0]); err != nil {
				log.Fatal(err)
			}
		} else {
			// If returned data is slice with multiple data objects -> return response with JSON object
			if err := json.NewEncoder(w).Encode(rates); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func GetLocalJsonData(w http.ResponseWriter, rateType, dateFrom, dateTo string, historical bool) {
	if rates, err := core.GetLocalCurrenciesJson(rateType, dateFrom, dateTo, historical); err != nil {
		returnError(w, http.StatusBadRequest, err.Error())
		return
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(rates); err != nil {
			log.Fatal(err)
		}
	}
}

func GetMetalJsonData(w http.ResponseWriter, symbols []string, weight, dateFrom, dateTo string, historical bool) {
	if rates, err := core.GetMetalsJson(symbols, dateFrom, dateTo, historical); err != nil {
		returnError(w, http.StatusBadRequest, err.Error())
		return
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		if (weight == "") {
			if err := json.NewEncoder(w).Encode(rates); err != nil {
				log.Fatal(err)
			}
		} else {
			// calculate weight
			rates := utils.CalculateWeight(rates)
			if err := json.NewEncoder(w).Encode(rates); err != nil {
				log.Fatal(err)
			}
		}
	}
}


func returnError(w http.ResponseWriter, status int, description string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	err := core.ErrorJson {
		Type 		: "error",
		Description : description,
	}

	if err := json.NewEncoder(w).Encode(err); err != nil {
		log.Fatal(err)
	}
}

func checkRequestDateParam(w http.ResponseWriter, r *http.Request) bool {
	request, _ := url.ParseQuery(r.URL.RawQuery)

	dateFrom := request.Get("from")
	dateTo := request.Get("to")
	if dateFrom == "" || dateTo == "" {
		returnError(w, http.StatusBadRequest, core.ERR_NOT_SPECIFIED_DATE)
		return false
	}

	reg := regexp.MustCompile(`^[0-9]{4}(\/|\-)[0-9]{1,2}(\/|\-)[0-9]{1,2}$`)

	if !(reg.MatchString(dateFrom) && reg.MatchString(dateTo)) {
		// not correct date format
		returnError(w, http.StatusBadRequest, core.ERR_INVALID_DATE_FORMAT)
		return false
	}

	return true
}

