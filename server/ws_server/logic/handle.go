package logic

import (
	"context"
	"encoding/json"
	"github.com/qsock/qf/net/ws"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qim/lib/constdef"
	"github.com/qsock/qim/lib/method"
	"github.com/qsock/qim/lib/proto/msg"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/stream"
)

func OnConnect(s *ws.Session) {
	req := new(msg.SessConnectReq)
	req.SessId = s.GetId()
	req.ServerKey = GetRegistKey()
	resp := new(ret.BytesResp)
	if err := qgrpc.Call(context.Background(), method.MsgSessConnect, req, resp); err != nil {
		qlog.Error(req, err, s.GetId())
		_ = s.Close()
		return
	}
	_ = s.Write(resp.Val)
}

func OnDisConnect(s *ws.Session) {
	qlog.Info("closed", s.GetId(), s.Keys)
}

// 收到消息
func OnMessage(s *ws.Session, b []byte) {
	// 收到json消息
	m := new(constdef.JsonRet)
	if err := json.Unmarshal(b, m); err != nil {
		qlog.Error(s.GetId(), string(b), s.Keys)
		_ = s.Close()
		return
	}
	switch m.T {
	case stream.StreamType_Ping:
		{
			OnPing(s)
		}
	}
	qlog.Info("closed", s.GetId(), s.Keys, string(b))
}

func OnSent(s *ws.Session, b []byte) {
	qlog.Info("sent", s.GetId(), s.Keys, string(b))
}

func OnPing(s *ws.Session) {
	b, _ := json.Marshal(&constdef.JsonRet{T: stream.StreamType_Pong})
	_ = s.Write(b)
}

func OnError(s *ws.Session, err error) {
	qlog.Error(s.GetId(), s.Keys, err)
	_ = s.Close()
}

func OnClose(s *ws.Session, t int, str string) {
	qlog.Info("closed", s.GetId(), s.Keys, t, str)
	req := new(msg.UserClosedReq)
	req.SessId = s.GetId()
	resp := new(ret.EmptyResp)
	if err := qgrpc.Call(context.Background(), method.MsgUserClosed, req, resp); err != nil {
		qlog.Error(req, err, s.GetId())
	}
}
