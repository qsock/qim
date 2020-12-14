package logic

import (
	"context"
	"encoding/json"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qim/lib/proto/event"
	"github.com/qsock/qim/lib/proto/msg"
	"github.com/qsock/qim/lib/proto/stream"
)

// 发送好友 的hint信息
func onFriendAdd(e *event.Friend) {
	req := new(msg.MsgReq)
	m := new(stream.MsgModel)
	m.ChatType = stream.ChatType_ChatTypeSingle
	m.SenderId = e.OperatorId
	m.RecvId = e.RecverId
	m.MsgType = stream.MsgType_MsgTypeHint
	m.Hint = &stream.HintMsg{T: stream.HintType_BecomeFriends, Content: "你们已经成为好友了，可以愉快的聊天了"}
	req.M = m
	_, err := Msg(context.Background(), req)
	if err != nil {
		qlog.Error(err, req)
	}
}

func onFriendDel(e *event.Friend) {
	req := new(msg.SysMsgReq)
	m := new(stream.SysMsgModel)
	m.SenderId = e.OperatorId
	m.RecverId = e.RecverId
	m.MsgType = stream.SysMsgType_Command
	m.Command = &stream.SysCommandMsg{T: stream.SysCommandType_SysCommandFriendDel,
		OperatorId: e.OperatorId, RecverId: e.RecverId}

	req.M = m
	_, err := SysMsg(context.Background(), req)
	if err != nil {
		qlog.Error(err, req)
	}
}

func onFriendMarkname(e *event.Friend) {
	req := new(msg.SysMsgReq)
	m := new(stream.SysMsgModel)
	m.SenderId = e.OperatorId
	m.RecverId = e.RecverId
	m.MsgType = stream.SysMsgType_Command
	mam := map[string]interface{}{"user_id": e.RecverId}
	b, _ := json.Marshal(mam)
	m.Command = &stream.SysCommandMsg{
		T:          stream.SysCommandType_SysCommandFriendMarkname,
		OperatorId: e.OperatorId, RecverId: e.OperatorId, Content: e.Content, Extra: string(b)}
	req.M = m
	_, err := SysMsg(context.Background(), req)
	if err != nil {
		qlog.Error(err, req)
	}
}

func onFriendApply(e *event.Friend) {
	req := new(msg.SysMsgReq)
	m := new(stream.SysMsgModel)
	m.SenderId = e.OperatorId
	m.RecverId = e.RecverId
	m.MsgType = stream.SysMsgType_Command
	m.Command = &stream.SysCommandMsg{
		T:          stream.SysCommandType_SysCommandFriendApply,
		OperatorId: e.OperatorId, RecverId: e.RecverId, Content: e.Content}
	req.M = m
	_, err := SysMsg(context.Background(), req)
	if err != nil {
		qlog.Error(err, req)
	}
}

func onFriendApplyReject(e *event.Friend) {
}
