package wx

import (
	"sync"
	"wx-proxy-service/internal/config"
)

type WxAppIdConfigMgr struct {
	mapWxAppConfig sync.Map
}

func NewWxAppConfigMgr(config config.Config) *WxAppIdConfigMgr {

	WxApp := new(WxAppIdConfigMgr)

	for _, AppIdInfo := range config.WxAppInfo {
		WxApp.mapWxAppConfig.Store(AppIdInfo.AppId, AppIdInfo.Appsecret)
	}

	return WxApp
}

func (s *WxAppIdConfigMgr) GetAppSecret(appId string) string {
	info, ok := s.mapWxAppConfig.Load(appId)
	if !ok {
		return ""
	}
	return info.(string)
}
