package logic

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/config/kvconfig"
	"github.com/qsock/qim/lib/cachename"
	"github.com/qsock/qim/lib/dao"
)

func GetUserSeqId(ctx context.Context, userId int64) (int64, error) {
	kvConn := dao.GetKvConn(kvconfig.KvDefault)
	cacheKey := cachename.RedisPassportSeq(userId)
	seqId, err := kvConn.Get(cacheKey).Int64()
	if seqId > 0 {
		return seqId, nil
	}
	if err != nil && err != redis.Nil {
		qlog.Ctx(ctx).Error(userId)
	}
	seqId, err = getUserSeqIdFromDb(ctx, userId)
	if err != nil {
		return 0, err
	}
	_ = kvConn.Set(cacheKey, seqId, 60*60)
	return seqId, nil
}

func getUserSeqIdFromDb(ctx context.Context, userId int64) (int64, error) {
	var seqId int64
	if err := dao.GetConn(dbconfig.DbPassport).QueryRow("select seq_id from `seq`"+
		" where user_id=? limit 1",
		userId).Scan(&seqId); err != nil {
		return 0, err
	}
	return seqId, nil
}
