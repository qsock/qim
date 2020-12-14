package mq

var (
	// prod addr of kafka
	KaAddr = []string{"127.0.0.1:9092"}
)

var (
	// dev addr of kafka
	KaDevAddr = []string{"127.0.0.1:9092"}
)

const (
	ConsumserMsg   = "msg_server_cg"
	ConsumserEvent = "event_server_cg"
)

const (
	// 日志上报的topic
	TopicLogTrace = "log_trace"
	TopicEvent    = "event"
	TopicIm       = "im"
	TopicPush     = "push"
)

const (
	// http的日志上报
	ELogTraceHttp = "http"
	ELogTraceGrpc = "grpc"
)

const (
	// 新用户注册
	EEventFresher = "fresher"

	EEventFriendApply       = "friend_apply"
	EEventFriendApplyReject = "friend_apply_reject"
	//EEventFriendApplyAgree  = "friend_apply_agree"
	EEventFriendApplyDel    = "friend_apply_del"
	EEventFriendApplyIgnore = "friend_apply_ignore"

	EEventGroupApply       = "group_apply"
	EEventGroupApplyReject = "group_apply_reject"
	//EEventGroupApplyAgree     = "group_apply_agree"
	EEventGroupApplyDel       = "group_apply_del"
	EEventGroupApplyIgnore    = "group_apply_ignore"
	EEventGroupManager        = "group_manager"
	EEventGroupNewMember      = "group_new_member"
	EEventGroupDelMember      = "group_del_member"
	EEventGroupLeaveMember    = "group_leave_member"
	EEventGroupUpdateName     = "group_update_name"
	EEventGroupUpdateNotice   = "group_update_notice"
	EEventGroupUpdateAvatar   = "group_update_avatar"
	EEventGroupUpdateJointype = "group_update_jointype"
	EEventGroupUpdateMute     = "group_update_mute"
	EEventGroupMuteone        = "group_mute_one"
	EEventGroupDismiss        = "group_dismiss"
	EEventGroupTransfer       = "group_transfer"
	EEventGroupBlock          = "group_block"

	// 成功添加好友
	EEventFriendAdd      = "friend_add"
	EEventFriendDel      = "friend_del"
	EEventFriendMarkname = "friend_mark_name"
)

const (
	EImNew      = "new"
	EImSys      = "sys"
	EImChatRoom = "chat_room"
)

const (
	EPushNewMsg = "new_msg"
)
