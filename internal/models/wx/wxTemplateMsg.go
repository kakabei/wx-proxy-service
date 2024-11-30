package wx

import (
	"context"
	"fmt"
	"time"

	"github.com/levigross/grequests"
	"github.com/zeromicro/go-zero/core/logx"
)

// 查询公众号的 access_token
// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html

func SendTemplateMessage(ctx context.Context, requestId, appId, accessToken, data string) (err error) {

	wxUrl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", accessToken)

	logx.Infof("[%s] SendTemplateMessage appId[%s] accessToken:%s data : %s", requestId, appId, accessToken, data)

	// 请求参数
	ro := grequests.RequestOptions{
		RequestTimeout: time.Duration(3000) * time.Millisecond,
		JSON:           data,
	}

	// 发送请求
	grsp, err := grequests.Post(wxUrl, &ro)
	if err != nil {
		logx.Errorf("[%s] SendTemplateMessage appId[%s] accessToken:%s data : %s, err:%+v", requestId, appId, accessToken, data, err)
		return
	}

	if int(grsp.StatusCode/100) != 2 {
		logx.Errorf("[%s] SendTemplateMessage appId[%s] accessToken:%s data : %s, grsp:%+v", requestId, appId, accessToken, data, grsp)
		return
	}

	logx.Infof("[%s] SendTemplateMessage appId[%s] grsp : %+v", requestId, appId, grsp)

	return nil
}
