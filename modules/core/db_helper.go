package core

import (
	"github.com/defaulteg/api/database"
	"errors"
	"database/sql"
	"fmt"
)


// Gets list of metals with info: metal id, name, selector, site
// Gets all sources for category and fill selector slice for each
func GetSources(category string) ([]Source, error) {

	sources := make([]Source, 0)

	stmtOut, err := database.Instance.Query(`
		SELECT site, id
		FROM sources
		WHERE type_id = (
		SELECT id
		FROM type
		WHERE name = ?)`, category)

	if err != nil {
		return nil, errors.New("Cannot execute query: 'Get source list'")
	}
	defer stmtOut.Close()

	for stmtOut.Next() {

		var source Source

		err := stmtOut.Scan(&source.Site, &source.Id)
		if err != nil {
			return nil, errors.New("Error in query 'Get source list': Cannot get next row")
		}

		// For cryptocurrencies selectors must be omitted
		if category != CryptocurrencyCategory {
			if err := source.getSelectors(); err != nil {
				return nil, err
			}
		}

		sources = append(sources, source)
	}

	if err = stmtOut.Err(); err != nil {
		return nil, err
	}

	return sources, nil
}


// Gets source site and selector for metal by id
func (s *Source) getSelectors() error {

	stmtOut, err := database.Instance.Query(`
		SELECT name, element_id
		FROM selectors
		LEFT JOIN sources
		ON selectors.source_id = sources.id
		WHERE sources.id = ?`, s.Id)

	if err != nil {
		return errors.New("Cannot execute query: 'Get selector list'")
	}

	defer stmtOut.Close()

	for stmtOut.Next() {
		var temp Selector

		err := stmtOut.Scan(&temp.Name, &temp.ElementId)
		if err != nil {
			return errors.New("Error in query 'Get selector list': Cannot get next row")
		}

		s.Selectors = append(s.Selectors, temp)
	}

	return nil
}

func PushToDatabase(sources []Source, tableName string) error {
	//push newly fetched rates to database

	sqlStr := "INSERT INTO " + tableName + "(element_id, rate) VALUES " //metal_id and rate
	var values = []interface{}{}

	for _, source := range sources {
		for _, selector := range source.Selectors {
			sqlStr += "(?, ?),"
			values = append(values, selector.ElementId, selector.Rate)
		}
	}

	if err := pushValuesToDatabase(values, sqlStr); err != nil {
		return err
	}

	return nil
}

func PushToDatabaseCryptoData(rates []CryptoRate, tableName string) error {

	sqlStr := "INSERT INTO " + tableName + "(element_id, rate, pair_id, market_id) VALUES "
	var values = []interface{}{}

	for _, rate := range rates {
		sqlStr += "(?, ?, ?, ?),"
		values = append(values, rate.Id, rate.Rate, rate.PairId, rate.MarketId)
	}

	if err := pushValuesToDatabase(values, sqlStr); err != nil {
		return err
	}

	return nil
}

func PushToDatabaseUnique(rates []LocalRate, tableName string) error {

	sqlStr := "INSERT INTO " + tableName + "(element_id, rate_buy, rate_sell, bank_id)" + " VALUES "
	var values = []interface{}{}

	for _, rate := range rates {
		sqlStr += "(?, ?, ?, ?),"
		values = append(values, rate.ElementId, rate.RateBuy, rate.RateSell, rate.BankId)
	}

	if err := pushValuesToDatabase(values, sqlStr); err != nil {
		return err
	}

	return nil
}

func pushValuesToDatabase(values []interface{}, sqlStr string) error {
	sqlStr = sqlStr[0:len(sqlStr)-1]
	stmt, err := database.Instance.Prepare(sqlStr)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(values...); err != nil {
		return err
	}

	return nil
}


func GetElementIds(elementTypeId int, queryType int, table string) (map[string]int, error) {

	// Fetch currency names and ids to MAP
	var dataMap map[string]int
	dataMap = make(map[string]int)

	var stmtOut *sql.Rows
	var err error

	if queryType == QueryBase {
		stmtOut, err = database.Instance.Query("SELECT id, name FROM " + table)
	} else if queryType == QueryWhere {
		stmtOut, err = database.Instance.Query("SELECT id, name FROM elements WHERE type_id = ?", elementTypeId)
	} else {
		// Gets compliant map table like Map[BANK_ID] = SELECTOR.NAME
		stmtOut, err = database.Instance.Query(`
			SELECT bank_id, name FROM selectorAndBankMatchers
			LEFT JOIN selectors
			ON selectorAndBankMatchers.selector_id = selectors.id`)
	}

	if err != nil {
		return nil, err
	}
	defer stmtOut.Close()

	// For each query response row
	for stmtOut.Next() {

		var id int
		var name string

		err = stmtOut.Scan(&id, &name)
		if err != nil {
			return nil, err
		}

		// Save data name and id. e.g. currencyMap["USD"] = 15
		dataMap[name] = id
	}

	return dataMap, nil
}

