package cryptocurrency

import (
	"github.com/defaulteg/api/modules/core"
	"os/exec"
	"strings"
	//"fmt"
	"errors"
	//"github.com/defaulteg/api/utils"
	//"time"
	"sync"
	"runtime"
)

//TODO: Get historical rates
//TODO: migrate db - add all selectors and elements from config?
//TODO: use chan model instead??

var (
	rates []core.CryptoRate							// List of fetched rates
	pairMap, coinMap, marketMap map[string]int		// Map of id's
	wg sync.WaitGroup								// Wait for all workers to fetch rates
	mutex = &sync.Mutex{}							// Simple mutex for db push
)

func Fetch() error {

	if sources, err := core.GetSources(core.CryptocurrencyCategory); err != nil {
		return err
	} else {

		// Get coint id's
		coinMap, err = core.GetElementIds(core.CryptocurrencyTypeId, core.QueryWhere, "")        	//map with [coinName] = id
		if err != nil {
			return err
		}

		// Get market id's
		marketMap, err = core.GetElementIds(core.CryptocurrencyTypeId, core.QueryBase, "markets")      //map with [marketName] = id
		if err != nil {
			return err
		}

		// Get pair id's
		pairMap, err = core.GetElementIds(core.CryptocurrencyTypeId, core.QueryBase, "pairs")         //map with [pairName] = id
		if err != nil {
			return err
		}

		rates = make([]core.CryptoRate, 0)

		// For each source fetch rates to cryptoRate object slice
		for _, source := range sources {
			wg.Add(1)
			go fetchWorker(source)
		}

		// Wait for all workers
		wg.Wait()

		// Push sources to database
		if err := core.PushToDatabaseCryptoData(rates, core.CryptocurrencyTable); err != nil {
			return err
		}
	}

	return nil
}

func fetchWorker(source core.Source) error {

	//defer utils.ExecutionTime("fetcher", time.Now())

	defer wg.Done()

	// Fetch data from source
	cmd := exec.Command("phantomjs", core.PathToCryptocurrencyFetcher, source.Site)

	if res, err := cmd.Output(); err != nil {
		return errors.New("Cannot execute phantomjs command. " + err.Error());
	} else {
		pjsOutputTemp := strings.TrimSpace(string(res))              // Trim all spaces from phantomjs response
		pjsOutput := strings.Split(pjsOutputTemp, ",")               // Split response string to slice with "," divider

		for i := 1; i < len(pjsOutput); i += 3 {
			var rate core.CryptoRate
			rate.Name = pjsOutput[0]                    			// Zero output is cryptocurrency name
			rate.Id = coinMap[pjsOutput[0]]                			// Id of cryptocurrency name
			rate.MarketId = marketMap[pjsOutput[i]]        			// Id of market
			rate.PairId = pairMap[pjsOutput[i + 1]]        			// Id of pair
			rate.Rate = pjsOutput[i + 2]                			// Rate of cryptocurrency

			mutex.Lock()
			rates = append(rates, rate)
			mutex.Unlock()

			runtime.Gosched()
		}

	}

	return nil
}



