package logic

import (
	"context"
	"fmt"
	"github.com/qiniu/api.v7/v7/storage"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qim/config/common"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/lib/constdef"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/proto/file"
	"github.com/qsock/qim/lib/proto/ret"
	"github.com/qsock/qim/lib/tablename"
	"github.com/qsock/qim/server/file_server/config"
	"path"
)

func GetUploadToken(ctx context.Context, req *file.GetUploadTokenReq) (*file.GetUploadTokenResp, error) {
	resp := new(file.GetUploadTokenResp)
	if config.GetConfig().FileType == "qiniu" {
		qiniuToken := GetQiniuToken(req.Path)
		resp.Path = req.Path
		resp.Tokens = map[int32]string{int32(file.UploadType_UploadQiniu): qiniuToken}
	}
	return resp, nil
}

func GetUserFile(ctx context.Context, req *file.GetUserFileReq) (*file.GetUserFileResp, error) {
	ssql := "select id,url,path,created_on from " + tablename.FileTable(req.UserId, config.GetEnv()) +
		" where user_id=? "
	condis := []interface{}{req.UserId}
	if len(req.Path) > 0 {
		ssql += " and path=?"
		condis = append(condis, req.Path)
	}
	ssql += " order by id desc limit %d offset %d"
	ssql = fmt.Sprintf(ssql, req.PageSize, req.Page*req.PageSize)
	rows, err := dao.GetConn(dbconfig.DbFile).Query(ssql, condis...)
	if err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	resp := new(file.GetUserFileResp)
	defer rows.Close()
	items := make([]*file.UserFile, 0)
	for rows.Next() {
		item := new(file.UserFile)
		item.UserId = req.UserId
		if err := rows.Scan(&item.Id, &item.Url, &item.Path, &item.CreatedOn); err != nil {
			qlog.Ctx(ctx).Error(err, req)
			continue
		}
		items = append(items, item)
	}
	resp.Files = items
	return resp, nil
}

func UserUploadSucceed(ctx context.Context, req *file.UserUploadSucceedReq) (*ret.EmptyResp, error) {
	isql := "insert into " + tablename.FileTable(req.UserId, config.GetEnv()) +
		" (user_id,url,path,created_on) values(?,?,?,unix_timestamp())"
	if _, err := dao.GetConn(dbconfig.DbFile).Exec(isql, req.UserId, req.Url, req.Path); err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	resp := new(ret.EmptyResp)
	return resp, nil
}

func UploadFileByUrl(ctx context.Context, req *file.UploadFileByUrlReq) (*file.UploadFileByUrlResp, error) {
	resp := new(file.UploadFileByUrlResp)
	var err error
	if config.GetConfig().FileType == "qiniu" {
		resp.Url, err = upload2Qiniu(ctx, req.Url, req.Path, req.UserId)
		if err != nil {
			return nil, err
		}
	}
	if _, err := UserUploadSucceed(ctx, &file.UserUploadSucceedReq{UserId: req.UserId, Url: resp.Url, Path: req.Path}); err != nil {
		return nil, err
	}
	return resp, nil
}

func upload2Qiniu(ctx context.Context, url, p string, userId int64) (string, error) {
	bucketManager := storage.NewBucketManager(Credential(), &storage.Config{UseHTTPS: true})
	ret, err := bucketManager.FetchWithoutKey(url, common.GetQiniuBucketByPath(p))
	if err != nil {
		qlog.Ctx(ctx).Error(url, p, userId, err)
		return "", err
	}
	fileUrl := path.Join(common.GetQiniuUrlByPath(p), ret.Key)
	return fileUrl, nil
}

func GetSysAvatars(ctx context.Context, req *file.GetSysAvatarsReq) (*file.GetSysAvatarsResp, error) {
	ssql := "select url from " + tablename.FileTable(0, config.GetEnv()) +
		" where user_id=0 and path=?"
	rows, err := dao.GetConn(dbconfig.DbFile).Query(ssql, constdef.FilePathAvatar)
	if err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	resp := new(file.GetSysAvatarsResp)
	defer rows.Close()
	items := make([]string, 0)
	for rows.Next() {
		var str string
		if err := rows.Scan(&str); err != nil {
			qlog.Ctx(ctx).Error(err, req)
			continue
		}
		items = append(items, str)
	}
	resp.Avatars = items
	return resp, nil
}
