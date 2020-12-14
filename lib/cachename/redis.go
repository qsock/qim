package cachename

import "strconv"

func RedisPassportSms(tel string) string {
	return "passport:sms:" + tel
}

func RedisPassportSeq(userId int64) string {
	return "passport:seq:" + strconv.FormatInt(userId, 10)
}

func RedisPassportRefreshToken(userId int64) string {
	return "passport:token:" + strconv.FormatInt(userId, 10)
}

func RedisUserInfo(userId int64) string {
	return "user:info:" + strconv.FormatInt(userId, 10)
}

func RedisUserWs(userId int64) string {
	return "user:ws" + strconv.FormatInt(userId, 10)
}

func RedisUserBan(userId int64) string {
	return "user:ban:" + strconv.FormatInt(userId, 10)
}

func RedisGroupInfo(userId int64) string {
	return "group:info:" + strconv.FormatInt(userId, 10)
}

func RedisChatRoomMsg(chatId string) string {
	return "chatroom:msg:" + chatId
}
