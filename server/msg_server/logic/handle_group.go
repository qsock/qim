package logic

import (
	"context"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qim/lib/gclient"
	"github.com/qsock/qim/lib/proto/event"
	"github.com/qsock/qim/lib/proto/msg"
	"github.com/qsock/qim/lib/proto/stream"
	"time"
)

func onGroupApply(e *event.GroupNewApply) {
	req := new(msg.SysMsgReq)
	m := new(stream.SysMsgModel)
	m.SenderId = e.OperatorId
	m.RecverId = e.RecverId
	m.MsgType = stream.SysMsgType_Command
	m.Command = &stream.SysCommandMsg{
		T:          stream.SysCommandType_SysCommandGroupApply,
		OperatorId: e.OperatorId, RecverId: e.GroupId, Content: e.Content}

	req.M = m
	_, err := SysMsg(context.Background(), req)
	if err != nil {
		qlog.Error(err, req)
	}
}

func onGroupManager(e *event.GroupManager) {
	req := new(msg.MsgReq)
	m := new(stream.MsgModel)
	groupInfo, err := GetGroupInfo(context.TODO(), e.GroupId)
	if err != nil {
		qlog.Error(e, err)
		return
	}
	m.ChatType = groupInfo.T
	m.SenderId = e.UserId
	m.RecvId = e.GroupId
	m.MsgType = stream.MsgType_MsgTypeHint
	var t = stream.HintType_GroupCharger
	if !e.IsManager {
		t = stream.HintType_GroupChargerCancel
	}
	lableowner, err := gclient.FullUserNameLabel(context.Background(), e.UserId)
	if err != nil {
		qlog.Error(err, e)
		return
	}
	lableManaer, err := gclient.FullUserNameLabel(context.Background(), e.UserId)
	if err != nil {
		qlog.Error(err, e)
		return
	}
	m.Hint = &stream.HintMsg{T: t, Content: lableowner + "认命了" + lableManaer + "为管理员"}
	req.M = m
	if _, err := Msg(context.Background(), req); err != nil {
		qlog.Error(err, req)
	}
}

func onGroupNewMember(e *event.GroupMember) {
	req := new(msg.MsgReq)
	m := new(stream.MsgModel)
	groupInfo, err := GetGroupInfo(context.TODO(), e.GroupId)
	if err != nil {
		qlog.Error(e, err)
		return
	}
	m.ChatType = groupInfo.T
	m.SenderId = e.OperatorId
	m.RecvId = e.GroupId
	m.MsgType = stream.MsgType_MsgTypeHint
	var label string
	if e.OperatorId > 0 {
		lstr, err := gclient.FullUserNameLabel(context.Background(), e.OperatorId)
		if err != nil {
			qlog.Error(err, e)
			return
		}
		label += lstr + "邀请"
	}
	lstr, err := gclient.FullUserNameLabel(context.Background(), e.UserId)
	if err != nil {
		qlog.Error(err, e)
		return
	}
	label += lstr + "加入了群组"
	m.Hint = &stream.HintMsg{T: stream.HintType_GroupJoin, Content: label}
	req.M = m
	if _, err := Msg(context.Background(), req); err != nil {
		qlog.Error(err, req)
	}
}

func onGroupDelMember(e *event.GroupMember) {
	req := new(msg.MsgReq)
	m := new(stream.MsgModel)
	groupInfo, err := GetGroupInfo(context.TODO(), e.GroupId)
	if err != nil {
		qlog.Error(e, err)
		return
	}
	m.ChatType = groupInfo.T
	m.SenderId = e.OperatorId
	m.RecvId = e.GroupId
	m.MsgType = stream.MsgType_MsgTypeHint
	var label string
	lstr, err := gclient.FullUserNameLabel(context.Background(), e.OperatorId)
	if err != nil {
		qlog.Error(err, e)
		return
	}
	label += lstr + "将"
	lstr, err = gclient.FullUserNameLabel(context.Background(), e.UserId)
	if err != nil {
		qlog.Error(err, e)
		return
	}
	label += lstr + "踢出了群组"
	m.Hint = &stream.HintMsg{T: stream.HintType_GroupLeave, Content: label}
	req.M = m
	if _, err := Msg(context.Background(), req); err != nil {
		qlog.Error(err, req)
	}
}

