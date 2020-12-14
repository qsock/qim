package ginproxy

import (
	"github.com/gin-gonic/gin"
	"github.com/qsock/qf/util/snowflake"
	"github.com/qsock/qim/lib/proto/model"
	"strconv"
	"time"
)

// 需要在nginx设置
//server{
//set $trace_id "${request_id}";
//if ($http_x_atrace_id != "" ){
//set $trace_id "${http_x_atrace_id}";
//}
//add_header trace_id $trace_id;
//#配置反向代理时使用
//proxy_set_header x-atrace-id $trace_id;
//...
//}
func Parse(c *gin.Context) {
	meta := new(model.RequestMeta)
	meta.TraceId = c.GetHeader("trace_id")
	if len(meta.TraceId) == 0 {
		meta.TraceId = strconv.FormatInt(snowflake.NextId(), 10)
	}
	{
		h := c.GetHeader("device")
		d, _ := strconv.ParseInt(h, 10, 32)
		meta.Device = model.Device(d)
	}
	meta.AppName = c.GetHeader("appname")
	meta.AppVersion = c.GetHeader("ver")
	meta.DeviceId = c.GetHeader("device-id")
	meta.UserIp = c.ClientIP() //这里需要设置nginx把用户ip传进来
	meta.Lat = c.GetHeader("lat")
	meta.Lng = c.GetHeader("lng")
	meta.CreatedOn = time.Now().UnixNano()
	c.Set(kMeta, meta)
}
