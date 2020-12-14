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
	"github.com/qsock/qim/lib/ginproxy"
	"github.com/qsock/qim/lib/metrics/rpc"
	"github.com/qsock/qim/lib/proto/ws"
	"github.com/qsock/qim/lib/types"
	"github.com/qsock/qim/server/ws_server/config"
	"google.golang.org/grpc"
	"net"
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

func main() {
	if *version {
		fmt.Println("buildtime:", buildtime)
		fmt.Println("githash:  ", githash)
		return
	}

	if err := config.Init(*configFile); err != nil {
		panic(err)
	}

	// 初始化日志
	if err := qlog.OpenKv(config.GetConfig().LogType, config.GetConfig().Log); err != nil {
		panic(err)
	}

	defer qlog.Close()

	// 初始化呢kafka服务
	if err := ka.Init(types.GetKaAddrs(config.GetEnv())); err != nil {
		panic(err)
	}
	// 启动grpc
	{
		ln, err := net.Listen("tcp", config.GetAddr())
		if err != nil {
			panic(err)
		}

		s := grpc.NewServer(rpc.LogMiddleware())
		ws.RegisterWsServer(s, new(Server))
		go s.Serve(ln)
		defer s.GracefulStop()
		if err := qgrpc.Init(config.GetOp()); err != nil {
			panic(err)
		}
	}

	{
		{
			e := gin.New()
			// 拦截panic的中间件，无法拦截协程里面的panic，只有主协程panic
			e.Use(ginproxy.Recovery)
			//切割路由
			e.Use(ginproxy.Parse)
			e.Use(ginproxy.AccessLog)
			// 初始化路由
			SetRoute(e)

			//初始化http server
			s := &http.Server{
				Addr:              fmt.Sprintf(":%d", config.GetWs().Port),
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

	}
}
