package test

import (
	"wx-proxy-service/internal/config"
	"wx-proxy-service/internal/svc"
)

func createMockServiceContext() *svc.ServiceContext {
	c := config.Config{
		WxAppInfo: []config.WxAppConfig{
			{
				AppId:     "wx3cc8fd6963e31a32",
				Appsecret: "27d428d407******************",
			},
		},
	}
	return svc.NewServiceContext(c)
}
