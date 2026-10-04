package domain

import "strings"

// Supported currency constants
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

// IsSupportedCurrency returns true if the uppercase currency string is supported.
func IsSupportedCurrency(currency string) bool {
	_, ok := supportedCurrencies[strings.ToUpper(strings.TrimSpace(currency))]
	return ok
}
