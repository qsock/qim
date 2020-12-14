package dbconfig

import (
	"github.com/qsock/qf/store/db"
)

const (
	// id db
	DbIdA string = "DbIdA"
	DbIdB string = "DbIdB"
	DbIdC string = "DbIdC"
	DbIdD string = "DbIdD"

	// passport db
	DbPassport string = "DbPassport"
	//DbPassportR1 string = "DbPassportR1"

	// file db
	DbFile string = "DbFile"
	//DbFileR1 string = "DbFileR1"

	// user db
	DbUser string = "DbUser"
	//DbUserR1 string = "DbUserR1"

	DbEvent string = "DbEvent"
	//DbEventR1 string = "DbEventR1"

	DbMsgShard0 string = "DbMsgShard0"
	DbMsgShard1 string = "DbMsgShard1"
	DbMsgShard2 string = "DbMsgShard2"
	DbMsgShard3 string = "DbMsgShard3"

	//DbMsgShard0R1 string = "DbMsgShard0R1"
	//DbMsgShard1R1 string = "DbMsgShard1R1"
	//DbMsgShard2R1 string = "DbMsgShard2R1"
	//DbMsgShard3R1 string = "DbMsgShard3R1"
)

// prod config map
var ConfigMap = map[string]db.Config{
	DbIdA: DbIdAConfig,
	DbIdB: DbIdBConfig,
	DbIdC: DbIdCConfig,
	DbIdD: DbIdDConfig,

	DbPassport: DbPassportConfig,
	//DbPassportR1: DbPassportR1Config,

	DbFile: DbFileConfig,
	//DbFileR1: DbFileR1Config,

	DbEvent: DbEventConfig,
	//DbEventR1: DbEventR1Config,

	DbUser: DbUserConfig,
	//DbUserR1: DbUserR1Config,

	DbMsgShard0: DbMsgShard0Config,
	//DbMsgShard0R1: DbMsgShard0R1Config,
	DbMsgShard1: DbMsgShard1Config,
	//DbMsgShard1R1: DbMsgShard1R1Config,
	DbMsgShard2: DbMsgShard2Config,
	//DbMsgShard2R1: DbMsgShard2R1Config,
	DbMsgShard3: DbMsgShard3Config,
	//DbMsgShard3R1: DbMsgShard3R1Config,
}

// dev config map
var ConfigDevMap = map[string]db.Config{
	DbIdA: DbIdADevConfig,
	DbIdB: DbIdBDevConfig,
	DbIdC: DbIdCDevConfig,
	DbIdD: DbIdDDevConfig,

	DbPassport:  DbPassportDevConfig,
	DbFile:      DbFileDevConfig,
	DbEvent:     DbEventDevConfig,
	DbUser:      DbUserDevConfig,
	DbMsgShard0: DbMsgShard0DevConfig,
}
