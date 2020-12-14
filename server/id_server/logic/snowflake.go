package logic

import (
	"context"
	"github.com/qsock/qf/util/snowflake"
	"github.com/qsock/qim/lib/proto/id"
)

func SnowflakeIdToTime(ctx context.Context, req *id.SnowflakeIdToTimeReq) (*id.SnowflakeIdToTimeResp, error) {
	resp := new(id.SnowflakeIdToTimeResp)
	resp.UnixTime = snowflake.ToTimeUnix(req.Id)
	return resp, nil
}

func GenSnowflakeId(ctx context.Context, req *id.GenSnowflakeIdReq) (*id.GenSnowflakeIdResp, error) {
	resp := new(id.GenSnowflakeIdResp)
	resp.Id = snowflake.NextId()
	return resp, nil
}
