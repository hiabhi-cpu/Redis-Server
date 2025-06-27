package main

import (
	"reflect"
	"testing"
)

func TestRESPStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []RespValue
	}{
		{
			name:  "Null Bulk String",
			input: "$-1\r\n",
			expected: []RespValue{
				{Type: NullType, Value: "nil"},
			},
		},
		{
			name:  "Single Bulk String in Array",
			input: "*1\r\n$4\r\nping\r\n",
			expected: []RespValue{
				{
					Type: ArrayType,
					Value: []RespValue{
						{Type: BulkStringType, Value: "ping"},
					},
				},
			},
		},
		{
			name:  "Two Bulk Strings in Array (echo hello world)",
			input: "*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n",
			expected: []RespValue{
				{
					Type: ArrayType,
					Value: []RespValue{
						{Type: BulkStringType, Value: "echo"},
						{Type: BulkStringType, Value: "hello world"},
					},
				},
			},
		},
		{
			name:  "Two Bulk Strings in Array (get key)",
			input: "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n",
			expected: []RespValue{
				{
					Type: ArrayType,
					Value: []RespValue{
						{Type: BulkStringType, Value: "get"},
						{Type: BulkStringType, Value: "key"},
					},
				},
			},
		},
		{
			name:  "Simple OK",
			input: "+OK\r\n",
			expected: []RespValue{
				{Type: SimpleStringType, Value: "OK"},
			},
		},
		{
			name:  "Error message",
			input: "-Error message\r\n",
			expected: []RespValue{
				{Type: ErrorType, Value: "Error message"},
			},
		},
		{
			name:  "Empty Bulk String",
			input: "$0\r\n\r\n",
			expected: []RespValue{
				{Type: BulkStringType, Value: ""},
			},
		},
		{
			name:  "Simple String: hello world",
			input: "+hello world\r\n",
			expected: []RespValue{
				{Type: SimpleStringType, Value: "hello world"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := De_serialise(tt.input)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}
