package main

import (
	"flag"
	"fmt"
	"net/http"

	"wx-proxy-service/internal/config"
	"wx-proxy-service/internal/handler"
	"wx-proxy-service/internal/svc"
	"wx-proxy-service/internal/types"

	"wx-proxy-service/internal/middleware"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/wx-proxy-service-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)

	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		switch e := err.(type) {
		case *types.HTTPCommonHead:
			return http.StatusOK, e.Data()
		default:
			return http.StatusInternalServerError, nil
		}
	})

	server.Use(
		middleware.LogHandle,
	)

	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
