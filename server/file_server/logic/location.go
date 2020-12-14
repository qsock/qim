package logic

import (
	"context"
	"github.com/qsock/qf/qlog"
	lcache "github.com/qsock/qf/store/cache"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/lib/cachename"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/proto/file"
	"github.com/qsock/qim/lib/proto/model"
)

var (
	cache *lcache.Cache
)

func init() {
	cache = lcache.New(1000)
}

func GetProvinceAndCity(ctx context.Context, req *file.GetProvinceAndCityReq) (*file.GetProvinceAndCityResp, error) {
	resp := new(file.GetProvinceAndCityResp)
	provinces, err := GetLvlocation(ctx, 0)
	if err != nil {
		return nil, err
	}
	cities, err := GetLvlocation(ctx, 1)
	if err != nil {
		return nil, err
	}
	resp.Provinces, resp.Cities = provinces, cities
	return resp, nil
}

func GetLvlocation(ctx context.Context, lv int) ([]*model.Cnarea2019, error) {
	cacheKey := cachename.MemoryLocationLevelCache(lv)
	val, ok := cache.Get(cacheKey)
	if ok {
		return val.([]*model.Cnarea2019), nil
	}
	items := make([]*model.Cnarea2019, 0)
	ssql := "select id,level,parent_code,area_code,zip_code," +
		"city_code,name,short_name,merger_name,pinyin," +
		"lng,lat from cnarea_2019 where level=?"
	rows, err := dao.GetConn(dbconfig.DbFile).Query(ssql, lv)
	if err != nil {
		qlog.Ctx(ctx).Error(lv)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		item := new(model.Cnarea2019)
		if err := rows.Scan(&item.Id, &item.Level, &item.ParentCode, &item.AreaCode, &item.ZipCode,
			&item.CityCode, &item.Name, &item.ShortName, &item.MergerName, &item.Pinyin,
			&item.Lng, &item.Lat); err != nil {
			qlog.Ctx(ctx).Error(lv, err)
			continue
		}
		items = append(items, item)
	}
	cache.SetEx(cacheKey, items, 360)
	return items, nil
}
