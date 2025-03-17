package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func ParseQueryParamInt(r *http.Request, dest *int, paramName string, required bool) error {
	param := r.URL.Query().Get(paramName)
	if param == "" {
		if required {
			return fmt.Errorf("query param %s is required", paramName)
		}
		return nil
	}

	parsedValue, err := strconv.Atoi(param)
	if err != nil {
		return fmt.Errorf("invalid value '%s': expected an integer", param)
	}

	*dest = parsedValue
	return nil
}

func ParseUUID(str, paramName string) (uuid.UUID, error) {
	if str == "" {
		return uuid.UUID{}, fmt.Errorf("path parameter %s is required", paramName)
	}

	parsedUUID, err := uuid.Parse(str)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid path parameter %s: expected a valid UUID", paramName)
	}

	return parsedUUID, nil
}

func ParseIsoDate(dateStr string, dest *time.Time) error {
	result, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return fmt.Errorf("invalid date format: expected %s, got %s", "YYYY-MM-DD", dateStr)
	}

	*dest = result
	return nil
}
