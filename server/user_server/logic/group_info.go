package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/config/kvconfig"
	"github.com/qsock/qim/lib/cachename"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/proto/user"
	"github.com/qsock/qim/lib/util"
	"strings"
	"time"
)

func ClearGroupCache(ctx context.Context, groupId int64) {
	kvconn := dao.GetKvConn(kvconfig.KvDefault)
	kvconn.Del(cachename.RedisGroupInfo(groupId))
}

// 通过id获取群组信息
func GetGroupById(ctx context.Context, groupId int64) (*user.GroupInfo, error) {
	infos, err := GetGroupByIds(ctx, []int64{groupId})
	if err != nil {
		return nil, err
	}
	if len(infos) == 0 {
		return nil, sql.ErrNoRows
	}
	return infos[0], nil
}

func GetGroupByIds(ctx context.Context, groupIds []int64) ([]*user.GroupInfo, error) {
	infos, _ := getGroupInfoByIdsOnCache(ctx, groupIds)
	if len(infos) == len(groupIds) {
		return infos, nil
	}

	if infos == nil {
		infos = make([]*user.GroupInfo, 0)
	}

	sids := make([]int64, 0)
	for _, id := range groupIds {
		var flag bool
		for _, info := range infos {
			if info.Id == id {
				flag = true
			}
		}
		if !flag {
			sids = append(sids, id)
		}
	}

	sinfos, err := getGroupInfoByIdsOnDb(ctx, sids)
	if err != nil {
		return nil, err
	}

	pipl := dao.GetKvConn(kvconfig.KvDefault).Pipeline()
	for _, info := range sinfos {
		key := cachename.RedisGroupInfo(info.Id)
		b, _ := json.Marshal(info)
		pipl.Set(key, string(b), time.Second*86400)
	}
	if _, err := pipl.Exec(); err != nil {
		qlog.Ctx(ctx).Error(sids, err)
	}
	infos = append(infos, sinfos...)
	return infos, nil
}

func getGroupInfoByIdsOnCache(ctx context.Context, groupIds []int64) ([]*user.GroupInfo, error) {
	pipl := dao.GetKvConn(kvconfig.KvDefault).Pipeline()
	for _, id := range groupIds {
		key := cachename.RedisGroupInfo(id)
		pipl.Get(key)
	}
	cmds, err := pipl.Exec()
	if err != nil {
		qlog.Ctx(ctx).Error(groupIds, err)
		return nil, err
	}
	infos := make([]*user.GroupInfo, 0)
	for _, cmd := range cmds {
		val := cmd.(*redis.StringCmd)
		var info *user.GroupInfo
		if err := json.Unmarshal([]byte(val.String()), info); err != nil {
			qlog.Ctx(ctx).Error(groupIds, cmd, val.String(), err)
			continue
		}
		if info.GetId() > 0 {
			infos = append(infos, info)
		}
	}
	return infos, nil
}

func getGroupInfoByIdsOnDb(ctx context.Context, groupIds []int64) ([]*user.GroupInfo, error) {
	ssql := "select id,name,mute_until,notice,created_on," +
		"max_member_ct,current_ct,avatar,join_type,deleted_on,t from `group` where id in (?) "
	items := make([]*user.GroupInfo, 0)
	rows, err := dao.GetConn(dbconfig.DbUser).Query(ssql, util.Int64sToStr(groupIds))
	if err != nil {
		qlog.Ctx(ctx).Error(groupIds, err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := new(user.GroupInfo)
		var avatar string
		if err := rows.Scan(&item.Id, &item.Name, &item.MuteUtil, &item.Notice, &item.CreatedOn,
			&item.MaxMemberCt, &item.CurrentCt, &avatar, &item.JoinType, &item.DeletedOn, &item.T); err != nil {
			qlog.Ctx(ctx).Error(groupIds, err)
			return nil, err
		}
		avatars := strings.Split(avatar, ",")
		item.Avatars = avatars
		items = append(items, item)
	}
	return items, nil
}
