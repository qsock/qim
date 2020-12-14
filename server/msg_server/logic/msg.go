package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qf/store/ka"
	"github.com/qsock/qim/config/mq"
	"github.com/qsock/qim/lib/codes"
	"github.com/qsock/qim/lib/constdef"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/gclient"
	"github.com/qsock/qim/lib/method"
	levent "github.com/qsock/qim/lib/proto/event"
	"github.com/qsock/qim/lib/proto/msg"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/stream"
	"github.com/qsock/qim/lib/proto/user"
	"github.com/qsock/qim/lib/tablename"
	"github.com/qsock/qim/server/msg_server/config"
	"time"
)

func SessConnect(ctx context.Context, req *msg.SessConnectReq) (*ret.BytesResp, error) {
	m := &constdef.JsonRet{T: stream.StreamType_ConnectS2C}
	m.Data = &stream.ConnectMsgModel{Key: req.ServerKey, Uuid: req.SessId}
	b, _ := json.Marshal(m)
	return &ret.BytesResp{Val: b}, nil
}

func genChatId(ctx context.Context, senderId, recverId int64, chatType stream.ChatType) (string, error) {
	arrs := make([]int64, 0, 3)
	arrs = append(arrs, int64(chatType))
	if chatType == stream.ChatType_ChatTypeSingle {
		if senderId > recverId {
			arrs = append(arrs, recverId, senderId)
		} else {
			arrs = append(arrs, senderId, recverId)
		}
	} else if chatType == stream.ChatType_ChatTypeGroup ||
		chatType == stream.ChatType_ChatTypeRoom {
		arrs = append(arrs, 0, recverId)
	} else {
		return "", nil
	}
	return hashIdCodec.EncodeInt64(arrs)
}

func exactChatId(ctx context.Context, chatId string) (stream.ChatType, int64, int64, error) {
	vals, err := hashIdCodec.DecodeInt64WithError(chatId)
	if err != nil {
		qlog.Ctx(ctx).Error(chatId, err)
		return stream.ChatType_ChatTypeFalse, 0, 0, err
	}
	if len(vals) != 3 {
		qlog.Ctx(ctx).Error(chatId, vals)
		return stream.ChatType_ChatTypeFalse, 0, 0, errors.New("chatId invalid")
	}

	return stream.ChatType(vals[0]), vals[1], vals[2], nil
}

func onRecvSingleMsg(ctx context.Context, m *stream.MsgModel) (*msg.MsgResp, error) {
	flag, err := isFriend(ctx, m.SenderId, m.RecvId)
	if err != nil {
		return nil, err
	}
	resp := new(msg.MsgResp)

	if !flag {
		qlog.Ctx(ctx).Error(m, flag)
		resp.Err = codes.Error(codes.ErrorUserIsFriend)
		return resp, nil
	}

	// 保存消息
	if err := saveMsg(ctx, m); err != nil {
		return nil, err
	}
	e := new(levent.NewMsg)
	e.ChatId = m.ChatId
	e.MsgId = m.MsgId
	ka.TopicEvent(mq.TopicIm, mq.EImNew, e)
	resp.ChatId = m.ChatId
	resp.MsgId = m.MsgId
	return resp, nil
}

