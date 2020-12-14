package ginproxy

import (
	"github.com/gin-gonic/gin"
	"github.com/qsock/qim/lib/codes"
	"github.com/qsock/qim/lib/proto/errmsg"
	"net/http"
)

type Resp struct {
	// 返回的错误码
	Code int32 `json:"code"`
	// 返回的错误信息
	Message string `json:"message,omitempty"`
	// 返回的内容
	Data interface{} `json:"data,omitempty"`
}

func Ret(c *gin.Context, resp interface{}) {
	c.JSON(http.StatusOK, Resp{Data: resp})
}

func OK(c *gin.Context) {
	Ret(c, "OK")
}

func FormatError(c *gin.Context, e *errmsg.ErrMsg) {
	c.Header("x-time-elapsed", "10")
	c.JSON(http.StatusOK, Resp{Code: e.Code, Message: e.Message})
	c.Abort()
}

// 参数错误
func ParameterError(c *gin.Context) {
	FormatError(c, codes.Error(codes.ErrorParameter))
}

func ServerError(c *gin.Context) {
	FormatError(c, codes.Error(codes.ErrorCommon))
}
