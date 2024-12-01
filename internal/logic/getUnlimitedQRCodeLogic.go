package logic

import (
	"context"

	"wx-proxy-service/internal/common"
	"wx-proxy-service/internal/models/wx"
	"wx-proxy-service/internal/svc"
	"wx-proxy-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUnlimitedQRCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUnlimitedQRCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUnlimitedQRCodeLogic {
	return &GetUnlimitedQRCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUnlimitedQRCodeLogic) GetUnlimitedQRCode(req *types.GetUnlimitedQRCodeReq) (resp *types.GetUnlimitedQRCodeResp, err error) {
	requestId := common.GetRequstId(l.ctx)
	resp = new(types.GetUnlimitedQRCodeResp)
	resp.Ret = types.CommonRet{Code: 0, Msg: "OK", RequestID: requestId}

	if req.Source == "" || req.AppId == "" || req.Scene == "" {
		return nil, types.NewResultError(requestId, types.HttpCheckParamError)
	}

	appSecret := l.svcCtx.WxAppIdConfigMgr.GetAppSecret(req.AppId)
	if len(appSecret) == 0 {
		l.Logger.Errorf("[%s] AppSecret on in configration from[%s] err", requestId, req.Source)
		return nil, types.NewResultError(requestId, types.HttpAppSecretErr)
	}

	tokenInfo, err := l.svcCtx.WxOfficailAccountMgr.GetWxAccessToken(requestId, req.AppId, appSecret)
	if err != nil {
		l.Logger.Errorf("[%s] wx.GetWxAccessToken err : %+v", requestId, err)
		return nil, types.NewResultError(requestId, types.HttpGetAccessTokenErr)
	}

	codeInfo, err := wx.GetUnlimitedQRCode(requestId, req.AppId, tokenInfo.AccessToken, req.Scene)
	if err != nil {
		l.Logger.Errorf("[%s] wx.GetUnlimitedQRCode err : %+v", requestId, err)
		return nil, types.NewResultError(requestId, types.HttpGetUnlimitedQRCodeErr)
	}

	resp.Body = types.GetUnlimitedQRCode{
		FlowId:   req.FlowId,
		AppId:    req.AppId,
		QRBuffer: codeInfo.EncodedData,
	}

	return resp, nil
}
