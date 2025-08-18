package main

import (
	"flag"
	"fmt"

	"shorturl/internal/config"
	"shorturl/internal/handler"
	"shorturl/internal/svc"
	"shorturl/pkg/base62"
	"shorturl/pkg/blacklist"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/shortener-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// 设置base62字符集
	base62.SetCharacterSet(c.Base62CharacterSet)
	// 设置黑名单
	blacklist.SetBlackListSet(c.BlackList)
	
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
