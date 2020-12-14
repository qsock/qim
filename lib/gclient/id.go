package gclient

import (
	"context"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qim/lib/method"
	"github.com/qsock/qim/lib/proto/id"
	"github.com/qsock/qim/lib/proto/ret"
)

func RegistKey(ctx context.Context, key string, offset int64, size int32) error {
	if len(key) == 0 || offset == 0 || size == 0 {
		return nil
	}
	req := new(id.RegistKeyReq)
	req.Keys = []*id.KeyItem{{Key: key, Offset: offset, Size_: size}}
	resp := new(ret.EmptyResp)
	if err := qgrpc.CallIn500ms(ctx, method.IdRegistKey, req, resp); err != nil {
		qlog.Get().Ctx(ctx).Error(key, offset, size, err)
		return err
	}
	return nil
}

func GenDbId(ctx context.Context, key string) int64 {
	if len(key) == 0 {
		return 0
	}
	req := new(id.GenDbIdReq)
	req.Key = key
	resp := new(id.GenDbIdResp)
	if err := qgrpc.CallIn200ms(ctx, method.IdGenDbId, req, resp); err != nil {
		qlog.Get().Ctx(ctx).Error(key, err)
		return 0
	}
	return resp.Id
}

func GenSnowflake() int64 {
	req := new(id.GenSnowflakeIdReq)
	resp := new(id.GenSnowflakeIdResp)
	if err := qgrpc.CallIn100ms(context.Background(), method.IdGenSnowflakeId, req, resp); err != nil {
		qlog.Get().Logger().Error(err)
		return 0
	}
	return resp.Id
}

func TraceToTime(itemId int64) int64 {
	req := new(id.SnowflakeIdToTimeReq)
	req.Id = itemId
	resp := new(id.SnowflakeIdToTimeResp)
	if err := qgrpc.CallIn100ms(context.Background(), method.IdSnowflakeIdToTime, req, resp); err != nil {
		qlog.Get().Logger().Error(err)
		return 0
	}
	return resp.UnixTime
}
