package test

import (
	"context"
	"testing"
	"wx-proxy-service/internal/logic"
	"wx-proxy-service/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestGetWxQrcode(t *testing.T) {
	ctx := context.Background()
	svcCtx := createMockServiceContext()

	tests := []struct {
		name    string
		req     *types.GetWxQrcodeReq
		wantErr bool
	}{
		{
			name: "正常请求",
			req: &types.GetWxQrcodeReq{
				AppId:  "wx3cc8fd6963e31a32",
				Source: "test",
				Scene:  "test_scene",
			},
			wantErr: false,
		},
		{
			name: "参数缺失-AppId为空",
			req: &types.GetWxQrcodeReq{
				Source: "test",
				Scene:  "test_scene",
			},
			wantErr: true,
		},
		{
			name: "参数缺失-Source为空",
			req: &types.GetWxQrcodeReq{
				AppId: "wx3cc8fd6963e31a32",
				Scene: "test_scene",
			},
			wantErr: true,
		},
		{
			name: "场景值为空-应该成功",
			req: &types.GetWxQrcodeReq{
				AppId:  "wx3cc8fd6963e31a32",
				Source: "test",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logic.NewGetWxQrcodeLogic(ctx, svcCtx)
			resp, err := l.GetWxQrcode(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, resp)
		})
	}
}
