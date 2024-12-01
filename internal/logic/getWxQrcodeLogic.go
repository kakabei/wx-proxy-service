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
		l.Logger.Errorf("[%s] appid[%s] wx.GetWxAccessToken err : %+v", requestId, req.AppId, err)
		return nil, types.NewResultError(requestId, types.HttpGetAccessTokenErr)
	}

	actionName := "QR_STR_SCENE"
	ticketInfo, err := l.svcCtx.WxOfficailAccountMgr.GetWxQrCodeTicket(requestId, req.AppId, tokenInfo.AccessToken, actionName, req.Scene, int(req.ExpireSeconds))
	if err != nil {
		l.Logger.Errorf("[%s] appid[%s] wx.GetWxQrCodeTicket err : %+v", requestId, req.AppId, err)
		return nil, types.NewResultError(requestId, types.HttpGetUnlimitedQRCodeErr)
	}

	qrCodeInfo, err := l.svcCtx.WxOfficailAccountMgr.GetWxQrCodeWithParameters(requestId, req.AppId, ticketInfo.Ticket)
	if err != nil {
		l.Logger.Errorf("[%s] wx.GetWxQrCodeWithParameters err : %+v", requestId, err)
		return nil, types.NewResultError(requestId, types.HttpGetUnlimitedQRCodeErr)
	}

	resp.Body = types.GetWxQrcode{
		FlowId:        req.FlowId,
		AppId:         req.AppId,
		Ticket:        ticketInfo.Ticket,
		QRBuffer:      qrCodeInfo.EncodedData,
		ExpireSeconds: int32(ticketInfo.ExpireSeconds),
	}

	return resp, nil
}
