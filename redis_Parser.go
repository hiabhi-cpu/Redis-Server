package main

import (
	"errors"
	"strconv"
	"strings"
)

func De_serialise(cmd string) (RespValue, error) {
	if !isValidCmd(cmd) {
		return RespValue{}, errors.New("Not valid strings")
	}
	cmds := strings.Split(cmd, "\r\n")
	// fmt.Println(len(cmds))
	// fmt.Println(isValidCmd(cmd))
	i := 0
	var res RespValue
	stringLineCnt := 0

	var currChar byte
	if len(cmds[i]) != 0 {
		currChar = cmds[i][0]
	} else {
		currChar = ' '
	}
	if currChar == '+' {
		if stringLineCnt > 0 {
			return RespValue{}, errors.New("Has multiple strings in same line")
		}
		stringLineCnt++
		temRes, err := parseString(cmds[i])
		if err != nil {
			return RespValue{}, errors.New("Need to have correct bulk string")
		}
		res = RespValue{
			Type:  SimpleStringType,
			Value: temRes,
		}
	} else if currChar == '$' {
		if cmds[i][1:] == "-1" {
			return RespValue{
				Type:  NullType,
				Value: "nil",
			}, nil
		}
		if i+1 >= len(cmds)-1 {
			return RespValue{}, errors.New("Need to have correct bulk string")
		}
		tempRes, err := parseBulkString(cmds[i], cmds[i+1])
		if err != nil {
			return RespValue{}, errors.New("Need to have correct bulk string")
		}
		res = RespValue{
			Type:  BulkStringType,
			Value: tempRes,
		}
	} else if currChar == '-' {
		// return
	}

	// fmt.Println(res)
	return res, nil
}

func isValidCmd(cmd string) bool {
	// fmt.Println("Is valid checking")
	cmds := strings.Split(cmd, "\r\n")

	return cmds[len(cmds)-1] == ""
}

func parseString(cmd string) (string, error) {
	// fmt.Println(cmd[1:])
	return cmd[1:], nil
}

func parseBulkString(size, cmd string) (string, error) {
	// fmt.Println(size)
	intSize, err := strconv.Atoi(size[1:])
	if err != nil {
		return "", getError("No correct int in bulk")
	}
	if intSize != len(cmd) {
		return "", getError("Bulk String length not correctly given")
	}
	// fmt.Println(len(cmd))
	return cmd, nil
}

func getError(str string) error {
	return errors.New(str)
}
