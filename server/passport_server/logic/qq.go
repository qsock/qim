package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/lib/codes"
	"github.com/qsock/qim/lib/constdef"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/gclient"
	"github.com/qsock/qim/lib/method"
	"github.com/qsock/qim/lib/proto/file"
	"github.com/qsock/qim/lib/proto/model"
	"github.com/qsock/qim/lib/proto/passport"
	"github.com/qsock/qim/lib/proto/user"
	"github.com/qsock/qim/server/passport_server/config"
	"io/ioutil"
	"net/http"
)

func Qqlogin(ctx context.Context, req *passport.QqloginReq) (*passport.LoginResp, error) {
	resp := new(passport.LoginResp)
	qqInfo, err := GetQQSimpleUserInfo(ctx, req.Openid, req.Token)
	if err != nil {
		return nil, err
	}
	if qqInfo.Ret != 0 {
		qlog.Ctx(ctx).Error(req, err, qqInfo)
		resp.Err = codes.Error(codes.ErrorPassportQq)
		return resp, nil
	}
	defer func() {
		isql := "insert into qqinfo (openid,nickname,figureurl,figureurl_1," +
			"figureurl_2,figureurl_qq_1,figureurl_qq_2,gender,is_yellow_vip," +
			"vip,yellow_vip_level,level,is_yellow_year_vip,created_on)values(?,?,?,?," +
			"?,?,?,?,?," +
			"?,?,?,?,unix_timestamp()) on duplicate key update nickname=?,figureurl=?,figureurl_1=?," +
			"figureurl_2=?,figureurl_qq_1=?,figureurl_qq_2=?,gender=?,is_yellow_vip=?," +
			"vip=?,yellow_vip_level=?,level=?,is_yellow_year_vip=?"
		var gender model.Gender
		if qqInfo.Gender == "男" {
			gender = model.Gender_GenderMale
		} else if qqInfo.Gender == "女" {
			gender = model.Gender_GenderFemale
		}
		if _, err := dao.GetConn(dbconfig.DbPassport).Exec(isql, req.Openid, qqInfo.Nickname, qqInfo.Figureurl, qqInfo.Figureurl1,
			qqInfo.Figureurl2, qqInfo.FigureurlQq1, qqInfo.FigureurlQq2, gender, qqInfo.IsYellowVip,
			qqInfo.Vip, qqInfo.YellowVipLevel, qqInfo.Level, qqInfo.IsYellowYearVip,
			qqInfo.Nickname, qqInfo.Figureurl, qqInfo.Figureurl1,
			qqInfo.Figureurl2, qqInfo.FigureurlQq1, qqInfo.FigureurlQq2, gender, qqInfo.IsYellowVip,
			qqInfo.Vip, qqInfo.YellowVipLevel, qqInfo.Level, qqInfo.IsYellowYearVip,
		); err != nil {
			qlog.Ctx(ctx).Error(req, err, qqInfo)
		}
	}()
	{
		var userId int64
		ssql := "select user_id from qq where openid=? limit 1"
		if err := dao.GetConn(dbconfig.DbPassport).QueryRow(ssql, req.Openid).Scan(&userId); err != nil && err != sql.ErrNoRows {
			qlog.Ctx(ctx).Error(err, req)
			return nil, err
		}
		// 如果已存在
		if userId > 0 {
			return doLogin(ctx, userId, req.Meta)
		}
	}

	return doQqRegister(ctx, qqInfo, req.Openid, req.Meta)
}

func doQqRegister(ctx context.Context, info *QQUserInfo, openId string, meta *model.RequestMeta) (*passport.LoginResp, error) {
	avatar := info.FigureurlQq2
	if len(avatar) == 0 {
		avatar = info.Figureurl2
	}
	if len(avatar) == 0 {
		avatar = info.FigureurlQq1
	}

	// 创建id
	userId := gclient.GenDbId(ctx, config.GetConfig().IdKey)
	avatar, err := exchangeAvatar(ctx, avatar, userId)
	if err != nil {
		qlog.Ctx(ctx).Error(info, meta, err, openId)
		return nil, err
	}
	name := info.Nickname
	var gender model.Gender
	if info.Gender == "男" {
		gender = model.Gender_GenderMale
	} else if info.Gender == "女" {
		gender = model.Gender_GenderFemale
	}

	// 先去创建用户，多次创建不用更新
	{
		creq := new(user.CreateReq)
		creq.UserId = userId
		creq.Avatar = avatar
		creq.Name = name
		creq.Gender = gender
		if err := createUser(ctx, creq); err != nil {
			qlog.Ctx(ctx).Error(info, meta, err, openId)
			return nil, err
		}
	}
	isql := "insert into qq (user_id,openid) values(?,?)"
	if _, err := dao.GetConn(dbconfig.DbPassport).Exec(isql, userId, openId); err != nil {
		qlog.Ctx(ctx).Error(info, meta, err, openId)
		return nil, err
	}
	_ = doSeq(ctx, userId)
	return doLogin(ctx, userId, meta)
}

func exchangeAvatar(ctx context.Context, avatar string, userId int64) (string, error) {
	creq := new(file.UploadFileByUrlReq)
	creq.UserId = userId
	creq.Path = constdef.FilePathAvatar
	creq.Url = avatar
	cresp := new(file.UploadFileByUrlResp)
	if err := qgrpc.Call(ctx, method.FileUploadFileByUrl, creq, cresp); err != nil {
		qlog.Ctx(ctx).Error(creq, err)
		return "", err
	}
	return cresp.Url, nil
}

type QQUserInfo struct {
	Ret             int    `json:"ret"`
	Msg             string `json:"msg"`
	Nickname        string `json:"nickname"`
	Figureurl       string `json:"figureurl"`
	Figureurl1      string `json:"figureurl_1"`
	Figureurl2      string `json:"figureurl_2"`
	FigureurlQq1    string `json:"figureurl_qq_1"`
	FigureurlQq2    string `json:"figureurl_qq_2"`
	Gender          string `json:"gender"`
	IsYellowVip     string `json:"is_yellow_vip"`
	Vip             string `json:"vip"`
	YellowVipLevel  string `json:"yellow_vip_level"`
	Level           string `json:"level"`
	IsYellowYearVip string `json:"is_yellow_year_vip"`
}

func GetQQSimpleUserInfo(ctx context.Context, openid, token string) (*QQUserInfo, error) {
	wurl := "https://graph.qq.com/user/get_user_info?access_token=%s&oauth_consumer_key=%s&openid=%s"
	wurl = fmt.Sprintf(wurl, token, config.GetConfig().QqApp.AppId, openid)
	resp, err := http.Get(wurl)
	if err != nil {
		qlog.Ctx(ctx).Error(token, openid, err)
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		qlog.Ctx(ctx).Error(token, openid, err)
		return nil, err
	}
	qlog.Ctx(ctx).Info(openid, token, string(b))
	m := new(QQUserInfo)
	err = json.Unmarshal(b, m)
	if err != nil {
		qlog.Ctx(ctx).Error(token, openid, err)
		return nil, err
	}
	return m, nil
}
