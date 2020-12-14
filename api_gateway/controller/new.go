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
)

// @summary 新的朋友&群组列表
// @description 新的朋友&群组列表
// @accept  json
// @tags    friend
// @produce json
// @param   x-token header string true "校验的header" required
// @param   page query int true "总页数，从0开始"
// @param   page_size query int false "每页数量，默认30条"
// @router  /friend/new/apply/list [GET]
// @success 200 {object} user.NewApplyListResp "请求返回"
func NewApplyList(c *gin.Context) {
	// @param   entity body user.FriendDelReq true "请求参数" required
	req := new(user.NewApplyUserListReq)
	req.Page = ginproxy.GetInt32Query(c, "page", 0)
	req.PageSize = ginproxy.GetInt32Query(c, "page_size", 30)
	req.UserId = ginproxy.GetUserId(c)
	ctx := ginproxy.GetCtx(c)
	resp := new(user.NewApplyListResp)
	err := qgrpc.Call(ctx, method.UserNewApplyUserList, req, resp)
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

// @summary 申请添加新好友
// @description 申请添加新好友
// @accept  json
// @tags    friend
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.FriendNewApplyReq true "请求参数" required
// @router  /friend/new/apply [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func FriendNewApply(c *gin.Context) {
	req := new(user.FriendNewApplyReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if req.FriendId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserFriendNewApply, req, resp)
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

// @summary 同意添加好友
// @description 同意添加好友
// @accept  json
// @tags    friend
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.FriendNewAgreeReq true "请求参数" required
// @router  /friend/new/agree [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func FriendNewAgree(c *gin.Context) {
	req := new(user.FriendNewAgreeReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if req.FriendId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserFriendNewAgree, req, resp)
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

// @summary 拒绝添加好友
// @description 拒绝添加好友
// @accept  json
// @tags    friend
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.FriendNewRejectReq true "请求参数" required
// @router  /friend/new/reject [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func FriendNewReject(c *gin.Context) {
	req := new(user.FriendNewRejectReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if req.FriendId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserFriendNewReject, req, resp)
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

// @summary 删除新好友列表
// @description 删除新好友列表的其中一个内容
// @accept  json
// @tags    friend
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.FriendNewDelReq true "请求参数" required
// @router  /friend/new/del [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func FriendNewDel(c *gin.Context) {
	req := new(user.FriendNewDelReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if req.Id <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserFriendNewDel, req, resp)
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

// @summary 忽略这个好友
// @description 忽略这个好友的申请
// @accept  json
// @tags    friend
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.FriendNewIgnoreReq true "请求参数" required
// @router  /friend/new/ignore [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func FriendNewIgnore(c *gin.Context) {
	req := new(user.FriendNewIgnoreReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if req.FriendId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserFriendNewDel, req, resp)
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
