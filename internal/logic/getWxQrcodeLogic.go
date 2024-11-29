package logic

import (
	"context"

	"wx-proxy-service/internal/common"
	"wx-proxy-service/internal/svc"
	"wx-proxy-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWxQrcodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWxQrcodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWxQrcodeLogic {
	return &GetWxQrcodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWxQrcodeLogic) GetWxQrcode(req *types.GetWxQrcodeReq) (resp *types.GetWxQrcodeResp, err error) {
	requestId := common.GetRequstId(l.ctx)
	resp = new(types.GetWxQrcodeResp)

	resp.Ret = types.CommonRet{Code: 0, Msg: "OK", RequestID: requestId}

	if req.AppId == "" || req.Source == "" {
		return nil, types.NewResultError(requestId, types.HttpCheckParamError)
	}

	AppSecret := l.svcCtx.WxAppIdConfigMgr.GetAppSecret(req.AppId)
	if len(AppSecret) == 0 {
		l.Logger.Errorf("[%s] AppSecret on in configration from[%s] err", requestId, req.Source)
		return nil, types.NewResultError(requestId, types.HttpAppSecretErr)
	}

	miniAppsAccessToken, err := l.svcCtx.WxAppMgr.GetMiniAppAccessToken(requestId, req.AppId, AppSecret)
	if err != nil {
		l.Logger.Errorf("[%s] wx.GetMiniAppAccessToken err : %+v", requestId, err)
		return nil, types.NewResultError(requestId, types.HttpGetAccessTokenErr)
	}

	codeInfo, err := l.svcCtx.WxAppMgr.GetUnlimitedQRCode(requestId, req.AppId, miniAppsAccessToken.AccessToken, req.Scene)
	if err != nil {
		l.Logger.Errorf("[%s] wx.GetWebAccessToken err : %+v", requestId, err)
		return nil, types.NewResultError(requestId, types.HttpGetUnlimitedQRCodeErr)
	}

	resp.Body.AppId = req.AppId
	resp.Body.QRBuffer = codeInfo.EncodedData
	return
}
