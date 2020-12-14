package logic

import (
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qim/server/ws_server/config"
)

const (
	CreatedOn = "created_on"
)

func GetRegistKey() string {
	return qgrpc.GetRegisterKey(config.GetOp().ServerName)
}
