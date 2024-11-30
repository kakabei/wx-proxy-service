package test

import (
	"context"
	"testing"
	"wx-proxy-service/internal/logic"
	"wx-proxy-service/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestSendWxTemplateMsg(t *testing.T) {
	ctx := context.Background()
	svcCtx := createMockServiceContext()

	tests := []struct {
		name    string
		req     *types.SendWxTemplateMsgReq
		wantErr bool
	}{
		{
			name: "正常请求",
			req: &types.SendWxTemplateMsgReq{
				Source: "test",
				FlowId: "test_flow_id",
				AppId:  "wx3cc8fd6963e31a32",
				OpenId: "oFks96q7XErZCVZfbjAkFc__-8fc",
				Data:   `{"first":{"value":"测试消息"},"keyword1":{"value":"测试内容"},"remark":{"value":"测试备注"}}`,
			},
			wantErr: false,
		},
		{
			name: "参数缺失-AppId为空",
			req: &types.SendWxTemplateMsgReq{
				Source: "test",
				FlowId: "test_flow_id",
				OpenId: "oFks96q7XErZCVZfbjAkFc__-8fc",
				Data:   `{"first":{"value":"测试消息"}}`,
			},
			wantErr: true,
		},
		{
			name: "参数缺失-Source为空",
			req: &types.SendWxTemplateMsgReq{
				FlowId: "test_flow_id",
				AppId:  "wx3cc8fd6963e31a32",
				OpenId: "oFks96q7XErZCVZfbjAkFc__-8fc",
				Data:   `{"first":{"value":"测试消息"}}`,
			},
			wantErr: true,
		},
		{
			name: "参数缺失-OpenId为空",
			req: &types.SendWxTemplateMsgReq{
				Source: "test",
				FlowId: "test_flow_id",
				AppId:  "wx3cc8fd6963e31a32",
				Data:   `{"first":{"value":"测试消息"}}`,
			},
			wantErr: true,
		},
		{
			name: "参数缺失-Data为空",
			req: &types.SendWxTemplateMsgReq{
				Source: "test",
				FlowId: "test_flow_id",
				AppId:  "wx3cc8fd6963e31a32",
				OpenId: "test_openid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logic.NewSendWxTemplateMsgLogic(ctx, svcCtx)
			resp, err := l.SendWxTemplateMsg(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, resp)
		})
	}
}
