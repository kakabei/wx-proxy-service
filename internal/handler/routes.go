// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"wx-proxy-service/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/getwxuserinfo",
				Handler: GeWxUserInfoHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/getwxqrcode",
				Handler: GetWxQrcodeHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/getuserphone",
				Handler: GetUserPhoneHandler(serverCtx),
			},
		},
		rest.WithPrefix("/v1/service/wx"),
	)
}