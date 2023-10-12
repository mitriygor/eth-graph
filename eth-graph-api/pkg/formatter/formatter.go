package formatter

import (
	"fmt"
	"strconv"
	"strings"
)

// FormatUSD takes a string representing a numerical value and formats it
// into a USD currency format (e.g., $1,234,567.89).
//
// Parameters:
//   - `numberStr`: A string representing the numerical value to be formatted.
//     This input string should be capable of being parsed into a float64.
//
// Returns:
// - A string representing the input number formatted in USD currency format.
// - An error, which will be non-nil if the input string cannot be parsed into a number.
//
// Example usage:
//
//	formatted, err := FormatUSD("1234567.891234")
//
// If err is nil, formatted will hold the string "$1,234,567.89".
func FormatUSD(numberStr string) (string, error) {
	number, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		return "", err
	}

	usd := fmt.Sprintf("$%.2f", number)

	parts := strings.Split(usd, ".")
	integerPart := parts[0]

	for i := len(integerPart) - 3; i > 1; i -= 3 {
		integerPart = integerPart[:i] + "," + integerPart[i:]
	}

	return integerPart + "." + parts[1], nil
}
