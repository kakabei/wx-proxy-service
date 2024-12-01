# 微信代理服务 (WeChat Proxy Service)

## 项目简介

微信代理服务是一个基于 go-zero 框架开发的微服务，用于处理微信公众号的各种接口调用。该服务封装了微信公众号的常用功能，提供统一的 HTTP 接口，简化了微信公众号的开发流程。

## 功能特性

- 获取微信二维码
  - 支持临时二维码生成
  - 支持永久二维码生成
- 微信用户信息获取
  - 获取用户基本信息
  - 获取用户手机号
- 模板消息
  - 支持发送微信模板消息
- 消息回调处理
  - 支持配置多环境回调处理
  - 支持自定义回调 URL
  - 支持事件过滤
  - 支持 OpenID 白名单

## 快速开始

### 环境要求

- Go 1.16 或以上版本
- go-zero 1.5.0 或以上版本 (推荐使用最新的稳定版本)
  - 安装命令：`go get -u github.com/zeromicro/go-zero`
  - 文档参考：[go-zero 官方文档](https://go-zero.dev/cn/)
- Redis
- 微信公众号账号

### 安装部署

1. 克隆项目

```bash
git clone [项目地址]
cd wx-proxy-service
```

2. 修改配置文件

```yaml
# etc/wx-proxy-api.yaml
Name: wx-proxy-api
Host: 0.0.0.0
Port: 36012

Redis:
  Host: 127.0.0.1:6379
  Type: node
  Pass: ""

WxAppInfo:
  - AppID: "your_app_id"
    AppSecret: "your_app_secret"

WxMsgMgr:
  AllowMsgEvent: ["SCAN","subscribe"]
  WxOpenIdList:
    - EnvName: "dev"
      AllowMsgEvent: ["SCAN","subscribe"]
      HandleUrl: http://your-callback-url
      OpenIdList: ["allowed_open_id"]
```

3. 运行服务

```bash
go run wx-proxy.go -f etc/wx-proxy-api.yaml
```

## API 接口说明

所有接口均以 `/v1/service/wx` 为基础路径。以下是各接口的详细说明和测试示例：

### 获取微信二维码

#### 临时二维码

```bash
curl --location 'http://127.0.0.1:36012/v1/service/wx/getwxqrcode' \
--header 'Content-Type: application/json' \
--data '{
    "source": "test",
    "flow_id": "test_123",
    "app_id": "wx3cc8fd6963e31a32",
    "scene": "test_scene",
    "expire_seconds": 3600
}'
```

#### 永久二维码

```bash
curl --location 'http://127.0.0.1:36012/v1/service/wx/getunlimitedqrcode' \
--header 'Content-Type: application/json' \
--data '{
    "source": "test",
    "flow_id": "test_123",
    "app_id": "wx3cc8fd6963e31a32",
    "scene": "test_scene"
}'
```

### 获取用户信息

```bash
curl --location 'http://127.0.0.1:36012/v1/service/wx/getwxuserinfo' \
--header 'Content-Type: application/json' \
--data '{
    "source": "test",
    "app_id": "wx3cc8fd6963e31a32",
    "code": "获取的网页授权code"
}'
```

### 获取用户手机号

```bash
curl --location 'http://127.0.0.1:36012/v1/service/wx/getuserphone' \
--header 'Content-Type: application/json' \
--data '{
    "source": "test",
    "app_id": "wx3cc8fd6963e31a32",
    "code": "获取手机号授权code"
}'
```

### 发送模板消息

```bash
curl --location 'http://127.0.0.1:36012/v1/service/wx/sendwxtemplatemsg' \
--header 'Content-Type: application/json' \
--data '{
    "source": "test",
    "app_id": "wx3cc8fd6963e31a32",
    "template_id": "模板ID",
    "touser": "用户OpenID",
    "data": {
        "first": {
            "value": "消息标题"
        },
        "keyword1": {
            "value": "内容1"
        },
        "keyword2": {
            "value": "内容2"
        },
        "remark": {
            "value": "备注信息"
        }
    }
}'
```

### 微信消息回调接口

#### 验证签名（配置回调时使用）

```bash
curl --location 'http://127.0.0.1:36012/v1/service/wx/wxmsg?signature=签名&timestamp=时间戳&nonce=随机串&echostr=随机串'
```

#### 接收消息通知

```bash
# 此接口由微信服务器调用，无需手动测试
# POST http://127.0.0.1:36012/v1/service/wx/wxmsg
```

### 响应示例

#### 成功响应

```json
{
    "ret": {
        "code": 0,
        "msg": "OK",
        "request_id": "wx_123456789"
    },
    "body": {
        // 接口相关的返回数据
    }
}
```

#### 错误响应

```json
{
    "ret": {
        "code": 3584002,
        "msg": "参数错误",
        "request_id": "wx_123456789"
    }
}
```

注意：

1. 所有请求都需要设置 `Content-Type: application/json` 头
2. 请确保 `app_id` 在配置文件中已正确配置
3. `source` 字段用于标识调用来源，建议使用有意义的标识
4. `flow_id` 用于跟踪请求，建议使用唯一标识

## 消息回调配置

服务支持配置多环境的消息回调处理：

```yaml
WxMsgMgr:
  AllowMsgEvent: ["SCAN","subscribe"]  # 全局允许的事件类型
  WxOpenIdList:
    - EnvName: "dev"                   # 环境名称
      AllowMsgEvent: ["SCAN"]          # 环境特定允许的事件
      HandleUrl: http://callback-url    # 回调地址
      OpenIdList: ["allowed_open_id"]   # 允许的OpenID列表
```

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 3584002 | 参数错误 |
| 3584004 | 获取 access_token 失败 |
| 3584008 | 获取用户信息失败 |

## 开发说明

### 项目结构

```
wx-proxy-service/
├── etc/                # 配置文件
├── internal/           # 内部实现
│   ├── config/        # 配置定义
│   ├── handler/       # 请求处理
│   ├── logic/         # 业务逻辑
│   ├── svc/           # 服务上下文
│   └── types/         # 类型定义
├── test/              # 测试用例
└── wx-proxy.go        # 主程序
```

### 测试

```bash
go test -v ./test/...
```

## 贡献指南

1. Fork 本仓库
2. 创建您的特性分支 (git checkout -b feature/AmazingFeature)
3. 提交您的更改 (git commit -m 'Add some AmazingFeature')
4. 推送到分支 (git push origin feature/AmazingFeature)
5. 打开一个 Pull Request

## 许可证

[MIT License](LICENSE)

## 联系方式

如有问题，请提交 Issue 或联系维护团队。
