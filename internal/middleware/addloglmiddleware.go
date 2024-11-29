package middleware

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
	"wx-proxy-service/internal/common"
	"wx-proxy-service/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

// Handle 中间件函数
func LogHandle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		requestId := common.CreatRequestId()
		r = r.Clone(context.WithValue(r.Context(), types.CtxValRequestId{}, requestId))

		startTime := time.Now()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logx.WithContext(r.Context()).Errorf("Failed to read request body: %v", err)
		}

		// 创建一个新的请求主体用于后续读取
		buf := bufferPool.Get().(*bytes.Buffer)
		buf.Reset()
		buf.Write(body)
		r.Body = ioutil.NopCloser(buf)

		defer func() {
			if buf != nil {
				buf.Reset()
				bufferPool.Put(buf)
			}
		}()

		logx.WithContext(r.Context()).Debugf("[%s][Request]: %s %s %+v %s", requestId, r.Method, r.RequestURI, r.Header, body)

		recorder := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			body:           make([]byte, 0),
		}
		next(recorder, r)

		logx.WithContext(r.Context()).Debugf("[%s][Response]: %s %s %s cost:%d", requestId, r.Method, r.RequestURI, string(recorder.body), time.Since(startTime).Milliseconds())

	}
}

// 自定义的 ResponseWriter
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

// WriteHeader 重写 WriteHeader 方法，捕获状态码
func (r *responseRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// 重写 Write 方法，捕获响应数据
func (r *responseRecorder) Write(body []byte) (int, error) {
	r.body = body
	return r.ResponseWriter.Write(body)
}
