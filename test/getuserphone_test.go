package test

import (
	"encoding/json"
	"net/http"
	"testing"
	"wx-proxy-service/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestGetUserPhone(t *testing.T) {
	client := NewHTTPTestClient()

	tests := []struct {
		name       string
		req        *types.GetUserPhoneReq
		wantStatus int
		wantErr    bool
	}{
		{
			name: "正常请求",
			req: &types.GetUserPhoneReq{
				Source: "test",
				AppId:  "wx3cc8fd6963e31a32",
				Code:   "061ePf1004WLgT1n9x000iUw1E3ePf1t",
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := client.DoRequest(t, "POST", "/v1/service/wx/getuserphone", tt.req)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			var result types.GetUserPhoneResp
			err := json.Unmarshal(body, &result)
			assert.NoError(t, err)

			if tt.wantErr {
				assert.NotEqual(t, 0, result.Ret.Code)
			} else {
				assert.Equal(t, 0, result.Ret.Code)
				assert.NotEmpty(t, result.Body.PhoneNumber)
			}
		})
	}
}
