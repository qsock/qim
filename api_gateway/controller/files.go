package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qim/lib/ginproxy"
	"github.com/qsock/qim/lib/method"
	"github.com/qsock/qim/lib/proto/file"
)

// @summary 获取所有系统头像
// @description 获取所有系统头像
// @accept  json
// @tags    file
// @produce json
// @router  /file/sys/avatars [GET]
// @success 200 {object} file.GetSysAvatarsResp "请求返回"
func SysAvatars(c *gin.Context) {
	req := new(file.GetSysAvatarsReq)
	resp := new(file.GetSysAvatarsResp)
	ctx := ginproxy.GetCtx(c)
	err := qgrpc.Call(ctx, method.FileGetSysAvatars, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	ginproxy.Ret(c, resp)
}