func onRecvGroupMsg(ctx context.Context, m *stream.MsgModel, ids ...[]int64) (*msg.MsgResp, error) {
	resp := new(msg.MsgResp)
	groupInfo, err := GetGroupInfo(ctx, m.RecvId)
	if err != nil {
		return nil, err
	}
	if groupInfo.DeletedOn > 0 {
		memberInfo, err := GetGroupMemberInfo(ctx, m.RecvId, m.SenderId)
		if err != nil {
			if err == sql.ErrNoRows {
				resp.Err = codes.Error(codes.ErrorUserGroupNotMember)
				return resp, nil
			}
			return nil, err
		}
		if memberInfo.IsBlocked {
			resp.Err = codes.Error(codes.ErrorUserGroupBeenBlock)
			return resp, nil
		}
		if memberInfo.MuteUntil > time.Now().Unix() {
			resp.Err = codes.Error(codes.ErrorUserGroupBeenMute)
			return resp, nil
		}
		if groupInfo.MuteUtil > time.Now().Unix() {
			resp.Err = codes.Error(codes.ErrorUserGroupBeenMute)
			return resp, nil
		}
	} else if groupInfo.DeletedOn == 0 && len(ids) > 0 {
		// 如果是发送群组解散的消息
	} else {
		resp.Err = codes.Error(codes.ErrorMsgInvalid)
		return resp, nil
	}
	// 保存消息
	if err := saveMsg(ctx, m); err != nil {
		return nil, err
	}
	e := new(levent.NewMsg)
	e.ChatId = m.ChatId
	e.MsgId = m.MsgId
	if len(ids) > 0 {
		e.Ids = ids[0]
	}
	ka.TopicEvent(mq.TopicIm, mq.EImNew, e)
	resp.ChatId = m.ChatId
	resp.MsgId = m.MsgId
	return resp, nil
}

func saveMsg(ctx context.Context, m *stream.MsgModel) error {
	shardName := tablename.ImMsgShard(m.ChatId, config.GetEnv())
	tname := tablename.ImMsg(m.ChatId, config.GetEnv())
	imshard := dao.GetConn(shardName)
	imIsql := "insert into " + tname +
		" (msg_id,chat_id,sender_id,recver_id,chat_type," +
		"msg_type,status,created_on,content)" +
		"values (?,?,?,?,?," +
		"?,?,?,?)"
	b, _ := json.Marshal(m)
	_, err := imshard.Exec(imIsql,
		m.MsgId, m.ChatId, m.SenderId, m.RecvId, m.ChatType,
		m.MsgType, m.Status, m.CreatedOn, string(b),
	)
	if err != nil {
		qlog.Ctx(ctx).Error(err, string(b))
		return nil
	}
	return err
}

func isGroupMemberBeenMute(ctx context.Context, userId int64, groupId int64) (bool, error) {
	req := new(user.IsGroupMemberBeenMuteReq)
	req.UserId = userId
	req.GroupId = groupId
	resp := new(ret.BoolResp)
	if err := qgrpc.Call(ctx, method.UserIsGroupMemberBeenMute, req, resp); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return false, err
	}
	return resp.Flag, nil
}

func GetGroupInfo(ctx context.Context, groupId int64) (*user.GroupInfo, error) {
	req := new(user.GroupInfoReq)
	req.GroupId = groupId
	resp := new(user.GroupInfoResp)
	if err := qgrpc.Call(ctx, method.UserGroupInfoById, req, resp); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	return resp.Info, nil
}

func GetGroupMemberInfo(ctx context.Context, groupId, memberId int64) (*user.GroupMember, error) {
	req := new(user.GroupMemberByIdsReq)
	req.GroupId = groupId
	req.UserIds = []int64{memberId}
	resp := new(user.GroupMemberByIdsResp)
	if err := qgrpc.Call(ctx, method.UserGroupMemberByIds, req, resp); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	if len(resp.GetMembers()) == 0 {
		qlog.Ctx(ctx).Error(req, sql.ErrNoRows)
		return nil, sql.ErrNoRows
	}
	return resp.Members[0], nil
}

func isGroupMember(ctx context.Context, userId int64, groupId int64) (bool, error) {
	req := new(user.IsGroupMemberReq)
	req.UserId = userId
	req.GroupId = groupId
	resp := new(ret.BoolResp)
	if err := qgrpc.Call(ctx, method.UserIsGroupMember, req, resp); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return false, err
	}
	return resp.Flag, nil
}

