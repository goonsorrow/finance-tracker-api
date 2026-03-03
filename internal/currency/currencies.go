package currency

import "strings"

var SupportedCurrencies = []string{"USD", "EUR", "RUB", "GBP", "JPY"}

func IsSupported(code string) bool {
	for _, c := range SupportedCurrencies {
		if c == strings.ToUpper(code) {
			return true
		}
	}
	return false
}
