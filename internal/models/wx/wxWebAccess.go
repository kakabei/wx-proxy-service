package wx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"wx-proxy-service/internal/common"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"
)

// from： https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html
// 网页版本 Access 信息
type WebAccessToken struct {
	AccessToken    string `json:"access_token"`    // 网页授权接口调用凭证,注意：此access_token与基础支持的access_token不同
	ExpiresIn      int64  `json:"expires_in"`      // access_token接口调用凭证超时时间，单位（秒)
	RefreshToken   string `json:"refresh_token"`   // 用户刷新access_token
	OpenId         string `json:"openid"`          // 用户唯一标识，请注意，在未关注公众号时，用户访问公众号的网页，也会产生一个用户和公众号唯一的OpenID
	Scope          string `json:"scope"`           // 用户授权的作用域，使用逗号（,）分隔
	IsSnapshotUser int32  `json:"is_snapshotuser"` // 是否为快照页模式虚拟账号，只有当用户是快照页模式虚拟账号时返回，值为1
	UnionID        string `json:"unionid"`         // 用户统一标识（针对一个微信开放平台账号下的应用，同一用户的 unionid 是唯一的），只有当scope为"snsapi_userinfo"时返回
}

// 微信用户信息
type WebWxUserInfo struct {
	OpenId     string   `json:"openid"`     // 用户唯一标识，请注意，在未关注公众号时，用户访问公众号的网页，也会产生一个用户和公众号唯一的OpenID
	Nickname   string   `json:"nickname"`   // 用户昵称
	Sex        int      `json:"sex"`        // 用户的性别，值为1时是男性，值为2时是女性，值为0时是未知
	Province   string   `json:"province"`   // 用户个人资料填写的省份
	City       string   `json:"city"`       // 普通用户个人资料填写的城市
	Country    string   `json:"country"`    // 国家，如中国为CN
	HeadImgUrl string   `json:"headimgurl"` // 用户头像，最后一个数值代表正方形头像大小（有0、46、64、96、132数值可选，0代表640*640正方形头像），用户没有头像时该项为空。若用户更换头像，原有头像URL将失效。
	Privilege  []string `json:"privilege"`  // 用户特权信息，json 数组，如微信沃卡用户为（chinaunicom）
	UnionID    string   `json:"unionid"`    // 只有在用户将公众号绑定到微信开放平台账号后，才会出现该字段。
}

type CheckWebAccessTokenRet struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// https://developers.weixin.qq.com/doc/offiaccount/OA_Web_Apps/Wechat_webpage_authorization.html
// access_token 考虑 缓存到 redis
func GetWebAccessToken(flowId string, appId string, appSecret string, code string) (webAccessToken WebAccessToken, err error) {

	wxUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code", appId, appSecret, code)

	resp, err := httpc.Do(context.Background(), http.MethodGet, wxUrl, nil)
	if err != nil || resp.Status != "200 OK" {
		logx.Errorf("[%s] GetWebAccessToken httpc.Do err, appId[%s] %+v  resp:%s", flowId, appId, err, common.ToJSON(resp))
		return
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("[%s] GetWebAccessToken io.ReadAll err, appId[%s] err:%+v ", flowId, appId, err)
		return
	}

	logx.Infof("[%s] GetWebAccessToken: appId[%s]  Body : %+v", flowId, appId, string(bytes))

	var m map[string]interface{}
	err = json.Unmarshal([]byte(bytes), &m)
	if err != nil {
		logx.Errorf("[%s] GetWebAccessToken json.Unmarshal err , appId[%s] %+v  resp:%+v", flowId, appId, err, string(bytes))
		return
	}

	if errcode, ok := m["errcode"]; ok && errcode.(float64) != 0 {
		err = errors.New(string(bytes))
		logx.Errorf("[%s] GetWebAccessToken: appId[%s]  m : %+v", flowId, appId, m)
		return
	}

	webAccessToken.AccessToken = m["access_token"].(string)
	webAccessToken.ExpiresIn = int64(m["expires_in"].(float64))
	webAccessToken.RefreshToken = m["refresh_token"].(string)
	webAccessToken.OpenId = m["openid"].(string)
	webAccessToken.Scope = m["scope"].(string)

	if value, ok := m["is_snapshotuser"]; ok {
		webAccessToken.IsSnapshotUser = int32(value.(float64))
	}

	return
}

