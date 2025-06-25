package main

import "fmt"

func main() {
	fmt.Println("Hello")
	redis_Parser("$2\r\nhi\r\n")
}
