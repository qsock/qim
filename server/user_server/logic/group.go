package logic

import (
	"context"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/store/ka"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/config/mq"
	"github.com/qsock/qim/lib/codes"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/gclient"
	"github.com/qsock/qim/lib/proto/event"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/user"
	"github.com/qsock/qim/lib/tablename"
	"github.com/qsock/qim/lib/util"
	"github.com/qsock/qim/server/user_server/config"
	"strings"
)

func GroupCreate(ctx context.Context, req *user.GroupCreateReq) (*user.GroupCreateResp, error) {
	resp := new(user.GroupCreateResp)
	req.MemberIds = append(req.MemberIds, req.UserId)
	ids := []int64{req.UserId}
	ids = append(ids, req.MemberIds...)
	ids = util.UniqueInt64s(ids)

	if len(ids) == 1 {
		resp.Err = codes.Error(codes.ErrorUserGroupMemberMustBiggerThanOne)
		return resp, nil
	}
	groupId := gclient.GenDbId(ctx, config.GetConfig().IdKey)
	names := make([]string, 0)
	avatars := make([]string, 0)

	if len(req.Name) == 0 ||
		len(req.Avatar) == 0 {
		ids2 := append([]int64{}, ids...)
		if len(ids2) > 9 {
			ids2 = ids[:9]
		}
		infos, err := GetUserInfoByIds(ctx, ids2)
		if err != nil {
			qlog.Ctx(ctx).Error(req)
			return nil, err
		}

		for _, id := range ids {
			for _, info := range infos {
				if info.UserId == id {
					names = append(names, info.Name)
					avatars = append(avatars, info.Avatar)
				}
			}
		}
	}
	if len(req.Name) == 0 {
		req.Name = strings.Join(names, ",")
		if len(req.Name) > 32 {
			req.Name = req.Name[0:32]
		}
	}
	if len(req.Avatar) == 0 {
		req.Avatar = strings.Join(avatars, ",")
	}
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
		isql := "insert into `group` (id, name,created_on,max_member_ct,current_ct," +
			"avatar,join_type,t) values(?,?,unix_timestamp(),?,?," +
			"?,?,?)"
		_, err = tx.Exec(isql, groupId, req.Name, req.MaxMemberCt, len(ids),
			req.Avatar, user.GroupJoinType_GroupJoinVerify, req.T)
		if err != nil {
			qlog.Ctx(ctx).Error(err, req)
			return nil, err
		}
	}
	{
		isql := "insert into " + tablename.UserGroupMember(groupId, config.GetEnv()) + "(group_id,user_id,user_role,created_on) values "
		condis := make([]interface{}, 0)
		for _, memberId := range ids {
			if memberId == req.UserId {
				condis = append(condis, groupId, memberId, user.GroupRoleType_GroupRoleOwner)
			} else {
				condis = append(condis, groupId, memberId, user.GroupRoleType_GroupRoleNormal)
			}
			isql += "(?,?,?,unix_timestamp()),"
		}
		isql = isql[:len(isql)-1]
		_, err = tx.Exec(isql, condis...)
		if err != nil {
			qlog.Ctx(ctx).Error(err, req, ids)
			return nil, err
		}
	}
	{
		isql := "insert into user_group " +
			" (user_id,group_id) values "
		condis := make([]interface{}, 0)
		for _, memberId := range ids {
			condis = append(condis, memberId, groupId)
			isql += "(?,?),"
		}
		isql = isql[:len(isql)-1]
		_, err = tx.Exec(isql, condis...)
		if err != nil {
			qlog.Ctx(ctx).Error(err, req, ids)
			return nil, err
		}
	}

	resp.GroupId = groupId
	return resp, nil
}

func GroupInfoById(ctx context.Context, req *user.GroupInfoReq) (*user.GroupInfoResp, error) {
	info, err := GetGroupById(ctx, req.GroupId)
	if err != nil {
		return nil, err
	}
	resp := new(user.GroupInfoResp)
	resp.Info = info
	return resp, nil
}

func GroupInfoByIds(ctx context.Context, req *user.GroupInfosReq) (*user.GroupInfosResp, error) {
	infos, err := GetGroupByIds(ctx, req.GroupIds)
	if err != nil {
		return nil, err
	}
	resp := new(user.GroupInfosResp)
	resp.Infos = infos
	return resp, nil
}

