package common

import (
	"context"
	"fmt"
	"time"
	"wx-proxy-service/internal/config"
	"wx-proxy-service/internal/types"

	"github.com/levigross/grequests"
	"github.com/zeromicro/go-zero/core/logx"
)

func GeWxOpenIdConfig(config config.Config, openId string) (_ *config.WxOpenIdListConfig, err error) {
	if len(config.WxMsgMgr.WxOpenIdList) == 0 {
		return nil, fmt.Errorf("没有配置WxOpenIdList")
	}
	for _, openIdConf := range config.WxMsgMgr.WxOpenIdList {

		if StringInArray(openId, openIdConf.OpenIdList) {
			return &openIdConf, nil
		}
	}
	// 没有配置返回最后一个
	return &config.WxMsgMgr.WxOpenIdList[len(config.WxMsgMgr.WxOpenIdList)-1], nil
}

func PosWxMsg(ctx context.Context, requestId string, req *types.ReceiveWxMsgReq, url string) (_ *grequests.Response, err error) {
	ro := grequests.RequestOptions{
		RequestTimeout: time.Duration(3000) * time.Millisecond,
		XML:            req,
	}

	// 发送请求
	grsp, err := grequests.Post(url, &ro)
	if err != nil {
		logx.WithContext(ctx).Errorf("[%s] PosWxMsg err. url[%s] grsp:%s", requestId, url, grsp.String())
		return nil, err
	}

	logx.WithContext(ctx).Debugf("[%s] PosWxMsg grsp:%s", requestId, grsp.String())

	return grsp, nil

}
