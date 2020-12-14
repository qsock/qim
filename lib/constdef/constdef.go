package constdef

import (
	"github.com/qsock/qim/lib/proto/stream"
)

const (
	FilePathAvatar = "avatar"
)

const (
	HeaderToken = "x-token"
)

type JsonRet struct {
	T stream.StreamType `json:"stream_type"`
	// 返回的内容
	Data interface{} `json:"data,omitempty"`
}
