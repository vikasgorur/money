package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func parseNumber(s string) float64 {
	var value float64
	fmt.Sscanf(s, "%f", &value)
	return value
}

const (
	Lakh  = 100000.0
	Crore = 10000000.0
)

type InrAmount struct {
	Value float64
}

const (
	Million  = 1000000.0
	Billion  = 1000000000.0
	Trillion = 1000000000000.0
)

type UsdAmount struct {
	Value float64
}

type Amount interface {
	Convert(usdToInr float64) Amount
	FormatValue() string
}

func NewInrAmount(s string) *InrAmount {
	units := regexp.MustCompile(`lakh|crore`)
	unit := units.FindString(s)
	number := parseNumber(s)

	switch unit {
	case "lakh":
		return &InrAmount{Value: number * Lakh}
	case "crore":
		return &InrAmount{Value: number * Crore}
	default:
		return &InrAmount{Value: number}
	}
}

func (amount *InrAmount) Convert(usdToInr float64) Amount {
	return &UsdAmount{Value: amount.Value / usdToInr}
}

func (amount *InrAmount) FormatValue() string {
	if v := amount.Value / Crore; v >= 1.0 {
		return fmt.Sprintf("₹ %.1f crore", v)
	} else if v := amount.Value / Lakh; v >= 1.0 {
		return fmt.Sprintf("₹ %.1f lakh", v)
	} else {
		return fmt.Sprintf("₹ %.1f", amount.Value)
	}
}

func NewUsdAmount(s string) *UsdAmount {
	units := regexp.MustCompile("million|billion|trillion")
	unit := units.FindString(s)
	number := parseNumber(s)

	switch unit {
	case "million":
		return &UsdAmount{Value: number * Million}
	case "billion":
		return &UsdAmount{Value: number * Billion}
	case "trillion":
		return &UsdAmount{Value: number * Trillion}
	default:
		return &UsdAmount{Value: number}
	}
}

func (amount *UsdAmount) Convert(usdToInr float64) Amount {
	return &InrAmount{Value: amount.Value * usdToInr}
}

func (amount *UsdAmount) FormatValue() string {
	if v := amount.Value / Trillion; v >= 1.0 {
		return fmt.Sprintf("$ %.1f trillion", v)
	} else if v := amount.Value / Billion; v >= 1.0 {
		return fmt.Sprintf("$ %.1f billion", v)
	} else if v := amount.Value / Million; v >= 1.0 {
		return fmt.Sprintf("$ %.1f million", v)
	} else {
		return fmt.Sprintf("$ %.1f", amount.Value)
	}
}

var inrSignifiers = regexp.MustCompile(`lakh|crore|arab|rs|inr|₹|rupee`)
var usdSignifiers = regexp.MustCompile(`million|billion|trillion|\$|usd|dollar`)

func parseAmount(s string) (Amount, error) {
	if inrSignifiers.MatchString(s) {
		return NewInrAmount(s), nil
	} else if usdSignifiers.MatchString(s) {
		return NewUsdAmount(s), nil
	} else {
		return nil, errors.New("no currency recognized")
	}
}

type FixerResponse struct {
	Rates map[string]float64 `json:"rates"`
}

func getUsdToInr() float64 {
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

func main() {
	input := strings.Join(os.Args[1:], " ")
	amount, err := parseAmount(input)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(amount.Convert(getUsdToInr()).FormatValue())
}
