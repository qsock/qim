package dbconfig

import "github.com/qsock/qf/store/db"

const (
	// 默认最大的shard数量
	MsgTotalShard = 4
)

var (
	DbMsgShard0DevConfig = db.Config{
		Addr:         "127.0.0.1:3306",
		User:         "root",
		Pwd:          "123456",
		Db:           "msg_shard0",
		MaxOpenConns: 5,
		MaxIdleConns: 3,
		Options:      "interpolateParams=true",
	}
)

var (
	DbMsgShard0Config = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "msg_shard0",
		MaxOpenConns:    40,
		MaxIdleConns:    40,
		MaxConnLifeTime: 3600,
	}

	DbMsgShard0R1Config = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "msg_shard0",
		MaxOpenConns:    40,
		MaxIdleConns:    40,
		MaxConnLifeTime: 3600,
	}
	DbMsgShard1Config = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "msg_shard1",
		MaxOpenConns:    40,
		MaxIdleConns:    40,
		MaxConnLifeTime: 3600,
	}

	DbMsgShard1R1Config = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "msg_shard1",
		MaxOpenConns:    40,
		MaxIdleConns:    40,
		MaxConnLifeTime: 3600,
	}
	/******shard2****/
	DbMsgShard2Config = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "msg_shard2",
		MaxOpenConns:    40,
		MaxIdleConns:    40,
		MaxConnLifeTime: 3600,
	}

	DbMsgShard2R1Config = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "msg_shard2",
		MaxOpenConns:    40,
		MaxIdleConns:    40,
		MaxConnLifeTime: 3600,
	}
	/******shard3****/
	DbMsgShard3Config = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "msg_shard3",
		MaxOpenConns:    40,
		MaxIdleConns:    40,
		MaxConnLifeTime: 3600,
	}

	DbMsgShard3R1Config = db.Config{
		Addr:            "127.0.0.1:3306",
		User:            "root",
		Pwd:             "123456",
		Db:              "msg_shard3",
		MaxOpenConns:    40,
		MaxIdleConns:    40,
		MaxConnLifeTime: 3600,
	}
)
