package main

import (
	"context"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/user"
	"github.com/qsock/qim/server/user_server/logic"
)

type Server struct{}

func (*Server) Ping(ctx context.Context, req *ret.NoArgs) (*ret.NoArgs, error) {
	return new(ret.NoArgs), nil
}

func (server *Server) Create(ctx context.Context, req *user.CreateReq) (*ret.EmptyResp, error) {
	return logic.Create(ctx, req)
}

func (server *Server) Infos(ctx context.Context, req *user.InfosReq) (*user.InfosResp, error) {
	return logic.Infos(ctx, req)
}

func (server *Server) Info(ctx context.Context, req *user.InfoReq) (*user.InfoResp, error) {
	return logic.Info(ctx, req)
}

func (server *Server) NewApplyUserList(ctx context.Context, req *user.NewApplyUserListReq) (*user.NewApplyListResp, error) {
	return logic.NewApplyUserList(ctx, req)
}

func (server *Server) Lastactive(ctx context.Context, req *user.LastactiveReq) (*ret.EmptyResp, error) {
	return logic.Lastactive(ctx, req)
}

func (server *Server) Update(ctx context.Context, req *user.UpdateReq) (*ret.EmptyResp, error) {
	return logic.Update(ctx, req)
}

// 好友相关
func (server *Server) IsFriend(ctx context.Context, req *user.IsFriendReq) (*ret.BoolResp, error) {
	return logic.IsFriend(ctx, req)
}

func (server *Server) FriendIds(ctx context.Context, req *user.FriendIdsReq) (*user.FriendIdsResp, error) {
	return logic.FriendIds(ctx, req)
}

func (server *Server) FriendByIds(ctx context.Context, req *user.FriendByIdsReq) (*user.FriendByIdsResp, error) {
	return logic.FriendByIds(ctx, req)
}

func (server *Server) FriendsByUid(ctx context.Context, req *user.FriendsByUidReq) (*user.FriendsByUidResp, error) {
	return logic.FriendsByUid(ctx, req)
}

func (server *Server) FriendDel(ctx context.Context, req *user.FriendDelReq) (*ret.EmptyResp, error) {
	return logic.FriendDel(ctx, req)
}

func (server *Server) FriendMarknameUpdate(ctx context.Context, req *user.FriendMarknameUpdateReq) (*ret.EmptyResp, error) {
	return logic.FriendMarknameUpdate(ctx, req)
}

func (server *Server) FriendNewApply(ctx context.Context, req *user.FriendNewApplyReq) (*ret.EmptyResp, error) {
	return logic.FriendNewApply(ctx, req)
}

func (server *Server) FriendNewAgree(ctx context.Context, req *user.FriendNewAgreeReq) (*ret.EmptyResp, error) {
	return logic.FriendNewAgree(ctx, req)
}

func (server *Server) FriendNewReject(ctx context.Context, req *user.FriendNewRejectReq) (*ret.EmptyResp, error) {
	return logic.FriendNewReject(ctx, req)
}

func (server *Server) FriendNewDel(ctx context.Context, req *user.FriendNewDelReq) (*ret.EmptyResp, error) {
	return logic.FriendNewDel(ctx, req)
}

func (server *Server) FriendNewIgnore(ctx context.Context, req *user.FriendNewIgnoreReq) (*ret.EmptyResp, error) {
	return logic.FriendNewIgnore(ctx, req)
}

// 黑名单部分
func (server *Server) BlacklistAdd(ctx context.Context, req *user.BlacklistAddReq) (*ret.EmptyResp, error) {
	return logic.BlacklistAdd(ctx, req)
}

func (server *Server) BlacklistDel(ctx context.Context, req *user.BlacklistDelReq) (*ret.EmptyResp, error) {
	return logic.BlacklistDel(ctx, req)
}

func (server *Server) Blacklist(ctx context.Context, req *user.BlacklistReq) (*user.BlacklistResp, error) {
	return logic.Blacklist(ctx, req)
}

// 群组部分
func (server *Server) IsGroupMember(ctx context.Context, req *user.IsGroupMemberReq) (*ret.BoolResp, error) {
	return logic.IsGroupMember(ctx, req)
}

func (server *Server) IsGroupMemberBeenMute(ctx context.Context, req *user.IsGroupMemberBeenMuteReq) (*ret.BoolResp, error) {
	return logic.IsGroupMemberBeenMute(ctx, req)
}

func (server *Server) GroupMemberIds(ctx context.Context, req *user.GroupMemberIdsReq) (*user.GroupMemberIdsResp, error) {
	return logic.GroupMemberIds(ctx, req)
}

func (server *Server) GroupMembersByGroupId(ctx context.Context, req *user.GroupMembersByGroupIdReq) (*user.GroupMembersByGroupIdResp, error) {
	return logic.GroupMembersByGroupId(ctx, req)
}

func (server *Server) GroupCreate(ctx context.Context, req *user.GroupCreateReq) (*user.GroupCreateResp, error) {
	return logic.GroupCreate(ctx, req)
}

