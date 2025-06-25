package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestParseSimpleStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected string
		hasError bool
	}{
		{
			name:     "OK reply",
			input:    []byte("+OK\r\n"),
			expected: "OK",
		},
		{
			name:     "hello reply",
			input:    []byte("+hello\r\n"),
			expected: "hello",
		},
		{
			name:     "empty reply",
			input:    []byte("+\r\n"),
			expected: "",
		},
		{
			name:     "no CRLF",
			input:    []byte("+hello"),
			hasError: true,
		},
		{
			name:     "multiple lines",
			input:    []byte("+hello\r\n+world\r\n"),
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseRESP(tt.input)

			if tt.hasError && err == nil {
				t.Fatalf("expected error but got none for input %q", tt.input)
			}
			if !tt.hasError && err != nil {
				t.Fatalf("unexpected error for input %q: %v", tt.input, err)
			}
			if !tt.hasError && result != tt.expected {
				t.Fatalf("expected %q but got %q for input %q", tt.expected, result, tt.input)
			}
		})
	}
}

// Example Parser (you'd have this implemented in your resp.go)
func ParseRESP(data []byte) (string, error) {
	if len(data) == 0 || data[0] != '+' {
		return "", fmt.Errorf("invalid simple string")
	}
	// Find the end of the line
	end := bytes.Index(data, []byte("\r\n"))
	if end == -1 {
		return "", fmt.Errorf("no CRLF found")
	}
	// Extract the payload
	return string(data[1:end]), nil
}
