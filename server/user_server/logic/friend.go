package logic

import (
	"context"
	"fmt"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/store/ka"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/config/mq"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/proto/event"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/user"
	"github.com/qsock/qim/lib/util"
)

func FriendIds(ctx context.Context, req *user.FriendIdsReq) (*user.FriendIdsResp, error) {
	ssql := "select friend_id from friends where user_id=?"
	rows, err := dao.GetConn(dbconfig.DbUser).Query(ssql, req.UserId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	ids := make([]int64, 0)
	defer rows.Close()
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			continue
		}
		ids = append(ids, id)
	}
	resp := new(user.FriendIdsResp)
	resp.Ids = ids
	return resp, nil
}

func FriendByIds(ctx context.Context, req *user.FriendByIdsReq) (*user.FriendByIdsResp, error) {
	infos, err := GetUserInfoByIds(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	finfos := make([]*user.FriendItem, 0)
	for _, info := range infos {
		m := new(user.FriendItem)
		m.User = info
		m.UserId = info.UserId
		finfos = append(finfos, m)
	}
	resp := new(user.FriendByIdsResp)
	ssql := "select created_on,friend_id,markname from friends " +
		"where user_id=? and friend_id in (?)"
	rows, err := dao.GetConn(dbconfig.DbUser).Query(ssql, req.UserId, util.Int64sToStr(req.Ids))
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var (
			createdOn int64
			friendId  int64
			markName  string
		)
		if err := rows.Scan(&createdOn, &friendId, &markName); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			continue
		}
		for _, info := range finfos {
			if info.User.UserId == friendId {
				info.MarkName = markName
				info.FriendTime = createdOn
			}
		}
	}
	resp.Items = finfos
	return resp, nil
}

func FriendsByUid(ctx context.Context, req *user.FriendsByUidReq) (*user.FriendsByUidResp, error) {
	resp := new(user.FriendsByUidResp)
	csql := "select count(1) from friends where user_id=?"

	if err := dao.GetConn(dbconfig.DbUser).QueryRow(csql, req.UserId).Scan(&resp.Total); err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	if resp.Total <= req.PageSize*req.Page {
		return resp, nil
	}

	items := make([]*user.FriendItem, 0)
	ids := make([]int64, 0)
	ssql := "select created_on,friend_id,markname from friends " +
		"where user_id=? limit %d offset %d"
	ssql = fmt.Sprintf(ssql, req.PageSize, req.PageSize*req.Page)
	rows, err := dao.GetConn(dbconfig.DbUser).Query(ssql, req.UserId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		item := new(user.FriendItem)
		if err := rows.Scan(&item.FriendTime, &item.UserId, &item.MarkName); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			continue
		}
		ids = append(ids, item.UserId)
		items = append(items, item)
	}
	infos, err := GetUserInfoByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		for _, info := range infos {
			if item.UserId == info.UserId {
				item.User = info
			}
		}
	}
	resp.Items = items
	return resp, nil
}

// 删除好友
func FriendDel(ctx context.Context, req *user.FriendDelReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	dsql := "delete from friends " +
		"where (user_id=? and friend_id=?) or (user_id=? and friend_id=?)"
	result, err := dao.GetConn(dbconfig.DbUser).Exec(dsql, req.UserId, req.FriendId, req.FriendId, req.UserId)
	if err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}
	ka.TopicEvent(mq.TopicEvent, mq.EEventFriendDel, event.Friend{OperatorId: req.UserId, RecverId: req.FriendId})
	return resp, nil
}

func FriendMarknameUpdate(ctx context.Context, req *user.FriendMarknameUpdateReq) (*ret.EmptyResp, error) {
	usql := "update friends set markname=? where user_id=? and friend_id=?"
	result, err := dao.GetConn(dbconfig.DbUser).Exec(usql, req.Name, req.UserId, req.FriendId)
	if err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	resp := new(ret.EmptyResp)
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}
	ka.TopicEvent(mq.TopicEvent, mq.EEventFriendMarkname, event.Friend{OperatorId: req.UserId, RecverId: req.FriendId, Content: req.Name})
	return resp, nil
}

func isFriend(ctx context.Context, userId, friendId int64) (bool, error) {
	ssql := "select count(1) from friends where user_id=? and friend_id=? limit 1"
	var flag bool
	if err := dao.GetConn(dbconfig.DbUser).QueryRow(ssql, userId, friendId).Scan(&flag); err != nil {
		qlog.Ctx(ctx).Error(userId, friendId, err)
		return true, err
	}
	return flag, nil
}

// 好友相关
func IsFriend(ctx context.Context, req *user.IsFriendReq) (*ret.BoolResp, error) {
	flag, err := isFriend(ctx, req.UserId, req.FriendId)
	if err != nil {
		return nil, err
	}
	resp := new(ret.BoolResp)
	resp.Flag = flag
	return resp, nil
}