func (server *Server) GroupInfoById(ctx context.Context, req *user.GroupInfoReq) (*user.GroupInfoResp, error) {
	return logic.GroupInfoById(ctx, req)
}

func (server *Server) GroupInfoByIds(ctx context.Context, req *user.GroupInfosReq) (*user.GroupInfosResp, error) {
	return logic.GroupInfoByIds(ctx, req)
}

func (server *Server) GroupsByUid(ctx context.Context, req *user.GroupsByUidReq) (*user.GroupsByUidResp, error) {
	return logic.GroupsByUid(ctx, req)
}

func (server *Server) GroupAppointManager(ctx context.Context, req *user.GroupAppointManagerReq) (*ret.EmptyResp, error) {
	return logic.GroupAppointManager(ctx, req)
}

func (server *Server) GroupManagerList(ctx context.Context, req *user.GroupManagerListReq) (*user.GroupManagerListResp, error) {
	return logic.GroupManagerList(ctx, req)
}

func (server *Server) IsGroupManager(ctx context.Context, req *user.IsGroupManagerReq) (*ret.BoolResp, error) {
	return logic.IsGroupManager(ctx, req)
}

func (server *Server) GroupMemberByIds(ctx context.Context, req *user.GroupMemberByIdsReq) (*user.GroupMemberByIdsResp, error) {
	return logic.GroupMemberByIds(ctx, req)
}

func (server *Server) GroupMemberAdd(ctx context.Context, req *user.GroupMemberAddReq) (*ret.EmptyResp, error) {
	return logic.GroupMemberAdd(ctx, req)
}

func (server *Server) GroupMemberDel(ctx context.Context, req *user.GroupMemberDelReq) (*ret.EmptyResp, error) {
	return logic.GroupMemberDel(ctx, req)
}

func (server *Server) GroupLeave(ctx context.Context, req *user.GroupLeaveReq) (*ret.EmptyResp, error) {
	return logic.GroupLeave(ctx, req)
}

func (server *Server) GroupDismiss(ctx context.Context, req *user.GroupDismissReq) (*ret.EmptyResp, error) {
	return logic.GroupDismiss(ctx, req)
}

func (server *Server) GroupBlock(ctx context.Context, req *user.GroupBlockReq) (*ret.EmptyResp, error) {
	return logic.GroupBlock(ctx, req)
}

func (server *Server) GroupBlockList(ctx context.Context, req *user.GroupBlockListReq) (*user.GroupBlockListResp, error) {
	return logic.GroupBlockList(ctx, req)
}

func (server *Server) GroupJoin(ctx context.Context, req *user.GroupJoinReq) (*ret.EmptyResp, error) {
	return logic.GroupJoin(ctx, req)
}

func (server *Server) GroupJoinAgree(ctx context.Context, req *user.GroupJoinAgreeReq) (*ret.EmptyResp, error) {
	return logic.GroupJoinAgree(ctx, req)
}

func (server *Server) GroupJoinReject(ctx context.Context, req *user.GroupJoinRejectReq) (*ret.EmptyResp, error) {
	return logic.GroupJoinReject(ctx, req)
}

func (server *Server) GroupJoinIgnore(ctx context.Context, req *user.GroupJoinIgnoreReq) (*ret.EmptyResp, error) {
	return logic.GroupJoinIgnore(ctx, req)
}
func (server *Server) GroupTransfer(ctx context.Context, req *user.GroupTransferReq) (*ret.EmptyResp, error) {
	return logic.GroupTransfer(ctx, req)
}

func (server *Server) GroupMute(ctx context.Context, req *user.GroupMuteReq) (*ret.EmptyResp, error) {
	return logic.GroupMute(ctx, req)
}

func (server *Server) GroupMuteUser(ctx context.Context, req *user.GroupMuteUserReq) (*ret.EmptyResp, error) {
	return logic.GroupMuteUser(ctx, req)
}

func (server *Server) GroupMuteList(ctx context.Context, req *user.GroupMuteListReq) (*user.GroupMuteListResp, error) {
	return logic.GroupMuteList(ctx, req)
}

func (server *Server) GroupUpdateName(ctx context.Context, req *user.GroupUpdateNameReq) (*ret.EmptyResp, error) {
	return logic.GroupUpdateName(ctx, req)
}

func (server *Server) GroupUpdateNotice(ctx context.Context, req *user.GroupUpdateNoticeReq) (*ret.EmptyResp, error) {
	return logic.GroupUpdateNotice(ctx, req)
}

func (server *Server) GroupUpdateAvatar(ctx context.Context, req *user.GroupUpdateAvatarReq) (*ret.EmptyResp, error) {
	return logic.GroupUpdateAvatar(ctx, req)
}

func (server *Server) GroupUpdateJoinType(ctx context.Context, req *user.GroupUpdateJoinTypeReq) (*ret.EmptyResp, error) {
	return logic.GroupUpdateJoinType(ctx, req)
}
