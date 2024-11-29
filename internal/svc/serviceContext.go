package svc

import (
	"wx-proxy-service/internal/config"
	"wx-proxy-service/internal/models/wx"
)

type ServiceContext struct {
	Config           config.Config
	WxAppIdConfigMgr *wx.WxAppIdConfigMgr
	WxAppMgr         *wx.WxAppMgr
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:           c,
		WxAppIdConfigMgr: wx.NewWxAppConfigMgr(c),
		WxAppMgr:         wx.NewWxAppMgr(c),
	}
}
