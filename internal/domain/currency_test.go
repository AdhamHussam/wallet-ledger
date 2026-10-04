package domain

import "testing"

func TestIsSupportedCurrency(t *testing.T) {
	testCases := []struct {
		currency string
		expected bool
	}{
		{"USD", true},
		{"usd", true},
		{"EUR", true},
		{"eur", true},
		{"GBP", true},
		{"CAD", true},
		{"JPY", true},
		{"AUD", true},
		{"CHF", true},
		{"XYZ", false},
		{"", false},
		{"USDT", false},
	}

	for _, tc := range testCases {
		got := IsSupportedCurrency(tc.currency)
		if got != tc.expected {
			t.Errorf("IsSupportedCurrency(%q) = %v; want %v", tc.currency, got, tc.expected)
		}
	}
}
