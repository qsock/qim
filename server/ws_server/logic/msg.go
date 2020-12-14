package logic

import (
	"context"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qim/lib/codes"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/ws"
)

func Exchange(ctx context.Context, req *ws.ExchangeReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	oldSess := WsServer.GetSessionById(req.Uuid)
	if oldSess == nil {
		qlog.Ctx(ctx).Error(req, "not found")
		resp.Err = codes.Error(codes.ErrorSessNotExists)
		return resp, nil
	}
	WsServer.ReplaceSessionById(oldSess, req.SessId)
	return resp, nil
}

func CloseUser(ctx context.Context, req *ws.CloseUserReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	sess := WsServer.GetSessionById(req.SessId)
	if sess == nil {
		qlog.Ctx(ctx).Error(req, "not found")
		resp.Err = codes.Error(codes.ErrorSessNotExists)
		return resp, nil
	}
	if err := sess.CloseWithMsg(req.Content); err != nil {
		qlog.Ctx(ctx).Error(req, err, string(req.Content))
		return nil, err
	}
	return new(ret.EmptyResp), nil
}

func IsSessOnline(ctx context.Context, req *ws.IsSessOnlineReq) (*ret.BoolResp, error) {
	resp := new(ret.BoolResp)
	sess := WsServer.GetSessionById(req.SessId)
	if sess != nil && !sess.IsClosed() {
		resp.Flag = true
	}
	return resp, nil
}
func Msg(ctx context.Context, req *ws.MsgReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	sess := WsServer.GetSessionById(req.SessId)
	if sess == nil {
		qlog.Ctx(ctx).Error(req.SessId, "not found", string(req.Content))
		resp.Err = codes.Error(codes.ErrorSessNotExists)
		return resp, nil
	}
	if err := sess.Write(req.Content); err != nil {
		qlog.Ctx(ctx).Error(req.SessId, err, string(req.Content))
		return nil, err
	}
	return resp, nil
}

func AllMsg(ctx context.Context, req *ws.AllMsgReq) (*ret.EmptyResp, error) {
	if err := WsServer.Broadcast(req.Content); err != nil {
		qlog.Ctx(ctx).Error(err, string(req.Content))
		return nil, err
	}
	return &ret.EmptyResp{}, nil
}