func onGroupLeaveMember(e *event.GroupMember) {
	req := new(msg.MsgReq)
	m := new(stream.MsgModel)
	groupInfo, err := GetGroupInfo(context.TODO(), e.GroupId)
	if err != nil {
		qlog.Error(e, err)
		return
	}
	m.ChatType = groupInfo.T
	m.SenderId = e.OperatorId
	m.RecvId = e.GroupId
	m.MsgType = stream.MsgType_MsgTypeHint
	var label string
	lstr, err := gclient.FullUserNameLabel(context.Background(), e.UserId)
	if err != nil {
		qlog.Error(err, e)
		return
	}
	label += lstr + "离开了群组"
	m.Hint = &stream.HintMsg{T: stream.HintType_GroupLeave, Content: label}
	req.M = m
	if _, err := Msg(context.Background(), req); err != nil {
		qlog.Error(err, req)
	}
}

func onGroupUpdateName(e *event.GroupUpdate) {
	req := new(msg.MsgReq)
	m := new(stream.MsgModel)
	groupInfo, err := GetGroupInfo(context.TODO(), e.GroupId)
	if err != nil {
		qlog.Error(e, err)
		return
	}
	m.ChatType = groupInfo.T
	m.SenderId = e.UserId
	m.RecvId = e.GroupId
	m.MsgType = stream.MsgType_MsgTypeHint
	var label string
	lstr, err := gclient.FullUserNameLabel(context.Background(), e.UserId)
	if err != nil {
		qlog.Error(err, e)
		return
	}
	label += lstr + "修改了群名称为:" + e.Str
	m.Hint = &stream.HintMsg{T: stream.HintType_GroupUpdateName, Content: label}
	req.M = m
	if _, err := Msg(context.Background(), req); err != nil {
		qlog.Error(err, req)
	}
}

func onGroupUpdateNotice(e *event.GroupUpdate) {
	req := new(msg.MsgReq)
	m := new(stream.MsgModel)
	groupInfo, err := GetGroupInfo(context.TODO(), e.GroupId)
	if err != nil {
		qlog.Error(e, err)
		return
	}
	m.ChatType = groupInfo.T
	m.SenderId = e.UserId
	m.RecvId = e.GroupId
	m.MsgType = stream.MsgType_MsgTypeHint
	var label string
	lstr, err := gclient.FullUserNameLabel(context.Background(), e.UserId)
	if err != nil {
		qlog.Error(err, e)
		return
	}
	label += lstr + "修改了群通知为:" + e.Str
	m.Hint = &stream.HintMsg{T: stream.HintType_GroupUpdateNotice, Content: label}
	req.M = m
	if _, err := Msg(context.Background(), req); err != nil {
		qlog.Error(err, req)
	}
}

func onGroupUpdateAvatar(e *event.GroupUpdate) {
	req := new(msg.MsgReq)
	m := new(stream.MsgModel)
	groupInfo, err := GetGroupInfo(context.TODO(), e.GroupId)
	if err != nil {
		qlog.Error(e, err)
		return
	}
	m.ChatType = groupInfo.T
	m.SenderId = e.UserId
	m.RecvId = e.GroupId
	m.MsgType = stream.MsgType_MsgTypeHint
	var label string
	lstr, err := gclient.FullUserNameLabel(context.Background(), e.UserId)
	if err != nil {
		qlog.Error(err, e)
		return
	}
	label += lstr + "修改了群头像"
	m.Hint = &stream.HintMsg{T: stream.HintType_GroupUpdateAvatar, Content: label}
	req.M = m
	if _, err := Msg(context.Background(), req); err != nil {
		qlog.Error(err, req)
	}
}

func onGroupUpdateMute(e *event.GroupUpdate) {
	req := new(msg.MsgReq)
	m := new(stream.MsgModel)
	groupInfo, err := GetGroupInfo(context.TODO(), e.GroupId)
	if err != nil {
		qlog.Error(e, err)
		return
	}
	m.ChatType = groupInfo.T
	m.SenderId = e.UserId
	m.RecvId = e.GroupId
	m.MsgType = stream.MsgType_MsgTypeHint
	var label string
	lstr, err := gclient.FullUserNameLabel(context.Background(), e.UserId)
	if err != nil {
		qlog.Error(err, e)
		return
	}
	if e.Time == 0 {
		label += lstr + "结束了全员禁言"
	} else {
		label += lstr + "开启了全员禁言，直到" + time.Unix(e.Time, 0).Local().Format("2006-01-02 15-04-05") + "结束"
	}

	m.Hint = &stream.HintMsg{T: stream.HintType_GroupMuteAll, Content: label}
	req.M = m
	if _, err := Msg(context.Background(), req); err != nil {
		qlog.Error(err, req)
	}
}

