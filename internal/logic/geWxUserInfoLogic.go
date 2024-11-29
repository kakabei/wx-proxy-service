package logic

import (
	"context"

	"wx-proxy-service/internal/common"
	"wx-proxy-service/internal/models/wx"
	"wx-proxy-service/internal/svc"
	"wx-proxy-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GeWxUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGeWxUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GeWxUserInfoLogic {
	return &GeWxUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GeWxUserInfoLogic) GeWxUserInfo(req *types.GetWxUserInfoReq) (resp *types.GetWxUserInfoResp, err error) {
	requestId := common.GetRequstId(l.ctx)
	resp = new(types.GetWxUserInfoResp)
	resp.Ret = types.CommonRet{Code: 0, Msg: "OK", RequestID: requestId}

	if req.Source == "" || req.AppId == "" || req.Code == "" {
		return nil, types.NewResultError(requestId, types.HttpCheckParamError)
	}

	AppSecret := l.svcCtx.WxAppIdConfigMgr.GetAppSecret(req.AppId)
	if len(AppSecret) == 0 {
		l.Logger.Errorf("[%s] AppSecret on in configration from[%s] err", requestId, req.Source)
		return nil, types.NewResultError(requestId, types.HttpAppSecretErr)
	}
	webAccessToken, err := wx.GetWebAccessToken(requestId, req.AppId, AppSecret, req.Code)
	if err != nil {
		l.Logger.Errorf("[%s] GetWebAccessToken err : %+v", requestId, err)
		return nil, types.NewResultError(requestId, types.HttpGetAccessTokenErr)
	}

	l.Logger.Infof("[%s] wx.GetWebAccessToken: from[%s] webAccessToken:%s", requestId, req.Source, common.ToJSON(webAccessToken))

	wxUserInfo, err := wx.GetWxUserinfo(requestId, webAccessToken.OpenId, webAccessToken.AccessToken)
	if err != nil {
		l.Logger.Errorf("[%s] webAccessToken : %+v", requestId, webAccessToken)
		return nil, types.NewResultError(requestId, types.HttpWxUserInfoErr)
	}

	l.Logger.Infof("[%s] wx.GetWxUserinfo: from[%s]  wxUserInfo %s", requestId, req.Source, common.ToJSON(wxUserInfo))

	resp.Body.AppId = wxUserInfo.UnionID
	resp.Body.OpenId = wxUserInfo.OpenId
	resp.Body.Nickname = wxUserInfo.Nickname
	resp.Body.Sex = wxUserInfo.Sex
	resp.Body.HeadImgUrl = wxUserInfo.HeadImgUrl
	resp.Body.UnionID = wxUserInfo.UnionID
	resp.Body.Province = wxUserInfo.Province
	resp.Body.City = wxUserInfo.City
	resp.Body.Country = wxUserInfo.Country

	return
}
