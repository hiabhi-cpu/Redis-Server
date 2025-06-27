package main

import "fmt"

func main() {
	fmt.Println("Hello")
	res, err := De_serialise("*2\r\n$3\r\nget\r\n$3\r\nkey\r\n")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}
