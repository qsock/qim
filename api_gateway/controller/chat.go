package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qim/lib/ginproxy"
	"github.com/qsock/qim/lib/method"
	"github.com/qsock/qim/lib/proto/msg"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/util"
	"strings"
)

// @summary 标记会话已读
// @description 标记会话已读
// @accept  json
// @tags    chat
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body msg.MarkChatReadReq true "请求参数" required
// @router  /chat/mark-read [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func MarkChatRead(c *gin.Context) {
	req := new(msg.MarkChatReadReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if len(req.ChatId) < 2 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.MsgMarkChatRead, req, resp)
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

// @summary 会话置顶
// @description 会话置顶
// @accept  json
// @tags    chat
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body msg.ChatAheadReq true "请求参数" required
// @router  /chat/ahead [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func ChatAhead(c *gin.Context) {
	req := new(msg.ChatAheadReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if len(req.ChatId) < 2 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.MsgChatAhead, req, resp)
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

// @summary 创建会话
// @description 创建会话
// @accept  json
// @tags    chat
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body msg.ChatTouchReq true "请求参数" required
// @router  /chat/touch [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func ChatTouch(c *gin.Context) {
	req := new(msg.ChatTouchReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if req.RecverId < 2 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.MsgChatTouch, req, resp)
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

// @summary 移除会话
// @description 移除会话
// @accept  json
// @tags    chat
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body msg.ChatRemoveReq true "请求参数" required
// @router  /chat/remove [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func ChatRemove(c *gin.Context) {
	req := new(msg.ChatRemoveReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if len(req.ChatId) < 2 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.MsgChatRemove, req, resp)
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

// @summary 会话免打扰
// @description 会话免打扰
// @accept  json
// @tags    chat
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body msg.ChatMuteReq true "请求参数" required
// @router  /chat/mute [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func ChatMute(c *gin.Context) {
	req := new(msg.ChatMuteReq)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if len(req.ChatId) < 2 {
		ginproxy.ParameterError(c)
		return
	}
	req.UserId = ginproxy.GetUserId(c)
	resp := new(ret.EmptyResp)
	err := qgrpc.Call(ctx, method.MsgChatMute, req, resp)
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

// @summary 会话的id列表
// @description 会话的id列表,因为会话是个经常变化的属性，所以每次进入app，只需要拉取会话id，跟本地对比
// @accept  json
// @tags    chat
// @produce json
// @param   x-token header string true "校验的header" required
// @router  /chat/ids [GET]
// @success 200 {object} msg.ChatIdsResp "请求返回"
func ChatIds(c *gin.Context) {
	req := new(msg.ChatIdsReq)
	ctx := ginproxy.GetCtx(c)
	req.UserId = ginproxy.GetUserId(c)
	resp := new(msg.ChatIdsResp)
	err := qgrpc.Call(ctx, method.MsgChatIds, req, resp)
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

// @summary 查看自己的会话列表
// @description 查看自己的会话列表
// @accept  json
// @tags    chat
// @produce json
// @param   x-token header string true "校验的header" required
// @param   page query int true "总页数，从0开始"
// @param   page_size query int false "每页数量，默认200条"
// @router  /chat/list [GET]
// @success 200 {object} msg.ChatByUidsResp "请求返回"
func ChatByUids(c *gin.Context) {
	req := new(msg.ChatByUidsReq)
	req.PageSize = ginproxy.GetInt32Query(c, "page_size", 200)
	req.Page = ginproxy.GetInt32Query(c, "page", 0)
	ctx := ginproxy.GetCtx(c)
	req.UserId = ginproxy.GetUserId(c)
	resp := new(msg.ChatByUidsResp)
	err := qgrpc.Call(ctx, method.MsgChatByUids, req, resp)
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

// @summary 通过会话id获取会话
// @description 通过会话id获取会话
// @accept  json
// @tags    chat
// @produce json
// @param   x-token header string true "校验的header" required
// @param   ids query string true "需要查看的id列表,逗号分割所有的id,例如:1234,12511,1100"
// @router  /chat/by-ids [GET]
// @success 200 {object} msg.ChatByIdsResp "请求返回"
func ChatByIds(c *gin.Context) {
	req := new(msg.ChatByIdsReq)
	req.Ids = strings.Split(c.Query("ids"), ",")
	ctx := ginproxy.GetCtx(c)
	req.UserId = ginproxy.GetUserId(c)
	if len(req.Ids) == 0 {
		ginproxy.ParameterError(c)
		return
	}
	resp := new(msg.ChatByIdsResp)
	err := qgrpc.Call(ctx, method.MsgChatByIds, req, resp)
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

// @summary 查找聊天记录id
// @description 查找聊天记录id
// @accept  json
// @tags    chat
// @produce json
// @param   x-token header string true "校验的header" required
// @param   chat_id query string true "会话id"
// @param   page query int true "总页数，从0开始"
// @param   page_size query int false "每页数量，默认200条"
// @router  /chat/record/ids [GET]
// @success 200 {object} msg.ChatRecordIdsResp "请求返回"
func ChatRecordIds(c *gin.Context) {
	req := new(msg.ChatRecordIdsReq)
	req.PageSize = ginproxy.GetInt32Query(c, "page_size", 200)
	req.Page = ginproxy.GetInt32Query(c, "page", 0)
	req.ChatId = c.Query("chat_id")
	ctx := ginproxy.GetCtx(c)
	req.UserId = ginproxy.GetUserId(c)
	if len(req.ChatId) == 0 {
		ginproxy.ParameterError(c)
		return
	}
	resp := new(msg.ChatRecordIdsResp)
	err := qgrpc.Call(ctx, method.MsgChatRecordIds, req, resp)
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

// @summary 通过聊天记录id查找聊天记录
// @description 通过聊天记录id查找聊天记录
// @accept  json
// @tags    chat
// @produce json
// @param   x-token header string true "校验的header" required
// @param   chat_id query string true "会话id"
// @param   ids query string true "需要查看的id列表,逗号分割所有的id,例如:1234,12511,1100"
// @router  /chat/record/by-ids [GET]
// @success 200 {object} msg.ChatRecordByIdsResp "请求返回"
func ChatRecordByIds(c *gin.Context) {
	req := new(msg.ChatRecordByIdsReq)
	req.Ids = util.StrToInt64s(c.Query("ids"))
	req.ChatId = c.Query("chat_id")
	ctx := ginproxy.GetCtx(c)
	req.UserId = ginproxy.GetUserId(c)
	if len(req.ChatId) == 0 ||
		len(req.Ids) == 0 {
		ginproxy.ParameterError(c)
		return
	}
	resp := new(msg.ChatRecordByIdsResp)
	err := qgrpc.Call(ctx, method.MsgChatRecordByIds, req, resp)
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
