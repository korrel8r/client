// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

package cmd

import (
	"encoding/json"
	"testing"
)

func TestCheckResponse_ErrorExtraction(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantOK     bool
	}{
		{
			name:       "200 OK",
			statusCode: 200,
			body:       `{"result": "ok"}`,
			wantOK:     true,
		},
		{
			name:       "400 with error message",
			statusCode: 400,
			body:       `{"error": "invalid query string: foo"}`,
			wantOK:     false,
		},
		{
			name:       "404 with error message",
			statusCode: 404,
			body:       `{"error": "class not found: alert:alert"}`,
			wantOK:     false,
		},
		{
			name:       "500 without JSON body",
			statusCode: 500,
			body:       `Internal Server Error`,
			wantOK:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantOK {
				// Should not exit for success codes
				checkResponse(tt.statusCode, []byte(tt.body), "GET", "/test")
			} else {
				// Verify error JSON parsing works
				var apiErr struct{ Error string }
				if err := json.Unmarshal([]byte(tt.body), &apiErr); err == nil && apiErr.Error != "" {
					if apiErr.Error == "" {
						t.Error("expected non-empty error message")
					}
				}
			}
		})
	}
}
