package tablename

import (
	"github.com/qsock/qf/encrypt"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/lib/types"
	"strconv"
)

const (
	ImMsgTableCount = 256
)

func getImShardNum(id int64) int64 {
	offset := id % (dbconfig.MsgTotalShard * ImMsgTableCount)
	return offset / ImMsgTableCount
}

func ImMsgShard(chatId string, env string) string {
	if env == types.EnvDev {
		return dbconfig.DbMsgShard0
	}
	// DbMsgShard0
	hashId := encrypt.HashInt(chatId)
	return "DbMsgShard" + strconv.FormatInt(getImShardNum(hashId), 10)
}

func ImMsgShardR1(chatId string, env string) string {
	if env == types.EnvDev {
		return dbconfig.DbMsgShard0
	}
	return ImMsgShard(chatId, env) + "R1"
}

func GetImTableNum(id int64) int64 {
	offset := id % (dbconfig.MsgTotalShard * ImMsgTableCount)
	return offset % ImMsgTableCount
}

func ImMsg(chatId string, env string) string {
	if env == types.EnvDev {
		return "im_msg_0"
	}
	hashId := encrypt.HashInt(chatId)
	return "im_msg_" + strconv.FormatInt(GetImTableNum(hashId), 10)
}

func SysMsgShard(recverId int64, env string) string {
	if env == types.EnvDev {
		return dbconfig.DbMsgShard0
	}
	return "DbMsgShard" + strconv.FormatInt(recverId%dbconfig.MsgTotalShard, 10)
}

func SysMsgShardR1(recverId int64, env string) string {
	if env == types.EnvDev {
		return dbconfig.DbMsgShard0
	}
	return SysMsgShard(recverId, env) + "R1"
}

func ChatlistShard(userId int64, env string) string {
	if env == types.EnvDev {
		return dbconfig.DbMsgShard0
	}
	return "DbMsgShard" + strconv.FormatInt(userId%dbconfig.MsgTotalShard, 10)
}

func ChatlistShardR1(userId int64, env string) string {
	if env == types.EnvDev {
		return dbconfig.DbMsgShard0
	}
	return ChatlistShard(userId, env) + "R1"
}
