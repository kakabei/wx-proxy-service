package types

type ErrorResp struct {
	Ret  HTTPCommonHead `json:"ret"`
	Body interface{}    `json:"body"`
}

type CodeErrorResponse struct {
	Ret  HTTPCommonHead `json:"ret"`
	Body interface{}    `json:"body"`
}

type HTTPCommonHead struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

const defaultCode = 35840001

var (
	HttpSucc                   HTTPCommonHead = HTTPCommonHead{Code: 0, Msg: "OK"}
	HttpCheckParamError        HTTPCommonHead = HTTPCommonHead{Code: 3584002, Msg: "参数错误"}
	HttpJosnMarshalErr         HTTPCommonHead = HTTPCommonHead{Code: 3584003, Msg: "json序列化失败"}
	HttpGetAccessTokenErr      HTTPCommonHead = HTTPCommonHead{Code: 3584003, Msg: "查询token失败"}
	HttpAppSecretErr           HTTPCommonHead = HTTPCommonHead{Code: 3584003, Msg: "AppSecret错误"}
	HttpWxUserInfoErr          HTTPCommonHead = HTTPCommonHead{Code: 3584003, Msg: "查询用户信息"}
	HttpGetUnlimitedQRCodeErr  HTTPCommonHead = HTTPCommonHead{Code: 3584003, Msg: "查询小程序二维码失败"}
	HttpGetUserPphoneNumberErr HTTPCommonHead = HTTPCommonHead{Code: 3584003, Msg: "查询用户手机号失败"}
)

/////

func NewResultError(requestId string, ret HTTPCommonHead) error {
	ret.RequestID = requestId
	return &ret
}

func NewCodeError(requestId string, code int, msg string) error {
	return &HTTPCommonHead{RequestID: requestId, Code: code, Msg: msg}
}

func NewDefaulResultError(ret HTTPCommonHead) error {
	return &ret
}

func NewDefaultError(requestId string, msg string) error {
	return NewCodeError(requestId, defaultCode, msg)
}

func (e *HTTPCommonHead) Error() string {
	return e.Msg
}

func (e *HTTPCommonHead) Data() *CodeErrorResponse {
	return &CodeErrorResponse{
		Ret: *e,
	}
}
