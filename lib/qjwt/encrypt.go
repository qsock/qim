package qjwt

import (
	"context"
	"encoding/hex"
	"github.com/qsock/qf/encrypt"
	"github.com/qsock/qf/qlog"
)

func aesEncrypt(ctx context.Context, b []byte) (string, error) {
	encryptVal, err := encrypt.AesCbcEncrypt(b, []byte(j.cfg.AesKey), []byte(j.cfg.AesIv))
	if err != nil {
		return "", err
	}
	str := hex.EncodeToString(encryptVal)
	return str, nil
}

func aesDecrypt(ctx context.Context, sign string) ([]byte, error) {
	encryptVals, err := hex.DecodeString(sign)
	if err != nil {
		qlog.Get().Logger().Ctx(ctx).Error(sign, err)
		return nil, err
	}
	b, err := encrypt.AesCbcDecrypt(encryptVals, []byte(j.cfg.AesKey), []byte(j.cfg.AesIv))
	if err != nil {
		qlog.Get().Logger().Ctx(ctx).Error(sign, err)
		return nil, err
	}
	return b, nil
}
