package logic

import (
	"github.com/gin-gonic/gin"
	"github.com/qsock/qf/net/ws"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qim/lib/ginproxy"
	"time"
)

var (
	WsServer *ws.Server
)

func init() {
	WsServer = ws.New()
	WsServer.HandleConnect(OnConnect)
	WsServer.HandleDisconnect(OnDisConnect)
	WsServer.HandleMessage(OnMessage)
	WsServer.HandleSent(OnSent)
	WsServer.HandleError(OnError)
	WsServer.HandleClose(OnClose)
}

func HandleWs(c *gin.Context) {
	ctx := ginproxy.GetCtx(c)
	if err := WsServer.HandleRequestWithKeys(c.Writer, c.Request, map[string]interface{}{CreatedOn: time.Now().Unix()}); err != nil {
		qlog.Ctx(ctx).Error(err)
	}
}
