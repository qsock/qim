package kvconfig

import "github.com/qsock/qf/store/kv"

var (
	//正式环境的codis配置
	KvDefaultConfig = &kv.Config{
		Addrs:       []string{"127.0.0.1:2379"},
		PoolSize:    200,
		ReadTimeout: 200,
	}
)

var (
	//测试环境的codis配置
	KvDefaultDevConfig = &kv.Config{
		Addrs:       []string{"127.0.0.1:6379"},
		PoolSize:    5,
		ReadTimeout: 200,
	}
)
