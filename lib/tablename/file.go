package tablename

import (
	"github.com/qsock/qim/lib/types"
	"strconv"
)

func FileTable(userId int64, env string) string {
	if env == types.EnvDev {
		return "file_0"
	}
	return "file_" + strconv.FormatInt(userId%64, 10)
}
