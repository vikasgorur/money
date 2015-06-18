package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/vikasgorur/money"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no input given")
		os.Exit(1)
	}

	input := strings.Join(os.Args[1:], " ")
	amount, err := money.ParseAmount(input)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(amount.Convert(money.GetUsdToInr()).FormatValue())
}