func GetJsonRates(queryCountType, queryType int, currencies []string, dateFrom, dateTo string) ([]JsonRatesData, error) {

	particularCurrencies := ""
	whereParam := ""

	// Get info about particular currencies
	if queryCountType == QueryGetParticular {
		particularCurrencies = " AND e.name IN ("

		for _, currency := range currencies {
			particularCurrencies += "'" + currency + "',"
		}
		// Trim last ","
		particularCurrencies = particularCurrencies[0:len(particularCurrencies)-1]
		particularCurrencies += ")"

	}

	if queryType == QueryGetLatest {
		whereParam = `b.created_at =
					  (
						SELECT DISTINCT max(created_at)
						FROM baseCurrencyRates b
					  )`
	}

	if queryType == QueryGetHistorical {
		whereParam = "b.created_at BETWEEN '" + dateFrom + "' AND '"+ dateTo + " 23:59:59'" //ORDER BY b.created_at DESC
	}


	// Prepare query string
	str := `
		SELECT e.name, b.rate, b.created_at
		FROM baseCurrencyRates b
		JOIN elements e
		ON b.element_id = e.id
		WHERE ` + whereParam + particularCurrencies

	// Run query and get data
	if queryType == QueryGetHistorical {
		// If historical sequence is fetched sort by date. (If QueryGetHistorical params)
		if jsonData, err := invokeQueryForSliceResult(str); err != nil {
		return jsonData, err
	} else {
		return jsonData, nil
	}

	} else {
		jsonDataSlice := make([]JsonRatesData, 0)
		if jsonData, err := invokeQueryForSingleResult(str); err != nil {
			return jsonDataSlice, err
		} else {

			jsonDataSlice = append(jsonDataSlice, jsonData)
			return jsonDataSlice, nil
		}
	}


	/*
	SELECT e.name, b.rate, b.created_at
	FROM baseCurrencyRates b
	JOIN elements e
	ON b.element_id = e.id
	WHERE b.created_at = ANY(
		SELECT MAX(created_at)
		FROM baseCurrencyRates
		GROUP BY element_id
	)
	 */

}

func invokeQueryForSliceResult(sqlStr string) ([]JsonRatesData, error) {
	stmtOut, err := database.Instance.Query(sqlStr)

	fmt.Println(sqlStr)

	jsonDataSlice := make([]JsonRatesData, 0)
	var jsonData JsonRatesData

	if err != nil {
		return jsonDataSlice, err
	}
	defer stmtOut.Close()

	prevDate := ""
	rates := make([]JsonRate, 0)

	for stmtOut.Next() {
		var rate JsonRate

		err := stmtOut.Scan(&rate.Name, &rate.Rate, &rate.Date)
		if err != nil {
			return jsonDataSlice, err
		}

		if rate.Date != prevDate && prevDate != "" {
			// Make jsonObject of rates for particular date
			jsonData.Rates = rates
			jsonData.Date = prevDate
			jsonDataSlice = append(jsonDataSlice, jsonData)
			// Clear slice
			rates = nil
		}
		rates = append(rates, rate)

		prevDate = rate.Date
	}

	jsonData.Rates = rates
	jsonData.Date = prevDate
	jsonDataSlice = append(jsonDataSlice, jsonData)

	return jsonDataSlice, nil

}


func invokeQueryForSingleResult(sqlStr string) (JsonRatesData, error) {
	stmtOut, err := database.Instance.Query(sqlStr)

	var jsonData JsonRatesData
	if err != nil {
		return jsonData, err
	}
	defer stmtOut.Close()

	rates := make([]JsonRate, 0)

	var rate JsonRate
	for stmtOut.Next() {

		err := stmtOut.Scan(&rate.Name, &rate.Rate, &rate.Date)
		if err != nil {
			return jsonData, err
		}

		rates = append(rates, rate)
	}

	jsonData.Rates = rates
	jsonData.Date = rate.Date

	return jsonData, nil
}

func GetAllCurrencySymbols() ([]string, error) {
	symbols := make([]string, 0)

	stmtOut, err := database.Instance.Query(`
		SELECT name FROM api.elements
		WHERE type_id = 3
	`)

	if err != nil {
		return nil, errors.New("Cannot execute query: 'Get all symbols'")
	}
	defer stmtOut.Close()

	for stmtOut.Next() {

		var symbol string

		err := stmtOut.Scan(&symbol)
		if err != nil {
			return nil, errors.New("Error in query 'Get all symbols list': Cannot get next row")
		}

		symbols = append(symbols, symbol)
	}

	if err = stmtOut.Err(); err != nil {
		return nil, err
	}

	return symbols, nil;
}

