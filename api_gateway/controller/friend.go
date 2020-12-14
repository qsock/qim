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

// @summary 更新备注名称
// @description 更新备注名称
// @accept  json
// @tags    friend
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.FriendMarknameUpdateReq true "请求参数" required
// @router  /friend/mark-name [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func FriendMarknameUpdate(c *gin.Context) {
	req := new(user.FriendMarknameUpdateReq)
	resp := new(ret.EmptyResp)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if len(req.Name) == 0 || req.FriendId <= 0 {
		ginproxy.ParameterError(c)
		return
	}
	err := qgrpc.Call(ctx, method.UserFriendMarknameUpdate, req, resp)
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

// @summary 好友id列表
// @description 好友id列表，每次进入app的时候调用一次，跟本地的db进行比对，如果有多余或者少的，就需要在本地处理，如果id比本地多了，则需要删除，如果id比本地少了，就需要添加
// @accept  json
// @tags    friend
// @produce json
// @param   x-token header string true "校验的header" required
// @router  /friend/list [GET]
// @success 200 {object} user.FriendIdsResp "请求返回"
func FriendIds(c *gin.Context) {
	req := new(user.FriendIdsReq)
	resp := new(user.FriendIdsResp)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	err := qgrpc.Call(ctx, method.UserFriendIds, req, resp)
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

// @summary 通过好友id，获取好友信息
// @description 通过好友id，获取好友信息
// @accept  json
// @tags    friend
// @produce json
// @param   x-token header string true "校验的header" required
// @param   ids query string true "需要查看的id列表,逗号分割所有的id,例如:1234,12511,1100"
// @router  /friend/by-ids [GET]
// @success 200 {object} user.FriendByIdsResp "请求返回"
func FriendByIds(c *gin.Context) {
	req := new(user.FriendByIdsReq)
	ids := c.Query("ids")
	req.Ids = util.StrToInt64s(ids)
	req.UserId = ginproxy.GetUserId(c)
	if len(req.Ids) == 0 {
		ginproxy.ParameterError(c)
		return
	}

	resp := new(user.FriendByIdsResp)
	ctx := ginproxy.GetCtx(c)
	err := qgrpc.Call(ctx, method.UserFriendByIds, req, resp)
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

// @summary 获取自己的好友列表
// @description 获取自己的好友列表
// @accept  json
// @tags    friend
// @produce json
// @param   x-token header string true "校验的header" required
// @param   page query int true "总页数，从0开始"
// @param   page_size query int false "每页数量，默认30条"
// @router  /friend/list [GET]
// @success 200 {object} user.FriendsByUidResp "请求返回"
func FriendsByUid(c *gin.Context) {
	req := new(user.FriendsByUidReq)
	req.Page = ginproxy.GetInt32Query(c, "page", 0)
	req.PageSize = ginproxy.GetInt32Query(c, "page_size", 30)
	req.UserId = ginproxy.GetUserId(c)
	ctx := ginproxy.GetCtx(c)
	resp := new(user.FriendsByUidResp)
	err := qgrpc.Call(ctx, method.UserFriendsByUid, req, resp)
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

// @summary 删除好友
// @description 删除好友
// @accept  json
// @tags    friend
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.FriendDelReq true "请求参数" required
// @router  /friend/del [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func FriendDel(c *gin.Context) {
	req := new(user.FriendDelReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	if req.FriendId <= 0 {
		ginproxy.ParameterError(c)
		return
	}

	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.UserFriendDel, req, resp)
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
