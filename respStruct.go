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
