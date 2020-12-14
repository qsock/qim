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
	"github.com/qsock/qim/lib/tablename"
	"github.com/qsock/qim/lib/util"
	"github.com/qsock/qim/server/user_server/config"
	"strings"
	"time"
)

func becomeGroupMember(ctx context.Context, groupId int64, userId int64) error {
	return becomeGroupMembers(ctx, groupId, []int64{userId})
}

func becomeGroupMembers(ctx context.Context, groupId int64, userIds []int64) error {
	info, err := GetGroupById(ctx, groupId)
	if err != nil {
		return err
	}
	if info.CurrentCt < 9 && len(info.Avatars) > 1 {
		uids := append([]int64{}, userIds...)
		if len(uids) > 7 {
			uids = uids[:7]
		}
		userInfos, err := GetUserInfoByIds(ctx, uids)
		if err != nil {
			return err
		}
		for _, id := range uids {
			if len(info.Avatars) >= 9 {
				break
			}
			for _, user := range userInfos {
				if id == user.UserId {
					info.Avatars = append(info.Avatars, user.Avatar)
				}
			}
		}
	}
	userIds = util.UniqueInt64s(userIds)
	tx, err := dao.GetConn(dbconfig.DbUser).Begin()
	if err != nil {
		qlog.Ctx(ctx).Error(err, groupId, userIds)
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	{
		isql := "insert into user_group (user_id,group_id) values "
		condis := make([]interface{}, 0)
		for _, uid := range userIds {
			isql += "(?,?),"
			condis = append(condis, uid, groupId)
		}
		isql = isql[:len(isql)-1]
		_, err = tx.Exec(isql, condis...)
		if err != nil {
			qlog.Ctx(ctx).Error(err, groupId, userIds)
			return err
		}
	}

	{
		isql := "insert into " + tablename.UserGroupMember(groupId, config.GetEnv()) +
			" (user_id, group_id, created_on) values "
		condis := make([]interface{}, 0)
		for _, uid := range userIds {
			isql += "(?,?,unix_timestamp()),"
			condis = append(condis, uid, groupId)
		}
		isql = isql[:len(isql)-1]
		_, err = tx.Exec(isql, condis...)
		if err != nil {
			qlog.Ctx(ctx).Error(err, groupId, userIds)
			return err
		}
	}

	{
		var ct = len(userIds)
		usql := "update group set current_ct=current_ct+%d,avatar=? where id=?"
		usql = fmt.Sprintf(usql, ct)
		_, err = tx.Exec(usql, groupId, strings.Join(info.Avatars, ","))
		if err != nil {
			qlog.Ctx(ctx).Error(err, groupId, userIds)
			return err
		}
	}

	ClearGroupCache(ctx, groupId)
	return nil
}

// 是否是群组成员
func IsGroupMember(ctx context.Context, req *user.IsGroupMemberReq) (*ret.BoolResp, error) {
	ssql := "select count(1) from user_group where user_id=? and group_id=? limit 1"
	resp := new(ret.BoolResp)
	if err := dao.GetConn(dbconfig.DbUser).QueryRow(ssql, req.UserId, req.GroupId).Scan(&resp.Flag); err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	return resp, nil
}

func getMemberById(ctx context.Context, userId, groupId int64) (*user.GroupMember, error) {
	infos, err := getMemberByIds(ctx, []int64{userId}, groupId)
	if err != nil {
		return nil, err
	}
	if len(infos) == 0 {
		return nil, sql.ErrNoRows
	}
	return infos[0], nil
}

func getMemberByIds(ctx context.Context, userIds []int64, groupId int64) ([]*user.GroupMember, error) {
	if len(userIds) == 0 {
		return nil, nil
	}
	mems := make([]*user.GroupMember, 0)
	ids := make([]int64, 0)
	tname := tablename.UserGroupMember(groupId, config.GetEnv())
	ssql := "select user_id,user_role,mark_name,mute_until,not_disturb,created_on from " + tname +
		" where group_id=? and user_id in (?)"
	rows, err := dao.GetConn(dbconfig.DbUser).Query(ssql, groupId, util.Int64sToStr(userIds))
	if err != nil {
		qlog.Ctx(ctx).Error(userIds, groupId, err)
		return nil, err
	}
	for rows.Next() {
		m := new(user.GroupMember)
		if err := rows.Scan(&m.UserId, &m.RoleType, &m.MarkName, &m.MuteUntil, &m.NotDisturb, &m.CreatedOn); err != nil {
			qlog.Ctx(ctx).Error(userIds, groupId, err)
			continue
		}
		ids = append(ids, m.UserId)
		mems = append(mems, m)
	}
	rows.Close()
	userInfos, err := GetUserInfoByIds(ctx, ids)
	if err != nil {
		qlog.Ctx(ctx).Error(userIds, groupId, err)
		return nil, err
	}
	for _, m := range mems {
		for _, info := range userInfos {
			if info.UserId == m.UserId {
				m.User = info
			}
		}
	}
	return mems, nil
}

// 是否被禁言了
func IsGroupMemberBeenMute(ctx context.Context, req *user.IsGroupMemberBeenMuteReq) (*ret.BoolResp, error) {
	resp := new(ret.BoolResp)
	member, err := getMemberById(ctx, req.UserId, req.GroupId)
	if err != nil {
		return nil, err
	}
	if member.RoleType == user.GroupRoleType_GroupRoleOwner ||
		member.RoleType == user.GroupRoleType_GroupRoleManager {
		resp.Flag = false
		return resp, nil
	}
	if member.MuteUntil > time.Now().Unix() {
		resp.Flag = true
		return resp, nil
	}
	groupInfo, err := GetGroupById(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}
	if groupInfo.MuteUtil > time.Now().Unix() {
		resp.Flag = true
		return resp, nil
	}
	return resp, nil
}

func GroupMemberIds(ctx context.Context, req *user.GroupMemberIdsReq) (*user.GroupMemberIdsResp, error) {
	tname := tablename.UserGroupMember(req.GroupId, config.GetEnv())
	ssql := "select user_id from " + tname + " where group_id=?"
	ids := make([]int64, 0)
	resp := new(user.GroupMemberIdsResp)
	rows, err := dao.GetConn(dbconfig.DbUser).Query(ssql, req.GroupId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			continue
		}
		ids = append(ids, id)
	}
	resp.Ids = ids
	return resp, nil
}

