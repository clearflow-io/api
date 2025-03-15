package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// ParseJSON parses the JSON body into the provided struct and handles type mismatch errors.
// If strict is true, it disallows unknown fields in the JSON.
func ParseJSON(r *http.Request, dest any, strict bool) ([]string, error) {
	decoder := json.NewDecoder(r.Body)
	if strict {
		decoder.DisallowUnknownFields() // Enforce strict decoding
	}

	if err := decoder.Decode(dest); err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError
		if errors.As(err, &unmarshalTypeError) {
			return []string{
				fmt.Sprintf("Invalid type for field '%s': expected '%s', got '%s'", unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Value),
			}, err
		}

		// Handle other JSON decoding errors
		return []string{err.Error()}, err
	}

	return nil, nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteJSONError(w http.ResponseWriter, status int, errors []string) error {
	return WriteJSON(w, status, ErrorResponse{Errors: errors})
}

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
