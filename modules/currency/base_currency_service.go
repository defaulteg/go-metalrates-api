package currency

import (
	"github.com/defaulteg/api/modules/core"
	"os"
	"encoding/csv"
	"bufio"
	"github.com/defaulteg/api/utils"
	"strings"
)

func FetchBase() error {

	// Get id of elements from database for further push to it
	var currencyMap map[string]int
	var err error

	if currencyMap, err = core.GetElementIds(core.BaseCurrencyTypeId, core.QueryWhere, ""); err != nil {
		return err
	}

	// Download csv rates data file
	if err := utils.DownloadFile(core.PathToSaveCurrencyZipCsv, core.CurrencyDownloadUrl); err != nil {
		return err
	}

	// Unzip csv file
	if err := utils.UnzipFile(core.PathToSaveCurrencyZipCsv, core.PathToUnzipCurrencyCsv); err != nil {
		return err
	}

	// Open csv
	f, err := os.Open(core.PathToCurrencyCsv)
	if err != nil {
		return err
	}

	// Read data from csv to slice of rate objects
	rates := make([]core.Selector, 0)
	reader := csv.NewReader(bufio.NewReader(f))

	// Read two records to two slices
	currencies, _ := reader.Read()
	values, _ := reader.Read()

	var rate core.Selector

	// Append all elements from records to rate slice. Zero element is date
	// Selector = next currency element with CUR, RATE, ID
	for i := 1; i < len(currencies) - 1; i++ {
		rate.Name = strings.TrimSpace(currencies[i])
		rate.Rate = strings.TrimSpace(values[i])
		rate.ElementId = currencyMap[rate.Name]

		rates = append(rates, rate)
	}

	// Make pseudo sources slice with selector slice of currency rates for pushing into database
	var s core.Source
	sources := make([]core.Source, 0)
	s.Selectors = rates
	s.Site = "csv"
	sources = append(sources, s)

	// Push to database
	if err := core.PushToDatabase(sources, core.BaseCurrencyTable); err != nil {
		return err
	}

	return nil
}

