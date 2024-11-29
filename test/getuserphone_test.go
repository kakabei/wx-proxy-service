package test

import (
	"context"
	"testing"
	"wx-proxy-service/internal/logic"
	"wx-proxy-service/internal/types"

	"github.com/stretchr/testify/assert"
)

func TestGetUserPhone(t *testing.T) {
	ctx := context.Background()
	svcCtx := createMockServiceContext()

	tests := []struct {
		name    string
		req     *types.GetUserPhoneReq
		wantErr bool
	}{
		{
			name: "正常请求",
			req: &types.GetUserPhoneReq{
				AppId:  "wx3cc8fd6963e31a32",
				Source: "test",
				Code:   "0611AjGa1GqOBI0xgMHa1ewUsg11AjG7",
			},
			wantErr: false,
		},
		{
			name: "参数缺失-AppId为空",
			req: &types.GetUserPhoneReq{
				Source: "test",
				Code:   "0611AjGa1GqOBI0xgMHa1ewUsg11AjG7",
			},
			wantErr: true,
		},
		{
			name: "参数缺失-Source为空",
			req: &types.GetUserPhoneReq{
				AppId: "wx3cc8fd6963e31a32",
				Code:  "0611AjGa1GqOBI0xgMHa1ewUsg11AjG7",
			},
			wantErr: true,
		},
		{
			name: "参数缺失-Code为空",
			req: &types.GetUserPhoneReq{
				AppId:  "wx3cc8fd6963e31a32",
				Source: "test",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logic.NewGetUserPhoneLogic(ctx, svcCtx)
			resp, err := l.GetUserPhone(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, resp)
		})
	}
}
