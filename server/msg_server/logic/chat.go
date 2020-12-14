package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qf/store/ka"
	"github.com/qsock/qim/config/mq"
	"github.com/qsock/qim/lib/codes"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/gclient"
	"github.com/qsock/qim/lib/method"
	"github.com/qsock/qim/lib/proto/event"
	"github.com/qsock/qim/lib/proto/msg"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/stream"
	"github.com/qsock/qim/lib/proto/user"
	"github.com/qsock/qim/lib/tablename"
	"github.com/qsock/qim/server/msg_server/config"
	"strings"
	"time"
)

func getMemberIdsByChatId(ctx context.Context, chatId string) ([]int64, error) {
	t, senderId, recverId, err := exactChatId(ctx, chatId)
	if err != nil {
		qlog.Ctx(ctx).Error(chatId)
		return nil, err
	}

	ids := make([]int64, 0)
	if t == stream.ChatType_ChatTypeSingle {
		ids = append(ids, senderId, recverId)
	} else if t == stream.ChatType_ChatTypeGroup || t == stream.ChatType_ChatTypeRoom {
		req := new(user.GroupMemberIdsReq)
		req.GroupId = recverId
		resp := new(user.GroupMemberIdsResp)
		if err := qgrpc.Call(ctx, method.UserGroupMemberIds, req, resp); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			return nil, err
		}
		ids = resp.Ids
	}
	return ids, nil
}

// 标记会话已读
func MarkChatRead(ctx context.Context, req *msg.MarkChatReadReq) (*ret.EmptyResp, error) {
	shardName := tablename.ChatlistShard(req.UserId, config.GetEnv())
	usql := "update chat_list set updated_on=unix_timestamp(),unread_ct=0,read_last_msg_id=last_msg_id " +
		"where user_id=? and chat_id=?"
	result, err := dao.GetConn(shardName).Exec(usql, req.UserId, req.ChatId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	resp := new(ret.EmptyResp)
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}
	ssql := "select read_last_msg_id from " +
		"where user_id=? and chat_id=? limit 1"
	var lastMsgId int64
	if err := dao.GetConn(shardName).QueryRow(ssql, req.UserId, req.ChatId).Scan(&lastMsgId); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}

	memberIds, err := getMemberIdsByChatId(ctx, req.ChatId)
	if err != nil {
		return nil, err
	}
	sysMsg := new(stream.SysMsgModel)
	sysMsg.MsgId = gclient.GenSnowflake()
	sysMsg.CreatedOn = time.Now().Unix()
	sysMsg.SenderId = req.UserId
	sysMsg.RecverId = req.UserId
	sysMsg.MsgType = stream.SysMsgType_Command
	commandMsg := new(stream.SysCommandMsg)
	commandMsg.RecverId = req.UserId
	commandMsg.OperatorId = req.UserId
	commandMsg.T = stream.SysCommandType_SysChatRead
	mam := map[string]interface{}{"chat_id": req.ChatId}
	b, _ := json.Marshal(mam)
	commandMsg.Extra = string(b)
	sysMsg.Command = commandMsg

	// 需要给所有人发个已读的系统消息
	go func() {
		for _, memberId := range memberIds {
			sysMsg.RecverId = memberId
			sysMsg.GetCommand().RecverId = memberId
			ka.TopicEvent(mq.TopicIm, mq.EImSys, sysMsg)
		}
	}()
	return resp, nil
}

func ChatAhead(ctx context.Context, req *msg.ChatAheadReq) (*ret.EmptyResp, error) {
	shardName := tablename.ChatlistShard(req.UserId, config.GetEnv())
	usql := "update chat_list set ahead_on=unix_timestamp(),updated_on=unix_timestamp() " +
		"where user_id=? and chat_id=?"
	result, err := dao.GetConn(shardName).Exec(usql, req.UserId, req.ChatId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	resp := new(ret.EmptyResp)
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}

	// 发个置顶的系统消息
	sysMsg := new(stream.SysMsgModel)
	sysMsg.MsgId = gclient.GenSnowflake()
	sysMsg.CreatedOn = time.Now().Unix()
	sysMsg.SenderId = req.UserId
	sysMsg.RecverId = req.UserId
	sysMsg.MsgType = stream.SysMsgType_Command
	commandMsg := new(stream.SysCommandMsg)
	commandMsg.RecverId = req.UserId
	commandMsg.OperatorId = req.UserId
	commandMsg.T = stream.SysCommandType_SysChatAhead
	if !req.IsAhead {
		commandMsg.T = stream.SysCommandType_SysChatAheadCancel
	}
	mam := map[string]interface{}{"chat_id": req.ChatId}
	b, _ := json.Marshal(mam)
	commandMsg.Extra = string(b)
	sysMsg.Command = commandMsg
	ka.TopicEvent(mq.TopicIm, mq.EImSys, sysMsg)
	return resp, nil
}

