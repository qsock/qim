package logic

import (
	"context"
	"database/sql"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/store/ka"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/config/mq"
	"github.com/qsock/qim/lib/codes"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/proto/event"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/user"
	"time"
)

// 申请加入群组
func GroupJoin(ctx context.Context, req *user.GroupJoinReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	{
		cresp, err := IsGroupMember(ctx, &user.IsGroupMemberReq{GroupId: req.GroupId, UserId: req.UserId})
		if err != nil {
			return nil, err
		}
		if cresp.Flag {
			resp.Err = codes.Error(codes.ErrorUserGroupAlreadyIn)
			return resp, nil
		}
	}
	// 是否被blocked
	{
		memInfo, err := getMemberById(ctx, req.UserId, req.GroupId)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if memInfo.GetUserId() > 0 && memInfo.IsBlocked {
			resp.Err = codes.Error(codes.ErrorUserGroupBeenBlock)
			return resp, nil
		}
	}

	info, err := GetGroupById(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}
	if info.CurrentCt >= info.MaxMemberCt {
		resp.Err = codes.Error(codes.ErrorUserGroupMaxMember)
		return resp, nil
	}

	if info.JoinType == user.GroupJoinType_GroupJoinNone {
		// 对方禁止添加
		resp.Err = codes.Error(codes.ErrorUserGroupForbidJoin)
		return resp, nil
	} else if info.JoinType == user.GroupJoinType_GroupJoinAnyone {
		if err := becomeGroupMember(ctx, req.GroupId, req.UserId); err != nil {
			return nil, err
		}
		// 加入新的群组
		ka.TopicEvent(mq.TopicEvent, mq.EEventGroupNewMember, event.GroupMember{OperatorId: 0,
			GroupId: req.GroupId, UserId: req.UserId})
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

		ssql := "select `ignore`,`updated_on` from new_apply where `a_id`=? and `apply_id`=? and `recver_id`=? limit 1"
		if err := dao.GetConn(dbconfig.DbUser).
			QueryRow(ssql, req.GroupId, req.UserId, req.GroupId).Scan(&ignore, &updatedOn); err != nil && err != sql.ErrNoRows {
			qlog.Ctx(ctx).Error(req, err)
			return nil, err
		}
		//一小时内不可以重复申请
		if updatedOn > time.Now().Unix()+360 {
			resp.Err = codes.Error(codes.ErrorUserApplyDenied)
			return resp, nil
		}

		if _, err := dao.GetConn(dbconfig.DbUser).Exec(isql,
			req.UserId, req.UserId, req.GroupId, user.NewApplyType_NewApplyGroup,
			req.Reason,
			req.Reason); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			return nil, err
		}

		// 如果被忽略,就只更新自己
		if ignore {
			return resp, nil
		}
	}

	if _, err := dao.GetConn(dbconfig.DbUser).Exec(isql,
		req.GroupId, req.UserId, req.GroupId, user.NewApplyType_NewApplyGroup,
		req.Reason,
		req.Reason); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}

	ka.TopicEvent(mq.TopicEvent, mq.EEventGroupApply, event.GroupNewApply{OperatorId: req.UserId,
		RecverId: req.GroupId, GroupId: req.GroupId, Content: req.Reason})
	return resp, nil
}

// 同意申请
func GroupJoinAgree(ctx context.Context, req *user.GroupJoinAgreeReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	info, err := GetGroupById(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}
	if info.CurrentCt >= info.MaxMemberCt {
		resp.Err = codes.Error(codes.ErrorUserGroupMaxMember)
		return resp, nil
	}

	memberInfo, err := getMemberById(ctx, req.UserId, req.GroupId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	if memberInfo.RoleType == user.GroupRoleType_GroupRoleNormal {
		resp.Err = codes.Error(codes.ErrorUserGroupHasNoRight)
		return resp, nil
	}

	usql := "update new_apply set status=?,updated_on=unix_timestamp(),operator_id=? " +
		"where apply_id=? and recver_id=? and apply_type=? and status=?"
	result, err := dao.GetConn(dbconfig.DbUser).Exec(usql, user.NewApplyStatus_NewApplySucceed, req.UserId,
		req.GroupId, req.UserId, user.NewApplyType_NewApplyGroup, user.NewApplyStatus_NewApply)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}

	if err := becomeGroupMember(ctx, req.GroupId, req.MemberId); err != nil {
		return nil, err
	}

	// 加入新的群组
	ka.TopicEvent(mq.TopicEvent, mq.EEventGroupNewMember, event.GroupMember{OperatorId: req.UserId,
		GroupId: req.GroupId, UserId: req.MemberId})
	return resp, nil
}

// 拒绝申请
func GroupJoinReject(ctx context.Context, req *user.GroupJoinRejectReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	memberInfo, err := getMemberById(ctx, req.UserId, req.GroupId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	if memberInfo.RoleType == user.GroupRoleType_GroupRoleNormal {
		resp.Err = codes.Error(codes.ErrorUserGroupHasNoRight)
		return resp, nil
	}

	usql := "update new_apply, updated_on=unixtimestamp(),operator_id=?,operator_reason=? " +
		"where apply_id=? and recver_id=? and apply_type=? and status=?"

	result, err := dao.GetConn(dbconfig.DbUser).Exec(usql, user.NewApplyStatus_NewApplyRejected, req.UserId, req.Reason,
		req.GroupId, req.UserId, user.NewApplyType_NewApplyGroup, user.NewApplyStatus_NewApply)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}
	ka.TopicEvent(mq.TopicEvent, mq.EEventGroupApplyReject, event.GroupNewApply{OperatorId: req.UserId,
		RecverId: req.MemberId, GroupId: req.GroupId, Content: req.Reason})

	return resp, nil
}

// 忽略申请
func GroupJoinIgnore(ctx context.Context, req *user.GroupJoinIgnoreReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	dsql := "update new_apply set `ignore`=1,updated_on=unix_timestamp() ,operator_id=?" +
		"where a_id=? and apply_id=? and recver_id=? and apply_type=? and ignore=0"
	result, err := dao.GetConn(dbconfig.DbUser).Exec(dsql, req.UserId,
		req.GroupId, req.MemberId, req.UserId, user.NewApplyType_NewApplyUser)
	if err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}

	ka.TopicEvent(mq.TopicEvent, mq.EEventGroupApplyIgnore, event.GroupNewApply{OperatorId: req.UserId,
		RecverId: req.MemberId, GroupId: req.GroupId})
	return resp, nil
}
