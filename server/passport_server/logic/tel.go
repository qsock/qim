package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qf/util/coderand"
	"github.com/qsock/qf/util/uuid"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/config/kvconfig"
	"github.com/qsock/qim/lib/cachename"
	"github.com/qsock/qim/lib/codes"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/gclient"
	"github.com/qsock/qim/lib/method"
	"github.com/qsock/qim/lib/proto/file"
	"github.com/qsock/qim/lib/proto/model"
	"github.com/qsock/qim/lib/proto/passport"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/user"
	"github.com/qsock/qim/lib/util"
	"github.com/qsock/qim/server/passport_server/config"
)

func TelLogin(ctx context.Context, req *passport.TelLoginReq) (*passport.LoginResp, error) {
	resp := new(passport.LoginResp)
	cacheKey := cachename.RedisPassportSms(req.Tel)
	kvconn := dao.GetKvConn(kvconfig.KvDefault)
	result := kvconn.Get(cacheKey).Val()
	if len(result) == 0 {
		resp.Err = codes.Error(codes.ErrorPassportSmsCode)
		return resp, nil
	}
	smsM := new(passport.SmsModel)
	if err := json.Unmarshal([]byte(result), smsM); err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	if smsM.Code != req.Code {
		resp.Err = codes.Error(codes.ErrorPassportSmsCode)
		return resp, nil
	}
	clearSmsCache(ctx, req.Tel)
	// 查看手机是否存在
	{
		var userId int64
		ssql := "select user_id from mobile where tel=? limit 1"
		if err := dao.GetConn(dbconfig.DbPassport).QueryRow(ssql, req.Tel).Scan(&userId); err != nil && err != sql.ErrNoRows {
			qlog.Ctx(ctx).Error(err, req)
			return nil, err
		}
		// 如果已存在
		if userId > 0 {
			return doLogin(ctx, userId, req.Meta)
		}
	}

	return doTelRegister(ctx, req.Tel, req.Meta)
}

func doSeq(ctx context.Context, userId int64) error {
	isql := "insert into `seq` (seq_id,user_id) values(1,?)"
	_, err := dao.GetConn(dbconfig.DbPassport).Exec(isql, userId)
	if err != nil {
		qlog.Ctx(ctx).Error(err, userId)
		return err
	}
	return nil
}

func doTelRegister(ctx context.Context, tel string, meta *model.RequestMeta) (*passport.LoginResp, error) {
	// 创建id
	userId := gclient.GenDbId(ctx, config.GetConfig().IdKey)
	avatar, err := getRandomAvatar(ctx)
	if err != nil {
		qlog.Ctx(ctx).Error(err, tel, meta)
		return nil, err
	}
	name := util.HideStar(tel)

	// 先去创建用户，多次创建不用更新
	{
		creq := new(user.CreateReq)
		creq.UserId = userId
		creq.Avatar = avatar
		creq.Name = name
		if err := createUser(ctx, creq); err != nil {
			qlog.Ctx(ctx).Error(err, creq)
			return nil, err
		}
	}
	isql := "insert into mobile (user_id,tel) values(?,?)"
	if _, err := dao.GetConn(dbconfig.DbPassport).Exec(isql, userId, tel); err != nil {
		qlog.Ctx(ctx).Error(err, tel, meta)
		return nil, err
	}
	_ = doSeq(ctx, userId)
	return doLogin(ctx, userId, meta)
}

func doLogin(ctx context.Context, userId int64, meta *model.RequestMeta) (*passport.LoginResp, error) {
	resp := new(passport.LoginResp)

	cacheKey := cachename.RedisUserBan(userId)
	if dao.GetKvConn(kvconfig.KvDefault).Exists(cacheKey).Val() > 0 {
		resp.Err = codes.Error(codes.ErrorPassportBanned)
		return resp, nil
	}

	refreshToken := uuid.NewString()
	token, err := GenToken(ctx, userId, meta.Device, meta.UserIp)
	if err != nil {
		qlog.Ctx(ctx).Error(err, meta, userId)
		return nil, err
	}

	auth := new(model.UserAuth)
	auth.UserId = userId
	auth.Token = token
	auth.RefreshToken = refreshToken
	resp.Auth = auth
	saveRefreshToken(ctx, refreshToken, userId, meta.Device)
	return resp, nil
}

func clearSmsCache(ctx context.Context, tel string) error {
	cacheKey := cachename.RedisPassportSms(tel)
	if err := dao.GetKvConn(kvconfig.KvDefault).Del(cacheKey).Err(); err != nil {
		qlog.Ctx(ctx).Error(err, tel)
		return err
	}
	return nil
}

func saveRefreshToken(ctx context.Context, token string, userId int64, device model.Device) error {
	cacheKey := cachename.RedisPassportRefreshToken(userId)
	if err := dao.GetKvConn(kvconfig.KvDefault).HSet(cacheKey, device.String(), token).Err(); err != nil {
		qlog.Ctx(ctx).Error(err, token, userId, device)
		return err
	}
	return nil
}

func createUser(ctx context.Context, creq *user.CreateReq) error {
	cresp := new(ret.EmptyResp)
	if err := qgrpc.Call(ctx, method.UserCreate, creq, cresp); err != nil {
		qlog.Ctx(ctx).Error(creq, err)
		return err
	}
	return nil
}

func getRandomAvatar(ctx context.Context) (string, error) {
	creq := new(file.GetSysAvatarsReq)
	cresp := new(file.GetSysAvatarsResp)

	if err := qgrpc.Call(ctx, method.FileGetSysAvatars, creq, cresp); err != nil {
		qlog.Ctx(ctx).Error(creq, err)
		return "", err
	}
	if len(cresp.Avatars) == 0 {
		return "", nil
	}
	idx := coderand.Uint32(uint32(len(cresp.Avatars)))
	return cresp.Avatars[idx], nil
}
