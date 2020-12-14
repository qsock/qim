package logic

import (
	"context"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/config/kvconfig"
	"github.com/qsock/qim/lib/cachename"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/proto/passport"
	"github.com/qsock/qim/lib/proto/ret"
	"time"
)

func Ban(ctx context.Context, req *passport.BanReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	isql := "insert into ban(user_id,created_on,end_on) values(?,unix_timestamp(),?) on duplicate key update " +
		"end_on=?"
	_, err := dao.GetConn(dbconfig.DbPassport).Exec(isql, req.UserId, req.EndOn, req.EndOn)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	cacheKey := cachename.RedisUserBan(req.UserId)
	if err := dao.GetKvConn(kvconfig.KvDefault).Set(cacheKey, req.EndOn,
		time.Duration(req.EndOn-time.Now().Unix())*time.Second).Err(); err != nil {
		qlog.Ctx(ctx).Error(req, err)
	}
	return resp, nil
}

func UnBan(ctx context.Context, req *passport.UnBanReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)

	cacheKey := cachename.RedisUserBan(req.UserId)
	if err := dao.GetKvConn(kvconfig.KvDefault).Del(cacheKey).Err(); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	dsql := "delete from ban where user_id=?"
	_, err := dao.GetConn(dbconfig.DbPassport).Exec(dsql, req.UserId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	return resp, nil
}