func onGroupUpdateMuteOne(e *event.GroupMember) {
	req := new(msg.MsgReq)
	m := new(stream.MsgModel)
	groupInfo, err := GetGroupInfo(context.TODO(), e.GroupId)
	if err != nil {
		qlog.Error(e, err)
		return
	}
	m.ChatType = groupInfo.T
	m.SenderId = e.OperatorId
	m.RecvId = e.GroupId
	m.MsgType = stream.MsgType_MsgTypeHint

	var label string
	var t = stream.HintType_GroupMuteSomeOne
	if !e.Flag {
		t = stream.HintType_GroupMuteSomeOneCancel
	}

	if t == stream.HintType_GroupMuteSomeOne {
		{
			lstr, err := gclient.FullUserNameLabel(context.Background(), e.UserId)
			if err != nil {
				qlog.Error(err, e)
				return
			}
			label += lstr + "被"
		}
		{

			lstr, err := gclient.FullUserNameLabel(context.Background(), e.OperatorId)
			if err != nil {
				qlog.Error(err, e)
				return
			}
			if e.Time == 0 {
				label += lstr + "结束禁言"
			} else {
				label += lstr + "禁言直到" + time.Unix(e.Time, 0).Local().Format("2006-01-02 15-04-05") + "结束"
			}

		}
	} else {
		{
			lstr, err := gclient.FullUserNameLabel(context.Background(), e.UserId)
			if err != nil {
				qlog.Error(err, e)
				return
			}
			label += lstr + "被"
		}
		{

			lstr, err := gclient.FullUserNameLabel(context.Background(), e.OperatorId)
			if err != nil {
				qlog.Error(err, e)
				return
			}
			label += lstr + "取消了禁言"
		}
	}

	m.Hint = &stream.HintMsg{T: t, Content: label}
	req.M = m
	if _, err := Msg(context.Background(), req); err != nil {
		qlog.Error(err, req)
	}
}

func onGroupDismiss(e *event.GroupDismiss) {
	m := new(stream.MsgModel)
	groupInfo, err := GetGroupInfo(context.TODO(), e.GroupId)
	if err != nil {
		qlog.Error(e, err)
		return
	}
	m.ChatType = groupInfo.T
	m.SenderId = 0
	m.RecvId = e.GroupId
	m.MsgType = stream.MsgType_MsgTypeHint
	var label string
	lstr, err := gclient.FullUserNameLabel(context.Background(), e.OperatorId)
	if err != nil {
		qlog.Error(err, e)
		return
	}
	label += lstr + "解散了群组"
	m.Hint = &stream.HintMsg{T: stream.HintType_GroupDismiss, Content: label}

	_, _ = onRecvGroupMsg(context.Background(), m, e.Ids)
}

func onGroupTransfer(e *event.GroupMember) {
	req := new(msg.MsgReq)
	m := new(stream.MsgModel)
	groupInfo, err := GetGroupInfo(context.TODO(), e.GroupId)
	if err != nil {
		qlog.Error(e, err)
		return
	}
	m.ChatType = groupInfo.T
	m.SenderId = e.OperatorId
	m.RecvId = e.GroupId
	m.MsgType = stream.MsgType_MsgTypeHint
	var label string
	lstr, err := gclient.FullUserNameLabel(context.Background(), e.OperatorId)
	if err != nil {
		qlog.Error(err, e)
		return
	}
	label += lstr + "将群主转让给"
	lstr, err = gclient.FullUserNameLabel(context.Background(), e.UserId)
	if err != nil {
		qlog.Error(err, e)
		return
	}
	label += lstr
	m.Hint = &stream.HintMsg{T: stream.HintType_GroupTransfer, Content: label}
	req.M = m
	if _, err := Msg(context.Background(), req); err != nil {
		qlog.Error(err, req)
	}
}
