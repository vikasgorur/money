package money

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

func parseNumber(s string) (float64, error) {
	var value float64
	n, err := fmt.Sscanf(s, "%f", &value)

	if n != 1 || err != nil {
		return 0, errors.New("invalid number")
	} else {
		return value, nil
	}
}

const (
	lakh  = 100000.0
	crore = 10000000.0
)

type inrAmount struct {
	Value float64
}

const (
	million  = 1000000.0
	billion  = 1000000000.0
	trillion = 1000000000000.0
)

type usdAmount struct {
	Value float64
}

type Amount interface {
	Convert(usdToInr float64) Amount
	FormatValue() string
}

func newInrAmount(s string) (*inrAmount, error) {
	units := regexp.MustCompile(`lakh|crore`)
	unit := units.FindString(s)
	number, err := parseNumber(s)

	if err != nil {
		return nil, err
	}

	switch unit {
	case "lakh":
		return &inrAmount{Value: number * lakh}, nil
	case "crore":
		return &inrAmount{Value: number * crore}, nil
	default:
		return &inrAmount{Value: number}, nil
	}
}

func (amount *inrAmount) Convert(usdToInr float64) Amount {
	return &usdAmount{Value: amount.Value / usdToInr}
}

func (amount *inrAmount) FormatValue() string {
	if v := amount.Value / crore; v >= 1.0 {
		return fmt.Sprintf("₹ %.1f crore", v)
	} else if v := amount.Value / lakh; v >= 1.0 {
		return fmt.Sprintf("₹ %.1f lakh", v)
	} else {
		return fmt.Sprintf("₹ %.1f", amount.Value)
	}
}

func newUsdAmount(s string) (*usdAmount, error) {
	units := regexp.MustCompile("million|billion|trillion")
	unit := units.FindString(s)
	number, err := parseNumber(s)

	if err != nil {
		return nil, err
	}

	switch unit {
	case "million":
		return &usdAmount{Value: number * million}, nil
	case "billion":
		return &usdAmount{Value: number * billion}, nil
	case "trillion":
		return &usdAmount{Value: number * trillion}, nil
	default:
		return &usdAmount{Value: number}, nil
	}
}

func (amount *usdAmount) Convert(usdToInr float64) Amount {
	return &inrAmount{Value: amount.Value * usdToInr}
}

func (amount *usdAmount) FormatValue() string {
	if v := amount.Value / trillion; v >= 1.0 {
		return fmt.Sprintf("$ %.1f trillion", v)
	} else if v := amount.Value / billion; v >= 1.0 {
		return fmt.Sprintf("$ %.1f billion", v)
	} else if v := amount.Value / million; v >= 1.0 {
		return fmt.Sprintf("$ %.1f million", v)
	} else {
		return fmt.Sprintf("$ %.1f", amount.Value)
	}
}

var inrSignifiers = regexp.MustCompile(`lakh|crore|rs|inr|₹|rupee`)
var usdSignifiers = regexp.MustCompile(`million|billion|trillion|\$|usd|dollar`)

func ParseAmount(s string) (Amount, error) {
	if inrSignifiers.MatchString(s) {
		return newInrAmount(s)
	} else if usdSignifiers.MatchString(s) {
		return newUsdAmount(s)
	} else {
		return nil, errors.New("no currency recognized")
	}
}

type FixerResponse struct {
	Rates map[string]float64 `json:"rates"`
}

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
