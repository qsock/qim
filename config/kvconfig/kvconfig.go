package kvconfig

import "github.com/qsock/qf/store/kv"

const (
	KvDefault string = "KvDefault"
)

// prod config
var ConfigMap = map[string]*kv.Config{
	KvDefault: KvDefaultConfig,
}

// dev config
var ConfigDevMap = map[string]*kv.Config{
	KvDefault: KvDefaultDevConfig,
}
