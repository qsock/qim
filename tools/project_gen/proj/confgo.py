
content = '''package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/qsock/qf/net/ipaddr"
	"github.com/qsock/qf/qlog/types"
	"github.com/qsock/qf/service/qgrpc"
	"os"
)

var config *Config

type Config struct {
	Port int           `toml:"port"`
	Env  string        `toml:"env"`
	Op   *qgrpc.Config `toml:"op"`

	LogType  string                 `toml:"logtype"`
	LogLevel types.LEVEL            `toml:"loglevel"`
	Log      map[string]interface{} `toml:"log"`
}

func Init(file string) error {
	config = new(Config)

	if _, err := toml.DecodeFile(file, config); err != nil {
		return err
	}

	//注册的地址
	if config.Op.Addr == "" {
		config.Op.Addr = fmt.Sprintf("%s:%d", ipaddr.GetLocalIp(), config.Port)
	}

	_ = toml.NewEncoder(os.Stdout).Encode(config)
	return nil
}

func GetConfig() *Config {
	return config
}

func GetOp() *qgrpc.Config {
	if config != nil {
		return config.Op
	}
	return nil
}

func GetAddr() string {
	if config != nil {
		return fmt.Sprintf(":%d", config.Port)
	}
	return ""
}

func GetEnv() string {
	if config != nil {
		return config.Env
	}
	return ""
}
'''

def gen(name, srv_dir) :
    with open(srv_dir+"/config.go","w") as f:
        f.write(content)