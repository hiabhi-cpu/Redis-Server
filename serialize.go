package main

import "fmt"

func Serialise(resp []RespValue) (string, error) {
	str := ""
	for _, r := range resp {
		switch r.Type {
		case SimpleStringType:
			{
				str = str + "+" + r.Value.(string) + "\r\n"
			}
		case ErrorType:
			{
				str = str + "-" + r.Value.(string) + "\r\n"
			}
		case IntegerType:
			{
				str = str + ":" + fmt.Sprint(r.Value.(int)) + "\r\n"
			}
		case BulkStringType:
			{
				length := fmt.Sprint(len(r.Value.(string)))
				str = str + "$" + length + "\r\n" + r.Value.(string) + "\r\n"
			}
		case NullType:
			{
				str += "$-1\r\n"
			}
		case ArrayType:
			{
				s, err := Serialise(r.Value.([]RespValue))
				if err != nil {
					return "", GetError("Error in parsing array")
				}
				length := fmt.Sprint(len(r.Value.([]RespValue)))
				str = str + "*" + length + "\r\n" + s
			}
		}
	}
	// fmt.Println(str)

	return str, nil
}
