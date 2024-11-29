package handler

import (
	"net/http"

	"wx-proxy-service/internal/logic"
	"wx-proxy-service/internal/svc"
	"wx-proxy-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GeWxUserInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetWxUserInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.WithContext(r.Context()).Errorf("httpx.Parse url:%s err:%+v", r.RequestURI, err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGeWxUserInfoLogic(r.Context(), svcCtx)
		resp, err := l.GeWxUserInfo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
