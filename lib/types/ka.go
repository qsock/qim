package types

import (
	"github.com/qsock/qim/config/mq"
)

// 获取kafka的地址
func GetKaAddrs(env string) []string {
	if env == EnvDev {
		return mq.KaDevAddr
	}
	return mq.KaAddr
}
