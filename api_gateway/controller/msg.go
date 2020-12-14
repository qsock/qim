package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qim/lib/ginproxy"
	"github.com/qsock/qim/lib/method"
	"github.com/qsock/qim/lib/proto/msg"
)

// @summary 得到系统消息
// @description 得到系统消息
// @accept  json
// @tags    msg
// @produce json
// @param   x-token header string true "校验的header" required
// @param   page query int true "总页数，从0开始"
// @param   page_size query int false "每页数量，默认200条"
// @router  /msg/sys [GET]
// @success 200 {object} msg.GetSysMsgResp "请求返回"
func GetSysMsg(c *gin.Context) {
	req := new(msg.GetSysMsgReq)
	req.PageSize = ginproxy.GetInt32Query(c, "page_size", 200)
	req.Page = ginproxy.GetInt32Query(c, "page", 0)
	ctx := ginproxy.GetCtx(c)
	req.UserId = ginproxy.GetUserId(c)
	resp := new(msg.GetSysMsgResp)
	err := qgrpc.Call(ctx, method.MsgGetSysMsg, req, resp)
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

// @summary 发送消息
// @description 发送消息
// @accept  json
// @tags    msg
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body msg.MsgReq true "请求参数" required
// @router  /msg [POST]
// @success 200 {object} msg.MsgResp "请求返回"
func Msg(c *gin.Context) {
	req := new(msg.MsgReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if req.M == nil || req.M.RecvId < 1 {
		ginproxy.ParameterError(c)
		return
	}
	meta := ginproxy.GetMeta(c)
	req.M.SenderId = ginproxy.GetUserId(c)
	req.M.Device = meta.Device
	resp := new(msg.MsgResp)
	err := qgrpc.Call(ctx, method.MsgMsg, req, resp)
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

// @summary 撤回自己的消息
// @description 撤回自己的消息
// @accept  json
// @tags    msg
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body msg.RevertSelfMsgReq true "请求参数" required
// @router  /msg/revert [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func RevertSelfMsg(c *gin.Context) {
	req := new(msg.RevertSelfMsgReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if req.UserId < 2 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(msg.MsgResp)
	err := qgrpc.Call(ctx, method.MsgRevertSelfMsg, req, resp)
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

// @summary 管理员撤回消息
// @description 管理员撤回消息
// @accept  json
// @tags    msg
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body msg.ManagerChatMsgRevertReq true "请求参数" required
// @router  /msg/revert/by-manager [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func ManagerChatMsgRevert(c *gin.Context) {
	req := new(msg.ManagerChatMsgRevertReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if req.UserId < 2 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(msg.MsgResp)
	err := qgrpc.Call(ctx, method.MsgManagerChatMsgRevert, req, resp)
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

// @summary 交换消息服务器令牌
// @description 交换消息服务器令牌
// @accept  json
// @tags    msg
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body msg.ExchangeReq true "请求参数" required
// @router  /msg/exchange [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func Exchange(c *gin.Context) {
	req := new(msg.ExchangeReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if req.UserId < 2 || len(req.Key) < 2 || len(req.Uuid) < 2 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(msg.MsgResp)
	err := qgrpc.Call(ctx, method.MsgExchange, req, resp)
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
