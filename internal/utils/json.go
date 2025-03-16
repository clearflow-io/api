package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/sirupsen/logrus"
)

// ParseJSON parses the JSON body into the provided struct and handles type mismatch errors.
// If strict is true, it disallows unknown fields in the JSON.
func ParseJSON(r *http.Request, dest any, strict bool) error {
	decoder := json.NewDecoder(r.Body)
	if strict {
		decoder.DisallowUnknownFields() // Enforce strict decoding
	}

	if err := decoder.Decode(dest); err != nil {
		var unmarshalTypeError *json.UnmarshalTypeError
		if errors.As(err, &unmarshalTypeError) {
			return fmt.Errorf("invalid type for field '%s': expected '%s', got '%s'",
				unmarshalTypeError.Field, unmarshalTypeError.Type, unmarshalTypeError.Value)
		}

		// Handle other JSON decoding errors
		return err
	}

	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Handle nil slices and maps
	if v == nil {
		v = struct{}{} // Encode nil as an empty JSON object
	} else {
		// Use reflection to check for nil slices or maps
		val := reflect.ValueOf(v)
		if val.Kind() == reflect.Slice && val.IsNil() {
			v = []any{} // Encode nil slice as an empty JSON array
		} else if val.Kind() == reflect.Map && val.IsNil() {
			v = map[string]any{} // Encode nil map as an empty JSON object
		}
	}

	return json.NewEncoder(w).Encode(v)
}

func WriteJSONError(w http.ResponseWriter, status int, err error) error {
	logrus.Warnf("Writing JSON error: %s", err)
	return WriteJSON(w, status, ErrorResponse{Error: err.Error()})
}
