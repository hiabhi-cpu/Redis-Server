package main

import (
	"testing"
)

func TestDe_serialise(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		hasError bool
	}{
		// Simple Strings
		{"Simple OK", "+OK\r\n", "OK", false},
		{"Simple Hello", "+Hello\r\n", "Hello", false},
		{"Empty Simple", "+\r\n", "", false},
		{"Invalid Simple (no CRLF)", "+hello", "", true},

		// Bulk Strings
		{"Bulk Hello", "$5\r\nhello\r\n", "hello", false},
		{"Empty Bulk", "$0\r\n\r\n", "", false},
		{"Null Bulk", "$-1\r\n", "nil", false}, // assume your function returns "nil" for null bulk
		{"Invalid Bulk (no CRLF)", "$5\r\nhello", "", true},
		{"Invalid Bulk Length", "$abc\r\nhello\r\n", "", true},
		{"Short Bulk Content", "$5\r\nhi\r\n", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := De_serialise(tt.input)
			if tt.hasError {
				if err == nil {
					t.Errorf("expected error for input %q, got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %q: %v", tt.input, err)
				}
				if result.Value != tt.expected {
					t.Errorf("expected %q, got %q", tt.expected, result)
				}
			}
		})
	}
}
