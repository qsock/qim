package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qim/lib/ginproxy"
	"github.com/qsock/qim/lib/method"
	"github.com/qsock/qim/lib/proto/file"
)

// @summary 获取省市
// @description 获取省市
// @accept  json
// @tags    file
// @produce json
// @router  /file/locations [GET]
// @success 200 {object} file.GetProvinceAndCityResp "请求返回"
func Locations(c *gin.Context) {
	req := new(file.GetProvinceAndCityReq)
	resp := new(file.GetProvinceAndCityResp)
	ctx := ginproxy.GetCtx(c)
	err := qgrpc.Call(ctx, method.FileGetProvinceAndCity, req, resp)
	if err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ServerError(c)
		return
	}
	ginproxy.Ret(c, resp)
}
