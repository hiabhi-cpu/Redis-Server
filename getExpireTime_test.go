package main

import (
	"testing"
	"time"
)

func createBulk(value string) RespValue {
	return RespValue{Type: BulkStringType, Value: value}
}

func TestGetExpireTime(t *testing.T) {
	nowMilli := time.Now().UnixMilli()

	tests := []struct {
		name     string
		cmdArray []RespValue
		wantErr  bool
		validate func(expireTime int64) bool
	}{
		{
			name: "EX 10s",
			cmdArray: []RespValue{
				createBulk("SET"),
				createBulk("mykey"),
				createBulk("hello"),
				createBulk("EX"),
				createBulk("10"),
			},
			wantErr: false,
			validate: func(et int64) bool {
				diff := et - nowMilli
				return diff >= 9500 && diff <= 10500
			},
		},
		{
			name: "PX 5000ms",
			cmdArray: []RespValue{
				createBulk("SET"),
				createBulk("mykey"),
				createBulk("hello"),
				createBulk("PX"),
				createBulk("5000"),
			},
			wantErr: false,
			validate: func(et int64) bool {
				diff := et - nowMilli
				return diff >= 4800 && diff <= 5200
			},
		},
		{
			name: "EXAT future time",
			cmdArray: []RespValue{
				createBulk("SET"),
				createBulk("mykey"),
				createBulk("hello"),
				createBulk("EXAT"),
				createBulk("9999999999"), // year 2286
			},
			wantErr: false,
			validate: func(et int64) bool {
				return et == 9999999999*1000
			},
		},
		{
			name: "PXAT future time",
			cmdArray: []RespValue{
				createBulk("SET"),
				createBulk("mykey"),
				createBulk("hello"),
				createBulk("PXAT"),
				createBulk("1729785600123"),
			},
			wantErr: false,
			validate: func(et int64) bool {
				return et == 1729785600123
			},
		},
		{
			name: "Unknown option",
			cmdArray: []RespValue{
				createBulk("SET"),
				createBulk("mykey"),
				createBulk("hello"),
				createBulk("ABC"),
				createBulk("123"),
			},
			wantErr: true,
		},
		{
			name: "Missing expiration value",
			cmdArray: []RespValue{
				createBulk("SET"),
				createBulk("mykey"),
				createBulk("hello"),
				createBulk("PX"),
			},
			wantErr: true,
		},
		{
			name: "No expiration option provided",
			cmdArray: []RespValue{
				createBulk("SET"),
				createBulk("mykey"),
				createBulk("hello"),
			},
			wantErr: false,
			validate: func(et int64) bool {
				return et == 0 // no expiry
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expireTime, err := GetExpireTime(tt.cmdArray)

			if (err != nil) != tt.wantErr {
				t.Errorf("expected error = %v, got = %v (%v)", tt.wantErr, err != nil, err)
				return
			}
			if !tt.wantErr && tt.validate != nil {
				if !tt.validate(expireTime) {
					t.Errorf("unexpected expire time: got %v", expireTime)
				}
			}
		})
	}
}
