package logic

import (
	"context"

	"wx-proxy-service/internal/common"
	"wx-proxy-service/internal/svc"
	"wx-proxy-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReceiveWxMsgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewReceiveWxMsgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceiveWxMsgLogic {
	return &ReceiveWxMsgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ReceiveWxMsgLogic) ReceiveWxMsg(req *types.ReceiveWxMsgReq) (resp *types.ReceiveWxMsgResp, err error) {

	requestId := common.GetRequstId(l.ctx)
	resp = &types.ReceiveWxMsgResp{
		XmlData: string("success"),
	}

	if !common.StringInArray(req.Event, l.svcCtx.Config.WxMsgMgr.AllowMsgEvent) {
		l.Logger.Debugf("[%s]  event[%s]  on InArray:%s ", requestId, req.Event, common.ToJSON(l.svcCtx.Config.WxMsgMgr.AllowMsgEvent))
		return
	}

	openIdConf, err := common.GeWxOpenIdConfig(l.svcCtx.Config, req.FromUserName)
	if err != nil {
		l.Logger.Errorf("[%s] GeWxOpenIdConfig err:%+v", requestId, err)
		return
	}

	l.Logger.Infof("[%s] GeWxOpenIdConfig:%s", requestId, common.ToJSON(openIdConf))

	if !common.StringInArray(req.Event, openIdConf.AllowMsgEvent) {
		l.Logger.Errorf("[%s] Openid[%s] event[%s] on InArray:%s", requestId, req.FromUserName, req.Event, common.ToJSON(openIdConf))
		return
	}

	// 转发微信消息
	grsp, err := common.PosWxMsg(l.ctx, requestId, req, openIdConf.HandleUrl)
	if err != nil {
		l.Logger.Errorf("[%s] PosWxMsg error OpenId[%s] error:%s", requestId, req.FromUserName, err)
		return
	}

	resp.XmlData = grsp.String()
	return
}
