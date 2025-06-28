package main

import (
	"testing"
)

func TestSerialise(t *testing.T) {
	tests := []struct {
		name     string
		input    []RespValue
		expected string
	}{
		{
			name: "Simple String",
			input: []RespValue{
				{Type: SimpleStringType, Value: "OK"},
			},
			expected: "+OK\r\n",
		},
		{
			name: "Integer",
			input: []RespValue{
				{Type: IntegerType, Value: 42},
			},
			expected: ":42\r\n",
		},
		{
			name: "Error",
			input: []RespValue{
				{Type: ErrorType, Value: "Something went wrong"},
			},
			expected: "-Something went wrong\r\n",
		},
		{
			name: "Bulk String",
			input: []RespValue{
				{Type: BulkStringType, Value: "hello"},
			},
			expected: "$5\r\nhello\r\n",
		},
		{
			name: "Null",
			input: []RespValue{
				{Type: NullType, Value: "nil"},
			},
			expected: "$-1\r\n",
		},
		{
			name: "Array",
			input: []RespValue{
				{
					Type: ArrayType,
					Value: []RespValue{
						{Type: BulkStringType, Value: "get"},
						{Type: BulkStringType, Value: "key"},
					},
				},
			},
			expected: "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			output, err := Serialise(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if output != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, output)
			}
		})
	}
}
