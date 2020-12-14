package tablename

import (
	"github.com/qsock/qim/lib/types"
	"strconv"
)

func UserLastactive(userId int64, env string) string {
	if env == types.EnvDev {
		return "lastactive_0"
	}
	return "lastactive_" + strconv.FormatInt(userId%128, 10)
}

func UserGroupMember(groupId int64, env string) string {
	if env == types.EnvDev {
		return "group_member_0"
	}
	return "group_member_" + strconv.FormatInt(groupId%128, 10)
}
