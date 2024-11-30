package wx

const (
	WxAcccessTokenKey = "wx:appid:accesstoken"
)

// 小程序
type WxAccessToken struct {
	AccessToken string `json:"access_token"` // 小程序访问 access_token与基础支持的access_token不同
}

type GetUnlimitedQRCodeResp struct {
	AppId       string `json:"appid"`        //  小程序AppId
	EncodedData string `json:"encoded_data"` // 二维码	buffer
}
type UserPphoneNumberResp struct {
	AppId       string `json:"appid"`           // 小程序appid
	PhoneNumber string `json:"purePhoneNumber"` // 纯净手机号
}
