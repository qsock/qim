package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qf/store/ka"
	"github.com/qsock/qim/config/mq"
	e "github.com/qsock/qim/lib/proto/errmsg"
	"github.com/qsock/qim/lib/proto/event"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func LogMiddleware() grpc.ServerOption {
	return grpc.UnaryInterceptor(logInterceptor)
}

func logInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	ctx = grpcCtx2Ctx(ctx)

	defer func() {
		if e := recover(); e != nil {
			qlog.Get().Ctx(ctx).Fatal("Panic||err:%v||stack:%s", e, string(debug.Stack()))
		}
	}()

	// 执行
	resp, err := handler(ctx, req)

	// 忽略ping消息
	if strings.HasSuffix(info.FullMethod, "/Ping") {
		return resp, err
	}

	report := new(event.RpcReport)
	report.Method = info.FullMethod
	report.CreatedOn = start.UnixNano()
	report.EndOn = time.Now().UnixNano()
	report.Resp = JsonStr(resp)
	report.Req = JsonStr(req)
	if err != nil {
		report.Err = err.Error()
	}

	ka.TopicEvent(mq.TopicLogTrace, mq.ELogTraceGrpc, report)
	if err != nil {
		qlog.Get().Ctx(ctx).Error(report.Method, report.Req, report.Resp, err.Error())
		return resp, err
	}

	m := reflect.ValueOf(resp).MethodByName("GetErr")
	if m.IsValid() {
		var emptyArgs []reflect.Value
		v := m.Call(emptyArgs)
		if len(v) > 0 {
			if e, ok := v[0].Interface().(*e.ErrMsg); ok && e != nil {
				qlog.Get().Ctx(ctx).Error(report.Method, report.Req, report.Resp, e.GetMessage())
				return resp, err
			}
		}
	}

	// 打印只打印一部分，kafka要传全部的
	qlog.Get().Ctx(ctx).Info(report.Method, time.Since(start), Str(report.Req), Str(report.Resp))
	return resp, err
}

func grpcCtx2Ctx(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	v, ok := md[qgrpc.MetaKey]
	if !ok {
		return ctx
	}
	mam := map[string]string{}
	_ = json.Unmarshal([]byte(v[0]), &mam)
	return context.WithValue(ctx, qgrpc.MetaKey, mam)
}

func JsonStr(r interface{}) string {
	b, _ := json.Marshal(r)
	return string(b)
}

//每次只打印500个字
func Str(r interface{}) string {
	s := fmt.Sprintf("%v", r)
	idx := len(s)
	if len(s) > 512 {
		idx = 512
		return "Has" + strconv.Itoa(len(s)) + "More " + s[:idx]
	}
	return s
}
