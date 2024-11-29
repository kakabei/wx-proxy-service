package wx

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"wx-proxy-service/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest/httpc"
)

// 微信公众号
type WxOfficialAccountMgr struct {
	Redis *redis.Redis
}

func NewWxOfficialAccountMgr(appConf config.Config) *WxOfficialAccountMgr {

	WxApp := new(WxOfficialAccountMgr)

	WxApp.Redis = redis.New(appConf.Redis.Host, func(r *redis.Redis) {
		r.Type = appConf.Redis.Type
		r.Pass = appConf.Redis.Pass
	})

	return WxApp
}

type QRCodeTicketRequest struct {
	ExpireSeconds int    `json:"expire_seconds"`
	ActionName    string `json:"action_name"`
	ActionInfo    struct {
		Scene struct {
			SceneStr string `json:"scene_str"`
		} `json:"scene"`
	} `json:"action_info"`
}

type QRCodeTicketResponse struct {
	Ticket        string `json:"ticket"`
	ExpireSeconds int64  `json:"expire_seconds"`
	URL           string `json:"url"`
}

type QRCodeResponse struct {
	AppId       string `json:"appid"`        //  小程序AppId
	EncodedData string `json:"encoded_data"` // 二维码	buffer
}

// 查询公众号的 access_token
// https://developers.weixin.qq.com/doc/offiaccount/Basic_Information/Get_access_token.html
func (s *WxOfficialAccountMgr) GetOfficialAccountAccessToken(flowId string, appId string, AppSecret string) (miniAppsAccessToken MiniProgramAccessToken, err error) {

	key := fmt.Sprintf("%s%s", MiniAppsAcccessTokenKey, appId)
	value, err := s.Redis.Get(key)
	if err == nil && len(value) > 0 {
		logx.Infof("[%s] GetOfficialAccountAccessToken: ss.Redis.Get success. appid[%s] value: %+v", flowId, appId, value)
		miniAppsAccessToken.AccessToken = value
		return
	} else if err != nil || len(value) == 0 {
		logx.Errorf("[%s] GetOfficialAccountAccessToken: ss.Redis.Get err. appid[%s] value[%s] err:%+v", flowId, appId, value, err)
	}

	wxUrl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?appid=%s&secret=%s&grant_type=client_credential", appId, AppSecret)

	resp, err := httpc.Do(context.Background(), http.MethodGet, wxUrl, nil)
	if err != nil || resp.Status != "200 OK" {
		logx.Errorf("[%s] GetOfficialAccountAccessToken httpc.Do get err, appId[%s] %+v  resp:%+v", flowId, appId, err, resp)
		return
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("[%s] GetOfficialAccountAccessToken io.ReadAll err , appId[%s] err:%+v ", flowId, appId, err)
		return
	}

	logx.Infof("[%s] GetOfficialAccountAccessToken: appId[%s] Body : %+v", flowId, appId, string(bytes))

	var m map[string]interface{}
	err = json.Unmarshal([]byte(bytes), &m)
	if err != nil {
		logx.Errorf("[%s] GetOfficialAccountAccessToken json.Unmarshal err , appId[%s] %+v  resp:%+v", flowId, appId, err, string(bytes))

		return
	}

	if errcode, ok := m["errcode"]; ok && errcode.(float64) != 0 {
		err = errors.New(string(bytes))
		logx.Errorf("[%s] GetOfficialAccountAccessToken: appId[%s]  m : %+v", flowId, appId, m)
		return
	}

	miniAppsAccessToken.AccessToken = m["access_token"].(string)
	expiresIn := int64(m["expires_in"].(float64))

	// 写出入缓存
	if err := s.Redis.Setex(key, miniAppsAccessToken.AccessToken, int(expiresIn-10)); err != nil {
		logx.Errorf("[%s] GetOfficialAccountAccessToken: s.Redis.Setex err. appid[%s] err: %+v", flowId, appId, err)
	}

	return
}

