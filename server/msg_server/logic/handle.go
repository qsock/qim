package logic

import (
	"encoding/json"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/store/ka"
	"github.com/qsock/qim/config/mq"
	levent "github.com/qsock/qim/lib/proto/event"
	"github.com/qsock/qim/lib/proto/stream"
)

// 收到群系统消息
func HandleImMsg(b []byte) {
	e := &ka.EConsumer{}
	err := json.Unmarshal(b, e)
	if err != nil {
		qlog.Error(err, string(b))
		return
	}
	switch e.Type {
	case mq.EImNew:
		{
			m := &levent.NewMsg{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onNewMsg(m)
		}
	case mq.EImSys:
		{
			m := &stream.SysMsgModel{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onSysMsg(m)
		}
	}
}

func HandleServerEvent(b []byte) {
	qlog.Info(string(b))
	e := &ka.EConsumer{}
	err := json.Unmarshal(b, e)
	if err != nil {
		qlog.Error(err, string(b))
		return
	}
	switch e.Type {

	// 成功添加好友
	case mq.EEventFriendAdd:
		{
			m := &levent.Friend{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onFriendAdd(m)
		}
	case mq.EEventFriendDel:
		{
			m := &levent.Friend{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onFriendDel(m)
		}
	case mq.EEventFriendMarkname:
		{
			m := &levent.Friend{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onFriendMarkname(m)
		}
	case mq.EEventFresher:
		{
		}
	case mq.EEventFriendApply:
		{
			m := &levent.Friend{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onFriendApply(m)
		}
	case mq.EEventFriendApplyReject:
	case mq.EEventFriendApplyDel:
	case mq.EEventFriendApplyIgnore:
	case mq.EEventGroupApply:
		{
			m := &levent.GroupNewApply{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onGroupApply(m)
		}
	case mq.EEventGroupApplyReject:
	case mq.EEventGroupApplyIgnore:
	case mq.EEventGroupManager:
		{
			m := &levent.GroupManager{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onGroupManager(m)
		}
	case mq.EEventGroupNewMember:
		{
			// 新成员
			m := &levent.GroupMember{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onGroupNewMember(m)
		}
	case mq.EEventGroupDelMember:
		{
			// 删除成员
			m := &levent.GroupMember{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onGroupDelMember(m)
		}
	case mq.EEventGroupLeaveMember:
		{
			// 离开群组
			m := &levent.GroupMember{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onGroupLeaveMember(m)
		}
	case mq.EEventGroupUpdateName:
		{
			// 更新群组
			m := &levent.GroupUpdate{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onGroupUpdateName(m)
		}
	case mq.EEventGroupUpdateNotice:
		{
			// 更新群组
			m := &levent.GroupUpdate{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onGroupUpdateNotice(m)
		}
	case mq.EEventGroupUpdateAvatar:
		{
			// 更新群组
			m := &levent.GroupUpdate{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onGroupUpdateAvatar(m)
		}
	case mq.EEventGroupUpdateJointype:
	case mq.EEventGroupUpdateMute:
		{
			// 更新群组
			m := &levent.GroupUpdate{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onGroupUpdateMute(m)
		}
	case mq.EEventGroupMuteone:
		{
			// 更新群组
			m := &levent.GroupMember{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onGroupUpdateMuteOne(m)
		}
	case mq.EEventGroupDismiss:
		{
			m := &levent.GroupDismiss{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onGroupDismiss(m)
		}
	case mq.EEventGroupTransfer:
		{
			m := &levent.GroupMember{}
			if err := json.Unmarshal(e.Msg, m); err != nil {
				qlog.Error(err, string(b))
				return
			}
			onGroupTransfer(m)
		}
	}
}

func HandleLogTrace(b []byte) {
	qlog.Info(string(b))
}
