package dao

import (
	"github.com/go-redis/redis"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/store/kv"
	"github.com/qsock/qim/config/kvconfig"
	"github.com/qsock/qim/lib/types"
)

//获取KV链接池
func GetKvConn(name string) *redis.Client {
	kmu.Lock()
	defer kmu.Unlock()

	client, err := kv.GetRedisConn(name)
	if err == nil && client != nil {
		return client
	}

	m := kvconfig.ConfigMap
	if env == types.EnvDev {
		m = kvconfig.ConfigDevMap
	}

	kv.Add(name, m[name])
	client, err = kv.GetRedisConn(name)
	if err != nil {
		qlog.Get().Logger().Error(name, err)
	}
	return client
}
