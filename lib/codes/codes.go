package codes

import (
	"github.com/qsock/qim/lib/proto/errmsg"
)

const (
	ErrorCommon int32 = 1

	ErrorInternalDB    = 1011
	ErrorInternalRedis = 1021
	ErrorInternalId    = 1051
	ErrorInternalHttp  = 1052

	ErrorAuthInvalid = 4001
	ErrorAuthExpired = 4002

	ErrorParameter = 20001

	ErrorPassportBanned       = 30001
	ErrorPassportPwd          = 30002
	ErrorPassportNameRepeated = 30003
	ErrorPassportSmsToosoon   = 30004
	ErrorPassportSmsCode      = 30005
	ErrorPassportQq           = 30006

	ErrorRoomIn = 40003

	// session 不存在
	ErrorSessNotExists = 50001

	// 消息发送失败
	ErrorMsgSentFailed       = 60001
	ErrorMsgLogined          = 60002
	ErrorMsgInvalid          = 60003
	ErrorMsgMorethan2Minutes = 60004

	ErrorUserApplyDenied   = 60005
	ErrorUserApplyRepeated = 60006
	ErrorUserCannotAddSelf = 60007
	ErrorUserIsFriend      = 60008
	ErrorUserFriendNoneAdd = 60009

	ErrorUserGroupMemberMustBiggerThanOne  = 60010
	ErrorUserGroupAlreadyIn                = 60011
	ErrorUserGroupForbidJoin               = 60012
	ErrorUserGroupMaxMember                = 60013
	ErrorUserGroupHasNoRight               = 60014
	ErrorUserGroupCannotPointSelf          = 60015
	ErrorUserGroupUserLessThanTwo          = 60016
	ErrorUserGroupUserManagerCannotDeleted = 60017
	ErrorUserGroupOwnnerCannotLeave        = 60018
	ErrorUserGroupNotMember                = 60019
	ErrorUserGroupBeenMute                 = 60020
	ErrorUserGroupBeenBlock                = 60021
)

var errText = map[int32]string{
	ErrorCommon: "服务端错误",

	ErrorInternalDB:    "数据库发生错误",
	ErrorInternalRedis: "缓存库发生错误",
	ErrorInternalId:    "id生成错误",
	ErrorInternalHttp:  "http调用错误",

	ErrorAuthInvalid: "鉴权失败，请重新登陆",
	ErrorAuthExpired: "登陆已过期，请重新登陆",

	ErrorParameter: "参数错误",

	ErrorPassportBanned:       "您已被封禁,请联系客服解封",
	ErrorPassportPwd:          "用户名或密码错误",
	ErrorPassportNameRepeated: "账号不能重复",
	ErrorPassportSmsToosoon:   "消息发送过快，请稍后重试",
	ErrorPassportSmsCode:      "验证码错误",
	ErrorPassportQq:           "qq登陆失败",

	ErrorRoomIn: "您已经在别的房间了,不能重复加房间",

	ErrorSessNotExists: "长连接session不存在",

	ErrorMsgSentFailed:       "消息发送失败",
	ErrorMsgLogined:          "长连接服务登陆失败",
	ErrorMsgInvalid:          "发送非法消息",
	ErrorMsgMorethan2Minutes: "消息发送超过了2分钟，不可以撤回",

	ErrorUserApplyDenied:                   "好友申请已发出，但被对方拒收",
	ErrorUserApplyRepeated:                 "一小时内，不要重复申请添加好友",
	ErrorUserCannotAddSelf:                 "不可以添加自己做好友",
	ErrorUserIsFriend:                      "你们已经是好友了,不要重复申请",
	ErrorUserFriendNoneAdd:                 "对方设置了，禁止任何人添加",
	ErrorUserGroupMemberMustBiggerThanOne:  "创建群组至少需要两个人哦",
	ErrorUserGroupAlreadyIn:                "已经在群组中了，不要重复添加",
	ErrorUserGroupForbidJoin:               "群组禁止陌生人添加",
	ErrorUserGroupMaxMember:                "群组已达到人数限制",
	ErrorUserGroupHasNoRight:               "用户没有这个操作的权限",
	ErrorUserGroupCannotPointSelf:          "不可以指定自己进行操作",
	ErrorUserGroupUserLessThanTwo:          "群组不能少于2个人，请直接解散",
	ErrorUserGroupUserManagerCannotDeleted: "不可以直接删除群管理员",
	ErrorUserGroupOwnnerCannotLeave:        "群主，不可以直接离开",
	ErrorUserGroupNotMember:                "不是群成员，不能进行此操作",
	ErrorUserGroupBeenMute:                 "您已经被禁言",
	ErrorUserGroupBeenBlock:                "您已经被此群拉黑",
}

func Error(code int32) *errmsg.ErrMsg {
	return &errmsg.ErrMsg{
		Code:    code,
		Message: ErrorDesc(code),
	}
}

func ErrorDesc(code int32) string {
	e, ok := errText[code]
	if ok {
		return e
	}
	return "internal error"
}
