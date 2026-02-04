// Copyright: This file is part of korrel8r, released under https://github.com/korrel8r/korrel8r/blob/main/LICENSE

package cmd

import (
	"testing"
)

func TestParseHTTPContext(t *testing.T) {
	tests := []struct {
		name         string
		errorString  string
		wantMethod   string
		wantEndpoint string
	}{
		{
			name:         "GET /objects error",
			errorString:  "[GET /objects][404] GetObjects default {\"error\":\"invalid class name: foo\"}",
			wantMethod:   "GET",
			wantEndpoint: "/objects",
		},
		{
			name:         "POST /graphs/neighbors error",
			errorString:  "[POST /graphs/neighbors][400] PostGraphsNeighbors default {\"error\":\"invalid query\"}",
			wantMethod:   "POST",
			wantEndpoint: "/graphs/neighbors",
		},
		{
			name:         "POST /graphs/goals error",
			errorString:  "[POST /graphs/goals][404] PostGraphsGoals default {\"error\":\"class not found: alert:alert\"}",
			wantMethod:   "POST",
			wantEndpoint: "/graphs/goals",
		},
		{
			name:         "GET /domains error",
			errorString:  "[GET /domains][500] GetDomains default {\"error\":\"internal server error\"}",
			wantMethod:   "GET",
			wantEndpoint: "/domains",
		},
		{
			name:         "empty string",
			errorString:  "",
			wantMethod:   "",
			wantEndpoint: "",
		},
		{
			name:         "no brackets",
			errorString:  "some random error message",
			wantMethod:   "",
			wantEndpoint: "",
		},
		{
			name:         "incomplete bracket",
			errorString:  "[GET /objects",
			wantMethod:   "",
			wantEndpoint: "",
		},
		{
			name:         "no space separator",
			errorString:  "[GET][404] GetObjects",
			wantMethod:   "",
			wantEndpoint: "",
		},
		{
			name:         "malformed but has brackets",
			errorString:  "[INVALID]",
			wantMethod:   "",
			wantEndpoint: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMethod, gotEndpoint := parseHTTPContext(tt.errorString)
			if gotMethod != tt.wantMethod {
				t.Errorf("parseHTTPContext() method = %v, want %v", gotMethod, tt.wantMethod)
			}
			if gotEndpoint != tt.wantEndpoint {
				t.Errorf("parseHTTPContext() endpoint = %v, want %v", gotEndpoint, tt.wantEndpoint)
			}
		})
	}
}

func TestParseHTTPContext_RealSwaggerErrors(t *testing.T) {
	// Test with actual swagger-generated error formats
	tests := []struct {
		name         string
		errorString  string
		wantMethod   string
		wantEndpoint string
	}{
		{
			name: "typical swagger client error",
			errorString: `[GET /api/v1alpha1/objects][400] GetObjects default  {
				"error": "invalid class name: not-a-class"
			}`,
			wantMethod:   "GET",
			wantEndpoint: "/api/v1alpha1/objects",
		},
		{
			name:         "swagger not found error",
			errorString:  "[POST /api/v1alpha1/graphs/neighbors][404] PostGraphsNeighbors default ",
			wantMethod:   "POST",
			wantEndpoint: "/api/v1alpha1/graphs/neighbors",
		},
		{
			name:         "null payload error (500)",
			errorString:  "[POST /graphs/neighbours][500] PostGraphsNeighbours default null",
			wantMethod:   "POST",
			wantEndpoint: "/graphs/neighbours",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMethod, gotEndpoint := parseHTTPContext(tt.errorString)
			if gotMethod != tt.wantMethod {
				t.Errorf("parseHTTPContext() method = %v, want %v", gotMethod, tt.wantMethod)
			}
			if gotEndpoint != tt.wantEndpoint {
				t.Errorf("parseHTTPContext() endpoint = %v, want %v", gotEndpoint, tt.wantEndpoint)
			}
		})
	}
}
