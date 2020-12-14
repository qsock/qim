package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qf/store/ka"
	"github.com/qsock/qim/config/kvconfig"
	"github.com/qsock/qim/config/mq"
	"github.com/qsock/qim/lib/cachename"
	"github.com/qsock/qim/lib/constdef"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/method"
	levent "github.com/qsock/qim/lib/proto/event"
	"github.com/qsock/qim/lib/proto/msg"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/stream"
	"github.com/qsock/qim/lib/proto/ws"
	"github.com/qsock/qim/lib/tablename"
	"github.com/qsock/qim/server/msg_server/config"
	"strconv"
)

type TinySess struct {
	userId    int64
	serverKey string
}

func getSessionIdsByUserIds(ctx context.Context, userIds []int64) ([]*TinySess, error) {
	kvconn := dao.GetKvConn(kvconfig.KvDefault)
	ms := make([]*TinySess, 0)
	pipl := kvconn.Pipeline()

	for _, userId := range userIds {
		cacheKey := cachename.RedisUserWs(userId)
		pipl.Get(cacheKey)
	}
	cmds, err := pipl.Exec()
	if err != nil {
		qlog.Ctx(ctx).Error(userIds, err)
		return nil, err
	}

	for i := 0; i < len(cmds); i++ {
		cmd := cmds[i]
		userId := userIds[i]
		result := cmd.(*redis.StringCmd)
		m := new(TinySess)
		m.serverKey = result.String()
		m.userId = userId
	}
	return ms, nil
}

func onNewMsg(e *levent.NewMsg) {
	ctx := context.Background()
	mmsg, err := getMsgById(ctx, e.ChatId, e.MsgId)
	if err != nil {
		qlog.Error(err, e)
		return
	}

	var memberIds []int64
	if len(e.Ids) == 0 {
		ids, err := getMemberIdsByChatId(ctx, e.ChatId)
		if err != nil {
			qlog.Error(err, e)
			return
		}
		memberIds = ids
	} else {
		memberIds = e.Ids
	}
	ms, err := getSessionIdsByUserIds(ctx, memberIds)
	succeedIds := make(map[int64]bool)
	mm := new(msg.MsgResp)
	mm.MsgId = e.MsgId
	mm.ChatId = e.ChatId
	p := new(constdef.JsonRet)
	p.T = stream.StreamType_NewMsgC2S
	p.Data = mm
	for _, m := range ms {
		// 发送消息
		sessId := strconv.FormatInt(m.userId, 10)
		if sendMsg(ctx, m.serverKey, sessId, p) {
			succeedIds[m.userId] = true
		}
	}

	// 对每个人生成自己独立的会话
	for _, memberId := range memberIds {
		if mmsg.MsgType == stream.MsgType_MsgTypeFalse ||
			mmsg.MsgType == stream.MsgType_MsgTypeCommand {
			continue
		}
		var incr bool
		if memberId != mmsg.SenderId {
			incr = true
		}
		_ = createChat(ctx,
			memberId, mmsg.ChatId, mmsg.ChatType, incr, mmsg.MsgId)
	}

	if mmsg.MsgType != stream.MsgType_MsgTypeFalse &&
		mmsg.MsgType != stream.MsgType_MsgTypeCommand {
		// 标记已读
		updateRead(ctx, mmsg.SenderId, mmsg.MsgId, mmsg.ChatId)
	}

	// 发送push消息
	for _, id := range memberIds {
		if !succeedIds[id] {
			ka.TopicEvent(mq.TopicPush, mq.EPushNewMsg, e)
		}
	}
}

// 消息标记已读
func updateRead(ctx context.Context, userId int64, msgId int64, chatId string) {
	shardName := tablename.ChatlistShard(userId, config.GetEnv())
	usql := "update chat_list set read_last_msg_id=? " +
		"where user_id=? and chat_id=?"
	_, err := dao.GetConn(shardName).Exec(usql, msgId, userId, chatId)
	if err != nil {
		qlog.Ctx(ctx).Error(userId, msgId, chatId, err)
	}
}

func createChat(ctx context.Context,
	userId int64, chatId string, chatType stream.ChatType, incr bool, msgId int64) error {
	shardName := tablename.ChatlistShard(userId, config.GetEnv())
	isql := "insert into chat_list (" +
		"user_id,chat_id,chat_type,created_on,updated_on," +
		"unread_ct,last_msg_id,deleted_on) values " +
		"(?,?,?,unix_timestamp(),unix_timestamp()," +
		"?,?,0) on duplicate key " +
		"update updated_on=unix_timestamp(),unread_ct=%s,last_msg_id=?,deleted_on=0"
	if incr {
		isql = fmt.Sprintf(isql, "unread_ct+1")
	} else {
		isql = fmt.Sprintf(isql, "0")
	}
	var unreadCt int
	if incr {
		unreadCt = 1
	}

	_, err := dao.GetConn(shardName).Exec(isql,
		userId, chatId, chatType,
		unreadCt, msgId,
		msgId)
	if err != nil {
		qlog.Ctx(ctx).Error(userId, chatId, chatType, unreadCt, msgId, err)
	}
	return err
}

func sendMsg(ctx context.Context, serverName, sessId string, packet *constdef.JsonRet) bool {
	creq := new(ws.MsgReq)
	creq.SessId = sessId
	creq.Content, _ = json.Marshal(packet)
	cresp := new(ret.EmptyResp)
	if err := qgrpc.CallWithServerName(ctx, serverName, method.WsMsg, creq, cresp); err != nil {
		qlog.Ctx(ctx).Error(serverName, sessId, packet.T, err)
		return false
	}
	if cresp.GetErr() != nil {
		return false
	}
	return true
}

//TODO 系统消息，之后再完善
func onSysMsg(e *stream.SysMsgModel) {
	//ctx := context.Background()
}