func ChatTouch(ctx context.Context, req *msg.ChatTouchReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)

	chatId, err := genChatId(ctx, req.UserId, req.RecverId, req.Type)
	if err != nil {
		return nil, err
	}
	if req.Type == stream.ChatType_ChatTypeSingle {
		flag, err := isFriend(ctx, req.UserId, req.RecverId)
		if err != nil {
			return nil, err
		}
		if !flag {
			resp.Err = codes.Error(codes.ErrorParameter)
			return resp, nil
		}
	} else if req.Type == stream.ChatType_ChatTypeGroup {
		flag, err := isGroupMember(ctx, req.UserId, req.RecverId)
		if err != nil {
			return nil, err
		}
		if !flag {
			resp.Err = codes.Error(codes.ErrorParameter)
			return resp, nil
		}
	} else {
		resp.Err = codes.Error(codes.ErrorParameter)
		return resp, nil
	}

	shardName := tablename.ChatlistShard(req.UserId, config.GetEnv())
	isql := "insert into chat_list (user_id,chat_id,chat_type,created_on,updated_on," +
		"unread_ct,deleted_on) values (?,?,?,unix_timestamp(),unix_timestamp()," +
		"0,0) on duplicate " +
		"update updated_on=unix_timestamp(),unread_ct=0,read_last_msg_id=last_msg_id,deleted_on=0"
	if _, err := dao.GetConn(shardName).Exec(isql,
		req.UserId, chatId, req.Type); err != nil {
		qlog.Ctx(ctx).Error(req, chatId, err)
		return nil, err
	}

	// 如果是聊天室
	if req.Type == stream.ChatType_ChatTypeRoom {
		m := new(stream.MsgModel)
		m.ChatType = req.Type
		m.SenderId = req.UserId
		m.RecvId = req.RecverId
		m.MsgType = stream.MsgType_MsgTypeHint
		label, _ := gclient.FullUserNameLabel(ctx, m.SenderId)
		m.Hint = &stream.HintMsg{T: stream.HintType_GroupJoin, Content: label + "进入了房间"}
		if _, err := Msg(ctx, &msg.MsgReq{M: m}); err != nil {
			return nil, err
		}
		return resp, nil
	}
	// 发个系统消息
	{
		sysMsg := new(stream.SysMsgModel)
		sysMsg.CreatedOn = time.Now().Unix()
		sysMsg.SenderId = req.UserId
		sysMsg.RecverId = req.UserId
		sysMsg.MsgType = stream.SysMsgType_Command
		commandMsg := new(stream.SysCommandMsg)
		commandMsg.RecverId = req.UserId
		commandMsg.OperatorId = req.UserId
		commandMsg.T = stream.SysCommandType_SysChatTouch
		mam := map[string]interface{}{"chat_id": chatId}
		b, _ := json.Marshal(mam)
		commandMsg.Extra = string(b)
		sysMsg.Command = commandMsg
		ka.TopicEvent(mq.TopicIm, mq.EImSys, sysMsg)
	}
	return resp, nil
}