func GroupManagerList(ctx context.Context, req *user.GroupManagerListReq) (*user.GroupManagerListResp, error) {
	resp := new(user.GroupManagerListResp)
	ssql := "select user_id from " + tablename.UserGroupMember(req.GroupId, config.GetEnv()) +
		" where group_id=? and user_role in (%s)"
	ids := make([]int64, 0)
	ssql = fmt.Sprintf(ssql,
		util.Int32sToStr([]int32{int32(user.GroupRoleType_GroupRoleManager), int32(user.GroupRoleType_GroupRoleOwner)}))
	rows, err := dao.GetConn(dbconfig.DbUser).Query(ssql)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			continue
		}
		ids = append(ids, id)
	}
	rows.Close()
	mems, err := getMemberByIds(ctx, ids, req.GroupId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	resp.Members = mems
	return resp, nil
}

func GroupMemberByIds(ctx context.Context, req *user.GroupMemberByIdsReq) (*user.GroupMemberByIdsResp, error) {
	resp := new(user.GroupMemberByIdsResp)
	infos, err := getMemberByIds(ctx, req.UserIds, req.GroupId)
	if err != nil {
		return nil, err
	}
	resp.Members = infos
	return resp, nil
}

func GroupMemberAdd(ctx context.Context, req *user.GroupMemberAddReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	info, err := GetGroupById(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}
	if info.JoinType == user.GroupJoinType_GroupJoinNone ||
		info.JoinType == user.GroupJoinType_GroupJoinVerify {
		cresp, err := IsGroupManager(ctx, &user.IsGroupManagerReq{UserId: req.UserId, GroupId: req.GroupId})
		if err != nil {
			return nil, err
		}
		if !cresp.Flag {
			resp.Err = codes.Error(codes.ErrorUserGroupHasNoRight)
			return resp, nil
		}
	}

	if info.CurrentCt >= info.MaxMemberCt {
		resp.Err = codes.Error(codes.ErrorUserGroupMaxMember)
		return resp, nil
	}
	if err := becomeGroupMembers(ctx, req.GroupId, req.MemberIds); err != nil {
		return nil, err
	}
	for _, id := range req.MemberIds {
		// 加入新的群组
		ka.TopicEvent(mq.TopicEvent, mq.EEventGroupNewMember, event.GroupMember{OperatorId: req.UserId,
			GroupId: req.GroupId, UserId: id})
	}
	return resp, nil
}

