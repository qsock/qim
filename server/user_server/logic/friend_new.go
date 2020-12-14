package logic

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/store/ka"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/config/mq"
	"github.com/qsock/qim/lib/codes"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/proto/event"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/user"
	"github.com/qsock/qim/lib/util"
	"time"
)

func FriendNewIgnore(ctx context.Context, req *user.FriendNewIgnoreReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	dsql := "update new_apply set `ignore`=1,updated_on=unix_timestamp(),operator_id=? " +
		"where a_id=? and apply_id=? and recver_id=? and apply_type=? and ignore=0"
	result, err := dao.GetConn(dbconfig.DbUser).Exec(dsql, req.UserId,
		req.UserId, req.FriendId, req.UserId, user.NewApplyType_NewApplyUser)
	if err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}
	ka.TopicEvent(mq.TopicEvent, mq.EEventFriendApplyIgnore, event.Friend{OperatorId: req.UserId, RecverId: req.FriendId})
	return resp, nil
}

// 申请添加新好友
func FriendNewApply(ctx context.Context, req *user.FriendNewApplyReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	if req.UserId == req.FriendId {
		//不能添加自己
		resp.Err = codes.Error(codes.ErrorUserCannotAddSelf)
		return resp, nil
	}
	// 是否是好友
	if flag, err := isFriend(ctx, req.UserId, req.FriendId); err != nil {
		return nil, err
	} else if flag {
		resp.Err = codes.Error(codes.ErrorUserIsFriend)
		return resp, nil
	}

	// 检查用户是否被拉黑
	flag, err := isUserInBlackList(ctx, req.FriendId, req.UserId)
	if err != nil {
		return nil, err
	}
	if flag {
		// 已经被拉黑，不用给用户感知，发送就可以了
		resp.Err = codes.Error(codes.ErrorUserApplyDenied)
		return resp, nil
	}

	// 判断用户加好友的设置,任意添加好友，需要请求，禁止添加
	info, err := GetUserInfoById(ctx, req.FriendId)
	if err != nil {
		return nil, err
	}
	if info.AddFriendType == user.FriendAddType_FriendAddNone {
		// 对方禁止添加
		resp.Err = codes.Error(codes.ErrorUserFriendNoneAdd)
		return resp, nil
	} else if info.AddFriendType == user.FriendAddType_FriendAddAnyone {
		if err := becomeFriend(ctx, req.UserId, req.FriendId); err != nil {
			return nil, err
		}
		return resp, nil
	}
	isql := "insert into new_apply (a_id, apply_id, recver_id, apply_type, ct," +
		"`ignore`, created_on, updated_on,`reason`) values " +
		"(?,?,?,?,1," +
		"0,unix_timestamp(),unix_timestamp(),?) on duplicate key " +
		"update updated_on=unix_timestamp(),ct=ct+1,`reason`=?,deleted_on=0"

	{
		var (
			ignore    bool
			updatedOn int64
		)
		// 查看一下新好友列表，不能一个小时内重复申请
		ssql := "select `ignore`,`updated_on` from new_apply where `a_id`=? and `apply_id`=? and `recver_id`=? limit 1"
		if err := dao.GetConn(dbconfig.DbUser).
			QueryRow(ssql, req.FriendId, req.UserId, req.FriendId).Scan(&ignore, &updatedOn); err != nil && err != sql.ErrNoRows {
			qlog.Ctx(ctx).Error(req, err)
			return nil, err
		}
		//一小时内不可以重复申请
		if updatedOn > time.Now().Unix()+360 {
			resp.Err = codes.Error(codes.ErrorUserApplyDenied)
			return resp, nil
		}
		if _, err := dao.GetConn(dbconfig.DbUser).Exec(isql,
			req.UserId, req.UserId, req.FriendId, user.NewApplyType_NewApplyUser,
			req.Reason,
			req.Reason); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			return nil, err
		}

		// 不能被忽略
		if ignore {
			return resp, nil
		}
	}

	if _, err := dao.GetConn(dbconfig.DbUser).Exec(isql,
		req.FriendId, req.UserId, req.FriendId, user.NewApplyType_NewApplyUser,
		req.Reason,
		req.Reason); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}

	ka.TopicEvent(mq.TopicEvent, mq.EEventFriendApply, event.Friend{OperatorId: req.UserId, RecverId: req.FriendId, Content: req.Reason})
	return resp, nil
}

func FriendNewReject(ctx context.Context, req *user.FriendNewRejectReq) (*ret.EmptyResp, error) {
	usql := "update new_apply set status=?, updated_on=unixtimestamp(),operator_id=?,operator_reason=? " +
		"where apply_id=? and recver_id=? and apply_type=? and status=?"
	result, err := dao.GetConn(dbconfig.DbUser).Exec(usql, user.NewApplyStatus_NewApplyRejected, req.UserId, req.Reason,
		req.FriendId, req.UserId, user.NewApplyType_NewApplyUser, user.NewApplyStatus_NewApply)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}

	resp := new(ret.EmptyResp)
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}
	ka.TopicEvent(mq.TopicEvent, mq.EEventFriendApplyReject, event.Friend{OperatorId: req.UserId, RecverId: req.FriendId, Content: req.Reason})
	return resp, nil
}

