package dbconfig

import "github.com/qsock/qf/store/db"

var (
	DbEventDevConfig = db.Config{
		Addr:         "127.0.0.1:3306",
		User:         "root",
		Pwd:          "123456",
		Db:           "event",
		MaxOpenConns: 5,
		MaxIdleConns: 3,
		Options:      "interpolateParams=true",
	}
)

var (
	DbEventConfig = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "event",
		MaxOpenConns:    20,
		MaxIdleConns:    20,
		MaxConnLifeTime: 3600,
	}

	DbEventR1Config = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "event",
		MaxOpenConns:    20,
		MaxIdleConns:    20,
		MaxConnLifeTime: 3600,
	}
)
