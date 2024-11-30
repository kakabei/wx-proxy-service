package logic

import (
	"context"

	"wx-proxy-service/internal/common"
	"wx-proxy-service/internal/models/wx"
	"wx-proxy-service/internal/svc"
	"wx-proxy-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendWxTemplateMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSendWxTemplateMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendWxTemplateMsgLogic {
	return &SendWxTemplateMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SendWxTemplateMsgLogic) SendWxTemplateMsg(req *types.SendWxTemplateMsgReq) (resp *types.SendWxTemplateMsgResp, err error) {
	requestId := common.GetRequstId(l.ctx)
	resp = new(types.SendWxTemplateMsgResp)
	resp.Ret = types.CommonRet{Code: 0, Msg: "OK", RequestID: requestId}
	if req.Source == "" || req.AppId == "" || req.OpenId == "" || req.Data == "" {
		return nil, types.NewResultError(requestId, types.HttpCheckParamError)
	}

	AppSecret := l.svcCtx.WxAppIdConfigMgr.GetAppSecret(req.AppId)
	if len(AppSecret) == 0 {
		l.Logger.Errorf("[%s] AppSecret on in configration from[%s] err", requestId, req.Source)
		return nil, types.NewResultError(requestId, types.HttpAppSecretErr)
	}

	tokenInfo, err := l.svcCtx.WxOfficailAccountMgr.GetWxAccessToken(requestId, req.AppId, AppSecret)
	if err != nil {
		l.Logger.Errorf("[%s] wx.GetMiniAppAccessToken err : %+v", requestId, err)
		return nil, types.NewResultError(requestId, types.HttpGetAccessTokenErr)
	}

	err = wx.SendTemplateMessage(l.ctx, requestId, req.AppId, tokenInfo.AccessToken, req.Data)
	if err != nil {
		logx.Errorf("[%s] wx.SendTemplateMessage err : %+v", requestId, err)
		return nil, types.NewResultError(requestId, types.HttpSendTemplateMsgErr)
	}

	return
}
