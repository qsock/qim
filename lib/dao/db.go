package dao

import (
	"database/sql"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/store/db"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/lib/types"
)

//获取连接池
func GetConn(name string) *sql.DB {
	mu.Lock()
	defer mu.Unlock()

	conn, err := db.GetDB(name)
	if conn != nil && err == nil {
		return conn
	}

	m := dbconfig.ConfigMap
	if env == types.EnvDev {
		m = dbconfig.ConfigDevMap
	}

	db.Add(name, m[name])
	conn, err = db.GetDB(name)
	if err != nil {
		qlog.Get().Logger().Error(err, name)
	}
	return conn
}
