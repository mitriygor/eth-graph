package validator

import (
	"regexp"
	"strconv"
	"time"
)

// IsValidToken checks the validity of a token ID. A valid token ID should satisfy the
// following conditions:
// - Have a length between minIdLength and maxIdLength.
// - Consist only of alphanumeric characters.
//
// Parameters:
// - `id`: a string representing the token ID to be checked.
//
// Returns:
// - A boolean value indicating whether the token ID is valid.
func IsValidToken(id string) bool {
	const (
		minIdLength = 40
		maxIdLength = 60
	)

	if len(id) < minIdLength || len(id) > maxIdLength {
		return false
	}

	isAlphaNumeric := regexp.MustCompile(`^[a-zA-Z0-9]*$`).MatchString

	if !isAlphaNumeric(id) {
		return false
	}

	return true
}

// IsValidBlock validates the provided block string by ensuring it is a positive integer.
//
// Parameters:
// - `block`: a string representing a block number to be checked.
//
// Returns:
// - A boolean value indicating whether the block number is valid.
func IsValidBlock(block string) bool {
	num, err := strconv.Atoi(block)

	if err != nil {
		return false
	}

	if num <= 0 {
		return false
	}

	return true
}

// IsValidFirst checks the validity of a query limit value. A valid limit should:
// - Be between minFirst and maxFirst inclusive.
//
// Parameters:
// - `first`: an integer representing the query limit to be checked.
//
// Returns:
// - A boolean value indicating whether the query limit is valid.
func IsValidFirst(first int) bool {
	const (
		minFirst = 1
		maxFirst = 1000
	)

	if first < minFirst || first > maxFirst {
		return false
	}

	return true
}

// IsValidRange checks the validity of a time range. The following conditions should be met:
// - Both "from" and "to" can be successfully parsed into 64-bit integers.
// - "from" is less than "to".
// - "to" is less than or equal to the current Unix timestamp.
//
// Parameters:
// - `fromStr`: a string representing the start of the time range.
// - `toStr`: a string representing the end of the time range.
//
// Returns:
// - A boolean value indicating whether the time range is valid.
func IsValidRange(fromStr, toStr string) bool {
	from, err := strconv.ParseInt(fromStr, 10, 64)
	if err != nil {
		return false
	}

	to, err := strconv.ParseInt(toStr, 10, 64)
	if err != nil {
		return false
	}

	now := time.Now().Unix()

	return from < to && to <= now
}
