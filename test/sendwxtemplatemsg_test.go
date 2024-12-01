package test

import (
	"encoding/json"
	"net/http"
	"testing"
	"wx-proxy-service/internal/types"

	"github.com/stretchr/testify/assert"
)

type RecycleInstData struct {
	Time4 struct {
		Value string `json:"value"`
	} `json:"time4"`
	CharacterString2 struct {
		Value string `json:"value"`
	} `json:"character_string2"`
	CharacterString3 struct {
		Value string `json:"value"`
	} `json:"character_string3"`
	CharacterString5 struct {
		Value string `json:"value"`
	} `json:"character_string5"`
	CharacterString6 struct {
		Value string `json:"value"`
	} `json:"character_string6"`
}

type RecycleInst struct {
	Touser      string          `json:"touser"`
	TemplateID  string          `json:"template_id"`
	URL         string          `json:"url"`
	ClientMsgID string          `json:"client_msg_id"`
	Data        RecycleInstData `json:"data"`
}

func TestSendWxTemplateMsg(t *testing.T) {
	client := NewHTTPTestClient()

	// 构造模板消息数据
	templateData := RecycleInst{
		Touser:      "oFks96q7XErZCVZfbjAkFc__-8fc",
		TemplateID:  "392SnkFA1Gxz4hVKbRyqfT2Z_6ZdYYNqBg3H-KCynZM",
		URL:         "https://example.com",
		ClientMsgID: "test_msg_id",
		Data: RecycleInstData{
			Time4: struct {
				Value string `json:"value"`
			}{
				Value: "2024-12-01 13:33:00",
			},
			CharacterString2: struct {
				Value string `json:"value"`
			}{
				Value: "1000125",
			},
			CharacterString3: struct {
				Value string `json:"value"`
			}{
				Value: "卡卡西",
			},
			CharacterString5: struct {
				Value string `json:"value"`
			}{
				Value: "202412011234043-kesrelh",
			},
			CharacterString6: struct {
				Value string `json:"value"`
			}{
				Value: "测试数据",
			},
		},
	}

	dataBytes, err := json.Marshal(templateData)
	assert.NoError(t, err)

	tests := []struct {
		name       string
		req        *types.SendWxTemplateMsgReq
		wantStatus int
		wantErr    bool
	}{
		{
			name: "正常请求",
			req: &types.SendWxTemplateMsgReq{
				Source: "test",
				FlowId: "test_flow_id",
				AppId:  "wx3cc8fd6963e31a32",
				OpenId: "oFks96q7XErZCVZfbjAkFc__-8fc",
				Data:   string(dataBytes),
			},
			wantStatus: http.StatusOK,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, body := client.DoRequest(t, "POST", "/v1/service/wx/sendwxtemplatemsg", tt.req)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			var result types.SendWxTemplateMsgResp
			err := json.Unmarshal(body, &result)
			assert.NoError(t, err)

			if tt.wantErr {
				assert.NotEqual(t, 0, result.Ret.Code)
			} else {
				assert.Equal(t, 0, result.Ret.Code)
			}
		})
	}
}
