package main

import (
	"flag"
	"fmt"

	"lanshan11/internal/config"
	"lanshan11/internal/handler"
	"lanshan11/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf) // ✅ 用 RestConf 创建 REST 服务
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.RestConf.Host, c.RestConf.Port) // ✅ 从 RestConf 取 Host/Port
	server.Start()
}
