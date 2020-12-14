package ginproxy

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/store/ka"
	"github.com/qsock/qim/config/mq"
	"github.com/qsock/qim/lib/proto/event"
	"github.com/qsock/qim/lib/proto/model"
	"strings"
	"time"
)

func AccessLog(c *gin.Context) {
	path := c.Request.URL.String()
	if path == "/" || strings.HasPrefix(path, "/swagger") {
		c.Next()
		return
	}

	//记录总请求时间
	start := time.Now()
	headers, _ := json.Marshal(c.Request.Header)
	if strings.Contains(strings.ToLower(c.GetHeader("Content-Type")), "application/json") {
		CopyReqBody(c)
	}

	c.Writer.Header().Set("x-trace", GetTraceId(c))
	c.Next()

	report := new(event.HttpReport)
	report.Meta = GetMeta(c)
	report.Method = c.Request.Method
	report.StatusCode = int32(c.Writer.Status())
	report.EndOn = time.Now().UnixNano()
	report.Path = path
	report.Alias = GetAlias(c)
	report.Headers = string(headers)
	report.Req = string(GetReqBody(c))
	report.Resp = GetResp(c)
	if report.StatusCode != 404 && report.StatusCode != 204 {
		ka.TopicEvent(mq.TopicLogTrace, mq.ELogTraceHttp, report)
	}

	//状态码|请求所花的时间|path|method|headers|ip||
	ks := "%d|t:%s|path:%s|method:%s|" +
		"headers:%s|ip:%s|"
	kv := []interface{}{
		report.StatusCode, time.Since(start).String(), report.Path, report.Method,
		string(headers), report.Meta.UserIp,
	}

	if report.Meta.Device != model.Device_DeviceFalse {
		ks = ks + "device:%s|"
		kv = append(kv, report.Meta.Device.String())
	}
	if len(report.Meta.AppName) > 0 {
		ks = ks + "appname:%s|"
		kv = append(kv, report.Meta.AppName)
	}

	if len(report.Meta.AppVersion) > 0 {
		ks = ks + "appver:%s|"
		kv = append(kv, report.Meta.AppVersion)
	}
	if len(report.Meta.DeviceId) > 0 {
		ks = ks + "deviceId:%s|"
		kv = append(kv, report.Meta.DeviceId)
	}
	if len(report.Meta.Lat) > 0 {
		ks = ks + "lat:%s|"
		kv = append(kv, report.Meta.Lat)
	}
	if len(report.Meta.Lng) > 0 {
		ks = ks + "lng:%s|"
		kv = append(kv, report.Meta.Lng)
	}
	if report.Meta.UserId > 0 {
		ks = ks + "uid:%s|"
		kv = append(kv, report.Meta.UserId)
	}
	if len(report.Req) > 0 {
		ks = ks + "req=%+s|"
		kv = append(kv, report.Req)
	}
	if len(report.Resp) > 0 {
		ks = ks + "req=%+s|"
		kv = append(kv, report.Resp)
	}

	if report.StatusCode >= 400 {
		qlog.Get().Ctx(GetCtx(c)).Errorf(ks, kv...)
	} else {
		qlog.Get().Ctx(GetCtx(c)).Infof(ks, kv...)
	}
}
