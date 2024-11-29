package wx

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"wx-proxy-service/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpc"
)

type WxAppMgr struct {
	Redis *redis.Redis
}

func NewWxAppMgr(appConf config.Config) *WxAppMgr {

	WxApp := new(WxAppMgr)

	WxApp.Redis = redis.New(appConf.Redis.Host, func(r *redis.Redis) {
		r.Type = appConf.Redis.Type
		r.Pass = appConf.Redis.Pass
	})

	return WxApp
}

// 查询小程序的 access_token
func (s *WxAppMgr) GetMiniAppAccessToken(flowId string, appId string, AppSecret string) (miniAppsAccessToken MiniProgramAccessToken, err error) {

	//key := fmt.Sprintf("%s%s", MiniAppsAcccessTokenKey, appId)
	//value, err := s.Redis.Get(key)
	//if err == nil && len(value) > 0 {
	//	logx.Errorf("[%s] GetMiniAppAccessToken Redis.Get success. appid[%s] value: %+v", flowId, appId, value)
	//	miniAppsAccessToken.AccessToken = value
	//	return
	//} else if err != nil || len(value) == 0 {
	//	logx.Errorf("[%s] GetMiniAppAccessToken Redis.Get err. appid[%s] value[%s] err:%+v", flowId, appId, value, err)
	//}

	wxUrl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?appid=%s&secret=%s&grant_type=client_credential", appId, AppSecret)

	resp, err := httpc.Do(context.Background(), http.MethodGet, wxUrl, nil)
	if err != nil || resp.Status != "200 OK" {
		logx.Errorf("[%s] GetMiniAppAccessToken httpc.Do get err, appId[%s] %+v  resp:%+v", flowId, appId, err, resp)
		return
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("[%s] GetMiniAppAccessToken io.ReadAll err , appId[%s] err:%+v ", flowId, appId, err)
		return
	}

	logx.Infof("[%s] GetMiniAppAccessToken: appId[%s] Body : %+v", flowId, appId, string(bytes))

	var m map[string]interface{}
	err = json.Unmarshal([]byte(bytes), &m)
	if err != nil {
		logx.Errorf("[%s] GetMiniAppAccessToken json.Unmarshal err , appId[%s] %+v  resp:%+v", flowId, appId, err, string(bytes))

		return
	}

	if errcode, ok := m["errcode"]; ok && errcode.(float64) != 0 {
		err = errors.New(string(bytes))
		logx.Errorf("[%s] GetMiniAppAccessToken: appId[%s]  m : %+v", flowId, appId, m)
		return
	}

	miniAppsAccessToken.AccessToken = m["access_token"].(string)
	//expiresIn := int64(m["expires_in"].(float64))

	// 写入缓存
	//if err := s.Redis.Setex(key, miniAppsAccessToken.AccessToken, int(expiresIn-10)); err != nil {
	//	logx.Errorf("[%s] GetMiniAppAccessToken: s.Redis.Setex err. appid[%s] err: %+v", flowId, appId, err)
	//}

	return
}

// 查询用户手机号
func (s *WxAppMgr) GetUserPphoneNumber(flowId string, appId, accessToken string, code string) (userResp UserPphoneNumberResp, err error) {

	wxUrl := fmt.Sprintf("https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token=%s", accessToken)

	type wxGetUserPhoneNumberReq struct {
		Code string `json:"code"` // 手机号获取凭证
	}
	type wxUserPphoneNumberResp struct {
		Errcode   int      `json:"errcode"` // 错误码
		Errmsg    string   `json:"errmsg"`  // 错误信息
		PhoneInfo struct { // 用户手机号信息
			PhoneNumber     string   `json:"phoneNumber"`     // 手机号
			PurePhoneNumber string   `json:"purePhoneNumber"` // 纯净手机号
			CountryCode     string   `json:"countryCode"`     // 国家码
			Watermark       struct { // 水印
				Timestamp int    `json:"timestamp"` // 时间戳
				Appid     string `json:"appid"`     // 小程序appid
			} `json:"watermark"`
		} `json:"phone_info"`
	}
	req := wxGetUserPhoneNumberReq{
		Code: code,
	}
	data, err := httpc.Do(context.Background(), http.MethodPost, wxUrl, req)
	if err != nil || data.Status != "200 OK" {
		logx.Errorf("[%s] httpc.Do get err, appId[%s] %+v  data:%+v", flowId, appId, err, data)
		return
	}

	bytes, err := io.ReadAll(data.Body)
	if err != nil {
		logx.Errorf("[%s] io.ReadAll err , appId[%s] err:%+v ", flowId, appId, err)
		return
	}

	logx.Infof("[%s] GetUserPphoneNumber: appId[%s]  Body : %+v", flowId, appId, string(bytes))

	resp := new(wxUserPphoneNumberResp)
	err = json.Unmarshal([]byte(bytes), resp)
	if err != nil {
		logx.Errorf("[%s] json.Unmarshal err , appId[%s] %+v  resp:%+v", flowId, appId, err, string(bytes))

		return
	}

	if resp.Errcode != 0 {
		err = errors.New(string(bytes))
		logx.Errorf("[%s] GetUserPphoneNumber: appId[%s]  resp : %+v", flowId, appId, resp)
		return
	}

	userResp.AppId = appId
	userResp.PhoneNumber = resp.PhoneInfo.PurePhoneNumber
	return
}

// 创建小程序的二维码
func (s *WxAppMgr) GetUnlimitedQRCode(flowId string, appId, accessToken string, scene string) (QRResp GetUnlimitedQRCodeResp, err error) {

	wxUrl := fmt.Sprintf("https://api.weixin.qq.com/wxa/getwxacodeunlimit?access_token=%s", accessToken)

	type wxGetUnlimitedQRCodeReq struct {
		Scene string `json:"scene"` // 场景值为 设备Id
	}

	req := wxGetUnlimitedQRCodeReq{
		Scene: scene,
	}

	resp, err := httpc.Do(context.Background(), http.MethodPost, wxUrl, req)
	if err != nil || resp.Status != "200 OK" {
		logx.Errorf("[%s] httpc.Do get err, appId[%s] %+v  resp:%+v", flowId, appId, err, resp)
		return
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("[%s] io.ReadAll err , appId[%s] err:%+v", flowId, appId, err)
		return
	}

	if json.Valid(bytes) { // 如果是json则说明已经有错误码
		var m map[string]interface{}
		err = json.Unmarshal([]byte(bytes), &m)
		logx.Errorf("[%s] GetUnlimitedQRCode json.Unmarshal err, appId[%s]  edata:%+v err:%+v", flowId, appId, string(bytes), err)
		return
	}
	QRResp.AppId = appId
	encodeUrl := base64.StdEncoding.EncodeToString(bytes)
	if len(encodeUrl) != 0 {
		QRResp.EncodedData = fmt.Sprintf("data:image/jpeg;base64,%s", base64.StdEncoding.EncodeToString(bytes))
	}

	return
}
