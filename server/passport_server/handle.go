package main

import (
	"context"
	"github.com/qsock/qim/lib/proto/passport"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/server/passport_server/logic"
)

type Server struct{}

func (*Server) Ping(ctx context.Context, req *ret.NoArgs) (*ret.NoArgs, error) {
	return new(ret.NoArgs), nil
}

func (server *Server) Sms(ctx context.Context, req *passport.SmsReq) (*ret.EmptyResp, error) {
	return logic.Sms(ctx, req)
}

func (server *Server) TelLogin(ctx context.Context, req *passport.TelLoginReq) (*passport.LoginResp, error) {
	return logic.TelLogin(ctx, req)
}

func (server *Server) Qqlogin(ctx context.Context, req *passport.QqloginReq) (*passport.LoginResp, error) {
	return logic.Qqlogin(ctx, req)
}

func (server *Server) WxLogin(ctx context.Context, req *passport.WxLoginReq) (*passport.LoginResp, error) {
	return logic.WxLogin(ctx, req)
}

func (server *Server) Logout(ctx context.Context, req *passport.LogoutReq) (*ret.EmptyResp, error) {
	return logic.Logout(ctx, req)
}

func (server *Server) Auth(ctx context.Context, req *passport.AuthReq) (*ret.IntResp, error) {
	return logic.Auth(ctx, req)
}

func (server *Server) Refresh(ctx context.Context, req *passport.RefreshReq) (*ret.StringResp, error) {
	return logic.Refresh(ctx, req)
}

func (server *Server) Ban(ctx context.Context, req *passport.BanReq) (*ret.EmptyResp, error) {
	return logic.Ban(ctx, req)
}

func (server *Server) UnBan(ctx context.Context, req *passport.UnBanReq) (*ret.EmptyResp, error) {
	return logic.UnBan(ctx, req)
}
