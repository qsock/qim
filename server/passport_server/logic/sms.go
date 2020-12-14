package logic

import (
	"context"
	"encoding/json"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/util/coderand"
	"github.com/qsock/qim/config/kvconfig"
	"github.com/qsock/qim/lib/cachename"
	"github.com/qsock/qim/lib/codes"
	"github.com/qsock/qim/lib/dao"
	"github.com/qsock/qim/lib/proto/passport"
	"github.com/qsock/qim/lib/proto/ret"
	"time"
)

// 发送短信
// 天极流控是由第三方去处理就可以了，我们不用管理
func Sms(ctx context.Context, req *passport.SmsReq) (*ret.EmptyResp, error) {
	resp := new(ret.EmptyResp)
	cacheKey := cachename.RedisPassportSms(req.Tel)
	kvconn := dao.GetKvConn(kvconfig.KvDefault)
	result := kvconn.Get(cacheKey).Val()
	smsM := new(passport.SmsModel)
	if len(result) > 0 {
		if err := json.Unmarshal([]byte(result), smsM); err != nil {
			qlog.Ctx(ctx).Error(err, req)
			return nil, err
		}
	}

	// 60秒才能发一次
	if smsM.GetCreatedOn() > time.Now().Unix()+60 {
		resp.Err = codes.Error(codes.ErrorPassportSmsToosoon)
		return resp, nil
	}
	smsM.Code = coderand.Num(4)
	smsM.CreatedOn = time.Now().Unix()
	b, _ := json.Marshal(smsM)
	if err := kvconn.Set(cacheKey, string(b), time.Minute*5).Err(); err != nil {
		qlog.Ctx(ctx).Error(err, req)
		return nil, err
	}
	qlog.Ctx(ctx).Debug(req.Tel, smsM.Code)
	return resp, nil
}
