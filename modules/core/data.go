package core

// Source url from where data must be fetched. For each source there can be multiple selectors, where
// selector equals to "id name of html element" which value will be parsed and stored in selector object

type Source struct {
	Site      	string
	Id        	int
	Selectors 	[]Selector
}

// Consists information about individual rate element.
// For metals: id from html markup in "NAME" field. Source is site url where there id fields are placed
// For currencies: "NAME" field is currency CODE. Source have pseudo site file: "csv"

type Selector struct {
	Name    	string
	ElementId 	int
	Rate     	string
}

type JsonRate struct {
	Name    	string		`json:"name"`
	Rate     	string		`json:"rate"`
	Date     	string		`json:"-"`
}

type JsonRatesData struct {
	Base 		string		`json:"base"`
	Date 		string		`json:"date"`
	Rates 		[]JsonRate	`json:"rates"`
}

type JsonRatesDataLocal struct {
	Date 		string		`json:"date"`
	Rates 	[]JsonRateLocal	`json:"rates"`
}

type JsonRateLocal struct {
	Name    	string		`json:"name"`
	RateBuy     string		`json:"rateSell"`
	RateSell	string 		`json:"rateBuy"`
	Commission	string 		`json:"commission"`
	Date     	string		`json:"-"`
}


type ErrorJson struct {
	Type		string		`json:"type"`
	Description string		`json:"description"`
}

// Object that contains information about newly fetched cryptocurrency rate from source:
// Contains id of cryptocoin; cryptocoin name; id of market (e.g. Bittrex, Kraken, ...);
// id of traiding pair (e.g. BTC/USD, LTC/EUR/, ...); current exchange rate of this cryptocoin
//TODO: add market cap? trading volume?

type CryptoRate struct {
	Id       	int
	Name    	string
	MarketId 	int
	PairId   	int
	Rate     	string
}

type LocalRate struct {
	ElementId 	int
	BankId 		int
	RateBuy		string
	RateSell 	string
}

const (
	// Category names
	MetalCategory string = "metal"
	CryptocurrencyCategory string = "cryptocurrency"
	LocalCurrencyCategory string = "local_currency"

	// Table names
	MetalTable string = "MetalRates"
	CryptocurrencyTable string = "CryptocurrencyRates"
	BaseCurrencyTable string = "BaseCurrencyRates"
	LocalCurrencyTable string = "localCurrencyRates"

	// Paths to script
	PathToPageElementFetcher string = "./scripts/page_element_fetcher.js" 		//path to script - SEQUENTIALLY
	PathToPageFetcher string = "./scripts/page_fetcher.js" 					//path to script - DOWNLOAD PAGE
	PathToCryptocurrencyFetcher string = "./scripts/crypto_fetcher.js"

	// Paths for file manipulations (unzipping, downloading, accessing, etc.)
	CurrencyDownloadUrl string = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref.zip"
	PathToSaveCurrencyZipCsv string = "./modules/currency/data/eurofxref.zip"
	PathToUnzipCurrencyCsv string = "./modules/currency/data"
	PathToCurrencyCsv string = "./modules/currency/data/eurofxref.csv"

	// Type id's
	MetalTypeId int = 1
	CryptocurrencyTypeId int = 2
	BaseCurrencyTypeId int = 3
	LocalCurrencyTypeID int = 4

	// Query types
	QueryBase int = 1 			// get list of all table elemenets with query
	QueryWhere int = 2			// get specific list from table with WHERE clause
	QueryMatcher int = 3   		// Gets compliant map table like Map[BANK_ID] = SELECTOR.NAME

	QueryGetAll int = 4			// Get rates of currencies in any from ('RUB', 'JPY', ...)
	QueryGetParticular int = 5
	QueryGetHistorical int = 6
	QueryGetLatest int = 7

	// Types of json object
	JsonSliceObject int = 1
	JsonOneObject int = 2

	// Error messages for errorPage
	ERR_SPECIFY_DATE string = "you must specify date 'from' and date 'to'"
	ERR_NOT_SPECIFIED_DATE string = "'date from' and 'date to' are not specified"
	ERR_INVALID_DATE_FORMAT string = "invalid date format (e.g. 2016/03/13 or 2016-04-12)"
	ERR_INVALID_RATE_FORMAT string = "Invalid rate type (Valid values: sell/buy/all)"

)

