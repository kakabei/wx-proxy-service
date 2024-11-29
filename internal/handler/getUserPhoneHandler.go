package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"wx-proxy-service/internal/logic"
	"wx-proxy-service/internal/svc"
	"wx-proxy-service/internal/types"
)

func GetUserPhoneHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserPhoneReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetUserPhoneLogic(r.Context(), svcCtx)
		resp, err := l.GetUserPhone(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
