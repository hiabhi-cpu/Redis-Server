package main

type RespType int

const (
	SimpleStringType RespType = iota
	ErrorType
	IntegerType
	BulkStringType
	ArrayType
	NullType
)

type RespValue struct {
	Type  RespType `json:"type"`
	Value any      `json:"value"`
}

type Entry struct {
	Value  RespValue `json:"value"`
	Expire int64     `json:"expire"`
}