//  删除会话，单边删除
func ChatRemove(ctx context.Context, req *msg.ChatRemoveReq) (*ret.EmptyResp, error) {
	shardName := tablename.ChatlistShard(req.UserId, config.GetEnv())
	usql := "update chat_list set deleted_on=unix_timestamp(),updated_on=unix_timestamp() " +
		"where user_id=? and chat_id=?"
	result, err := dao.GetConn(shardName).Exec(usql, req.UserId, req.ChatId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	resp := new(ret.EmptyResp)
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}

	// 发个置顶的系统消息
	sysMsg := new(stream.SysMsgModel)
	sysMsg.MsgId = gclient.GenSnowflake()
	sysMsg.CreatedOn = time.Now().Unix()
	sysMsg.SenderId = req.UserId
	sysMsg.RecverId = req.UserId
	sysMsg.MsgType = stream.SysMsgType_Command
	commandMsg := new(stream.SysCommandMsg)
	commandMsg.RecverId = req.UserId
	commandMsg.OperatorId = req.UserId
	commandMsg.T = stream.SysCommandType_SysChatDeleted
	mam := map[string]interface{}{"chat_id": req.ChatId}
	b, _ := json.Marshal(mam)
	commandMsg.Extra = string(b)
	sysMsg.Command = commandMsg
	ka.TopicEvent(mq.TopicIm, mq.EImSys, sysMsg)
	return resp, nil
}

// 清理会话，删除所有聊天信息
func ChatClear(ctx context.Context, req *msg.ChatClearReq) (*ret.EmptyResp, error) {
	chatType, _, recverId, err := exactChatId(ctx, req.ChatId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}

	memberIds, err := getMemberIdsByChatId(ctx, req.ChatId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	var flag bool
	for _, memberId := range memberIds {
		chatListShard := tablename.ChatlistShard(memberId, config.GetEnv())
		dsql := "delete from chat_list where chat_id=?"
		result, err := dao.GetConn(chatListShard).Exec(dsql, memberId, req.ChatId)
		if err != nil {
			qlog.Ctx(ctx).Error(req, err)
			continue
		}
		if n, _ := result.RowsAffected(); n > 0 {
			flag = true
		}
	}
	resp := new(ret.EmptyResp)
	if !flag {
		return resp, nil
	}
	imMsgShard := tablename.ImMsgShard(req.ChatId, config.GetEnv())
	imtable := tablename.ImMsg(req.ChatId, config.GetEnv())
	if _, err := dao.GetConn(imMsgShard).Exec("update "+imtable+
		" set status=? where chat_id=?", stream.MsgStatus_MsgDeleted, req.ChatId); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}

	// 发个清理的消息
	mmsg := new(stream.MsgModel)
	mmsg.MsgId = gclient.GenSnowflake()
	mmsg.CreatedOn = time.Now().Unix()
	mmsg.ChatId = req.ChatId
	mmsg.RecvId = recverId
	mmsg.CreatedOn = time.Now().Unix()
	mmsg.Status = stream.MsgStatus_MsgNormal
	mmsg.ChatType = chatType
	mmsg.MsgType = stream.MsgType_MsgTypeCommand
	commandMsg := new(stream.CommandMsg)
	commandMsg.T = stream.CommandType_ChatClear
	commandMsg.RecverId = recverId

	mmsg.Command = commandMsg
	if err := saveMsg(ctx, mmsg); err != nil {
		return nil, err
	}
	l := new(event.NewMsg)
	l.ChatId = mmsg.ChatId
	l.MsgId = mmsg.MsgId
	ka.TopicEvent(mq.TopicIm, mq.EImNew, l)

	return resp, nil
}

// 得到会话的id
func ChatIds(ctx context.Context, req *msg.ChatIdsReq) (*msg.ChatIdsResp, error) {
	shardName := tablename.ChatlistShard(req.UserId, config.GetEnv())
	ids := make([]string, 0)
	rows, err := dao.GetConn(shardName).Query("select chat_id from chat_list where user_id=?", req.UserId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			continue
		}
		ids = append(ids, id)
	}
	resp := new(msg.ChatIdsResp)
	resp.Ids = ids
	return resp, nil
}

func getRecverIdByChatId(ctx context.Context, chatId string, userId int64) (int64, error) {
	t, id1, id2, err := exactChatId(ctx, chatId)
	if err != nil {
		qlog.Ctx(ctx).Error(chatId, userId, err)
		return 0, err
	}
	if t == stream.ChatType_ChatTypeRoom {
		return id2, nil
	} else if t == stream.ChatType_ChatTypeGroup {
		return id2, nil
	}
	if id1 != userId {
		return id1, nil
	}
	return id2, nil
}

