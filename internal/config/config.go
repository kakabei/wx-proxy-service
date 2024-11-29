package config

import "github.com/zeromicro/go-zero/rest"

type RedisCofnig struct {
	Host string
	Pass string
	Type string
}

type WxAppConfig struct {
	AppId     string
	Appsecret string
}

type Config struct {
	rest.RestConf
	Redis     RedisCofnig
	WxAppInfo []WxAppConfig
}