func becomeFriend(ctx context.Context, userId int64, friendId int64) error {
	isql := "insert into friends (user_id,friend_id,created_on) " +
		"values(?,?,unix_timestamp()),(?,?,unix_timestamp()) "
	_, err := dao.GetConn(dbconfig.DbUser).Exec(isql, userId, friendId, friendId, userId)
	if err != nil {
		qlog.Ctx(ctx).Error(userId, friendId, err)
		return err
	}
	ka.TopicEvent(mq.TopicEvent, mq.EEventFriendAdd, event.Friend{OperatorId: userId, RecverId: friendId})
	return nil
}

func FriendNewAgree(ctx context.Context, req *user.FriendNewAgreeReq) (*ret.EmptyResp, error) {
	usql := "update new_apply set status=?,updated_on=unix_timestamp(),operator_id=?  " +
		"where apply_id=? and recver_id=? and apply_type=? and status=?"
	result, err := dao.GetConn(dbconfig.DbUser).Exec(usql, user.NewApplyStatus_NewApplySucceed, req.UserId,
		req.FriendId, req.UserId, user.NewApplyType_NewApplyUser, user.NewApplyStatus_NewApply)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	resp := new(ret.EmptyResp)
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}

	if err := becomeFriend(ctx, req.UserId, req.FriendId); err != nil {
		return nil, err
	}
	//ka.TopicEvent(mq.TopicEvent, mq.EEventFriendApplyAgree, event.Friend{OperatorId: req.UserId, RecverId: req.FriendId})
	return resp, nil
}

func FriendNewDel(ctx context.Context, req *user.FriendNewDelReq) (*ret.EmptyResp, error) {
	usql := "update new_apply set deleted_on=unix_timestamp() where " +
		"id=? and a_id=?"
	result, err := dao.GetConn(dbconfig.DbUser).Exec(usql,
		req.Id, req.UserId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}

	resp := new(ret.EmptyResp)
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}
	ka.TopicEvent(mq.TopicEvent, mq.EEventFriendApplyDel, event.Friend{OperatorId: req.UserId, Id: req.Id})
	return resp, nil
}

func NewApplyUserList(ctx context.Context, req *user.NewApplyUserListReq) (*user.NewApplyListResp, error) {
	resp := new(user.NewApplyListResp)
	items := make([]*user.NewItem, 0)
	csql := "select count(1) from new_apply " +
		"where a_id=?"
	if err := dao.GetConn(dbconfig.DbUser).QueryRow(csql, req.UserId).Scan(&resp.Total); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	if req.PageSize*req.Page > resp.Total {
		return resp, nil
	}

	ssql := "select id, apply_id, recver_id, apply_type, ct," +
		"`ignore`,created_on,`status`,updated_on, `reason`, " +
		"`operator_id`,`operator_reason` from new_apply " +
		"where a_id=? and deleted_on=0 order by updated_on desc limit %d offset %d"
	ssql = fmt.Sprintf(ssql, req.PageSize, req.Page*req.PageSize)

	uids := make([]int64, 0)
	gids := make([]int64, 0)
	rows, err := dao.GetConn(dbconfig.DbUser).Query(ssql, req.UserId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	for rows.Next() {
		item := new(user.NewItem)
		if err := rows.Scan(&item.Id, &item.ApplyId, &item.RecverId, &item.ApplyType, &item.Ct,
			&item.Ignore, &item.CreatedOn, &item.Status, &item.UpdatedOn, &item.Reason,
			&item.OperatorId, &item.OperatorReason); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			continue
		}
		uids = append(uids, item.ApplyId)
		if item.OperatorId > 0 {
			uids = append(uids, item.OperatorId)
		}
		if item.ApplyType == user.NewApplyType_NewApplyUser {
			uids = append(uids, item.RecverId)
		} else if item.ApplyType == user.NewApplyType_NewApplyGroup {
			gids = append(gids, item.RecverId)
		} else {
			qlog.Ctx(ctx).Error(req, item)
			continue
		}
		items = append(items, item)
	}
	_ = rows.Close()

	uids = util.UniqueInt64s(uids)
	gids = util.UniqueInt64s(gids)

	infos, err := GetUserInfoByIds(ctx, uids)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	groups, err := GetGroupByIds(ctx, gids)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}

	for _, item := range items {
		for _, info := range groups {
			if item.ApplyType == user.NewApplyType_NewApplyGroup && info.Id == item.RecverId {
				item.RecvGroup = info
			}
		}
	}

	for _, item := range items {
		for _, info := range infos {
			if item.ApplyType == user.NewApplyType_NewApplyUser && info.UserId == item.RecverId {
				item.RecvUser = info
			}
			if info.UserId == item.ApplyId {
				item.ApplyUser = info
			}
			if info.UserId == item.OperatorId {
				item.OperatorUser = info
			}
		}
	}

	resp.Items = items
	return resp, nil
}
