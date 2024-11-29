package types

import (
	"strconv"
	"time"
)

func GenerateRspHead(Code int, Msg string) HTTPCommonHead {
	return HTTPCommonHead{Code: Code, Msg: Msg, RequestID: strconv.FormatInt(time.Now().Unix(), 10)}
}