func ChatByUids(ctx context.Context, req *msg.ChatByUidsReq) (*msg.ChatByUidsResp, error) {
	resp := new(msg.ChatByUidsResp)
	shardName := tablename.ChatlistShard(req.UserId, config.GetEnv())

	csql := "select count(1) from chat_list where user_id=?"
	if err := dao.GetConn(shardName).QueryRow(csql, req.UserId).Scan(&resp.Total); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	if resp.Total <= req.Page*req.PageSize {
		return resp, nil
	}

	items := make([]*msg.ChatItem, 0)
	ssql := "select chat_id,chat_type,ahead_on,updated_on,unread_ct," +
		"is_mute,read_last_msg_id,last_msg_id from chat_list where user_id=? " +
		"order by updated_on limit %d offset %d"
	ssql = fmt.Sprintf(ssql, req.PageSize, req.PageSize*req.Page)

	gids := make([]int64, 0)
	uids := make([]int64, 0)

	rows, err := dao.GetConn(shardName).Query(ssql, req.UserId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	for rows.Next() {
		item := new(msg.ChatItem)
		if err := rows.Scan(&item.ChatId, &item.T, &item.AheadOn, &item.UpdatedOn, &item.UnreadCt,
			&item.IsMute, &item.ReadLastMsgId, &item.LastMsgId); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			continue
		}
		item.RecverId, _ = getRecverIdByChatId(ctx, item.ChatId, req.UserId)
		if item.T == stream.ChatType_ChatTypeGroup {
			gids = append(gids, item.RecverId)
		} else if item.T == stream.ChatType_ChatTypeSingle {
			uids = append(uids, item.RecverId)
		} else {
			continue
		}
		items = append(items, item)
	}
	_ = rows.Close()

	uinfos, err := gclient.FriendByIds(ctx, uids, req.UserId)
	if err != nil {
		return nil, err
	}
	ginfos, err := gclient.Groups(ctx, gids)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if item.T == stream.ChatType_ChatTypeSingle {
			for _, uinfo := range uinfos {
				if uinfo.GetUser().GetUserId() == item.RecverId {
					if len(uinfo.MarkName) > 0 {
						item.Name = uinfo.MarkName
					} else {
						item.Name = uinfo.User.Name
					}
					item.Avatars = []string{uinfo.GetUser().GetAvatar()}
				}

			}
		} else if item.T == stream.ChatType_ChatTypeGroup {
			for _, ginfo := range ginfos {
				if ginfo.GetId() == item.RecverId {
					item.Name = ginfo.Name
					item.Avatars = ginfo.Avatars
				}
			}
		}
	}

	resp.Items = items
	return resp, nil
}

func ChatByIds(ctx context.Context, req *msg.ChatByIdsReq) (*msg.ChatByIdsResp, error) {
	shardName := tablename.ChatlistShard(req.UserId, config.GetEnv())
	items := make([]*msg.ChatItem, 0)
	ssql := "select chat_id,chat_type,ahead_on,updated_on,unread_ct," +
		"is_mute,read_last_msg_id,last_msg_id from chat_list where chat_id in (?) and user_id=?"
	gids := make([]int64, 0)
	uids := make([]int64, 0)

	rows, err := dao.GetConn(shardName).Query(ssql, strings.Join(req.Ids, ","), req.UserId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	for rows.Next() {
		item := new(msg.ChatItem)
		if err := rows.Scan(&item.ChatId, &item.T, &item.AheadOn, &item.UpdatedOn, &item.UnreadCt,
			&item.IsMute, &item.ReadLastMsgId, &item.LastMsgId); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			continue
		}
		item.RecverId, _ = getRecverIdByChatId(ctx, item.ChatId, req.UserId)
		if item.T == stream.ChatType_ChatTypeGroup ||
			item.T == stream.ChatType_ChatTypeRoom {
			gids = append(gids, item.RecverId)
		} else if item.T == stream.ChatType_ChatTypeSingle {
			uids = append(uids, item.RecverId)
		} else {
			continue
		}
		items = append(items, item)
	}
	_ = rows.Close()

	uinfos, err := gclient.FriendByIds(ctx, uids, req.UserId)
	if err != nil {
		return nil, err
	}
	ginfos, err := gclient.Groups(ctx, gids)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.T == stream.ChatType_ChatTypeSingle {
			for _, uinfo := range uinfos {
				if uinfo.GetUser().GetUserId() == item.RecverId {
					if len(uinfo.MarkName) > 0 {
						item.Name = uinfo.MarkName
					} else {
						item.Name = uinfo.User.Name
					}
					item.Avatars = []string{uinfo.GetUser().GetAvatar()}
				}

			}
		} else if item.T == stream.ChatType_ChatTypeGroup ||
			item.T == stream.ChatType_ChatTypeRoom {
			for _, ginfo := range ginfos {
				if ginfo.GetId() == item.RecverId {
					item.Name = ginfo.Name
					item.Avatars = ginfo.Avatars
				}
			}
		}
	}

	resp := new(msg.ChatByIdsResp)
	resp.Items = items
	return resp, nil
}

