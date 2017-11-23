package modules

import (
	"github.com/defaulteg/api/modules/currency"
	"github.com/defaulteg/api/modules/metal"
	"github.com/defaulteg/api/modules/cryptocurrency"
	"github.com/robfig/cron"
	"log"
)

// fetcher all metal, cash, crypto prices

func StartFetcherService() {
	cron := cron.New()
	// local currencies every hour
	cron.AddFunc("@hourly", startLocalCurrencyService)
	// crypto currencies every 15 min
	cron.AddFunc("@every 15m", startCryptocurrencyService)
	// base currencies every day at 16:00 on working days
	cron.AddFunc("0 0 16 * * MON-FRI", startBaseCurrencyService)
	// metals every 5 min
	cron.AddFunc("@every 5m", startMetalService)

	cron.Start()
}


func startLocalCurrencyService() {
	if err := currency.FetchLocal(); err != nil {
		log.Fatal(err)
	} else {
		log.Println(" -  Local currency rates fetched")
	}
}

func startBaseCurrencyService() {
	if err := currency.FetchBase(); err != nil {
		log.Fatal(err)
	} else {
		log.Println(" -  Base currency rates fetched")
	}
}

func startCryptocurrencyService() {
	if err := cryptocurrency.Fetch(); err != nil {
		log.Fatal(err)
	} else {
		log.Println(" -  Cryptocurrency rates fetched")
	}
}

func startMetalService() {
	if err := metal.Fetch(); err != nil {
		log.Fatal(err)
	} else {
		log.Println(" -  Metal rates fetched")
	}
}
