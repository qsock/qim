package logic

import (
	"context"
	"encoding/json"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qim/config/common"
	"github.com/qsock/qim/config/kvconfig"
	"github.com/qsock/qim/lib/cachename"
	"github.com/qsock/qim/lib/codes"
	"github.com/qsock/qim/lib/constdef"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/method"
	"github.com/qsock/qim/lib/proto/msg"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/stream"
	"github.com/qsock/qim/lib/proto/ws"
	"strconv"
)

func Exchange(ctx context.Context, req *msg.ExchangeReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	sessId := strconv.FormatInt(req.UserId, 10)

	cacheKey := cachename.RedisUserWs(req.UserId)
	oldKey := dao.GetKvConn(kvconfig.KvDefault).Get(cacheKey).String()
	// 关闭已经登陆的server
	if len(oldKey) > 0 {
		if _, err := CloseWithMsg(ctx, &msg.CloseWithMsgReq{UserId: req.UserId, Msg: common.AccountOfflineMsg}); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			return nil, err
		}
	}

	{
		creq := new(ws.ExchangeReq)
		creq.SessId = sessId
		creq.Uuid = req.Uuid
		cresp := new(ret.EmptyResp)
		if err := qgrpc.CallWithServerName(ctx, req.Key, method.WsExchange, creq, cresp); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			resp.Err = codes.Error(codes.ErrorMsgLogined)
			return resp, nil
		}
		if cresp.GetErr() != nil {
			qlog.Ctx(ctx).Error(req, cresp.GetErr())
			resp.Err = codes.Error(codes.ErrorMsgLogined)
			return resp, nil
		}
		if err := dao.GetKvConn(kvconfig.KvDefault).Set(cacheKey, req.Key, -1).Err(); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			resp.Err = codes.Error(codes.ErrorMsgLogined)
		}
	}
	return resp, nil
}

func CloseWithMsg(ctx context.Context, req *msg.CloseWithMsgReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	cacheKey := cachename.RedisUserWs(req.UserId)
	serverKey := dao.GetKvConn(kvconfig.KvDefault).Get(cacheKey).String()
	if len(serverKey) == 0 {
		qlog.Ctx(ctx).Error(req, serverKey)
		return resp, nil
	}

	sessId := strconv.FormatInt(req.UserId, 10)
	creq := new(ws.CloseUserReq)
	creq.SessId = sessId
	p := new(constdef.JsonRet)
	p.T = stream.StreamType_SysMsgS2C
	p.Data = &stream.SysMsgModel{Close: &stream.SysCloseMsg{Content: req.Msg, Pop: true}}
	creq.Content, _ = json.Marshal(p)
	cresp := new(ret.EmptyResp)
	if err := qgrpc.CallWithServerName(ctx, serverKey, method.WsCloseUser, creq, cresp); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	return resp, nil
}

func UserClosed(ctx context.Context, req *msg.UserClosedReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	userId, err := strconv.ParseInt(req.SessId, 10, 64)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}

	cacheKey := cachename.RedisUserWs(userId)
	if err := dao.GetKvConn(kvconfig.KvDefault).Del(cacheKey).Err(); err != nil {
		qlog.Ctx(ctx).Error(req, err, cacheKey)
		return nil, err
	}
	return resp, nil
}
