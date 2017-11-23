package utils

import (
	"github.com/defaulteg/api/modules/core"
	"strconv"
	"errors"
	"log"
)

func SetBaseCurrency(rate core.JsonRatesData, base string) (core.JsonRatesData, error) {
	// Change rates if not base currency
	if base != "EUR" {
		log.Print(base)
		baseRate := 0.0
		for _, rate := range rate.Rates {
			if rate.Name == base {
				baseRate, _ = strconv.ParseFloat(rate.Rate, 64)
				break
			}
		}

		if baseRate == 0 {
			return rate, errors.New("Invalid base currency")
		}

		for j, _ := range rate.Rates {
			currentRate := rate.Rates[j].Rate
			currentRateFloat, _ := strconv.ParseFloat(currentRate, 64)
			newRate := (1/baseRate) * currentRateFloat
			rate.Rates[j].Rate = strconv.FormatFloat(newRate, 'f', 6, 64)
		}
	}

	return rate, nil
}


