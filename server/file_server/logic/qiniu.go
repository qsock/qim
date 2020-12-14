package logic

import (
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/qsock/qim/config/common"
	"github.com/qsock/qim/server/file_server/config"
)

// 获得鉴权对象
func Credential() *qbox.Mac {
	ak, sk := common.GetQiniuConf(config.GetEnv())
	return qbox.NewMac(ak, sk)
}

func GetQiniuToken(path string) string {
	bucket := common.GetQiniuBucketByPath(path)
	policy := storage.PutPolicy{
		CallbackBodyType: "application/json",
		CallbackURL:      common.GetQiniuCallbackUrl(),
		Scope:            bucket,
		CallbackBody:     `{"key":"$(key)","hash":"$(etag)","fsize":$(fsize),"bucket":"$(bucket)","path":"$(x:path)","uid":"$(x:uid)"}`,
	}
	upToken := policy.UploadToken(Credential())
	return upToken
}