func GetLocalCurrenciesJson(rateType, dateFrom, dateTo string, historical bool) ([]JsonRatesDataLocal, error){

	jsonRatesSlice := make([]JsonRatesDataLocal, 0)
	var jsonData JsonRatesDataLocal

	sqlStr := `
			SELECT rate_buy, rate_sell, name, comission, created_at
			FROM localCurrencyRates
			INNER JOIN banks ON localCurrencyRates.bank_id = banks.id
		`
	if historical {
		sqlStr += "WHERE created_at BETWEEN '" + dateFrom + "' AND '"+ dateTo + " 23:59:59'"
	} else {
		sqlStr += `
			ORDER BY created_at DESC
			LIMIT 6
		`
	}

	stmtOut, err := database.Instance.Query(sqlStr)

	if err != nil {
		return jsonRatesSlice, err
	}
	defer stmtOut.Close()

	prevDate := ""
	rates := make([]JsonRateLocal, 0)

	for stmtOut.Next() {
		var rate JsonRateLocal

		err := stmtOut.Scan(&rate.RateBuy, &rate.RateSell, &rate.Name, &rate.Commission, &rate.Date)
		if err != nil {
			return jsonRatesSlice, err
		}

		if rate.Date != prevDate && prevDate != "" {
			// Make jsonObject of rates for particular date
			jsonData.Rates = rates
			jsonData.Date = prevDate
			jsonRatesSlice = append(jsonRatesSlice, jsonData)
			// Clear slice
			rates = nil
		}
		rates = append(rates, rate)

		prevDate = rate.Date
	}

	jsonData.Rates = rates
	jsonData.Date = prevDate
	jsonRatesSlice = append(jsonRatesSlice, jsonData)

	return jsonRatesSlice, nil;
}

func GetMetalsJson(symbols, dateFrom, dateTo string, historical bool) ([]JsonRatesData, error){

	jsonRatesSlice := make([]JsonRatesData, 0)
	var jsonData JsonRatesData

	//SELECT elements.name, metalRates.rate, metalRates.created_at FROM metalRates
	//INNER JOIN elements ON metalRates.element_id = elements.id
	//WHERE metalRates.created_at IN (
	//	SELECT DISTINCT max(metalRates.created_at) FROM metalRates
	//) AND elements.name IN ( 'gold', 'silver')
	//AND metalRates.created_at BETWEEN '2016-03-24' AND '2016-03-24 23:59:59';

	sqlStr := `
		SELECT elements.name, metalRates.rate, metalRates.created_at FROM metalRates
		INNER JOIN elements ON metalRates.element_id = elements.id
		WHERE metalRates.created_at IN (
			SELECT DISTINCT max(metalRates.created_at) FROM metalRates
		)
	`
	if historical {
		sqlStr += "WHERE created_at BETWEEN '" + dateFrom + "' AND '"+ dateTo + " 23:59:59'"
	} else {
		sqlStr += `
			ORDER BY created_at DESC
			LIMIT 6
		`
	}

	stmtOut, err := database.Instance.Query(sqlStr)

	if err != nil {
		return jsonRatesSlice, err
	}
	defer stmtOut.Close()

	prevDate := ""
	rates := make([]JsonRate, 0)

	for stmtOut.Next() {
		var rate JsonRate

		err := stmtOut.Scan(&rate.Name, &rate.Rate, &rate.Date)
		if err != nil {
			return jsonRatesSlice, err
		}

		if rate.Date != prevDate && prevDate != "" {
			// Make jsonObject of rates for particular date
			jsonData.Rates = rates
			jsonData.Date = prevDate
			jsonRatesSlice = append(jsonRatesSlice, jsonData)
			// Clear slice
			rates = nil
		}
		rates = append(rates, rate)

		prevDate = rate.Date
	}

	jsonData.Rates = rates
	jsonData.Date = prevDate
	jsonRatesSlice = append(jsonRatesSlice, jsonData)

	return jsonRatesSlice, nil;
}






//SELECT * FROM localCurrencyRates
//INNER JOIN banks ON localCurrencyRates.bank_id = banks.id
//ORDER BY created_at DESC
//LIMIT 6

//ALL
//SELECT name, rate_buy, rate_sell, comission  FROM localCurrencyRates
//INNER JOIN banks ON localCurrencyRates.bank_id = banks.id
//WHERE created_at IN (
//SELECT DISTINCT max(created_at) FROM localCurrencyRates
//)


//SELECT * FROM localCurrencyRates
//INNER JOIN banks ON localCurrencyRates.bank_id = banks.id
//WHERE created_at BETWEEN '2016-03-15' AND '2016-03-24 23:59:59'


