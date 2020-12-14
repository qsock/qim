package ginproxy

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qim/lib/proto/model"
	"io/ioutil"
	"strconv"
)

const (
	kReqBody = "req-body"
	kResp    = "resp"
	kAlias   = "alias"
	kMeta    = "meta"
)

func GetBoolQuery(c *gin.Context, query string, defaultValue bool) bool {
	v, ok := c.GetQuery(query)
	if !ok {
		return defaultValue
	}
	i, err := strconv.ParseBool(v)
	if err != nil {
		return defaultValue
	}
	return i
}

func GetInt32Query(c *gin.Context, query string, defaultValue int32) int32 {
	return int32(GetInt64Query(c, query, int64(defaultValue)))
}

func GetInt64Query(c *gin.Context, query string, defaultValue int64) int64 {
	v, ok := c.GetQuery(query)
	if !ok {
		return defaultValue
	}

	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return defaultValue

	}
	return i
}

func GetFloat64Query(c *gin.Context, query string, defaultValue float64) float64 {
	v, ok := c.GetQuery(query)
	if !ok {
		return defaultValue
	}

	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return defaultValue
	}
	return f
}

func GetInt64Path(c *gin.Context, key string) (int64, error) {
	v := c.Param(key)
	return strconv.ParseInt(v, 10, 64)
}

func GetInt32Path(c *gin.Context, key string) (int32, error) {
	v := c.Param(key)
	i, err := strconv.ParseInt(v, 10, 32)
	return int32(i), err
}

func GetTraceId(c *gin.Context) string {
	if o := GetMeta(c); o != nil {
		return o.TraceId
	}
	return ""
}

func GetUserId(c *gin.Context) int64 {
	if o := GetMeta(c); o != nil {
		return o.UserId
	}
	return 0
}

func GetDeviceId(c *gin.Context) string {
	if o := GetMeta(c); o != nil {
		return o.DeviceId
	}
	return ""
}

func GetAppName(c *gin.Context) string {
	if o := GetMeta(c); o != nil {
		return o.AppName
	}
	return ""
}

func GetAlias(c *gin.Context) string {
	if v, ok := c.Get(kAlias); ok {
		return v.(string)
	}
	return ""
}

func SetAlias(c *gin.Context, uri string) {
	c.Set(kAlias, uri)
}

func GetMeta(c *gin.Context) *model.RequestMeta {
	if v, ok := c.Get(kMeta); ok {
		if o, ok := v.(*model.RequestMeta); ok {
			return o
		}
	}
	return nil
}

func GetReqBody(c *gin.Context) []byte {
	if b, ok := c.Get(kReqBody); ok {
		return b.([]byte)
	}
	return nil
}

func CopyReqBody(c *gin.Context) {
	rawBytes, _ := ioutil.ReadAll(c.Request.Body)
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawBytes))
	c.Set(kReqBody, rawBytes)
}

func GetResp(c *gin.Context) string {
	return c.GetString(kResp)
}

func SetResp(c *gin.Context, val interface{}) {
	b, _ := json.Marshal(val)
	c.Set(kResp, string(b))
}

func GetCtx(c *gin.Context) context.Context {
	traceId := GetTraceId(c)
	return context.WithValue(context.Background(), qgrpc.MetaKey, map[string]string{"trace_id": traceId})
}
