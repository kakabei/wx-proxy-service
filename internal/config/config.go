package config

import "github.com/zeromicro/go-zero/rest"

type RedisConfig struct {
	Host string
	Pass string
	Type string
}

type WxAppConfig struct {
	AppId     string
	Appsecret string
}

type WxOpenIdListConfig struct {
	EnvName       string
	AllowMsgEvent []string
	HandleUrl     string
	OpenIdList    []string
}

type WxMsgMgrConfig struct {
	AllowMsgEvent []string
	WxOpenIdList  []WxOpenIdListConfig
}

type Config struct {
	rest.RestConf
	Redis     RedisConfig
	WxAppInfo []WxAppConfig
	WxMsgMgr  WxMsgMgrConfig
}
