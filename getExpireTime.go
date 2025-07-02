package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func GetExpireTime(cmdArray []RespValue) (int64, error) {
	var expireTime int64
	if len(cmdArray) > 3 {
		// fmt.Println(cmdArray)
		if len(cmdArray) < 5 {
			return 0, errors.New("SET requires a time for expiration")
		}
		option := strings.ToUpper(cmdArray[3].Value.(string))
		switch option {
		case "EX":
			tempSec, err := strconv.Atoi(cmdArray[4].Value.(string))
			if err != nil {
				return 0, errors.New("SET requires a time for expiration")
			}
			expireTime = time.Now().Add(time.Duration(tempSec) * time.Second).UnixMilli()

		case "PX":
			tempMilli, err := strconv.ParseInt(cmdArray[4].Value.(string), 10, 64)
			if err != nil {
				return 0, errors.New("SET requires a valid expiration time in milliseconds")

			}
			expireTime = time.Now().UnixMilli() + tempMilli

		case "EXAT":
			timestampSec, err := strconv.ParseInt(cmdArray[4].Value.(string), 10, 64)
			if err != nil {
				return 0, errors.New("SET requires valid EXAT time (UNIX seconds)")
			}
			expireTime = timestampSec * 1000 // Convert to milliseconds

		case "PXAT":
			timestampMilli, err := strconv.ParseInt(cmdArray[4].Value.(string), 10, 64)
			if err != nil {
				return 0, errors.New("SET requires valid PXAT time (UNIX ms)")
			}
			expireTime = timestampMilli

		default:
			return 0, errors.New("Unknown expiration option: must be EX or PX")

		}
	} else {
		expireTime = 0
	}
	return expireTime, nil
}
