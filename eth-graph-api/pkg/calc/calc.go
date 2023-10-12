package calc

import (
	"fmt"
	"math/big"
)

// SumNumbers takes a slice of strings, each representing a number,
// and returns the sum of all these numbers as a string.
// It utilizes big.Float to handle the arithmetic of potentially very large numbers
// and ensures accurate sum calculation without floating-point precision issues.
//
// Parameters:
//   - `numbers`: A slice of strings, where each string represents a number that
//     needs to be included in the sum.
//
// Returns:
//   - A string representing the sum of all the numbers in the input slice.
//   - An error, which will be non-nil if any of the input strings cannot be
//     parsed into a number.
//
// Example usage:
//
//	sum, err := SumNumbers([]string{"10", "20.5", "30.75"})
//
// In the example above, if err is nil, sum will be a string representing the
// sum of the numbers ["10", "20.5", "30.75"].
func SumNumbers(numbers []string) (string, error) {
	sum := new(big.Float)
	for _, numStr := range numbers {
		num, _, err := big.ParseFloat(numStr, 10, 0, big.ToNearestEven)
		if err != nil {
			return "", fmt.Errorf("invalid number: %s, error: %v", numStr, err)
		}
		sum.Add(sum, num)
	}

	return sum.Text('f', 10), nil
}
