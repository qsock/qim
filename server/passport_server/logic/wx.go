package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qim/config/dbconfig"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/gclient"
	"github.com/qsock/qim/lib/proto/model"
	"github.com/qsock/qim/lib/proto/passport"
	"github.com/qsock/qim/lib/proto/user"
	"github.com/qsock/qim/server/passport_server/config"
	"io/ioutil"
	"net/http"
	"strings"
)

func WxLogin(ctx context.Context, req *passport.WxLoginReq) (*passport.LoginResp, error) {
	acm, err := getWxAccessToken(ctx, req.Code)
	if err != nil {
		return nil, err
	}

	info, err := getWxUserInfo(ctx, acm.Openid, acm.AccessToken)
	if err != nil {
		return nil, err
	}
	defer func() {
		if len(info.Openid) == 0 {
			qlog.Ctx(ctx).Error(err, req, "nil")
			return
		}
		isql := "insert into wxinfo(openid,nickname,sex,province,city," +
			"country,headimgurl,privilege,unionid,created_on)values(?,?,?,?,?," +
			"?,?,?,?,unix_timestamp()) on duplicate key update nickname=?,sex=?,province=?,city=?," +
			"country=?,headimgurl=?,privilege=?,unionid=?"
		if _, err := dao.GetConn(dbconfig.DbPassport).Exec(isql, info.Openid, info.Nickname, info.Sex, info.Province, info.City,
			info.Country, info.Headimgurl, strings.Join(info.Privilege, ","), info.Unionid,
			info.Nickname, info.Sex, info.Province, info.City,
			info.Country, info.Headimgurl, strings.Join(info.Privilege, ","), info.Unionid,
		); err != nil {
			qlog.Ctx(ctx).Error(err, req, info)
		}
	}()

	{
		var userId int64
		ssql := "select user_id from wx where openid=? limit 1"
		if err := dao.GetConn(dbconfig.DbPassport).QueryRow(ssql, acm.Openid).Scan(&userId); err != nil && err != sql.ErrNoRows {
			qlog.Ctx(ctx).Error(err, req)
			return nil, err
		}
		// 如果已存在
		if userId > 0 {
			return doLogin(ctx, userId, req.Meta)
		}
	}

	return doWxRegister(ctx, info, acm, req.Meta)
}

func doWxRegister(ctx context.Context, info *WxUserInfo, acm *WxAccessToken, meta *model.RequestMeta) (*passport.LoginResp, error) {
	// 创建id
	userId := gclient.GenDbId(ctx, config.GetConfig().IdKey)
	creq := new(user.CreateReq)
	creq.UserId = userId
	if len(info.Openid) == 0 {
		creq.Name = acm.Openid
		creq.Avatar, _ = getRandomAvatar(ctx)
	} else {
		creq.Name = info.Nickname
		creq.Avatar, _ = exchangeAvatar(ctx, info.Headimgurl, userId)
		creq.Gender = model.Gender(info.Sex)
	}

	{
		if err := createUser(ctx, creq); err != nil {
			qlog.Ctx(ctx).Error(info, meta, err, acm)
			return nil, err
		}
	}
	isql := "insert into wx (user_id,openid,unionid) values(?,?,?)"
	if _, err := dao.GetConn(dbconfig.DbPassport).Exec(isql, userId, acm.Openid, acm.Unionid); err != nil {
		qlog.Ctx(ctx).Error(info, meta, err, acm)
		return nil, err
	}
	_ = doSeq(ctx, userId)
	return doLogin(ctx, userId, meta)
}

type WxAccessToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
}

func getWxAccessToken(ctx context.Context, code string) (*WxAccessToken, error) {
	wurl := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	wurl = fmt.Sprintf(wurl, config.GetConfig().WxApp.AppId, config.GetConfig().WxApp.Secret, code)
	resp, err := http.Get(wurl)
	if err != nil {
		qlog.Ctx(ctx).Error(err, code)
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		qlog.Ctx(ctx).Error(err, code)
		return nil, err
	}
	qlog.Ctx(ctx).Info(code, string(b))
	m := new(WxAccessToken)
	err = json.Unmarshal(b, m)
	if err != nil {
		qlog.Ctx(ctx).Error(err, code)
		return nil, err
	}
	return m, nil
}

type WxUserInfo struct {
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}

func getWxUserInfo(ctx context.Context, openId, accessToken string) (*WxUserInfo, error) {
	wurl := "https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s"
	wurl = fmt.Sprintf(wurl, accessToken, openId)
	resp, err := http.Get(wurl)
	if err != nil {
		qlog.Ctx(ctx).Error(err, accessToken, openId)
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		qlog.Ctx(ctx).Error(err, accessToken, openId)
		return nil, err
	}
	qlog.Ctx(ctx).Info(accessToken, openId, string(b))
	m := new(WxUserInfo)
	err = json.Unmarshal(b, m)
	if err != nil {
		qlog.Ctx(ctx).Error(err, accessToken, openId)
		return nil, err
	}
	return m, nil
}
