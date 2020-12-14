package logic

import (
	"context"
	"fmt"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/user"
)

func isUserInBlackList(ctx context.Context, userId, friendId int64) (bool, error) {
	ssql := "select count(1) from blacklist where user_id=? and black_user_id=? limit 1"
	var flag bool
	if err := dao.GetConn(dbconfig.DbUser).QueryRow(ssql, userId, friendId).Scan(&flag); err != nil {
		qlog.Ctx(ctx).Error(userId, friendId, err)
		return true, err
	}
	return flag, nil
}

// 黑名单部分
func BlacklistAdd(ctx context.Context, req *user.BlacklistAddReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	// 删除好友，并且加入黑名单
	_, err := FriendDel(ctx, &user.FriendDelReq{UserId: req.UserId, FriendId: req.BlackUserId})
	if err != nil {
		return nil, err
	}
	isql := "insert ignore into blacklist (user_id,black_user_id,created_on) values(?,?,unix_timestamp())"
	_, err = dao.GetConn(dbconfig.DbUser).Exec(isql, req.UserId, req.BlackUserId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	return resp, nil
}

func BlacklistDel(ctx context.Context, req *user.BlacklistDelReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	dsql := "delete from blacklist where user_id=?,black_user_id=?"
	_, err := dao.GetConn(dbconfig.DbUser).Exec(dsql, req.UserId, req.BlackUserId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	return resp, nil
}

func Blacklist(ctx context.Context, req *user.BlacklistReq) (*user.BlacklistResp, error) {
	resp := new(user.BlacklistResp)
	csql := "select count(1) from blacklist where user_id=?"
	if err := dao.GetConn(dbconfig.DbUser).QueryRow(csql, req.UserId).Scan(&resp.Total); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	if resp.Total < req.PageSize*req.Page {
		return resp, nil
	}

	ssql := "select black_user_id from blacklist where user_id=? order by created_on desc limit %d offset %d"
	ssql = fmt.Sprintf(ssql, req.PageSize, req.PageSize*req.Page)
	rows, err := dao.GetConn(dbconfig.DbUser).Query(ssql, req.UserId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			continue
		}
		ids = append(ids, id)
	}
	_ = rows.Close()
	infos, err := GetUserInfoByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	resp.Users = infos
	return resp, nil
}