// 消息记录id
func ChatRecordIds(ctx context.Context, req *msg.ChatRecordIdsReq) (*msg.ChatRecordIdsResp, error) {
	imMsgShard := tablename.ImMsgShard(req.ChatId, config.GetEnv())
	imtable := tablename.ImMsg(req.ChatId, config.GetEnv())
	csql := "select count(1) from " + imtable +
		" where chat_id=? and status=?"
	resp := new(msg.ChatRecordIdsResp)
	if err := dao.GetConn(imMsgShard).QueryRow(csql, req.ChatId, stream.MsgStatus_MsgNormal).Scan(&resp.Total); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	if resp.Total == 0 {
		return resp, nil
	}
	if req.Page*req.PageSize > int32(resp.Total) {
		return resp, nil
	}

	ssql := "select msg_id from " + imtable +
		" where chat_id=? and status=? order by msg_id desc limit %d offset %d"
	ssql = fmt.Sprintf(ssql, req.PageSize, req.Page*req.PageSize)
	ids := make([]int64, 0)
	rows, err := dao.GetConn(imMsgShard).Query(ssql, req.ChatId, stream.MsgStatus_MsgNormal)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			continue
		}
		ids = append(ids, id)
	}
	resp.Ids = ids
	return resp, nil
}

func ChatRecordByIds(ctx context.Context, req *msg.ChatRecordByIdsReq) (*msg.ChatRecordByIdsResp, error) {
	resp := new(msg.ChatRecordByIdsResp)
	{
		shardName := tablename.ChatlistShardR1(req.UserId, config.GetEnv())
		ssql := "select count(1) from chat_list where user_id=? and chat_id=? limit 1"
		var flag bool
		if err := dao.GetConn(shardName).QueryRow(ssql, req.UserId, req.ChatId).Scan(&flag); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			return nil, err
		}
		if !flag {
			resp.Err = codes.Error(codes.ErrorParameter)
			return resp, nil
		}
	}
	msgs, err := getMsgByIds(ctx, req.ChatId, req.Ids)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}

	resp.Items = msgs
	return resp, nil
}

func ChatMute(ctx context.Context, req *msg.ChatMuteReq) (*ret.EmptyResp, error) {
	shardName := tablename.ChatlistShard(req.UserId, config.GetEnv())
	usql := "update chat_list set is_mute=?,updated_on=unix_timestamp() " +
		"where user_id=? and chat_id=?"
	result, err := dao.GetConn(shardName).Exec(usql, req.IsMute, req.UserId, req.ChatId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	resp := new(ret.EmptyResp)
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}

	// 发个置顶的系统消息
	sysMsg := new(stream.SysMsgModel)
	sysMsg.MsgId = gclient.GenSnowflake()
	sysMsg.CreatedOn = time.Now().Unix()
	sysMsg.SenderId = req.UserId
	sysMsg.RecverId = req.UserId
	sysMsg.MsgType = stream.SysMsgType_Command
	commandMsg := new(stream.SysCommandMsg)
	commandMsg.RecverId = req.UserId
	commandMsg.OperatorId = req.UserId
	commandMsg.T = stream.SysCommandType_SysChatMute
	if !req.IsMute {
		commandMsg.T = stream.SysCommandType_SysChatMuteCancel
	}
	mam := map[string]interface{}{"chat_id": req.ChatId}
	b, _ := json.Marshal(mam)
	commandMsg.Extra = string(b)
	sysMsg.Command = commandMsg
	ka.TopicEvent(mq.TopicIm, mq.EImSys, sysMsg)
	return resp, nil
}
