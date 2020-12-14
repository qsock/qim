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

// @summary 查看用户信息
// @description 查看用户信息
// @accept  json
// @tags    user
// @produce json
// @param   x-token header string true "校验的header" required
// @param   id query int false "需要查看的用户id"
// @router  /user/info [GET]
// @success 200 {object} user.InfoResp "请求返回"
func UserInfo(c *gin.Context) {
	req := new(user.InfoReq)
	resp := new(user.InfoResp)
	ctx := ginproxy.GetCtx(c)
	req.UserId = ginproxy.GetInt64Query(c, "id", 0)
	if req.UserId == 0 {
		req.UserId = ginproxy.GetUserId(c)
	}
	err := qgrpc.Call(ctx, method.UserInfo, req, resp)
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

// @summary 用户上次活跃
// @description 用户上次活跃
// @accept  json
// @tags    user
// @produce json
// @param   x-token header string true "校验的header" required
// @router  /user/lastactive [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func UserLastactive(c *gin.Context) {
	req := new(user.LastactiveReq)
	req.Meta = ginproxy.GetMeta(c)
	resp := new(ret.EmptyResp)
	ctx := ginproxy.GetCtx(c)
	err := qgrpc.Call(ctx, method.UserLastactive, req, resp)
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

// @summary 更新用户信息
// @description 更新用户信息
// @accept  json
// @tags    user
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body user.UpdateReq true "请求参数" required
// @router  /user/update [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func UserUpdate(c *gin.Context) {
	req := new(user.UpdateReq)
	resp := new(ret.EmptyResp)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	err := qgrpc.Call(ctx, method.UserUpdate, req, resp)
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
