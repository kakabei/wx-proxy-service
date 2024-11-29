package common

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/levigross/grequests"
	"github.com/zeromicro/go-zero/core/logx"
)

func HttpPost(sessionId string, url string, timeoutInMS int, reqBody interface{}, rspBody interface{}) (errCode int, err error) {
	//global.VUA_LOG.Debug(fmt.Sprintf("[%v] url[%v], request[%+v]", sessionId, url, reqBody))
	logx.Debugf("[%v] url[%v], request:%s", sessionId, url, ToJSON(reqBody))
	// 请求参数
	ro := grequests.RequestOptions{
		RequestTimeout: time.Duration(timeoutInMS) * time.Millisecond,
		JSON:           reqBody,
	}

	// 发送请求
	grsp, e := grequests.Post(url, &ro)
	if e != nil {
		err = e
		errCode = -3001
		return
	}

	logx.Debugf("[%v] response:%s]", sessionId, ToJSON(grsp))
	return parseHttpResponse(grsp, rspBody)
}

func parseHttpResponse(resp *grequests.Response, body interface{}) (int, error) {
	if int(resp.StatusCode/100) != 2 {
		return -3001, fmt.Errorf("non-2xx status code[%v]", resp.StatusCode)
	}

	if body != nil {
		if err := json.Unmarshal(resp.Bytes(), body); err != nil {
			return -3001, fmt.Errorf("json Unmarshal failure, err[%v]", err)
		}
	}

	return 0, nil
}
