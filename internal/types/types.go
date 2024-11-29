// Code generated by goctl. DO NOT EDIT.
package types

type CommonRet struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

type GetWxUserInfoReq struct {
	Source string `json:"source"` // 业务来源
	AppId  string `json:"appid"`  // 微信公众号 appid
	Code   string `json:"code"`   // 用户临时凭证
}

type GetWxUserInfoBody struct {
	AppId          string `json:"appid"`            // 微信公众号 appid
	AccessToken    string `json:"access_token"`     // 该 公众号 AppId 获取到的凭证
	OpenId         string `json:"openid"`           // 用户的OpenId
	RefreshToken   string `json:"refresh_token"`    // 用户刷新access_token
	Nickname       string `json:"nickname"`         // 用户昵称
	Sex            int    `json:"sex"`              // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Province       string `json:"province"`         // 用户个人资料填写的省份
	City           string `json:"city"`             // 普通用户个人资料填写的城市
	Country        string `json:"country"`          // 国家，如中国为CN
	HeadImgUrl     string `json:"headimgurl"`       // 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	UnionID        string `json:"unionid"`          // 只有在用户将公众号绑定到微信开放平台账号后，才会出现该字段。
	IsSnapshotUser int32  `json:"is_snapshot_user"` // 否为快照页模式虚拟账号，只有当用户是快照页模式虚拟账号时返回，值为1
}

type GetWxUserInfoResp struct {
	Ret  CommonRet         `json:"ret"`
	Body GetWxUserInfoBody `json:"body"`
}

type GetWxQrcodeReq struct {
	Source        string `json:"source"`         // 业务来源
	Scene         string `json:"scene"`          // 场景值
	ExpireSeconds int32  `json:"expire_seconds"` // 二维码， 有效时间
	AppId         string `json:"appid"`          // 微信公众号 appid
}

type GetWxQrcodeBody struct {
	FlowId        string `json:"flow_id"`        // 流水ID
	AppId         string `json:"appid"`          // 小程序的 AppID
	Ticket        string `json:"ticket"`         // 二维码 ticket
	QRBuffer      string `json:"qr_buffer"`      // 图片二进制流
	ExpireSeconds int32  `json:"expire_seconds"` // 二维码， 有效时间
}

type GetWxQrcodeResp struct {
	Ret  CommonRet       `json:"ret"`
	Body GetWxQrcodeBody `json:"body"`
}

type GetUserPhoneReq struct {
	Source string `json:"source"` // 业务来源
	AppId  string `json:"appid"`  // 小程序的 AppID
	Code   string `json:"code"`   // 手机号获取凭证
}

type GetUserPhoneBody struct {
	AppId       string `json:"appid"`           // 小程序的 AppID
	PhoneNumber string `json:"purePhoneNumber"` // 手机号
}

type GetUserPhoneResp struct {
	Ret  CommonRet        `json:"ret"`
	Body GetUserPhoneBody `json:"body"`
}