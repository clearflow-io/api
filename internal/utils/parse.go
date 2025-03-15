package utils

import (
	"fmt"
	"net/http"
	"strconv"
)

// ParseQueryParamInt parses a query parameter string into the provided destination.
// Returns user-friendly error messages if parsing fails.
func ParseQueryParamInt(r *http.Request, paramName string, dest *int) ([]string, error) {
	param := r.URL.Query().Get(paramName)
	if param == "" {
		return nil, nil
	}

	parsedValue, err := strconv.Atoi(param)
	if err != nil {
		return []string{
			fmt.Sprintf("Invalid value '%s'. Expected an integer.", param),
		}, err
	}

	*dest = parsedValue
	return nil, nil
}

func IntPointerToValue(ptr *int) int {
	if ptr == nil {
		return 0 // Default value if nil
	}
	return *ptr
}
