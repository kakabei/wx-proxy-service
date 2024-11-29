package types

import (
	"encoding/xml"
)

var (
	Period int64 = 60 // 结算周期 即1分钟扣一次代币
)

type CtxValRequestId struct{}

var (
	DefautHeaderUrl string = "https://thirdwx.qlogo.cn/mmopen/vi_32/iavwib8v3s0vgVkva0rvx36NCibdfiaqqh5qgNYI8DiaB8Etb1xjJbv2N0ibzHtW4neguibnuow1RZzbOxgMLAwI3e1bQ/132"
	TimeFormat      string = "2006-01-02 15:04:05"

	WxQrCodeLoginCache string = "WxQrCodeLogin:"
)

const (
	PhoneModeLogin int = 0 // 登录
	PhoneModeBind  int = 1 // 绑定
	PhoneModeAlter int = 2 // 修改
)

type HttpCommonResponse struct {
	Head struct {
		Code      int    `json:"code"`
		Msg       string `json:"msg"`
		RequestID string `json:"request_id"`
	} `json:"ret"`
	Body interface{} `json:"body"`
}

type WxuseDataCache struct {
	AppId  string `json:"appid"`
	BizId  int64  `json:"biz_id"`
	Ticket string `json:"ticket"`
	Uid    uint64 `json:"uid"`
	UUID   string `json:"uuid"`
}

type ReceiveWxMsg struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
}
