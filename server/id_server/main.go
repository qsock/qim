package main

import (
	"flag"
	"fmt"
	"github.com/qsock/qf/qlog"
	_ "github.com/qsock/qf/qlog/qfilelog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qf/store/ka"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/metrics/rpc"
	"github.com/qsock/qim/lib/proto/id"
	"github.com/qsock/qim/lib/types"
	"github.com/qsock/qim/server/id_server/config"
	"github.com/qsock/qim/server/id_server/logic"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
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
	var err error
	if err := config.Init(*configFile); err != nil {
		panic(err)
	}

	// 初始化日志
	if err := qlog.OpenKv(config.GetConfig().LogType, config.GetConfig().Log); err != nil {
		panic(err)
	}
	defer qlog.Close()

	dao.SetEnv(config.GetEnv())
	logic.Init()
	// 初始化呢kafka服务
	if err := ka.Init(types.GetKaAddrs(config.GetEnv())); err != nil {
		panic(err)
	}

	ln, err := net.Listen("tcp", config.GetAddr())
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(rpc.LogMiddleware())
	id.RegisterIdServer(s, new(Server))
	go s.Serve(ln)
	defer s.GracefulStop()
	if err := qgrpc.Init(config.GetOp()); err != nil {
		panic(err)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGTERM, syscall.SIGINT)
	<-exit
}
