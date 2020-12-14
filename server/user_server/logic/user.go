package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/config/kvconfig"
	"github.com/qsock/qim/lib/cachename"
	"github.com/qsock/qim/lib/codes"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/proto/model"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/user"
	"github.com/qsock/qim/lib/tablename"
	"github.com/qsock/qim/lib/util"
	"github.com/qsock/qim/server/user_server/config"
	"time"
)

func Create(ctx context.Context, req *user.CreateReq) (*ret.EmptyResp, error) {
	isql := "insert into userinfo (user_id, nickname, avatar, gender, created_on, updated_on) " +
		"values(?,?,?,?,unix_timestamp(),unix_timestamp())"
	_, err := dao.GetConn(dbconfig.DbUser).Exec(isql,
		req.UserId, req.Name, req.Avatar, req.Gender)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	resp := new(ret.EmptyResp)
	return resp, nil
}

func Lastactive(ctx context.Context, req *user.LastactiveReq) (*ret.EmptyResp, error) {
	tname := tablename.UserLastactive(req.Meta.UserId, config.GetEnv())
	isql := "insert into " + tname +
		" (user_id,lat,lng,ip,device," +
		"device_id,app_ver,created_on) values (?,?,?,?,?," +
		"?,?,unix_timestamp())"
	_, err := dao.GetConn(dbconfig.DbUser).Exec(isql, req.Meta.UserId, req.Meta.Lat, req.Meta.Lng, req.Meta.UserIp, req.Meta.Device,
		req.Meta.DeviceId, req.Meta.AppVersion)
	if err != nil {
		qlog.Ctx(ctx).Error(req.Meta, err)
		return nil, err
	}

	resp := new(ret.EmptyResp)
	return resp, nil
}

func Update(ctx context.Context, req *user.UpdateReq) (*ret.EmptyResp, error) {
	usql := "update userinfo set updated_on=unix_timestamp(),"
	conds := make([]interface{}, 0)
	if len(req.Name) > 0 {
		usql += "nickname=?,"
		conds = append(conds, req.Name)
	}
	if len(req.Avatar) > 0 {
		usql += "avatar=?,"
		conds = append(conds, req.Avatar)
	}

	if req.Gender != model.Gender_GenderUnknown {
		usql += "gender=?,"
		conds = append(conds, req.Gender)
	}

	if req.Birthday > 0 {
		usql += "birthday=?,"
		conds = append(conds, req.Birthday)
	}

	if len(req.Brief) > 0 {
		usql += "brief=?,"
		conds = append(conds, req.Brief)
	}
	if req.AddFriendType != user.FriendAddType_FriendAddTypeFalse {
		usql += "add_friend_type=? "
		conds = append(conds, req.AddFriendType)
	}

	resp := new(ret.EmptyResp)
	if len(conds) == 0 {
		resp.Err = codes.Error(codes.ErrorParameter)
		return resp, nil
	}
	usql = usql[:len(usql)-1]
	usql += " where user_id=? "
	conds = append(conds, req.UserId)
	_, err := dao.GetConn(dbconfig.DbUser).Exec(usql, conds...)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}

	_ = clearUserCache(ctx, req.UserId)
	return resp, nil
}

func clearUserCache(ctx context.Context, userId int64) error {
	key := cachename.RedisUserInfo(userId)
	if err := dao.GetKvConn(kvconfig.KvDefault).Del(key).Err(); err != nil {
		qlog.Ctx(ctx).Error(userId, key, err)
		return err
	}
	return nil
}

func GetUserInfoById(ctx context.Context, id int64) (*user.UserInfo, error) {
	infos, err := GetUserInfoByIds(ctx, []int64{id})
	if err != nil {
		return nil, err
	}
	if len(infos) == 0 {
		return nil, nil
	}
	return infos[0], nil
}

func GetUserInfoByIds(ctx context.Context, ids []int64) ([]*user.UserInfo, error) {
	infos, _ := getUserInfoByIdsOnCache(ctx, ids)
	if len(infos) == len(ids) {
		return infos, nil
	}

	if infos == nil {
		infos = make([]*user.UserInfo, 0)
	}

	sids := make([]int64, 0)
	for _, id := range ids {
		var flag bool
		for _, info := range infos {
			if info.UserId == id {
				flag = true
			}
		}
		if !flag {
			sids = append(sids, id)
		}
	}

	sinfos, err := getUserInfoByIdsOnDb(ctx, sids)
	if err != nil {
		return nil, err
	}
	pipl := dao.GetKvConn(kvconfig.KvDefault).Pipeline()
	for _, info := range sinfos {
		key := cachename.RedisUserInfo(info.UserId)
		b, _ := json.Marshal(info)
		pipl.Set(key, string(b), time.Second*86400)
	}
	if _, err := pipl.Exec(); err != nil {
		qlog.Ctx(ctx).Error(sids, err)
	}
	infos = append(infos, sinfos...)
	return infos, nil
}

func getUserInfoByIdsOnCache(ctx context.Context, ids []int64) ([]*user.UserInfo, error) {
	pipl := dao.GetKvConn(kvconfig.KvDefault).Pipeline()
	for _, id := range ids {
		key := cachename.RedisUserInfo(id)
		pipl.Get(key)
	}
	cmds, err := pipl.Exec()
	if err != nil {
		qlog.Ctx(ctx).Error(ids, err)
		return nil, err
	}
	infos := make([]*user.UserInfo, 0)
	for _, cmd := range cmds {
		val := cmd.(*redis.StringCmd)
		if len(val.Val()) == 0 {
			continue
		}
		userInfo := new(user.UserInfo)
		if err := json.Unmarshal([]byte(val.Val()), userInfo); err != nil {
			qlog.Ctx(ctx).Error(ids, cmd, val.Val(), err)
			continue
		}
		if userInfo.GetUserId() > 0 {
			infos = append(infos, userInfo)
		}
	}
	return infos, nil
}

func getUserInfoByIdsOnDb(ctx context.Context, ids []int64) ([]*user.UserInfo, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	ssql := "select user_id,nickname,gender,avatar,birthday," +
		"brief,add_friend_type from userinfo where user_id in (%s)"
	ssql = fmt.Sprintf(ssql, util.Int64sToStr(ids))
	rows, err := dao.GetConn(dbconfig.DbUser).Query(ssql)
	if err != nil {
		qlog.Ctx(ctx).Error(ids, err)
		return nil, err
	}
	infos := make([]*user.UserInfo, 0)
	defer rows.Close()
	for rows.Next() {
		userInfo := new(user.UserInfo)
		if err := rows.Scan(&userInfo.UserId, &userInfo.Name, &userInfo.Gender, &userInfo.Avatar, &userInfo.Birthday,
			&userInfo.Brief, &userInfo.AddFriendType); err != nil {
			qlog.Ctx(ctx).Error(ids, err)
			continue
		}
		infos = append(infos, userInfo)
	}
	return infos, nil
}

func Infos(ctx context.Context, req *user.InfosReq) (*user.InfosResp, error) {
	resp := new(user.InfosResp)
	infos, err := GetUserInfoByIds(ctx, req.UserIds)
	if err != nil {
		return nil, err
	}
	resp.Users = infos
	return resp, nil
}

func Info(ctx context.Context, req *user.InfoReq) (*user.InfoResp, error) {
	resp := new(user.InfoResp)
	info, err := GetUserInfoById(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	resp.User = info
	return resp, nil
}
