package dbconfig

import "github.com/qsock/qf/store/db"

var (
	DbFileDevConfig = db.Config{
		Addr:         "127.0.0.1:3306",
		User:         "root",
		Pwd:          "123456",
		Db:           "file",
		MaxOpenConns: 5,
		MaxIdleConns: 3,
		Options:      "interpolateParams=true",
	}
)

var (
	DbFileConfig = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "file",
		MaxOpenConns:    20,
		MaxIdleConns:    20,
		MaxConnLifeTime: 3600,
	}

	DbFileR1Config = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "file",
		MaxOpenConns:    20,
		MaxIdleConns:    20,
		MaxConnLifeTime: 3600,
	}
)