func delGroupMembers(ctx context.Context, groupId int64, userIds []int64) error {
	userIds = util.UniqueInt64s(userIds)
	tx, err := dao.GetConn(dbconfig.DbUser).Begin()
	if err != nil {
		qlog.Ctx(ctx).Error(err, groupId, userIds)
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	{
		dsql := "delete from user_group where user_id in (?) and group_id=?"
		_, err = tx.Exec(dsql, util.Int64sToStr(userIds), groupId)
		if err != nil {
			qlog.Ctx(ctx).Error(err, groupId, userIds)
			return err
		}
	}

	{
		isql := "delete from " + tablename.UserGroupMember(groupId, config.GetEnv()) +
			" where group_id=? and user_id in (?) "
		_, err = tx.Exec(isql, groupId, util.Int64sToStr(userIds))
		if err != nil {
			qlog.Ctx(ctx).Error(err, groupId, userIds)
			return err
		}
	}

	{
		var ct = len(userIds)
		usql := "update group set current_ct=current_ct-%d where id=?"
		usql = fmt.Sprintf(usql, ct)
		_, err = tx.Exec(usql, groupId)
		if err != nil {
			qlog.Ctx(ctx).Error(err, groupId, userIds)
			return err
		}
	}

	ClearGroupCache(ctx, groupId)
	return nil
}

func GroupMemberDel(ctx context.Context, req *user.GroupMemberDelReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	info, err := GetGroupById(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}
	if info.CurrentCt <= 2 {
		resp.Err = codes.Error(codes.ErrorUserGroupUserLessThanTwo)
		return resp, nil
	}

	// 查看里面是否有管理员
	mems, err := getMemberByIds(ctx, req.MemberIds, req.GroupId)
	if err != nil {
		return nil, err
	}
	for _, mem := range mems {
		if mem.RoleType == user.GroupRoleType_GroupRoleManager || mem.RoleType == user.GroupRoleType_GroupRoleOwner {
			resp.Err = codes.Error(codes.ErrorUserGroupUserManagerCannotDeleted)
			return resp, nil
		}
	}

	{
		cresp, err := IsGroupManager(ctx, &user.IsGroupManagerReq{UserId: req.UserId, GroupId: req.GroupId})
		if err != nil {
			return nil, err
		}
		if !cresp.Flag {
			resp.Err = codes.Error(codes.ErrorUserGroupHasNoRight)
			return resp, nil
		}
	}

	for _, id := range req.MemberIds {
		// 删除成员
		ka.TopicEvent(mq.TopicEvent, mq.EEventGroupDelMember, event.GroupMember{OperatorId: req.UserId,
			GroupId: req.GroupId, UserId: id})
	}
	return resp, nil
}

// 离开群组
func GroupLeave(ctx context.Context, req *user.GroupLeaveReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	// 查看里面是否有管理员
	mem, err := getMemberById(ctx, req.UserId, req.GroupId)
	if err != nil {
		return nil, err
	}
	if mem.RoleType == user.GroupRoleType_GroupRoleOwner {
		resp.Err = codes.Error(codes.ErrorUserGroupOwnnerCannotLeave)
		return resp, nil
	}
	if err := delGroupMembers(ctx, req.GroupId, []int64{req.UserId}); err != nil {
		return nil, err
	}
	ka.TopicEvent(mq.TopicEvent, mq.EEventGroupLeaveMember, event.GroupMember{OperatorId: req.UserId,
		GroupId: req.GroupId, UserId: req.UserId})
	return resp, nil
}

func GroupDismiss(ctx context.Context, req *user.GroupDismissReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	mem, err := getMemberById(ctx, req.UserId, req.GroupId)
	if err != nil {
		return nil, err
	}
	if mem.RoleType == user.GroupRoleType_GroupRoleOwner {
		resp.Err = codes.Error(codes.ErrorUserGroupHasNoRight)
		return resp, nil
	}
	cresp, err := GroupMemberIds(ctx, &user.GroupMemberIdsReq{GroupId: req.GroupId})
	if err != nil {
		return nil, err
	}
	//ids,err:= G
	tx, err := dao.GetConn(dbconfig.DbUser).Begin()
	if err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		_ = tx.Commit()
	}()

	{
		dsql := "delete from user_group where group_id=?"
		_, err = tx.Exec(dsql, req.GroupId)
		if err != nil {
			qlog.Ctx(ctx).Error(err, req)
			return nil, err
		}
	}

	{
		isql := "delete from " + tablename.UserGroupMember(req.GroupId, config.GetEnv()) +
			" where group_id=? "
		_, err = tx.Exec(isql, req.GroupId)
		if err != nil {
			qlog.Ctx(ctx).Error(err, req)
			return nil, err
		}
	}

	{
		usql := "update group set deleted_on=unix_timestamp() where id=?"
		_, err = tx.Exec(usql, req.GroupId)
		if err != nil {
			qlog.Ctx(ctx).Error(err, req)
			return nil, err
		}
	}

	ClearGroupCache(ctx, req.GroupId)

	ka.TopicEvent(mq.TopicEvent, mq.EEventGroupDismiss, event.GroupDismiss{OperatorId: req.UserId,
		GroupId: req.GroupId, Ids: cresp.Ids})
	return resp, nil
}

