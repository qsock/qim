package dbconfig

import "github.com/qsock/qf/store/db"

var (
	DbUserDevConfig = db.Config{
		Addr:         "127.0.0.1:3306",
		User:         "root",
		Pwd:          "123456",
		Db:           "user",
		MaxOpenConns: 5,
		MaxIdleConns: 3,
		Options:      "interpolateParams=true",
	}
)

var (
	DbUserConfig = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "user",
		MaxOpenConns:    20,
		MaxIdleConns:    20,
		MaxConnLifeTime: 3600,
	}
	DbUserR1Config = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "user",
		MaxOpenConns:    20,
		MaxIdleConns:    20,
		MaxConnLifeTime: 3600,
	}
)
