package money

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

// parseNumber expects the string to contain a single
// floating point number and returns the same. If the
// string doesn't contain exactly one number, it returns an error.
func parseNumber(s string) (float64, error) {
	var value float64
	n, err := fmt.Sscanf(s, "%f", &value)

	if n != 1 || err != nil {
		return 0, errors.New("invalid number")
	} else {
		return value, nil
	}
}

// Converters have a method that given the exchange rate from
// USD to INR return a printable result in the other currency.
type Converter interface {
	Convert(usdToInr float64) fmt.Stringer
}

// Inr is an amount denominated in Rupees.
type Inr float64

const (
	lakh  = 100000.0
	crore = 10000000.0
)

// Convert returns the equivalent Usd for an Inr.
func (amount Inr) Convert(usdToInr float64) fmt.Stringer {
	return Usd(float64(amount) / usdToInr)
}

func (amount Inr) String() string {
	if v := amount / crore; v >= 1.0 {
		return fmt.Sprintf("₹ %.1f crore", v)
	} else if v := amount / lakh; v >= 1.0 {
		return fmt.Sprintf("₹ %.1f lakh", v)
	} else {
		return fmt.Sprintf("₹ %.1f", amount)
	}
}

// Usd is an amount denominated in dollars.
type Usd float64

const (
	million  = 1000000.0
	billion  = 1000000000.0
	trillion = 1000000000000.0
)

func (amount Usd) Convert(usdToInr float64) fmt.Stringer {
	return Inr(float64(amount) * usdToInr)
}

func (amount Usd) String() string {
	if v := amount / trillion; v >= 1.0 {
		return fmt.Sprintf("$ %.1f trillion", v)
	} else if v := amount / billion; v >= 1.0 {
		return fmt.Sprintf("$ %.1f billion", v)
	} else if v := amount / million; v >= 1.0 {
		return fmt.Sprintf("$ %.1f million", v)
	} else {
		return fmt.Sprintf("$ %.1f", amount)
	}
}

// Parser represents a value that knows how to parse strings that denote
// an amount in a particular currency. A Parser also provides a method
// that identifies if a string is an amount in its currency.
type Parser interface {
	Match(s string) bool
	Parse(s string) (Converter, error)
}

// InrParser parses amounts in INR.
type InrParser struct{}

func (p InrParser) Match(s string) bool {
	return regexp.MustCompile(`lakh|crore|rs|inr|₹|rupee`).MatchString(s)
}

func (p InrParser) Parse(s string) (Converter, error) {
	units := regexp.MustCompile(`lakh|crore`)
	unit := units.FindString(s)
	number, err := parseNumber(s)

	if err != nil {
		return Usd(0), err
	}

	switch unit {
	case "lakh":
		return Inr(number * lakh), nil
	case "crore":
		return Inr(number * crore), nil
	default:
		return Inr(number), nil
	}
}

// UsdParser parses amounts in USD.
type UsdParser struct{}

func (p UsdParser) Match(s string) bool {
	return regexp.MustCompile(`million|billion|trillion|\$|usd|dollar`).MatchString(s)
}

func (p UsdParser) Parse(s string) (Converter, error) {
	units := regexp.MustCompile("million|billion|trillion")
	unit := units.FindString(s)
	number, err := parseNumber(s)

	if err != nil {
		return Inr(0), err
	}

	switch unit {
	case "million":
		return Usd(number * million), nil
	case "billion":
		return Usd(number * billion), nil
	case "trillion":
		return Usd(number * trillion), nil
	default:
		return Usd(number), nil
	}
}

// ErrorParser matches any string, and its Parse() always returns an error.
type ErrorParser struct{}

func (p ErrorParser) Match(s string) bool {
	return true
}

func (p ErrorParser) Parse(s string) (Converter, error) {
	return nil, errors.New("could not parse: unknown currency")
}

// Parse tries to match the string against all available parsers
// and returns the parsed valued from the first one that matches.
// The search will always succeed because ErrorParser will always
// match.
func Parse(s string) (Converter, error) {
	parsers := []Parser{InrParser{}, UsdParser{}, ErrorParser{}}
	for _, p := range parsers {
		if p.Match(s) {
			return p.Parse(s)
		}
	}

	panic("none of the parsers matched!")
}

// FixerResponse is the response from api.fixer.io, which provides
// us the exchange rate.
type FixerResponse struct {
	Rates map[string]float64 `json:"rates"`
}

// GetUsdToInr fetches the exchange rate from fixer.io. To simplify
// error handling, it returns a default value if there was an error
// accessing the API.
func GetUsdToInr() float64 {
	defaultRate := 62.0

	r, err := http.Get("http://api.fixer.io/latest?base=USD")
	if err != nil {
		return defaultRate
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return defaultRate
	}

	var response FixerResponse
	if err = json.Unmarshal(body, &response); err != nil {
		return defaultRate
	}

	return response.Rates["INR"]
}
