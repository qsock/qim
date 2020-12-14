package logic

import (
	"bufio"
	"context"
	"database/sql"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/proto/id"
	"sync"
	"time"
)

var (
	readPool *sync.Pool = &sync.Pool{
		New: func() interface{} {
			return bufio.NewReaderSize(nil, 128)
		}}
	c       = make(chan []*DB, 4)
	chanMap = make(map[string]chan int64)
	lock    = new(sync.RWMutex)
)

type DB struct {
	idx int64
	db  *sql.DB
}

func Init() {
	initDb()
	initKeys()
}

func initDb() {
	dbA := dao.GetConn(dbconfig.DbIdA)
	dbB := dao.GetConn(dbconfig.DbIdB)
	dbC := dao.GetConn(dbconfig.DbIdC)
	dbD := dao.GetConn(dbconfig.DbIdD)

	// 将db依次入channel
	c <- []*DB{&DB{0, dbA}, &DB{1, dbB}, &DB{2, dbC}, &DB{3, dbD}}
	c <- []*DB{&DB{1, dbB}, &DB{2, dbC}, &DB{3, dbD}, &DB{0, dbA}}
	c <- []*DB{&DB{2, dbC}, &DB{3, dbD}, &DB{0, dbA}, &DB{1, dbB}}
	c <- []*DB{&DB{3, dbD}, &DB{0, dbA}, &DB{1, dbB}, &DB{2, dbC}}

	return
}

func initKeys() {
	items := make([]*id.KeyItem, 0)
	ssql := "select `k`,`offset`,`size` from `keys`"
	conn := dao.GetConn(dbconfig.DbIdA)
	rows, err := conn.Query(ssql)
	if err != nil {
		panic(err)
	}
	for rows.Next() {
		item := new(id.KeyItem)
		if err := rows.Scan(&item.Key, &item.Offset, &item.Size_); err != nil {
			panic(err)
		}
		items = append(items, item)
	}
	// 加入
	for _, item := range items {
		if err := genDb(context.Background(), item); err != nil {
			panic(err)
		}
	}
}

func gen(key string) {
	for {
		id := getId(key)
		if id != 0 {
			chanMap[key] <- id
		} else {
			//以免把数据库打死
			time.Sleep(time.Second)
		}
	}
}

func list() []*DB {
	l := <-c
	c <- l
	return l
}

func getId(key string) int64 {
	for _, db := range list() {
		ret, err := db.db.Exec("update id set id=last_insert_id(id) + 4 where k = ?", key)
		if err != nil {
			qlog.Error(key)
			continue
		}
		id, _ := ret.LastInsertId()
		if (id % 4) == db.idx {
			qlog.Info(key, db.idx, id)
			return id
		}
		qlog.Error("gen id", key, db.idx, id)
	}
	return 0
}
