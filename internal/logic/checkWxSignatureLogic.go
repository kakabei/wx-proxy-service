package logic

import (
	"context"

	"wx-proxy-service/internal/common"
	"wx-proxy-service/internal/svc"
	"wx-proxy-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckWxSignatureLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckWxSignatureLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckWxSignatureLogic {
	return &CheckWxSignatureLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckWxSignatureLogic) CheckWxSignature(req *types.CheckWxSignatureReq) (resp *types.CheckWxSignatureResp, err error) {

	requestId := common.GetRequstId(l.ctx)
	resp = new(types.CheckWxSignatureResp)

	l.Logger.Infof("[%s] GetWxMsg req :%+v ", requestId, req)

	// todo 验证签名
	resp.Echostr = req.Echostr

	return
}
