package wx

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpc"
)

// 查询用户手机号
func GetUserPhoneNumber(flowId string, appId, accessToken string, code string) (userResp UserPphoneNumberResp, err error) {

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
func GetUnlimitedQRCode(flowId string, appId, accessToken string, scene string) (QRResp GetUnlimitedQRCodeResp, err error) {

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
