package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qim/lib/codes"
	"github.com/qsock/qim/lib/constdef"
	"github.com/qsock/qim/lib/ginproxy"
	"github.com/qsock/qim/lib/method"
	"github.com/qsock/qim/lib/proto/passport"
	"github.com/qsock/qim/lib/proto/ret"
)

func Auth(c *gin.Context) {
	token := c.GetHeader(constdef.HeaderToken)
	if len(token) < 10 {
		ginproxy.FormatError(c, codes.Error(codes.ErrorAuthInvalid))
		return
	}
	req := new(passport.AuthReq)
	req.Meta = ginproxy.GetMeta(c)
	req.Token = token
	resp := new(ret.IntResp)
	ctx := ginproxy.GetCtx(c)
	err := qgrpc.Call(ctx, method.PassportAuth, req, resp)
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
	ginproxy.GetMeta(c).UserId = resp.Val
	c.Next()
}

// @summary 发送登陆短信
// @description 发送登陆短信
// @accept  json
// @tags    passport
// @produce json
// @param   entity body passport.SmsReq true "请求参数" required
// @router  /passport/login/sms [POST]
// @success 200 {object} ginproxy.Resp "请求返回"
func LoginSms(c *gin.Context) {
	req := new(passport.SmsReq)
	resp := new(ret.EmptyResp)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if len(req.Tel) == 0 {
		ginproxy.ParameterError(c)
		return
	}
	err := qgrpc.Call(ctx, method.PassportSms, req, resp)
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

// @summary 手机号登陆
// @description 手机号登陆
// @accept  json
// @tags    passport
// @produce json
// @param   entity body passport.TelLoginReq true "请求参数" required
// @router  /passport/login/tel [POST]
// @success 200 {object} passport.LoginResp "请求返回"
func LoginTel(c *gin.Context) {
	req := new(passport.TelLoginReq)
	resp := new(passport.LoginResp)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if len(req.Tel) == 0 ||
		len(req.Code) == 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.Meta = ginproxy.GetMeta(c)
	err := qgrpc.Call(ctx, method.PassportTelLogin, req, resp)
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

// @summary qq号登陆
// @description qq号登陆
// @accept  json
// @tags    passport
// @produce json
// @param   entity body passport.QqloginReq true "请求参数" required
// @router  /passport/login/qq [POST]
// @success 200 {object} passport.LoginResp "请求返回"
func LoginQq(c *gin.Context) {
	req := new(passport.QqloginReq)
	resp := new(passport.LoginResp)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if len(req.Token) == 0 ||
		len(req.Openid) == 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.Meta = ginproxy.GetMeta(c)
	err := qgrpc.Call(ctx, method.PassportQqlogin, req, resp)
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

// @summary 微信号登陆
// @description 微信号登陆
// @accept  json
// @tags    passport
// @produce json
// @param   entity body passport.WxLoginReq true "请求参数" required
// @router  /passport/login/wx [POST]
// @success 200 {object} passport.LoginResp "请求返回"
func LoginWx(c *gin.Context) {
	req := new(passport.WxLoginReq)
	resp := new(passport.LoginResp)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}
	if len(req.Code) == 0 {
		ginproxy.ParameterError(c)
		return
	}
	req.Meta = ginproxy.GetMeta(c)
	err := qgrpc.Call(ctx, method.PassportWxLogin, req, resp)
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

// @summary 刷新token
// @description 刷新token
// @accept  json
// @tags    passport
// @produce json
// @param   x-token header string true "校验的header" required
// @param   entity body passport.RefreshReq true "请求参数" required
// @router  /passport/refresh [POST]
// @success 200 {object} ret.StringResp "请求返回"
func RefreshToken(c *gin.Context) {
	req := new(passport.RefreshReq)
	resp := new(ret.StringResp)
	ctx := ginproxy.GetCtx(c)
	if err := c.ShouldBindWith(req, binding.JSON); err != nil {
		qlog.Ctx(ctx).Error(req, err)
		ginproxy.ParameterError(c)
		return
	}

	req.Token = c.GetHeader(constdef.HeaderToken)
	if len(req.RefreshToken) < 10 ||
		len(req.Token) < 10 {
		ginproxy.ParameterError(c)
		return
	}
	req.Meta = ginproxy.GetMeta(c)
	err := qgrpc.Call(ctx, method.PassportRefresh, req, resp)
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

// @summary 登出
// @description 登出
// @accept  json
// @tags    passport
// @produce json
// @param   x-token header string true "校验的header" required
// @router  /passport/logout [POST]
// @success 200 {object} ret.StringResp "请求返回"
func Logout(c *gin.Context) {
	req := new(passport.LogoutReq)
	resp := new(ret.StringResp)
	ctx := ginproxy.GetCtx(c)

	req.Meta = ginproxy.GetMeta(c)
	err := qgrpc.Call(ctx, method.PassportLogout, req, resp)
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
