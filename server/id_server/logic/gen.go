package logic

import (
	"context"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/lib/codes"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/proto/id"
	"github.com/qsock/qim/lib/proto/ret"
	"strings"
)

func RegistKey(ctx context.Context, req *id.RegistKeyReq) (*ret.EmptyResp, error) {
	for _, item := range req.Keys {
		if err := registKey(ctx, item); err != nil {
			return nil, err
		}
	}
	resp := new(ret.EmptyResp)
	return resp, nil
}

func registKey(ctx context.Context, item *id.KeyItem) error {
	lock.RLock()
	_, ok := chanMap[item.Key]
	if ok {
		lock.RUnlock()
		return nil
	}
	lock.RUnlock()
	{
		isql := "insert into `keys` (`k`,`offset`,`size`) values(?,?,?)"
		conn := dao.GetConn(dbconfig.DbIdA)
		if _, err := conn.Exec(isql, item.Key, item.Offset, item.Size_); err != nil {
			if !strings.Contains(err.Error(), "Duplicate") {
				qlog.Ctx(ctx).Error(err, item)
				return err
			}
		}
	}
	return genDb(ctx, item)
}

func genDb(ctx context.Context, item *id.KeyItem) error {
	lock.RLock()
	var ct int
	for _, db := range list() {
		ssql := "select count(1) from id where k=?"
		var flag bool
		if err := db.db.QueryRow(ssql, item.Key).Scan(&flag); err != nil {
			qlog.Ctx(ctx).Error(item, err)
			continue
		}
		if flag {
			ct++
		}
	}
	if ct != 0 && ct != 4 {
		qlog.Ctx(ctx).Fatal(item, ct)
	}
	lock.RUnlock()

	lock.Lock()
	defer lock.Unlock()

	if ct == 0 {
		for item.Offset%4 != 0 {
			item.Offset++
		}

		isql := "insert into id (`k`,`id`) values (?,?)"
		for _, db := range list() {
			dbOffset := item.Offset + db.idx
			qlog.Ctx(ctx).Info(db.idx, dbOffset, item)
			_, err := db.db.Exec(isql, item.Key, dbOffset)
			if err != nil {
				qlog.Ctx(ctx).Error(db.idx, dbOffset, item, err)
				return err
			}
		}
	}

	chanMap[item.Key] = make(chan int64, item.Size_)
	go gen(item.Key)
	return nil
}

func GenDbId(ctx context.Context, req *id.GenDbIdReq) (*id.GenDbIdResp, error) {
	resp := new(id.GenDbIdResp)
	lock.RLock()
	defer lock.RUnlock()
	chanId, ok := chanMap[req.Key]
	if !ok {
		resp.Err = codes.Error(codes.ErrorParameter)
		return resp, nil
	}
	id := <-chanId
	resp.Id = id
	return resp, nil
}
