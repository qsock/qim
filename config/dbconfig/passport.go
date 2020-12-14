package dbconfig

import "github.com/qsock/qf/store/db"

var (
	DbPassportDevConfig = db.Config{
		Addr:         "127.0.0.1:3306",
		User:         "root",
		Pwd:          "123456",
		Db:           "passport",
		MaxOpenConns: 5,
		MaxIdleConns: 3,
		Options:      "interpolateParams=true",
	}
)

var (
	DbPassportConfig = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "passport",
		MaxOpenConns:    20,
		MaxIdleConns:    20,
		MaxConnLifeTime: 3600,
	}

	DbPassportR1Config = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "passport",
		MaxOpenConns:    20,
		MaxIdleConns:    20,
		MaxConnLifeTime: 3600,
	}
)
