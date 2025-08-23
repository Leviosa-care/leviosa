package httpx

import (
	"fmt"
	"strconv"
	"strings"
)

// Helper functions (parseIntFormValue, respondWithError, respondWithJSON) remain the same.
// Helper function to parse int form values with basic error messages
func ParseIntFormValue(s, fieldName string) (int, error) {
	fmt.Printf("parsing '%s' that has the value %q\n", fieldName, s)
	s = strings.TrimSpace(s)
	fmt.Printf("parsing '%s' that has the value %q\n", fieldName, s)
	if s == "" { // For required numeric fields
		return 0, fmt.Errorf("'%s' is required", fieldName)
	}
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid %s format. Must be a number", fieldName)
	}
	return int(val), nil
}
