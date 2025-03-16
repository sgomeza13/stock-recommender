package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// CleanDecimal converts a string representing a decimal number to a float64,
// handling various formats including thousand separators and currency symbols.
func CleanDecimal(s string) (float64, error) {
	// Create a regular expression to match common currency symbols
	currencySymbols := regexp.MustCompile(`[$€£¥₹₽₩]`)

	// Remove currency symbols
	s = currencySymbols.ReplaceAllString(s, "")

	// Remove all commas to handle thousand separators
	s = strings.ReplaceAll(s, ",", "")

	// Some regions use spaces as thousand separators
	s = strings.ReplaceAll(s, " ", "")

	// Handle European format where comma is used as decimal point
	// If there's a comma and no period, replace it with a period
	if strings.Contains(s, ",") && !strings.Contains(s, ".") {
		s = strings.Replace(s, ",", ".", 1)
	}

	// Trim any remaining whitespace
	s = strings.TrimSpace(s)

	// If the string is empty after all the processing, return an error
	if s == "" {
		return 0, fmt.Errorf("empty value after cleaning")
	}

	// Try to parse the cleaned string as a float
	return strconv.ParseFloat(s, 64)
}
