package config

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
	// websocket 暴露出去的port
	Ws   *Ws           `toml:"ws"`
	Port int           `toml:"port"`
	Env  string        `toml:"env"`
	Op   *qgrpc.Config `toml:"op"`

	LogType  string                 `toml:"logtype"`
	LogLevel types.LEVEL            `toml:"loglevel"`
	Log      map[string]interface{} `toml:"log"`
}

type Ws struct {
	PublicAddr string `toml:"public_addr"`
	Port       int    `toml:"port"`
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

func GetWs() *Ws {
	if config != nil {
		return config.Ws
	}
	return nil
}
