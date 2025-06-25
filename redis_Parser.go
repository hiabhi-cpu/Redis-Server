package main

import (
	"errors"
	"fmt"
	"strings"
)

func redis_Parser(cmd string) (string, error) {
	if !isValidCmd(cmd) {
		return "", errors.New("Not valid strings")
	}
	cmds := strings.Split(cmd, "\r\n")
	fmt.Println(len(cmds))
	fmt.Println(isValidCmd(cmd))
	i := 0
	stringLineCnt := 0
	for i = 0; i < len(cmds)-1; i++ {
		if cmds[i][0] == '+' {
			if stringLineCnt > 0 {
				return "", errors.New("Has multiple strings in same line")
			}
			stringLineCnt++
			parseString(cmds[i])
		} else if cmds[i][0] == '*' {
			if i+1 >= len(cmds)-1 {
				return "", errors.New("Need to have correct bulk string")
			}
			parseBulkString(cmds[i], cmds[i+1])
		}
	}
	return "", nil
}

func isValidCmd(cmd string) bool {
	// fmt.Println("Is valid checking")
	cmds := strings.Split(cmd, "\r\n")

	return cmds[len(cmds)-1] == ""
}

func parseString(cmd string) (string, error) {
	fmt.Println(cmd[1:])
	return cmd[1:], nil
}

func parseBulkString(size, cmd string) (string, error) {
	return "", nil
}
