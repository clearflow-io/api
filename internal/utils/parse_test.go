package utils

import (
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParseQueryParamInt(t *testing.T) {
	tests := []struct {
		name      string
		paramName string
		query     string
		required  bool
		expected  int
		wantErr   bool
	}{
		{
			name:      "valid integer",
			paramName: "age",
			query:     "age=25",
			required:  true,
			expected:  25,
			wantErr:   false,
		},
		{
			name:      "missing optional",
			paramName: "age",
			query:     "",
			required:  false,
			expected:  0,
			wantErr:   false,
		},
		{
			name:      "missing required",
			paramName: "age",
			query:     "",
			required:  true,
			wantErr:   true,
		},
		{
			name:      "invalid integer",
			paramName: "age",
			query:     "age=abc",
			required:  true,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, _ := url.Parse("http://example.com?" + tt.query)
			r := &http.Request{URL: u}
			var dest int
			err := ParseQueryParamInt(r, &dest, tt.paramName, tt.required)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, dest)
			}
		})
	}
}

func TestParseUUID(t *testing.T) {
	validUUID := uuid.New()
	tests := []struct {
		name      string
		str       string
		paramName string
		wantErr   bool
	}{
		{
			name:      "valid uuid",
			str:       validUUID.String(),
			paramName: "id",
			wantErr:   false,
		},
		{
			name:      "empty string",
			str:       "",
			paramName: "id",
			wantErr:   true,
		},
		{
			name:      "invalid uuid",
			str:       "not-a-uuid",
			paramName: "id",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUUID(tt.str, tt.paramName)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.str, got.String())
			}
		})
	}
}

func TestParseIsoDate(t *testing.T) {
	tests := []struct {
		name    string
		dateStr string
		wantErr bool
	}{
		{
			name:    "valid date",
			dateStr: "2023-10-27",
			wantErr: false,
		},
		{
			name:    "invalid format",
			dateStr: "27-10-2023",
			wantErr: true,
		},
		{
			name:    "invalid date",
			dateStr: "2023-13-45",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dest time.Time
			err := ParseIsoDate(tt.dateStr, &dest)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.dateStr, dest.Format("2006-01-02"))
			}
		})
	}
}
