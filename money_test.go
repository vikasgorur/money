package money

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsdMatch(t *testing.T) {
	a := assert.New(t)
	var p UsdParser

	a.True(p.Match("$100"))

	// Test that match is case-insensitive
	a.True(p.Match("5 Dollars"))
	a.True(p.Match("5 dollar"))

	a.True(p.Match("9 million"))
	a.True(p.Match("10 Billion"))
	a.True(p.Match("10 trillion"))

	a.False(p.Match("₹100"))
	a.False(p.Match("78"))
	a.False(p.Match("7 brazilian"))
}

func TestInrMatch(t *testing.T) {
	a := assert.New(t)
	var p InrParser

	a.True(p.Match("2 rupees"))

	// Test that match is case-insensitive
	a.True(p.Match("Rs..67"))
	a.True(p.Match("rs..67"))

	a.True(p.Match("₹100"))
	a.True(p.Match("9 LAKH"))
	a.True(p.Match("50 crore"))

	a.False(p.Match("$100"))
}
