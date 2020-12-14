package password

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"github.com/qsock/qf/encrypt"
	"github.com/qsock/qf/util/uuid"
	"github.com/qsock/qim/config/common"
	"strings"
)

const (
	EncryptMethod = "sha512"
)

func ValidatePassword(password string, store string) bool {
	if len(password) <= 0 || len(store) <= 0 {
		return false
	}
	decode := strings.Split(store, "$")
	if len(decode) != 3 {
		return false
	}
	salt := decode[2]
	src := encryptPassword(password, salt)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return hmac.Equal(dst, []byte(decode[0]))
}

// 加密pwd
func ObscurePassword(password string) string {
	salt := uuid.NewString()
	pwd := encryptPassword(password, salt)
	dst := make([]byte, hex.EncodedLen(len(pwd)))
	hex.Encode(dst, pwd)
	return strings.Join([]string{string(dst), EncryptMethod, salt}, "$")
}

func encryptPassword(plainPassword, salt string) []byte {
	sum := encrypt.HmacSha512(salt, plainPassword)
	dst := make([]byte, hex.EncodedLen(len(sum)))
	hex.Encode(dst, sum)
	mac := hmac.New(sha512.New, []byte(common.PwdSeed))
	mac.Write(dst)
	return mac.Sum(nil)
}
