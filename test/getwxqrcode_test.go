package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"wx-proxy-service/internal/types"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestGetWxQrcode(t *testing.T) {
	client := NewHTTPTestClient()

	tests := []struct {
		name           string
		request        types.GetWxQrcodeReq
		expectedStatus int
		checkResponse  func(*testing.T, *types.GetWxQrcodeResp)
	}{
		{
			name: "successful_request",
			request: types.GetWxQrcodeReq{
				Source:        "test",
				FlowId:        "test_flow_123",
				AppId:         "wx3cc8fd6963e31a32",
				Scene:         "test_scene",
				ExpireSeconds: 3600,
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp *types.GetWxQrcodeResp) {
				assert.Equal(t, 0, resp.Ret.Code, "Response code should be 0")
				assert.Equal(t, "OK", resp.Ret.Msg, "Response message should be OK")
				assert.NotEmpty(t, resp.Ret.RequestID, "RequestID should not be empty")
				assert.NotEmpty(t, resp.Body.Ticket, "Ticket should not be empty")
				assert.NotEmpty(t, resp.Body.QRBuffer, "QRBuffer should not be empty")
				assert.Equal(t, "wx3cc8fd6963e31a32", resp.Body.AppId, "AppId should match")
				assert.Equal(t, "test_flow_123", resp.Body.FlowId, "FlowId should match")
				assert.Equal(t, int32(3600), resp.Body.ExpireSeconds, "ExpireSeconds should match")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Log request details
			reqJson, _ := json.MarshalIndent(tt.request, "", "  ")
			logx.Infof("Test case: %s, Request: %s", tt.name, string(reqJson))

			resp, body := client.DoRequest(t, http.MethodPost, "/v1/service/wx/getwxqrcode", tt.request)
			assert.Equal(t, tt.expectedStatus, resp.StatusCode, fmt.Sprintf("HTTP status code should be %d", tt.expectedStatus))

			// Log response details
			logx.Infof("Response status: %d, body: %s", resp.StatusCode, string(body))

			var response types.GetWxQrcodeResp
			err := json.Unmarshal(body, &response)
			if err != nil {
				t.Fatalf("Failed to unmarshal response: %v\nResponse body: %s", err, string(body))
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, &response)
			}
		})
	}
}
