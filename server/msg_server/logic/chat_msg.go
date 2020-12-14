package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/store/ka"
	"github.com/qsock/qim/config/mq"
	"github.com/qsock/qim/lib/codes"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/gclient"
	"github.com/qsock/qim/lib/proto/event"
	"github.com/qsock/qim/lib/proto/msg"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/stream"
	"github.com/qsock/qim/lib/tablename"
	"github.com/qsock/qim/lib/util"
	"github.com/qsock/qim/server/msg_server/config"
	"time"
)

func getMsgById(ctx context.Context, chatId string, msgId int64) (*stream.MsgModel, error) {
	items, err := getMsgByIds(ctx, chatId, []int64{msgId})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		qlog.Ctx(ctx).Error(msgId, chatId, err)
		return nil, sql.ErrNoRows
	}
	return items[0], nil
}

func getMsgByIds(ctx context.Context, chatId string, msgIds []int64) ([]*stream.MsgModel, error) {
	shardName := tablename.ImMsgShardR1(chatId, config.GetEnv())
	tname := tablename.ImMsg(chatId, config.GetEnv())
	ssql := "select content from " + tname + " where msg_id in (%s) and chat_id=? and status=?"
	ssql = fmt.Sprintf(ssql, util.Int64sToStr(msgIds))
	rows, err := dao.GetConn(shardName).Query(ssql, chatId, stream.MsgStatus_MsgNormal)
	if err != nil {
		qlog.Ctx(ctx).Error(msgIds, chatId, err)
		return nil, err
	}
	items := make([]*stream.MsgModel, 0)
	defer rows.Close()
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			qlog.Ctx(ctx).Error(msgIds, chatId, err)
			continue
		}

		m := new(stream.MsgModel)
		if err := json.Unmarshal([]byte(content), m); err != nil {
			qlog.Ctx(ctx).Error(msgIds, chatId, err, content)
			continue
		}
		items = append(items, m)
	}
	return items, nil
}

// 撤回自己的消息
func RevertSelfMsg(ctx context.Context, req *msg.RevertSelfMsgReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	chatType, _, recverId, err := exactChatId(ctx, req.ChatId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		resp.Err = codes.Error(codes.ErrorMsgMorethan2Minutes)
		return resp, nil
	}

	t := gclient.TraceToTime(req.MsgId)
	if (t - time.Now().Unix()) > 120 {
		qlog.Ctx(ctx).Error(req, err)
		resp.Err = codes.Error(codes.ErrorMsgMorethan2Minutes)
		return resp, nil
	}
	userLabel, err := gclient.FullUserNameLabel(ctx, req.UserId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}

	shardName := tablename.ImMsgShardR1(req.ChatId, config.GetEnv())
	tname := tablename.ImMsg(req.ChatId, config.GetEnv())

	usql := "update " + tname + " set status=? where msg_id=? and chat_id=? and sender_id=?"
	result, err := dao.GetConn(shardName).Exec(usql, stream.MsgStatus_MsgRevert, req.MsgId, req.ChatId, req.UserId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}

	mmsg := new(stream.MsgModel)
	mmsg.MsgId = gclient.GenSnowflake()
	mmsg.SenderId = req.UserId
	mmsg.RecvId = recverId
	mmsg.CreatedOn = time.Now().Unix()
	mmsg.Status = stream.MsgStatus_MsgNormal
	mmsg.ChatId = req.ChatId
	mmsg.MsgType = stream.MsgType_MsgTypeRevertMsg
	mmsg.ChatType = chatType
	mmsg.Device = req.GetMeta().GetDevice()

	revertMsg := new(stream.RevertMsg)
	revertMsg.OperatorId = req.UserId
	revertMsg.Content = userLabel + "撤回了一条消息"
	mmsg.Revert = revertMsg

	if err := saveMsg(ctx, mmsg); err != nil {
		return nil, err
	}
	e := new(event.NewMsg)
	e.ChatId = mmsg.ChatId
	e.MsgId = mmsg.MsgId
	ka.TopicEvent(mq.TopicIm, mq.EImNew, e)
	return resp, nil
}

// 管理撤回消息
func ManagerChatMsgRevert(ctx context.Context, req *msg.ManagerChatMsgRevertReq) (*ret.EmptyResp, error) {
	//TODO 检查是否管理员

	resp := new(ret.EmptyResp)
	chatType, _, recverId, err := exactChatId(ctx, req.ChatId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		resp.Err = codes.Error(codes.ErrorMsgMorethan2Minutes)
		return resp, nil
	}

	shardName := tablename.ImMsgShardR1(req.ChatId, config.GetEnv())
	tname := tablename.ImMsg(req.ChatId, config.GetEnv())
	var senderId int64
	ssql := "select sender_id from " + tname + " where msg_id=? and chat_id=? limit 1"
	if err := dao.GetConn(shardName).QueryRow(ssql, req.MsgId, req.ChatId).Scan(&senderId); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return resp, nil
	}

	if senderId == req.UserId {
		resp.Err = codes.Error(codes.ErrorParameter)
		return resp, nil
	}

	userLabel, err := gclient.FullUserNameLabel(ctx, senderId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	operatorLabel, err := gclient.FullUserNameLabel(ctx, req.UserId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}

	usql := "update " + tname + " set status=? where msg_id=? and chat_id=? and sender_id=?"
	result, err := dao.GetConn(shardName).Exec(usql, stream.MsgStatus_MsgRevert, req.MsgId, req.ChatId, req.UserId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}

	mmsg := new(stream.MsgModel)
	mmsg.MsgId = gclient.GenSnowflake()
	mmsg.SenderId = req.UserId
	mmsg.RecvId = recverId
	mmsg.CreatedOn = time.Now().Unix()
	mmsg.Status = stream.MsgStatus_MsgNormal
	mmsg.ChatId = req.ChatId
	mmsg.MsgType = stream.MsgType_MsgTypeRevertMsg
	mmsg.ChatType = chatType
	mmsg.Device = req.GetMeta().GetDevice()

	revertMsg := new(stream.RevertMsg)
	revertMsg.OperatorId = req.UserId
	revertMsg.Content = operatorLabel + "撤回了" + userLabel + "的一条消息"
	mmsg.Revert = revertMsg

	if err := saveMsg(ctx, mmsg); err != nil {
		return nil, err
	}
	e := new(event.NewMsg)
	e.ChatId = mmsg.ChatId
	e.MsgId = mmsg.MsgId
	ka.TopicEvent(mq.TopicIm, mq.EImNew, e)
	return resp, nil
}

func GetSysMsg(ctx context.Context, req *msg.GetSysMsgReq) (*msg.GetSysMsgResp, error) {
	//TODO 还需要思考
	panic("implement me")
}

func GetMemberIdByChatId(ctx context.Context, req *msg.GetMemberIdByChatIdReq) (*msg.GetMemberIdByChatIdResp, error) {
	ids, err := getMemberIdsByChatId(ctx, req.ChatId)
	if err != nil {
		return nil, err
	}
	return &msg.GetMemberIdByChatIdResp{Ids: ids}, nil
}
