package main

import (
	"flag"
	"fmt"

	"webtranshome/internal/config"
	"webtranshome/internal/handler"
	"webtranshome/internal/shared/response"
	"webtranshome/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
)

var configFile = flag.String("f", "etc/transhome-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	// Tambahkan rest.WithCors() agar frontend bisa menembak API ini
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	// Daftarkan global error handler untuk mengubah error menjadi format amplop
	httpx.SetErrorHandler(response.ErrorHandler)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