func RefreshWebAccessToken(flowId string, appId string, refreshToken string) (webAccessToken WebAccessToken, err error) {

	wxUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/refresh_token?appid=%s&grant_type=refresh_token&refresh_token=%s", appId, refreshToken)

	resp, err := httpc.Do(context.Background(), http.MethodGet, wxUrl, nil)
	if err != nil || resp.Status != "200 OK" {
		logx.Errorf("[%s] httpc.Do get err , appId[%s] %+v  resp:%+v", flowId, appId, err, resp)
		return
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("[%s] io.ReadAll err , appId[%s] err:%+v ", flowId, appId, err)
		return
	}

	logx.Infof("[%s] RefreshWebAccessToken: appId[%s]  Body : %+v", flowId, appId, string(bytes))

	var m map[string]interface{}
	err = json.Unmarshal([]byte(bytes), &m)
	if err != nil {
		logx.Errorf("[%s] json.Unmarshal err , appId[%s] %+v  resp:%+v", flowId, appId, err, string(bytes))

		return
	}

	if errcode, ok := m["errcode"]; ok && errcode.(float64) != 0 {
		err = errors.New(string(bytes))
		logx.Errorf("[%s] GetWebAccessToken: appId[%s]  m : %+v", flowId, appId, m)
		return
	}

	return
}

func CheckWebAccessToken(flowId, appId string, accessToken string) (accessTokenValid bool, err error) {

	wxUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/auth?access_token=%s&openid=%s", accessToken, appId)

	resp, err := httpc.Do(context.Background(), http.MethodGet, wxUrl, nil)
	if err != nil || resp.Status != "200 OK" {
		logx.Errorf("[%s] httpc.Do get err , appId[%s] %+v  resp:%+v", flowId, appId, err, resp)
		return
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("[%s] io.ReadAll err , appId[%s] err:%+v ", flowId, appId, err)
		return
	}

	logx.Infof("[%s] RefreshWebAccessToken: appId[%s]  Body : %+v", flowId, appId, string(bytes))

	var m map[string]interface{}
	err = json.Unmarshal([]byte(bytes), &m)
	if err != nil {
		logx.Errorf("[%s] json.Unmarshal err , appId[%s] %+v  resp:%+v", flowId, appId, err, string(bytes))

		return
	}

	errcode, ok := m["errcode"]
	if ok {
		if errcode.(float64) == 40003 {
			return false, nil
		} else if errcode.(float64) == 400 {
			return true, nil
		}
	}
	return
}

func GetWxUserinfo(flowId string, appId string, accessToken string) (webWxUserInfo WebWxUserInfo, err error) {

	wxUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s&lang=zh_CN", accessToken, appId)

	resp, err := httpc.Do(context.Background(), http.MethodGet, wxUrl, nil)
	if err != nil || resp.Status != "200 OK" {
		logx.Errorf("[%s] httpc.Do get err , appId[%s] %+v  resp:%+v", flowId, appId, err, resp)
		return
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("[%s] io.ReadAll err , appId[%s] err: %+v", flowId, appId, err)
		return
	}

	logx.Infof("[%s] GetWxUserinfo: appId[%s]  Body : %+v", flowId, appId, string(bytes))

	var m map[string]interface{}
	err = json.Unmarshal([]byte(bytes), &m)
	if err != nil {
		logx.Errorf("[%s] json.Unmarshal err , appId[%s] %+v  resp:%+v", flowId, appId, err, string(bytes))

		return
	}

	if errcode, ok := m["errcode"]; ok && errcode.(float64) != 0 {
		err = errors.New(string(bytes))
		logx.Errorf("[%s] GetWebAccessToken: appId[%s]  m : %+v", flowId, appId, m)
		return
	}

	webWxUserInfo.Nickname = m["nickname"].(string)
	webWxUserInfo.HeadImgUrl = m["headimgurl"].(string)

	return
}
