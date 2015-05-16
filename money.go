package main

import (
	"fmt"
	"os"
	"strings"
)

type Amount struct {
	Currency string
	Value    float64
}

var indianUnits = map[string]uint64{"arab": 1000000000, "crore": 10000000, "lakh": 100000}
var usUnits = map[string]uint64{"trillion": 1000000000000, "billion": 1000000000, "million": 1000000}

var multipliersFor = map[string]map[string]uint64{"inr": indianUnits, "usd": usUnits}

// return all keys of the map
func keys(m map[string]uint64) []string {
	ks := make([]string, 0)
	for k, _ := range m {
		ks = append(ks, k)
	}
	return ks
}

// true if s contains val, returns val
func contains(s []string, val string) (bool, string) {
	for _, u := range s {
		if strings.Contains(strings.ToLower(val), u) {
			return true, u
		}
	}
	return false, ""
}

var inrSignifiers = append(keys(indianUnits), []string{"rs", "inr", "₹", "rupee"}...)
var usdSignifiers = append(keys(usUnits), []string{"$", "usd", "dollar"}...)

func parseCurrency(s string) string {
	if c, _ := contains(inrSignifiers, s); c {
		return "inr"
	} else if c, _ := contains(usdSignifiers, s); c {
		return "usd"
	} else {
		return ""
	}
}

func parseMultiplier(s string) uint64 {
	if c, v := contains(keys(indianUnits), s); c {
		return indianUnits[v]
	}

	if c, v := contains(keys(usUnits), s); c {
		return usUnits[v]
	}

	return 1
}

func parseNumber(s string) float64 {
	var value float64
	fmt.Sscanf(s, "%f", &value)
	return value
}

func parse(s string) Amount {
	currency := parseCurrency(s)
	multiplier := parseMultiplier(s)
	number := parseNumber(s)

	return Amount{Currency: currency, Value: float64(multiplier) * number}
}

// convert from one currency to the other
func convert(amount Amount) Amount {
	switch amount.Currency {
	case "inr":
		return Amount{Currency: "usd", Value: amount.Value / 62.0}
	case "usd":
		return Amount{Currency: "inr", Value: amount.Value * 62.0}
	default:
		return Amount{}
	}
}

func symbolFor(amount Amount) string {
	switch amount.Currency {
	case "inr":
		return "₹"
	case "usd":
		return "$"
	default:
		return ""
	}
}

func otherCurrency(currency string) string {
	switch currency {
	case "inr":
		return "usd"
	case "usd":
		return "inr"
	default:
		return ""
	}
}

func humanDivisorFor(amount Amount) (string, uint64) {
	multipliers := multipliersFor[amount.Currency]
	for k, v := range multipliers {
		if amount.Value/float64(v) > 1 {
			return k, v
		}
	}

	return "", 1
}

func humanize(amount Amount) string {
	symbol := symbolFor(amount)
	unit, divisor := humanDivisorFor(amount)
	return fmt.Sprintf("%s%.1f %s", symbol, amount.Value/float64(divisor), unit)
}

func main() {
	amount := convert(parse(strings.Join(os.Args[1:], " ")))
	fmt.Println(humanize(amount))
}
