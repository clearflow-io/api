package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseJSON(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	tests := []struct {
		name    string
		body    string
		strict  bool
		wantErr bool
	}{
		{
			name:    "valid json",
			body:    `{"name": "John", "age": 30}`,
			strict:  false,
			wantErr: false,
		},
		{
			name:    "invalid json",
			body:    `{"name": "John", "age": "thirty"}`,
			strict:  false,
			wantErr: true,
		},
		{
			name:    "unknown field strict",
			body:    `{"name": "John", "age": 30, "unknown": true}`,
			strict:  true,
			wantErr: true,
		},
		{
			name:    "unknown field not strict",
			body:    `{"name": "John", "age": 30, "unknown": true}`,
			strict:  false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.body))
			var dest TestStruct
			err := ParseJSON(r, &dest, tt.strict)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		v          any
		expectBody string
	}{
		{
			name:       "write object",
			status:     http.StatusOK,
			v:          map[string]string{"message": "success"},
			expectBody: `{"message":"success"}`,
		},
		{
			name:       "write nil",
			status:     http.StatusOK,
			v:          nil,
			expectBody: `{}`,
		},
		{
			name:       "write nil slice",
			status:     http.StatusOK,
			v:          []string(nil),
			expectBody: `[]`,
		},
		{
			name:       "write nil map",
			status:     http.StatusOK,
			v:          map[string]string(nil),
			expectBody: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := WriteJSON(w, tt.status, tt.v)
			assert.NoError(t, err)
			assert.Equal(t, tt.status, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

			var expected, actual any
			json.Unmarshal([]byte(tt.expectBody), &expected)
			json.Unmarshal(w.Body.Bytes(), &actual)
			assert.Equal(t, expected, actual)
		})
	}
}
