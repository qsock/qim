package main

import (
	"flag"
	"fmt"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/gin-gonic/gin"
	"github.com/qsock/qf/qlog"
	_ "github.com/qsock/qf/qlog/qfilelog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qf/store/ka"
	"github.com/qsock/qim/api_gateway/config"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/ginproxy"
	"github.com/qsock/qim/lib/types"
	"net/http"
	"time"
)

var (
	configFile = flag.String("C", "", "config file")
	version    = flag.Bool("v", false, "version")
	buildtime  = "2018-01-01_00:00:00"
	githash    = "master"
)

func init() {
	flag.Parse()
}

// @title       IM平台
// @version     1.0
// @description IM平台
// @schemes     http
// @basePath /im
func main() {
	if *version {
		fmt.Println("buildtime:", buildtime)
		fmt.Println("githash:  ", githash)
		return
	}
	if err := config.Init(*configFile); err != nil {
		panic(err)
	}
	if err := qlog.OpenKv(config.GetConfig().LogType, config.GetConfig().Log); err != nil {
		panic(err)
	}
	defer qlog.Close()

	dao.SetEnv(config.GetEnv())

	// 初始化呢kafka服务
	if err := ka.Init(types.GetKaAddrs(config.GetEnv())); err != nil {
		panic(err)
	}

	if err := qgrpc.Init(config.GetOp()); err != nil {
		panic(err)
	}

	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	//切割路由,首先进行切割，保证这里面不会crash
	e.Use(ginproxy.Parse)
	// 拦截panic的中间件，无法拦截协程里面的panic，只有主协程panic
	e.Use(ginproxy.Recovery)
	e.Use(ginproxy.AccessLog)
	// 初始化路由
	SetRoute(e)

	//初始化http server
	s := &http.Server{
		Addr:              config.GetAddr(),
		Handler:           e,
		ReadTimeout:       60 * time.Second,
		ReadHeaderTimeout: 60 * time.Second,
		IdleTimeout:       300 * time.Second,
		WriteTimeout:      60 * time.Second,
	}

	// 优雅关闭
	if err := gracehttp.Serve(s); err != nil {
		qlog.Info("stop", err)
	}
}