func GroupTransfer(ctx context.Context, req *user.GroupTransferReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	mem, err := getMemberById(ctx, req.UserId, req.GroupId)
	if err != nil {
		return nil, err
	}
	if mem.RoleType == user.GroupRoleType_GroupRoleOwner {
		resp.Err = codes.Error(codes.ErrorUserGroupHasNoRight)
		return resp, nil
	}
	cresp, err := IsGroupMember(ctx, &user.IsGroupMemberReq{GroupId: req.GroupId, UserId: req.MemberId})
	if err != nil {
		return nil, err
	}
	if !cresp.Flag {
		resp.Err = codes.Error(codes.ErrorUserGroupNotMember)
		return resp, nil
	}

	conn := dao.GetConn(dbconfig.DbUser)

	usql := "update " + tablename.UserGroupMember(req.GroupId, config.GetEnv()) +
		" set user_role=? where group_id=? and user_id=?"
	{
		_, err = conn.Exec(usql, user.GroupRoleType_GroupRoleOwner, req.GroupId, req.MemberId)
		if err != nil {
			qlog.Ctx(ctx).Error(err, req)
			return nil, err
		}
	}
	{
		_, err = conn.Exec(usql, user.GroupRoleType_GroupRoleNormal, req.GroupId, req.UserId)
		if err != nil {
			qlog.Ctx(ctx).Error(err, req)
			return nil, err
		}
	}

	ClearGroupCache(ctx, req.GroupId)

	ka.TopicEvent(mq.TopicEvent, mq.EEventGroupTransfer, event.GroupMember{OperatorId: req.UserId,
		GroupId: req.GroupId, UserId: req.MemberId})
	return resp, nil
}

