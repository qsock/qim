package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qim/lib/ginproxy"
	"github.com/qsock/qim/lib/method"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/proto/user"
	"github.com/qsock/qim/lib/util"
)

// @summary 创建群组
// @description 创建群组
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupCreateReq true "请求参数" required
// @router  /group/create [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupCreate(c *gin.Context) {
	req := new(user.GroupCreateReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	req.MemberIds = util.UniqueInt64s(req.MemberIds)
	if len(req.MemberIds) < 2 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupCreate, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 个人群组列表
// @description 个人群组列表
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   page query int true "总页数，从0开始"
// @param   page_size query int false "每页数量，默认30条"
// @router  /user/groups [GET]
// @success 200 {object} user.GroupsByUidResp "请求返回"
func GroupsByUid(c *gin.Context) {
	req := new(user.GroupsByUidReq)
	ctx := ginproxy.GetCtx(c)
	req.Page = ginproxy.GetInt32Query(c, "page", 0)
	req.PageSize = ginproxy.GetInt32Query(c, "page_size", 30)
	req.UserId = ginproxy.GetUserId(c)
	resp := new(user.GroupsByUidResp)
	err := qgrpc.Call(ctx, method.UserGroupsByUid, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 通过id查询群组
// @description 通过id查询一个群组
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   group_id query int true "群组id"
// @router  /group/info [GET]
// @success 200 {object} user.GroupInfoResp "请求返回"
func GroupInfoById(c *gin.Context) {
	req := new(user.GroupInfoReq)
	ctx := ginproxy.GetCtx(c)
	req.GroupId = ginproxy.GetInt64Query(c, "group_id", 0)
	resp := new(user.GroupInfoResp)
	if req.GroupId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	err := qgrpc.Call(ctx, method.UserGroupInfoById, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 通过id查询群组列表
// @description 通过id查询多个群组
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   group_ids query string true "需要查看的id列表,逗号分割所有的id,例如:1234,12511,1100"
// @router  /group/infos [GET]
// @success 200 {object} user.GroupInfoResp "请求返回"
func GroupInfoByIds(c *gin.Context) {
	req := new(user.GroupInfosReq)
	ctx := ginproxy.GetCtx(c)
	ids := c.Query("group_ids")
	req.GroupIds = util.StrToInt64s(ids)
	if len(req.GroupIds) == 0 {
		ginproxy.ParameterError(c)
		return
	}
	resp := new(user.GroupInfoResp)
	err := qgrpc.Call(ctx, method.UserGroupInfoByIds, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 群主指定管理员
// @description 群主指定管理员
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupAppointManagerReq true "请求参数" required
// @router  /group/appoint/manager [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupAppointManager(c *gin.Context) {
	req := new(user.GroupAppointManagerReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	req.ManagerIds = util.UniqueInt64s(req.ManagerIds)
	if util.InArrayInt64(req.UserId, req.ManagerIds) ||
		len(req.ManagerIds) == 0 ||
		req.GroupId <= 0 {
		ginproxy.ParameterError(c)
		return
	}

	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupAppointManager, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 群组管理员列表
// @description 群组管理员列表
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   group_id query int true "群组id"
// @router  /group/managers [GET]
// @success 200 {object} user.GroupManagerListResp "请求返回"
func GroupManagerList(c *gin.Context) {
	req := new(user.GroupManagerListReq)
	ctx := ginproxy.GetCtx(c)
	req.GroupId = ginproxy.GetInt64Query(c, "group_id", 0)
	if req.GroupId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	resp := new(user.GroupManagerListResp)
	err := qgrpc.Call(ctx, method.UserGroupManagerList, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 获取群组成员id列表
// @description 获取群组成员id列表
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   group_id query int true "群组id"
// @router  /group/member/ids [GET]
// @success 200 {object} user.GroupMemberIdsResp "请求返回"
func GroupMemberIds(c *gin.Context) {
	req := new(user.GroupMemberIdsReq)
	ctx := ginproxy.GetCtx(c)
	req.GroupId = ginproxy.GetInt64Query(c, "group_id", 0)
	if req.GroupId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	resp := new(user.GroupMemberIdsResp)
	err := qgrpc.Call(ctx, method.UserGroupMemberIds, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 通过群组成员id获取成员
// @description 通过群组成员id获取成员
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   group_id query int true "群组id"
// @param   ids query string true "需要查看的成员id列表,逗号分割所有的id,例如:1234,12511,1100"
// @router  /group/member/by-ids [GET]
// @success 200 {object} user.GroupMemberByIdsResp "请求返回"
func GroupMemberByIds(c *gin.Context) {
	req := new(user.GroupMemberByIdsReq)
	ctx := ginproxy.GetCtx(c)
	req.GroupId = ginproxy.GetInt64Query(c, "group_id", 0)
	ids := c.Query("ids")
	req.UserIds = util.StrToInt64s(ids)
	if req.GroupId <= 0 ||
		len(req.UserIds) == 0 {
		ginproxy.ParameterError(c)
		return
	}
	resp := new(user.GroupMemberByIdsResp)
	err := qgrpc.Call(ctx, method.UserGroupMemberByIds, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)

}

// @summary 所在群的群成员
// @description 所在群的群成员
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   group_id query int true "群组id"
// @param   page query int true "总页数，从0开始"
// @param   page_size query int false "每页数量，默认30条"
// @router  /group/members [GET]
// @success 200 {object} user.GroupMemberIdsResp "请求返回"
func GroupMembersByGroupId(c *gin.Context) {
	req := new(user.GroupMembersByGroupIdReq)
	ctx := ginproxy.GetCtx(c)
	req.GroupId = ginproxy.GetInt64Query(c, "group_id", 0)
	req.UserId = ginproxy.GetUserId(c)
	req.Page = ginproxy.GetInt32Query(c, "page", 0)
	req.PageSize = ginproxy.GetInt32Query(c, "page_size", 0)
	if req.GroupId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	resp := new(user.GroupMembersByGroupIdResp)
	err := qgrpc.Call(ctx, method.UserGroupMembersByGroupId, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 添加群成员
// @description 添加群成员
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupMemberAddReq true "请求参数" required
// @router  /group/member/add [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupMemberAdd(c *gin.Context) {
	req := new(user.GroupMemberAddReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	req.MemberIds = util.UniqueInt64s(req.MemberIds)
	if len(req.MemberIds) < 1 || req.GroupId == 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupMemberAdd, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 删除群成员
// @description 添加群成员
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupMemberDelReq true "请求参数" required
// @router  /group/member/del [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupMemberDel(c *gin.Context) {
	req := new(user.GroupMemberDelReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	req.MemberIds = util.UniqueInt64s(req.MemberIds)
	if len(req.MemberIds) < 1 || req.GroupId == 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupMemberDel, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 离开群组
// @description 离开群组
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupLeaveReq true "请求参数" required
// @router  /group/leave [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupLeave(c *gin.Context) {
	req := new(user.GroupLeaveReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}

	if req.GroupId == 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupLeave, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 解散群组
// @description 解散群组
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupDismissReq true "请求参数" required
// @router  /group/dismiss [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupDismiss(c *gin.Context) {
	req := new(user.GroupDismissReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}

	if req.GroupId == 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupDismiss, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 申请加入群组
// @description 申请加入群组
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupJoinReq true "请求参数" required
// @router  /group/join/apply [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupJoin(c *gin.Context) {
	req := new(user.GroupJoinReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}

	if req.GroupId == 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupJoin, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 同意申请加入群组
// @description 同意申请加入群组
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupJoinAgreeReq true "请求参数" required
// @router  /group/join/agree [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupJoinAgree(c *gin.Context) {
	req := new(user.GroupJoinAgreeReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}

	if req.GroupId <= 0 || req.MemberId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupJoinAgree, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 拒绝申请加入群组
// @description 拒绝申请加入群组
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupJoinRejectReq true "请求参数" required
// @router  /group/join/reject [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupJoinReject(c *gin.Context) {
	req := new(user.GroupJoinRejectReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}

	if req.GroupId <= 0 || req.MemberId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupJoinReject, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 忽略申请加入群组
// @description 忽略申请加入群组
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupJoinIgnoreReq true "请求参数" required
// @router  /group/join/ignore [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupJoinIgnore(c *gin.Context) {
	req := new(user.GroupJoinIgnoreReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}

	if req.GroupId <= 0 || req.MemberId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupJoinIgnore, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 转让群组
// @description 转让群组
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupTransferReq true "请求参数" required
// @router  /group/transfer [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupTransfer(c *gin.Context) {
	req := new(user.GroupTransferReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}

	if req.GroupId <= 0 || req.MemberId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupTransfer, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 群组禁言
// @description 群组禁言
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupMuteReq true "请求参数" required
// @router  /group/mute [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupMute(c *gin.Context) {
	req := new(user.GroupMuteReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}

	if req.GroupId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupMute, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 拉黑某个人
// @description 拉黑某个人
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupBlockReq true "请求参数" required
// @router  /group/block [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupBlock(c *gin.Context) {
	req := new(user.GroupBlockReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}

	if req.GroupId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupBlock, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 禁言列表
// @description 禁言列表
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   group_id query int true "群组id"
// @router  /group/mute/list [GET]
// @success 200 {object} user.GroupMuteListResp "请求返回"
func GroupBlockList(c *gin.Context) {
	req := new(user.GroupBlockListReq)
	ctx := ginproxy.GetCtx(c)
	req.GroupId = ginproxy.GetInt64Query(c, "group_id", 0)

	if req.GroupId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	resp := new(user.GroupBlockListResp)
	err := qgrpc.Call(ctx, method.UserGroupBlockList, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 禁言某些人
// @description 禁言某些人
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupMuteUserReq true "请求参数" required
// @router  /group/mute/user [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupMuteUser(c *gin.Context) {
	req := new(user.GroupMuteUserReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}

	if req.GroupId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupMuteUser, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 禁言列表
// @description 禁言列表
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   group_id query int true "群组id"
// @router  /group/mute/list [GET]
// @success 200 {object} user.GroupMuteListResp "请求返回"
func GroupMuteList(c *gin.Context) {
	req := new(user.GroupMuteListReq)
	ctx := ginproxy.GetCtx(c)
	req.GroupId = ginproxy.GetInt64Query(c, "group_id", 0)

	if req.GroupId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(user.GroupMuteListResp)
	err := qgrpc.Call(ctx, method.UserGroupMuteList, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 更新群组名称
// @description 更新群组名称
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupUpdateNameReq true "请求参数" required
// @router  /group/update/name [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupUpdateName(c *gin.Context) {
	req := new(user.GroupUpdateNameReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}

	if req.GroupId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupUpdateName, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 更新群通知
// @description 更新群通知
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupUpdateNoticeReq true "请求参数" required
// @router  /group/update/notice [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupUpdateNotice(c *gin.Context) {
	req := new(user.GroupUpdateNoticeReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}

	if req.GroupId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupUpdateNotice, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 更新群头像
// @description 更新群头像
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupUpdateAvatarReq true "请求参数" required
// @router  /group/update/avatar [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupUpdateAvatar(c *gin.Context) {
	req := new(user.GroupUpdateAvatarReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}

	if req.GroupId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupUpdateAvatar, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}

// @summary 修改群加入类型
// @description 修改群加入类型
// @accept  json
// @tags    group
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.GroupUpdateJoinTypeReq true "请求参数" required
// @router  /group/update/jointype [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func GroupUpdateJoinType(c *gin.Context) {
	req := new(user.GroupUpdateJoinTypeReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}

	if req.GroupId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserGroupUpdateJoinType, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	if resp.GetErr() != nil {
		qlog.Ctx(ctx).Error(req, resp.GetErr())
		ginproxy.FormatError(c, resp.GetErr())
		return
	}
	ginproxy.Ret(c, resp)
}