func isChatRoomMember(ctx context.Context, userId int64, groupId int64) (bool, error) {
	req := new(user.IsGroupMemberReq)
	req.UserId = userId
	req.GroupId = groupId
	resp := new(ret.BoolResp)
	if err := qgrpc.Call(ctx, method.UserIsGroupMember, req, resp); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return false, err
	}
	return resp.Flag, nil
}

func isFriend(ctx context.Context, userId int64, friendId int64) (bool, error) {
	req := new(user.IsFriendReq)
	req.UserId = userId
	req.FriendId = friendId
	resp := new(ret.BoolResp)
	if err := qgrpc.Call(ctx, method.UserIsFriend, req, resp); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return false, err
	}
	return resp.Flag, nil
}

func Msg(ctx context.Context, req *msg.MsgReq) (*msg.MsgResp, error) {
	resp := new(msg.MsgResp)

	// 先保存原始消息
	m := req.M

	if m.RecvId == 0 ||
		m.ChatType == stream.ChatType_ChatTypeFalse ||
		m.MsgType == stream.MsgType_MsgTypeFalse ||
		m.SenderId == 0 {
		b, _ := json.Marshal(m)
		qlog.Ctx(ctx).Error(m, string(b))
		resp.Err = codes.Error(codes.ErrorMsgInvalid)
		return resp, nil
	}
	// 生成消息id
	msgId := gclient.GenSnowflake()
	m.MsgId = msgId
	m.Status = stream.MsgStatus_MsgNormal
	var err error
	m.CreatedOn = time.Now().Unix()
	m.ChatId, err = genChatId(ctx, m.SenderId, m.RecvId, m.ChatType)
	if err != nil {
		b, _ := json.Marshal(m)
		qlog.Ctx(ctx).Error(err, string(b))
		resp.Err = codes.Error(codes.ErrorMsgInvalid)
		return resp, nil
	}
	if m.ChatType == stream.ChatType_ChatTypeSingle {
		return onRecvSingleMsg(ctx, m)
	} else if m.ChatType == stream.ChatType_ChatTypeGroup || m.ChatType == stream.ChatType_ChatTypeRoom {
		return onRecvGroupMsg(ctx, m)
	}

	resp.ChatId = m.ChatId
	resp.MsgId = msgId
	return resp, nil
}

func SysMsg(ctx context.Context, req *msg.SysMsgReq) (*ret.IntResp, error) {
	resp := new(ret.IntResp)
	m := req.M
	if m == nil {
		resp.Err = codes.Error(codes.ErrorParameter)
		return resp, nil
	}
	msgId := gclient.GenSnowflake()
	m.MsgId = msgId
	m.CreatedOn = time.Now().Unix()
	m.Status = stream.MsgStatus_MsgNormal

	if m.NeedSave {
		if err := saveSysMsg(ctx, m); err != nil {
			return nil, err
		}
	}

	// 立即发送 & 需要发送推送通知的
	// 才需要长连接推送
	// 其他的等用户拉取即可
	// kafka消息的系统消息是全量推送的
	if m.SendOn == 0 && m.NeedPush {
		ka.TopicEvent(mq.TopicIm, mq.EImSys, m)
		return resp, nil
	}
	panic("implement me")
}

func saveSysMsg(ctx context.Context, m *stream.SysMsgModel) error {
	shardName := tablename.SysMsgShard(m.RecverId, config.GetEnv())
	shard := dao.GetConn(shardName)
	imIsql := "insert into sys_msg" +
		" (msg_id,sender_id,recver_id,created_on,msg_type," +
		"need_push,send_on,status,content)" +
		"values (?,?,?,?,?," +
		"?,?,?,?)"
	b, _ := json.Marshal(m)

	_, err := shard.Exec(imIsql,
		m.MsgId, m.SenderId, m.RecverId, m.CreatedOn, m.MsgType,
		m.NeedPush, m.SendOn, m.Status, string(b),
	)
	if err != nil {
		qlog.Ctx(ctx).Error(err, string(b))
		return nil
	}
	return err
}
