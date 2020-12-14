package ginproxy

import (
	"github.com/gin-gonic/gin"
	"github.com/qsock/qf/qlog"
	"net/http/httputil"
	"runtime/debug"
)

func Recovery(c *gin.Context) {
	ctx := GetCtx(c)
	defer func() {
		if err := recover(); err != nil {
			dumpStr, _ := httputil.DumpRequest(c.Request, false)
			qlog.Get().Ctx(ctx).Fatalf("panic||req:%s||err:%v||stack:%s", string(dumpStr), err, string(debug.Stack()))
			c.AbortWithStatus(500)
		}
	}()
	c.Next()
}