func GroupsByUid(ctx context.Context, req *user.GroupsByUidReq) (*user.GroupsByUidResp, error) {
	resp := new(user.GroupsByUidResp)
	ids := make([]int64, 0)
	ssql := "select group_id from user_group where user_id=?"
	rows, err := dao.GetConn(dbconfig.DbUser).Query(ssql, req.UserId)
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
	resp.Total = int32(len(ids))
	if req.Page*req.PageSize >= resp.Total {
		return resp, nil
	}
	start := int(req.PageSize * req.Page)
	end := start + int(req.PageSize)
	if end > int(resp.Total) {
		end = int(resp.Total)
	}

	infos, err := GetGroupByIds(ctx, ids[start:end])
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	resp.Infos = infos
	return resp, nil
}

// 指定管理员
func GroupAppointManager(ctx context.Context, req *user.GroupAppointManagerReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	memberInfo, err := getMemberById(ctx, req.UserId, req.GroupId)
	if err != nil {
		return nil, err
	}
	if util.InArrayInt64(req.UserId, req.ManagerIds) {
		resp.Err = codes.Error(codes.ErrorUserGroupCannotPointSelf)
		return resp, nil
	}
	if memberInfo.RoleType != user.GroupRoleType_GroupRoleOwner {
		resp.Err = codes.Error(codes.ErrorUserGroupHasNoRight)
		return resp, nil
	}

	var roleType = user.GroupRoleType_GroupRoleManager
	if !req.IsAppoint {
		roleType = user.GroupRoleType_GroupRoleNormal
	}

	tname := tablename.UserGroupMember(req.GroupId, config.GetEnv())
	usql := "update " + tname + " set user_role=? where group_id and user_id in (?)"
	result, err := dao.GetConn(dbconfig.DbUser).Exec(usql, roleType, util.Int64sToStr(req.ManagerIds))
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		return nil, err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}

	for _, managerId := range req.ManagerIds {
		ka.TopicEvent(mq.TopicEvent, mq.EEventGroupManager, event.GroupManager{UserId: req.UserId, GroupId: req.GroupId,
			ManagerId: managerId, IsManager: req.IsAppoint})
	}
	return resp, nil
}

func IsGroupManager(ctx context.Context, req *user.IsGroupManagerReq) (*ret.BoolResp, error) {
	memberInfo, err := getMemberById(ctx, req.UserId, req.GroupId)
	if err != nil {
		return nil, err
	}
	resp := new(ret.BoolResp)
	if memberInfo.RoleType == user.GroupRoleType_GroupRoleManager ||
		memberInfo.RoleType == user.GroupRoleType_GroupRoleOwner {
		resp.Flag = true
	}
	return resp, nil
}

func GroupUpdateName(ctx context.Context, req *user.GroupUpdateNameReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)

	cresp, err := IsGroupManager(ctx, &user.IsGroupManagerReq{UserId: req.UserId, GroupId: req.GroupId})
	if err != nil {
		return nil, err
	}
	if !cresp.Flag {
		resp.Err = codes.Error(codes.ErrorUserGroupHasNoRight)
		return resp, nil
	}
	usql := "update group set name=? where id=?"
	result, err := dao.GetConn(dbconfig.DbUser).Exec(usql, req.Name, req.GroupId)
	if err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}
	ClearGroupCache(ctx, req.GroupId)
	ka.TopicEvent(mq.TopicEvent, mq.EEventGroupUpdateName, event.GroupUpdate{GroupId: req.GroupId, Str: req.Name, UserId: req.UserId})
	return resp, nil
}

func GroupUpdateNotice(ctx context.Context, req *user.GroupUpdateNoticeReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)

	cresp, err := IsGroupManager(ctx, &user.IsGroupManagerReq{UserId: req.UserId, GroupId: req.GroupId})
	if err != nil {
		return nil, err
	}
	if !cresp.Flag {
		resp.Err = codes.Error(codes.ErrorUserGroupHasNoRight)
		return resp, nil
	}
	usql := "update group set notice=? where id=?"
	result, err := dao.GetConn(dbconfig.DbUser).Exec(usql, req.Notice, req.GroupId)
	if err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}
	ClearGroupCache(ctx, req.GroupId)
	ka.TopicEvent(mq.TopicEvent, mq.EEventGroupUpdateNotice, event.GroupUpdate{GroupId: req.GroupId, Str: req.Notice, UserId: req.UserId})
	return resp, nil
}

