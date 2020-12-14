package main

import (
	"context"
	"github.com/qsock/qim/lib/proto/id"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/server/id_server/logic"
)

type Server struct{}

func (server *Server) SnowflakeIdToTime(ctx context.Context, req *id.SnowflakeIdToTimeReq) (*id.SnowflakeIdToTimeResp, error) {
	return logic.SnowflakeIdToTime(ctx, req)
}

func (server *Server) GenSnowflakeId(ctx context.Context, req *id.GenSnowflakeIdReq) (*id.GenSnowflakeIdResp, error) {
	return logic.GenSnowflakeId(ctx, req)
}

func (server *Server) RegistKey(ctx context.Context, req *id.RegistKeyReq) (*ret.EmptyResp, error) {
	return logic.RegistKey(ctx, req)
}

func (server *Server) GenDbId(ctx context.Context, req *id.GenDbIdReq) (*id.GenDbIdResp, error) {
	return logic.GenDbId(ctx, req)
}

func (*Server) Ping(ctx context.Context, req *ret.NoArgs) (*ret.NoArgs, error) {
	return new(ret.NoArgs), nil
}
