package logic

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/qsock/qim/lib/proto/model"
	"github.com/qsock/qim/lib/proto/passport"
	"github.com/qsock/qim/lib/qjwt"
	"time"
)

func GenToken(ctx context.Context,
	userId int64,
	device model.Device,
	ip string) (string, error) {

	m := new(passport.JwtClaims)
	m.UserId = userId
	m.Device = device
	m.UserIp = ip
	var err error
	m.SeqId, err = GetUserSeqId(ctx, userId)
	if err != nil && err != redis.Nil {
		return "", err
	}
	b, _ := m.Marshal()
	return qjwt.CreateToken(ctx, b, time.Now().Unix()+86400)
}
