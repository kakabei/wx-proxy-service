package svc

import (
	"wx-proxy-service/internal/config"
	"wx-proxy-service/internal/models/wx"
)

type ServiceContext struct {
	Config               config.Config
	WxAppIdConfigMgr     *wx.WxAppIdConfigMgr
	WxOfficailAccountMgr *wx.WxOfficialAccountMgr
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:               c,
		WxAppIdConfigMgr:     wx.NewWxAppConfigMgr(c),
		WxOfficailAccountMgr: wx.NewWxOfficialAccountMgr(c),
	}
}
