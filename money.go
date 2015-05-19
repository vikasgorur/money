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

func NewInrAmount(s string) (*InrAmount, error) {
	units := regexp.MustCompile(`lakh|crore`)
	unit := units.FindString(s)
	number, err := parseNumber(s)

	if err != nil {
		return nil, err
	}

	switch unit {
	case "lakh":
		return &InrAmount{Value: number * Lakh}, nil
	case "crore":
		return &InrAmount{Value: number * Crore}, nil
	default:
		return &InrAmount{Value: number}, nil
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

func NewUsdAmount(s string) (*UsdAmount, error) {
	units := regexp.MustCompile("million|billion|trillion")
	unit := units.FindString(s)
	number, err := parseNumber(s)

	if err != nil {
		return nil, err
	}

	switch unit {
	case "million":
		return &UsdAmount{Value: number * Million}, nil
	case "billion":
		return &UsdAmount{Value: number * Billion}, nil
	case "trillion":
		return &UsdAmount{Value: number * Trillion}, nil
	default:
		return &UsdAmount{Value: number}, nil
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

var inrSignifiers = regexp.MustCompile(`lakh|crore|rs|inr|₹|rupee`)
var usdSignifiers = regexp.MustCompile(`million|billion|trillion|\$|usd|dollar`)

func parseAmount(s string) (Amount, error) {
	if inrSignifiers.MatchString(s) {
		return NewInrAmount(s)
	} else if usdSignifiers.MatchString(s) {
		return NewUsdAmount(s)
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
	if len(os.Args) < 2 {
		fmt.Println("no input given")
		os.Exit(1)
	}

	input := strings.Join(os.Args[1:], " ")
	amount, err := parseAmount(input)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(amount.Convert(getUsdToInr()).FormatValue())
}