func GroupUpdateAvatar(ctx context.Context, req *user.GroupUpdateAvatarReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)

	cresp, err := IsGroupManager(ctx, &user.IsGroupManagerReq{UserId: req.UserId, GroupId: req.GroupId})
	if err != nil {
		return nil, err
	}
	if !cresp.Flag {
		resp.Err = codes.Error(codes.ErrorUserGroupHasNoRight)
		return resp, nil
	}
	usql := "update group set avatar=? where id=?"
	result, err := dao.GetConn(dbconfig.DbUser).Exec(usql, req.Avatar, req.GroupId)
	if err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}
	ClearGroupCache(ctx, req.GroupId)
	ka.TopicEvent(mq.TopicEvent, mq.EEventGroupUpdateAvatar, event.GroupUpdate{GroupId: req.GroupId, Str: req.Avatar, UserId: req.UserId})
	return resp, nil
}

func GroupUpdateJoinType(ctx context.Context, req *user.GroupUpdateJoinTypeReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)

	cresp, err := IsGroupManager(ctx, &user.IsGroupManagerReq{UserId: req.UserId, GroupId: req.GroupId})
	if err != nil {
		return nil, err
	}
	if !cresp.Flag {
		resp.Err = codes.Error(codes.ErrorUserGroupHasNoRight)
		return resp, nil
	}
	usql := "update group set join_type=? where id=?"
	result, err := dao.GetConn(dbconfig.DbUser).Exec(usql, req.JoinType, req.GroupId)
	if err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}
	ClearGroupCache(ctx, req.GroupId)
	ka.TopicEvent(mq.TopicEvent, mq.EEventGroupUpdateJointype,
		event.GroupUpdate{GroupId: req.GroupId, Type: int32(req.JoinType), UserId: req.UserId})
	return resp, nil
}

func GroupMute(ctx context.Context, req *user.GroupMuteReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	cresp, err := IsGroupManager(ctx, &user.IsGroupManagerReq{UserId: req.UserId, GroupId: req.GroupId})
	if err != nil {
		return nil, err
	}
	if !cresp.Flag {
		resp.Err = codes.Error(codes.ErrorUserGroupHasNoRight)
		return resp, nil
	}
	usql := "update group set mute_util=? where id=?"
	result, err := dao.GetConn(dbconfig.DbUser).Exec(usql, req.MuteUntil, req.GroupId)
	if err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}
	ClearGroupCache(ctx, req.GroupId)
	ka.TopicEvent(mq.TopicEvent, mq.EEventGroupUpdateMute, event.GroupUpdate{GroupId: req.GroupId, Time: req.MuteUntil, UserId: req.UserId})
	return resp, nil
}

func GroupMuteUser(ctx context.Context, req *user.GroupMuteUserReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
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
	{
		mems, err := getMemberByIds(ctx, req.MemberId, req.UserId)
		if err != nil {
			return nil, err
		}
		for _, mem := range mems {
			if mem.RoleType == user.GroupRoleType_GroupRoleManager ||
				mem.RoleType == user.GroupRoleType_GroupRoleOwner {

				resp.Err = codes.Error(codes.ErrorUserGroupHasNoRight)
				return resp, nil
			}
		}
	}

	tname := tablename.UserGroupMember(req.GroupId, config.GetEnv())
	usql := "update " + tname + " set mute_util=? where group_id=? and user_id in (?)"
	result, err := dao.GetConn(dbconfig.DbUser).Exec(usql, req.MuteUntil, req.GroupId, util.Int64sToStr(req.MemberId))
	if err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	if n, _ := result.RowsAffected(); n == 0 {
		return resp, nil
	}

	for _, id := range req.MemberId {
		ka.TopicEvent(mq.TopicEvent, mq.EEventGroupMuteone,
			event.GroupMember{GroupId: req.GroupId, OperatorId: req.UserId, UserId: id,
				Time: req.MuteUntil})
	}
	return resp, nil
}

// 禁言列表
func GroupMuteList(ctx context.Context, req *user.GroupMuteListReq) (*user.GroupMuteListResp, error) {
	resp := new(user.GroupMuteListResp)
	tname := tablename.UserGroupMember(req.GroupId, config.GetEnv())
	ssql := "select user_id from " + tname + "where group_id=? and mute_until>unix_timestamp()"
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
