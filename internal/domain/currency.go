package domain

import "strings"

const (
	USD = "USD"
	EUR = "EUR"
	GBP = "GBP"
	CAD = "CAD"
	JPY = "JPY"
	AUD = "AUD"
	CHF = "CHF"
)

var supportedCurrencies = map[string]struct{}{
	USD: {},
	EUR: {},
	GBP: {},
	CAD: {},
	JPY: {},
	AUD: {},
	CHF: {},
}

func IsSupportedCurrency(currency string) bool {
	_, ok := supportedCurrencies[strings.ToUpper(strings.TrimSpace(currency))]
	return ok
}
