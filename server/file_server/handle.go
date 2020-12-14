package main

import (
	"context"
	"github.com/qsock/qim/lib/proto/file"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/server/file_server/logic"
)

type Server struct{}

func (*Server) Ping(ctx context.Context, req *ret.NoArgs) (*ret.NoArgs, error) {
	return new(ret.NoArgs), nil
}

func (server *Server) GetUploadToken(ctx context.Context, req *file.GetUploadTokenReq) (*file.GetUploadTokenResp, error) {
	return logic.GetUploadToken(ctx, req)
}

func (server *Server) GetUserFile(ctx context.Context, req *file.GetUserFileReq) (*file.GetUserFileResp, error) {
	return logic.GetUserFile(ctx, req)
}

func (server *Server) UploadFileByUrl(ctx context.Context, req *file.UploadFileByUrlReq) (*file.UploadFileByUrlResp, error) {
	return logic.UploadFileByUrl(ctx, req)
}

func (server *Server) UserUploadSucceed(ctx context.Context, req *file.UserUploadSucceedReq) (*ret.EmptyResp, error) {
	return logic.UserUploadSucceed(ctx, req)
}

func (server *Server) GetSysAvatars(ctx context.Context, req *file.GetSysAvatarsReq) (*file.GetSysAvatarsResp, error) {
	return logic.GetSysAvatars(ctx, req)
}

func (server *Server) GetProvinceAndCity(ctx context.Context, req *file.GetProvinceAndCityReq) (*file.GetProvinceAndCityResp, error) {
	return logic.GetProvinceAndCity(ctx, req)
}
