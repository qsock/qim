package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/qsock/qim/lib/ginproxy"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/ws"
	"github.com/qsock/qim/server/ws_server/logic"
)

type Server struct{}

func (*Server) Ping(ctx context.Context, req *ret.NoArgs) (*ret.NoArgs, error) {
	return new(ret.NoArgs), nil
}

func (s *Server) IsSessOnline(ctx context.Context, req *ws.IsSessOnlineReq) (*ret.BoolResp, error) {
	return logic.IsSessOnline(ctx, req)
}

func (s *Server) Msg(ctx context.Context, req *ws.MsgReq) (*ret.EmptyResp, error) {
	return logic.Msg(ctx, req)
}

func (s *Server) AllMsg(ctx context.Context, req *ws.AllMsgReq) (*ret.EmptyResp, error) {
	return logic.AllMsg(ctx, req)
}

func (s *Server) CloseUser(ctx context.Context, req *ws.CloseUserReq) (*ret.EmptyResp, error) {
	return logic.CloseUser(ctx, req)
}

func (s *Server) Exchange(ctx context.Context, req *ws.ExchangeReq) (*ret.EmptyResp, error) {
	return logic.Exchange(ctx, req)
}

func SetRoute(e *gin.Engine) {
	e.Any("/ping", ginproxy.OK)
	e.Any("/ls", logic.HandleWs)
}
