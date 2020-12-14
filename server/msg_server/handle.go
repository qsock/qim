package main

import (
	"context"
	"github.com/qsock/qim/lib/proto/msg"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/server/msg_server/logic"
)

type Server struct{}

func (*Server) Ping(ctx context.Context, req *ret.NoArgs) (*ret.NoArgs, error) {
	return new(ret.NoArgs), nil
}

func (server *Server) Msg(ctx context.Context, req *msg.MsgReq) (*msg.MsgResp, error) {
	return logic.Msg(ctx, req)
}

func (server *Server) SessConnect(ctx context.Context, req *msg.SessConnectReq) (*ret.BytesResp, error) {
	return logic.SessConnect(ctx, req)
}

func (server *Server) SysMsg(ctx context.Context, req *msg.SysMsgReq) (*ret.IntResp, error) {
	return logic.SysMsg(ctx, req)
}

func (server *Server) Exchange(ctx context.Context, req *msg.ExchangeReq) (*ret.EmptyResp, error) {
	return logic.Exchange(ctx, req)
}

func (server *Server) CloseWithMsg(ctx context.Context, req *msg.CloseWithMsgReq) (*ret.EmptyResp, error) {
	return logic.CloseWithMsg(ctx, req)
}

func (server *Server) UserClosed(ctx context.Context, req *msg.UserClosedReq) (*ret.EmptyResp, error) {
	return logic.UserClosed(ctx, req)
}

func (server *Server) MarkChatRead(ctx context.Context, req *msg.MarkChatReadReq) (*ret.EmptyResp, error) {
	return logic.MarkChatRead(ctx, req)
}

func (server *Server) ChatAhead(ctx context.Context, req *msg.ChatAheadReq) (*ret.EmptyResp, error) {
	return logic.ChatAhead(ctx, req)
}

func (server *Server) ChatTouch(ctx context.Context, req *msg.ChatTouchReq) (*ret.EmptyResp, error) {
	return logic.ChatTouch(ctx, req)
}

func (server *Server) ChatRemove(ctx context.Context, req *msg.ChatRemoveReq) (*ret.EmptyResp, error) {
	return logic.ChatRemove(ctx, req)
}

func (server *Server) ChatClear(ctx context.Context, req *msg.ChatClearReq) (*ret.EmptyResp, error) {
	return logic.ChatClear(ctx, req)
}

func (server *Server) ChatIds(ctx context.Context, req *msg.ChatIdsReq) (*msg.ChatIdsResp, error) {
	return logic.ChatIds(ctx, req)
}

func (server *Server) ChatByIds(ctx context.Context, req *msg.ChatByIdsReq) (*msg.ChatByIdsResp, error) {
	return logic.ChatByIds(ctx, req)
}

func (server *Server) ChatByUids(ctx context.Context, req *msg.ChatByUidsReq) (*msg.ChatByUidsResp, error) {
	return logic.ChatByUids(ctx, req)
}

func (server *Server) ChatRecordIds(ctx context.Context, req *msg.ChatRecordIdsReq) (*msg.ChatRecordIdsResp, error) {
	return logic.ChatRecordIds(ctx, req)
}

func (server *Server) ChatRecordByIds(ctx context.Context, req *msg.ChatRecordByIdsReq) (*msg.ChatRecordByIdsResp, error) {
	return logic.ChatRecordByIds(ctx, req)
}

func (server *Server) ChatMute(ctx context.Context, req *msg.ChatMuteReq) (*ret.EmptyResp, error) {
	return logic.ChatMute(ctx, req)
}

func (server *Server) RevertSelfMsg(ctx context.Context, req *msg.RevertSelfMsgReq) (*ret.EmptyResp, error) {
	return logic.RevertSelfMsg(ctx, req)
}

func (server *Server) ManagerChatMsgRevert(ctx context.Context, req *msg.ManagerChatMsgRevertReq) (*ret.EmptyResp, error) {
	return logic.ManagerChatMsgRevert(ctx, req)
}

func (server *Server) GetSysMsg(ctx context.Context, req *msg.GetSysMsgReq) (*msg.GetSysMsgResp, error) {
	return logic.GetSysMsg(ctx, req)
}

func (server *Server) GetMemberIdByChatId(ctx context.Context, req *msg.GetMemberIdByChatIdReq) (*msg.GetMemberIdByChatIdResp, error) {
	return logic.GetMemberIdByChatId(ctx, req)
}
