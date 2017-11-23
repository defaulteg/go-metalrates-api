package currency

import (
	"github.com/defaulteg/api/modules/core"
)

//TODO: add comission table with ranges

func FetchLocal() error {
	// Get sources from where will be local currency data fetched
	if sources, err := core.GetSources(core.LocalCurrencyCategory); err != nil {
		return err
	} else {
		var rates = make([]core.LocalRate, 0)
		for _, source := range sources {
			// Get buy and sell rates
			if err := source.GetRatesFromStaticSource(); err != nil {
				return err
			}
			// Get bank matcher for this rate (e.g. Map[SELECTOR_NAME] = bank_id
			selectorBankMatcherMap, err := core.GetElementIds(core.CryptocurrencyTypeId, core.QueryMatcher, "")
			if err != nil {
				return err
			}
			// Count of all fetched selectors with rates divided by two
			selectorHalfCount := len(source.Selectors) / 2;

			// Go through all selectors / 2, because for each rate object there are two types of rate: buy, sell.
			for i := 0; i < selectorHalfCount; i++ {
				var rate core.LocalRate

				rate.ElementId = source.Selectors[i].ElementId						// Id of element; e.g. 17 = USD, 18 = RUB...
				rate.BankId = selectorBankMatcherMap[source.Selectors[i].Name]		// Bank id from matcher map
				rate.RateBuy = source.Selectors[i].Rate								// Buy rate for currency on base = EUR. if element is USD, it means: EUR costs x.xxxx USD
				rate.RateSell = source.Selectors[i + selectorHalfCount].Rate		// Sell rate for currency on base = EUR

				rates = append(rates, rate)											// Append to rate list for further push to database
			}
		}
		if err := core.PushToDatabaseUnique(rates, core.LocalCurrencyTable); err != nil {
			return err
		}
	}

	return nil
}