func GroupMembersByGroupId(ctx context.Context, req *user.GroupMembersByGroupIdReq) (*user.GroupMembersByGroupIdResp, error) {
	resp := new(user.GroupMembersByGroupIdResp)
	cresp, err := IsGroupMember(ctx, &user.IsGroupMemberReq{UserId: req.UserId, GroupId: req.GroupId})
	if err != nil {
		return nil, err
	}
	if !cresp.Flag {
		resp.Err = codes.Error(codes.ErrorUserGroupNotMember)
		return resp, nil
	}

	tname := tablename.UserGroupMember(req.GroupId, config.GetEnv())
	csql := "select count(1) from " + tname + " where group_id=?"
	if err := dao.GetConn(dbconfig.DbUser).QueryRow(csql, req.GroupId).Scan(&resp.Total); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	if resp.Total <= req.Page*req.PageSize {
		return resp, nil
	}

	mems := make([]*user.GroupMember, 0)
	ids := make([]int64, 0)
	ssql := "select user_id,user_role,mark_name,mute_until,not_disturb,created_on from " + tname +
		" where group_id=? and blocked=0 limit %d offset %d"
	ssql = fmt.Sprintf(ssql, req.PageSize, req.Page*req.PageSize)
	rows, err := dao.GetConn(dbconfig.DbUser).Query(ssql, req.GroupId)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	for rows.Next() {
		m := new(user.GroupMember)
		if err := rows.Scan(&m.UserId, &m.RoleType, &m.MarkName, &m.MuteUntil, &m.NotDisturb, &m.CreatedOn); err != nil {
			qlog.Ctx(ctx).Error(req, err)
			continue
		}
		ids = append(ids, m.UserId)
		mems = append(mems, m)
	}
	rows.Close()
	userInfos, err := GetUserInfoByIds(ctx, ids)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err, ids)
		return nil, err
	}
	for _, m := range mems {
		for _, info := range userInfos {
			if info.UserId == m.UserId {
				m.User = info
			}
		}
	}
	resp.Members = mems
	return resp, nil
}

func GroupBlock(ctx context.Context, req *user.GroupBlockReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	cresp, err := IsGroupMember(ctx, &user.IsGroupMemberReq{UserId: req.UserId, GroupId: req.GroupId})
	if err != nil {
		return nil, err
	}
	var esql string
	if cresp.Flag {
		esql = "update " + tablename.UserGroupMember(req.GroupId, config.GetEnv()) +
			"blocked=1 where group_id=? and user_id=? and blocked=0"
		if !req.IsBlock {
			esql = "update " + tablename.UserGroupMember(req.GroupId, config.GetEnv()) +
				"blocked=0 where group_id=? and user_id=? and blocked=1"
		}
	} else {
		esql = "insert into " + tablename.UserGroupMember(req.GroupId, config.GetEnv()) +
			"(group_id,user_id,blocked,created_on) value (?,?,1,unix_timestamp())"
		if !req.IsBlock {
			esql = "delete from " + tablename.UserGroupMember(req.GroupId, config.GetEnv()) +
				" where group_id=? and user_id=? and blocked=1"
		}
	}

	if !req.IsBlock {
		esql = "delete from " + tablename.UserGroupMember(req.GroupId, config.GetEnv()) +
			" where group_id=? and user_id=? and blocked=1)"
	}
	result, err := dao.GetConn(dbconfig.DbUser).Exec(esql, req.GroupId, req.UserId)
	if err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}

	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}

	ka.TopicEvent(mq.TopicEvent, mq.EEventGroupBlock, event.GroupMember{OperatorId: req.UserId,
		GroupId: req.GroupId, UserId: req.MemberId, Flag: req.IsBlock})
	return resp, nil
}

func GroupBlockList(ctx context.Context, req *user.GroupBlockListReq) (*user.GroupBlockListResp, error) {
	resp := new(user.GroupBlockListResp)
	tname := tablename.UserGroupMember(req.GroupId, config.GetEnv())
	ssql := "select user_id from " + tname + "where group_id=? and `blocked`=1"
	rows, err := dao.GetConn(dbconfig.DbUser).Query(ssql, req.GroupId)
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

	members, err := getMemberByIds(ctx, ids, req.GroupId)
	if err != nil {
		return nil, err
	}

	resp.Members = members
	return resp, nil
}
