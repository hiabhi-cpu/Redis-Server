package main

import "fmt"

func main() {
	fmt.Println("Hello")
	_, err := De_serialise("$2\r\nOK\r\n")
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(str)
}
