package main

import (
	"errors"

	"strconv"
	"strings"
)

func De_serialise(cmd string) ([]RespValue, error) {
	if !isValidCmd(cmd) {
		return []RespValue{}, errors.New("Not valid strings")
	}
	cmds := strings.Split(cmd, "\r\n")
	if len(cmds) == 0 {
		return []RespValue{}, errors.New("Not valid strings")
	}
	i := 0
	res := make([]RespValue, 0)
	stringLineCnt := 0

	for i < len(cmds) {
		var currChar byte
		if len(cmds[i]) != 0 {
			currChar = cmds[i][0]
		} else {
			currChar = ' '
		}
		if currChar == '+' {
			if stringLineCnt > 0 {
				return []RespValue{}, errors.New("Has multiple strings in same line")
			}
			stringLineCnt++
			temRes, err := parseString(cmds[i])
			if err != nil {
				return []RespValue{}, errors.New("Need to have correct bulk string")
			}

			res = append(res, temRes)
		} else if currChar == '$' {
			if cmds[i][1:] == "-1" {
				return append(res, RespValue{
					Type:  NullType,
					Value: "nil",
				}), nil
			}
			if i+1 >= len(cmds)-1 {
				return []RespValue{}, errors.New("Need to have correct bulk string")
			}
			tempRes, err := parseBulkString(cmds[i], cmds[i+1])
			if err != nil {
				return []RespValue{}, errors.New("Need to have correct bulk string")
			}

			res = append(res, tempRes)
		} else if currChar == '-' {
			tempRes, err := parseError(cmds[i])
			if err != nil {
				return []RespValue{}, errors.New("Need to have correct bulk string")
			}
			res = append(res, tempRes)
		} else if currChar == ':' {
			tempRes, err := parseInt(cmds[i])
			if err != nil {
				return []RespValue{}, errors.New("Need to have correct bulk string")
			}
			res = append(res, tempRes)
		} else if currChar == '*' {
			n, tempRes, err := parseArray(cmds[i:])
			if err != nil {
				return []RespValue{}, errors.New("Error in array")
			}
			i += n
			res = append(res, tempRes)
		}
		i++
	}
	return res, nil
}

func isValidCmd(cmd string) bool {
	cmds := strings.Split(cmd, "\r\n")

	return cmds[len(cmds)-1] == ""
}

func parseString(cmd string) (RespValue, error) {
	return RespValue{
		Type:  SimpleStringType,
		Value: cmd[1:],
	}, nil
}

func parseBulkString(size, cmd string) (RespValue, error) {
	intSize, err := strconv.Atoi(size[1:])
	if err != nil {
		return RespValue{}, GetError("No correct int in bulk")
	}
	if intSize != len(cmd) {
		return RespValue{}, GetError("Bulk String length not correctly given")
	}
	// fmt.Println(len(cmd))
	return RespValue{
		Type:  BulkStringType,
		Value: cmd,
	}, nil
}

func parseError(cmd string) (RespValue, error) {
	return RespValue{
		Value: cmd[1:],
		Type:  ErrorType,
	}, nil
}

func parseInt(cmd string) (RespValue, error) {
	n, err := strconv.Atoi(cmd[1:])
	if err != nil {
		return RespValue{}, err
	}
	return RespValue{
		Value: n,
		Type:  IntegerType,
	}, nil
}

func parseArray(cmd []string) (int, RespValue, error) {
	n, err := strconv.Atoi(cmd[0][1:])
	if err != nil {
		return 0, RespValue{}, err
	}
	tempCmd := ""
	for i := 1; i < len(cmd); i++ {
		tempCmd = tempCmd + cmd[i] + "\r\n"
	}
	tempResVal, err := De_serialise(tempCmd)
	if err != nil {
		return 0, RespValue{}, err
	}
	res := RespValue{
		Value: tempResVal,
		Type:  ArrayType,
	}
	if n != len(tempResVal) {
		return 0, RespValue{}, GetError("Not correct array elements")
	}
	return n + 1, res, nil
}

func GetError(str string) error {
	return errors.New(str)
}
