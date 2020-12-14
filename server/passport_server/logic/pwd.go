package logic

import (
	"context"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/config/kvconfig"
	"github.com/qsock/qim/lib/cachename"
	"github.com/qsock/qim/lib/codes"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/proto/passport"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/qjwt"
)

func Logout(ctx context.Context, req *passport.LogoutReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	usql := "update `seq` set seq_id = seq_id+1 where user_id=?"
	if _, err := dao.GetConn(dbconfig.DbPassport).Exec(usql, req.Meta.UserId); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	clearCache(ctx, req.Meta.UserId)
	return resp, nil
}

func clearCache(ctx context.Context, userId int64) {
	cacheKey := cachename.RedisPassportRefreshToken(userId)
	seqKey := cachename.RedisPassportSeq(userId)

	if err := dao.GetKvConn(kvconfig.KvDefault).Del(cacheKey, seqKey).Err(); err != nil {
		qlog.Ctx(ctx).Error(userId, cacheKey, seqKey, err)
	}
}

func Auth(ctx context.Context, req *passport.AuthReq) (*ret.IntResp, error) {
	resp := new(ret.IntResp)
	val, err := qjwt.Parse(ctx, req.Token)
	if err != nil {
		qlog.Ctx(ctx).Error(req.Meta, err)
		if err == qjwt.ErrTokenExpired {
			resp.Err = codes.Error(codes.ErrorAuthExpired)
			return resp, nil
		} else if err == qjwt.ErrTokenInvalid {
			resp.Err = codes.Error(codes.ErrorAuthInvalid)
			return resp, nil
		}
		return nil, err
	}
	m := new(passport.JwtClaims)
	if err := m.Unmarshal(val); err != nil {
		qlog.Ctx(ctx).Error(req.Meta, err)
		resp.Err = codes.Error(codes.ErrorAuthInvalid)
		return resp, nil
	}
	//check seq
	seqId, err := GetUserSeqId(ctx, m.UserId)
	if err != nil {
		qlog.Ctx(ctx).Error(req.Meta, err)
		return nil, err
	}
	if m.SeqId != seqId {
		resp.Err = codes.Error(codes.ErrorAuthExpired)
		return resp, nil
	}
	cacheKey := cachename.RedisUserBan(m.UserId)
	if dao.GetKvConn(kvconfig.KvDefault).Exists(cacheKey).Val() > 0 {
		resp.Err = codes.Error(codes.ErrorPassportBanned)
		return resp, nil
	}
	resp.Val = m.UserId
	return resp, nil
}

func Refresh(ctx context.Context, req *passport.RefreshReq) (*ret.StringResp, error) {
	resp := new(ret.StringResp)
	val, err := qjwt.ParseTokenWithoutTime(ctx, req.Token)
	if err != nil {
		resp.Err = codes.Error(codes.ErrorAuthInvalid)
		return resp, nil
	}

	m := new(passport.JwtClaims)
	if err := m.Unmarshal(val); err != nil {
		qlog.Ctx(ctx).Error(req.Meta, err)
		resp.Err = codes.Error(codes.ErrorAuthInvalid)
		return resp, nil
	}

	cacheKey := cachename.RedisPassportRefreshToken(m.UserId)
	refreshToken := dao.GetKvConn(kvconfig.KvDefault).HGet(cacheKey, req.Meta.Device.String()).Val()
	if refreshToken != req.RefreshToken {
		qlog.Ctx(ctx).Error(req.Meta, refreshToken, m, cacheKey, req.Meta.Device.String())
		resp.Err = codes.Error(codes.ErrorAuthInvalid)
		return resp, nil
	}

	newToken, err := GenToken(ctx, m.UserId, req.Meta.Device, req.Meta.UserIp)
	if err != nil {
		qlog.Ctx(ctx).Error(req.Meta, err)
		return nil, err
	}
	resp.Str = newToken
	return resp, nil
}
