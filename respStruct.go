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
	Type  RespType
	Value any
}

type Entry struct {
	Value  RespValue
	Expire int64
}
