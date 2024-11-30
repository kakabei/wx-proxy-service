package logic

import (
	"context"

	"wx-proxy-service/internal/common"
	"wx-proxy-service/internal/models/wx"
	"wx-proxy-service/internal/svc"
	"wx-proxy-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserPhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserPhoneLogic {
	return &GetUserPhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserPhoneLogic) GetUserPhone(req *types.GetUserPhoneReq) (resp *types.GetUserPhoneResp, err error) {

	RequestId := common.CreatRequestId()
	resp = new(types.GetUserPhoneResp)
	resp.Ret = types.CommonRet{Code: 0, Msg: "OK", RequestID: RequestId}

	if req.AppId == "" || req.Source == "" || req.Code == "" {
		return nil, types.NewResultError(RequestId, types.HttpCheckParamError)
	}

	AppSecret := l.svcCtx.WxAppIdConfigMgr.GetAppSecret(req.AppId)
	if len(AppSecret) == 0 {
		logx.Errorf("[%s] AppSecret no in configration from[%s] err", RequestId, req.Source)
		return nil, types.NewResultError(RequestId, types.HttpAppSecretErr)
	}

	tokenInfo, err := l.svcCtx.WxOfficailAccountMgr.GetWxAccessToken(RequestId, req.AppId, AppSecret)
	if err != nil {
		logx.Errorf("[%s] wx.GetWxAccessToken err. appid[%s] err : %+v", RequestId, req.AppId, err)
		return nil, types.NewResultError(RequestId, types.HttpGetAccessTokenErr)
	}
	userResp, err := wx.GetUserPhoneNumber(RequestId, req.AppId, tokenInfo.AccessToken, req.Code)
	if err != nil {
		logx.Errorf("[%s] wx.GetUserPhoneNumber err. appid[%s] err : %+v", RequestId, req.AppId, err)
		return nil, types.NewResultError(RequestId, types.HttpGetUserPphoneNumberErr)
	}

	resp.Body.AppId = req.AppId
	resp.Body.PhoneNumber = userResp.PhoneNumber
	return
}
