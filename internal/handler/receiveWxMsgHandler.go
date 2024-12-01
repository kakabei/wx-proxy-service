package handler

import (
	"encoding/xml"
	"io"
	"net/http"

	"wx-proxy-service/internal/logic"
	"wx-proxy-service/internal/svc"
	"wx-proxy-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ReceiveWxMsgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ReceiveWxMsgReq
		if err := httpx.Parse(r, &req); err != nil {
			logx.WithContext(r.Context()).Errorf("httpx.Parse url:%s err:%+v", r.RequestURI, err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logx.WithContext(r.Context()).Errorf("Failed to read request body: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		if err := xml.Unmarshal(body, &req); err != nil {
			logx.WithContext(r.Context()).Errorf("xml.Unmarshal err. request body: %v", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewReceiveWxMsgLogic(r.Context(), svcCtx)

		resp, err := l.ReceiveWxMsg(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			w.Write([]byte(resp.XmlData))
		}
	}
}
