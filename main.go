package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func main() {
	// fmt.Println("Hello")
	// res, err := De_serialise("*2\r\n$3\r\nget\r\n$3\r\nkey\r\n")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(res)
	// strRes, err := Serialise(res)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
	// fmt.Println(strRes)
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Redis lite server is listening to port :6379....")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

var buffHash = make(map[string]Entry)
var mux sync.Mutex

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New client connected")

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}

		input := string(buf[:n])
		fmt.Println("Received:  ", input)

		res, err := De_serialise(input)
		if err != nil {
			fmt.Println("Error in deserialise")
			conn.Write([]byte("-ERR invalid input\r\n"))
			continue
		}
		res = res[:1]
		// Redis commands are sent as array
		if len(res) != 1 || res[0].Type != ArrayType {
			conn.Write([]byte("-ERR invalid command format\r\n"))
			continue
		}

		// Extract the command and args
		cmdArray := res[0].Value.([]RespValue)
		if len(cmdArray) == 0 {
			conn.Write([]byte("-ERR empty command\r\n"))
			continue
		}

		command := strings.ToUpper(cmdArray[0].Value.(string))
		// fmt.Println(res)
		fmt.Println(cmdArray)
		var reply RespValue
		mux.Lock()
		switch command {
		case "PING":
			reply = RespValue{Type: SimpleStringType, Value: "PONG"}
		case "ECHO":
			if len(cmdArray) < 2 {
				reply = RespValue{Type: ErrorType, Value: "ECHO requires a message"}
			} else {
				reply = RespValue{Type: BulkStringType, Value: cmdArray[1].Value.(string)}
			}
		case "SET":
			if len(cmdArray) < 3 {
				reply = RespValue{Type: ErrorType, Value: "SET requires a key value pair"}
			} else {
				expireTime, err := GetExpireTime(cmdArray)
				if err != nil {
					reply = RespValue{Type: ErrorType, Value: err}
					break
				}
				buffHash[cmdArray[1].Value.(string)] = Entry{Value: cmdArray[2], Expire: expireTime}
				reply = RespValue{Type: SimpleStringType, Value: "OK"}
			}
		case "GET":
			if len(cmdArray) < 2 {
				reply = RespValue{Type: ErrorType, Value: "GET requires a key"}
			} else {
				entry, err := buffHash[cmdArray[1].Value.(string)]
				if err == false {
					reply = RespValue{Type: ErrorType, Value: "Key not found"}
					// fmt.Println(buffHash)
				} else {
					if entry.Expire > 0 && entry.Expire < time.Now().UnixMilli() {
						delete(buffHash, cmdArray[1].Value.(string))
						reply = RespValue{Type: ErrorType, Value: "Key expired "}
						// fmt.Println(buffHash)
					} else {
						reply = RespValue{Type: BulkStringType, Value: entry.Value.Value.(string)}
					}

				}
			}
		default:
			reply = RespValue{Type: ErrorType, Value: "unknown command"}
		}
		mux.Unlock()

		serial, err := Serialise([]RespValue{reply})
		if err != nil {
			fmt.Println("Error in serializing")
			continue
		}
		conn.Write([]byte(serial))
	}

}
