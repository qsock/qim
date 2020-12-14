package dbconfig

import (
	"github.com/qsock/qf/store/db"
)

var (
	DbIdADevConfig = db.Config{
		Addr:         "127.0.0.1:3306",
		User:         "root",
		Pwd:          "123456",
		Db:           "db_a",
		MaxOpenConns: 2,
		MaxIdleConns: 1,
	}

	DbIdBDevConfig = db.Config{
		Addr:         "127.0.0.1:3306",
		User:         "root",
		Pwd:          "123456",
		Db:           "db_b",
		MaxOpenConns: 2,
		MaxIdleConns: 1,
	}

	DbIdCDevConfig = db.Config{
		Addr:         "127.0.0.1:3306",
		User:         "root",
		Pwd:          "123456",
		Db:           "db_c",
		MaxOpenConns: 2,
		MaxIdleConns: 1,
	}

	DbIdDDevConfig = db.Config{
		Addr:         "127.0.0.1:3306",
		User:         "root",
		Pwd:          "123456",
		Db:           "db_d",
		MaxOpenConns: 2,
		MaxIdleConns: 1,
	}
)

var (
	DbIdAConfig = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "db_a",
		MaxOpenConns:    4,
		MaxIdleConns:    2,
		MaxConnLifeTime: 3600,
	}

	DbIdBConfig = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "db_b",
		MaxOpenConns:    4,
		MaxIdleConns:    2,
		MaxConnLifeTime: 3600,
	}

	DbIdCConfig = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "db_c",
		MaxOpenConns:    4,
		MaxIdleConns:    2,
		MaxConnLifeTime: 3600,
	}

	DbIdDConfig = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "db_d",
		MaxOpenConns:    4,
		MaxIdleConns:    2,
		MaxConnLifeTime: 3600,
	}
)