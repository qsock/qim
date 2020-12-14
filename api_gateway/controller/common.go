package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/qsock/qim/lib/ginproxy"
	"net/url"
	"time"
)

// @summary 得到系统时间。
// @tags    common
// @produce json
// @router  /time [GET]
// @success 200 {object} ginproxy.Resp "resp"
func ServerTime(c *gin.Context) {
	ginproxy.Ret(c, gin.H{"t": time.Now().Unix()})
}

// @summary 得到proto文件
// @description 请用  https://tool.chinaz.com/tools/urlencode.aspx 进行文本的decode
// @tags    common
// @produce json
// @param   name query string true "文件名:stream,errmsg,model"
// @router  /proto [GET]
// @success 200 {object} ginproxy.Resp "resp"
func Proto(c *gin.Context) {
	name := c.Query("name")
	var content string
	if name == "stream" {
		content = streamProto
	} else if name == "errmsg" {
		content = errmsgProto
	} else if name == "model" {
		content = modelProto
	}
	ginproxy.Ret(c, gin.H{"proto": url.PathEscape(content)})
}
