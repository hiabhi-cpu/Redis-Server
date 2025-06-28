package main

import (
	"reflect"
	"testing"
)

func TestDe_serialise_Errors(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []RespValue
		wantErr bool
	}{
		{
			name:    "Bulk string with incorrect length",
			input:   "$4\r\nabc\r\n", // Declares length 4, gives 3
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Array with fewer elements than declared",
			input:   "*3\r\n$3\r\nget\r\n$3\r\nkey\r\n", // Declares 3 elements, gives 2
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Integer with non-numeric value",
			input:   ":abc\r\n",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Bulk string with missing data line",
			input:   "$3\r\n", // Missing value line
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Error string missing \\r\\n",
			input:   "-Error message", // Incomplete message
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Null bulk string with missing pair",
			input:   "*2\r\n$-1\r\n$3\r\nkey\r\n", // Should have another value
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Array with missing size value",
			input:   "*\r\n$3\r\nget\r\n", // Invalid array size
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Two simple strings in single line",
			input:   "+OK\r\n+Second\r\n", // Your rule: only one simple string per line
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Bulk string claims 0 length but gives 1 char",
			input:   "$0\r\na\r\n",
			want:    nil,
			wantErr: true,
		},
		// {
		// 	name:    "Empty input",
		// 	input:   "",
		// 	want:    nil,
		// 	wantErr: false,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := De_serialise(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("De_serialise() error = %v, wantErr = %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("De_serialise() = %v, want = %v", got, tt.want)
			}
		})
	}
}
