package qjwt

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/qsock/qf/qlog"
	"time"
)

var (
	ErrInvalidConf  = errors.New("conf is not valid")
	ErrTokenExpired = errors.New("token has expired")
	ErrTokenInvalid = errors.New("token is invalid")
	j               *Jwt
)

type Jwt struct {
	cfg *Config
}

// payload
type CustomClaims struct {
	jwt.StandardClaims
	Sign string `json:"sign"`
}

func Init(c *Config) error {
	j = new(Jwt)
	if len(c.AesKey) != 32 {
		return ErrInvalidConf
	}
	if len(c.AesIv) != 16 {
		return ErrInvalidConf
	}
	if len(c.Signkey) == 0 {
		return ErrInvalidConf
	}
	j.cfg = c
	return nil
}

// CreateToken 生成一个token
func CreateToken(ctx context.Context, b []byte, expiresAt int64) (string, error) {
	signStr, err := aesEncrypt(ctx, b)
	if err != nil {
		qlog.Get().Ctx(ctx).Error(err, string(b))
		return "", err
	}
	if expiresAt == 0 {
		// 默认过期时间一天
		expiresAt = time.Now().Unix() + 86400
	}
	claims := CustomClaims{
		jwt.StandardClaims{
			ExpiresAt: expiresAt,
			Issuer:    j.cfg.Signkey, //签名的发行者
		},
		signStr,
	}
	return jwtCreateToken(ctx, claims)
}

func Parse(ctx context.Context, tokenString string) ([]byte, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.cfg.Signkey), nil
	})

	if err != nil {
		qlog.Get().Ctx(ctx).Error(err, tokenString)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			}
		}
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		qlog.Get().Ctx(ctx).Error(err, tokenString)
		return nil, ErrTokenInvalid
	}
	sign := claims.Sign
	if len(sign) == 0 {
		qlog.Get().Ctx(ctx).Error(err, tokenString, claims)
		return nil, ErrTokenInvalid
	}
	return aesDecrypt(ctx, sign)
}

func ParseTokenWithoutTime(ctx context.Context, tokenString string) ([]byte, error) {
	jwt.TimeFunc = func() time.Time {
		return time.Unix(0, 0)
	}

	if len(tokenString) < 10 {
		qlog.Get().Ctx(ctx).Error(tokenString)
		return nil, ErrTokenInvalid
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.cfg.Signkey), nil
	})
	if err != nil {
		qlog.Get().Ctx(ctx).Error(tokenString, err)
		return nil, err
	}
	claims, ok := token.Claims.(*CustomClaims)

	if !ok {
		qlog.Get().Ctx(ctx).Error(tokenString, ok)
		return nil, ErrTokenInvalid
	}
	sign := claims.Sign
	if len(sign) == 0 {
		qlog.Get().Ctx(ctx).Error(err, tokenString, claims)
		return nil, ErrTokenInvalid
	}
	return aesDecrypt(ctx, sign)
}

func jwtCreateToken(ctx context.Context, claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(j.cfg.Signkey))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