// 生成带参数二维码
// https://developers.weixin.qq.com/doc/offiaccount/Account_Management/Generating_a_Parametric_QR_Code.html
// 查询 ticket
func (s *WxOfficialAccountMgr) GetWxQrCodeTicket(flowId, appId, accessToken string, actionName, sceneStr string, expireSeconds int) (ticketInfo QRCodeTicketResponse, err error) {

	wxUrl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%s", accessToken)
	Req := QRCodeTicketRequest{
		ExpireSeconds: expireSeconds,
		ActionName:    actionName,
	}
	Req.ActionInfo.Scene.SceneStr = sceneStr

	logx.Infof("[%s] GetWxQrCodeTicket appId[%s]  Req : %+v", flowId, appId, Req)

	resp, err := httpc.Do(context.Background(), http.MethodPost, wxUrl, Req)
	if err != nil || resp.Status != "200 OK" {
		logx.Errorf("[%s]  GetWxQrCodeTicket httpc.Do post err, url[%s] appId[%s] %+v  resp:%+v", flowId, wxUrl, appId, err, resp)
		return
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("[%s] GetWxQrCodeTicket io.ReadAll err, url[%s] appId[%s] err:%+v ", flowId, wxUrl, appId, err)
		return
	}

	logx.Infof("[%s] GetWxQrCodeTicket appId[%s]  Body : %+v", flowId, appId, string(bytes))

	var m map[string]interface{}
	err = json.Unmarshal([]byte(bytes), &m)
	if err != nil {
		logx.Errorf("[%s] json.Unmarshal err , appId[%s] %+v  resp:%+v", flowId, appId, err, string(bytes))
		return
	}

	if errcode, ok := m["errcode"]; ok && errcode.(float64) != 0 {
		err = errors.New(string(bytes))
		logx.Errorf("[%s] GetWxQrCodeTicket: appId[%s]  m : %+v", flowId, appId, m)
		return
	}
	if ExpireSeconds, ok := m["expire_seconds"]; ok {
		ticketInfo.ExpireSeconds = int64(ExpireSeconds.(float64))
	}

	ticketInfo.Ticket = m["ticket"].(string)
	ticketInfo.URL = m["url"].(string)

	return ticketInfo, nil
}

// 请求公众号二维码
func (s *WxOfficialAccountMgr) GetWxQrCodeWithParameters(flowId, appId, ticket string) (qrCodeInfo QRCodeResponse, err error) {
	ticket, _ = url.QueryUnescape(ticket)
	wxUrl := fmt.Sprintf("https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=%s", ticket)

	resp, err := httpc.Do(context.Background(), http.MethodGet, wxUrl, nil)
	if err != nil || resp.Status != "200 OK" {
		logx.Errorf("[%s] GetWxQrCodeWithParameters httpc.Do get err, appId[%s] %+v  resp:%+v", flowId, appId, err, resp)
		return
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		logx.Errorf("[%s] GetWxQrCodeWithParameters io.ReadAll err , appId[%s] err:%+v ", flowId, appId, err)
		return
	}

	if json.Valid(bytes) { // 如果是json则说明已经有错误码
		var m map[string]interface{}
		err = json.Unmarshal([]byte(bytes), &m)
		logx.Errorf("[%s] GetWxQrCodeWithParameters json.Unmarshal err, appId[%s]  edata:%+v err:%+v ", flowId, appId, string(bytes), err)
		return
	}
	qrCodeInfo.AppId = appId
	encodeUrl := base64.StdEncoding.EncodeToString(bytes)
	if len(encodeUrl) == 0 {
		logx.Errorf("[%s] GetWxQrCodeWithParameters json.Unmarshal err, appId[%s]  edata:%+v err:%+v ", flowId, appId, string(bytes), err)
		return qrCodeInfo, errors.New("encodeUrl err")
	}
	qrCodeInfo.EncodedData = fmt.Sprintf("data:image/jpeg;base64,%s", encodeUrl)
	return qrCodeInfo, nil
}
