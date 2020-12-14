package main

import (
	"github.com/gin-gonic/gin"
	r "github.com/qsock/qim/api_gateway/controller"
	_ "github.com/qsock/qim/api_gateway/docs"
	"github.com/qsock/qim/lib/ginproxy"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetRoute(e *gin.Engine) {
	e.Any("/", ginproxy.OK)
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	qiyou := e.Group("/im")
	{
		qiyou.GET("/proto", r.Proto)
		qiyou.GET("/time", r.ServerTime)
		Passport(qiyou.Group("/passport"))
		User(qiyou.Group("/user"))
		Friend(qiyou.Group("/friend"))
		File(qiyou.Group("/file"))
		Chat(qiyou.Group("/chat"))
		Message(qiyou.Group("/msg"))
		Group(qiyou.Group("/group"))
		ChatRoom(qiyou.Group("/chatroom"))
	}
}

// 账户相关
func Passport(g *gin.RouterGroup) {
	// 手机号登陆
	g.POST("/login/tel", r.LoginTel)
	// 微信 登陆
	g.POST("/login/wx", r.LoginWx)
	// qq 登陆
	g.POST("/login/qq", r.LoginQq)
	// sms 登陆短信
	g.POST("/login/sms", r.LoginSms)
	g.POST("/refresh", r.RefreshToken)
	// 登出用户
	g.POST("/logout", r.Auth, r.Logout)
}

// 用户相关
func User(g *gin.RouterGroup) {
	// 得到用户信息
	g.GET("/info", r.Auth, r.UserInfo)
	g.POST("/lastactive", r.Auth, r.UserLastactive)
	g.POST("/update", r.Auth, r.UserUpdate)
}

func Friend(g *gin.RouterGroup) {
	g.POST("/mark-name", r.Auth, r.FriendMarknameUpdate)
	g.GET("/ids", r.Auth, r.FriendIds)
	g.GET("/by-ids", r.Auth, r.FriendByIds)
	g.GET("/list", r.Auth, r.FriendsByUid)
	g.POST("/del", r.Auth, r.FriendDel)

	g.GET("/new/apply/list", r.Auth, r.NewApplyList)
	g.POST("/new/apply", r.Auth, r.FriendNewApply)
	g.POST("/new/agree", r.Auth, r.FriendNewAgree)
	g.POST("/new/reject", r.Auth, r.FriendNewReject)
	g.POST("/new/del", r.Auth, r.FriendNewDel)
	g.POST("/new/ignore", r.Auth, r.FriendNewIgnore)
}

func Chat(g *gin.RouterGroup) {
	g.POST("/mark-read", r.Auth, r.MarkChatRead)
	g.POST("/ahead", r.Auth, r.ChatAhead)
	g.POST("/touch", r.Auth, r.ChatTouch)
	g.POST("/remove", r.Auth, r.ChatRemove)
	g.POST("/mute", r.Auth, r.ChatMute)

	g.GET("/ids", r.Auth, r.ChatIds)
	g.GET("/by-ids", r.Auth, r.ChatByIds)
	g.GET("/list", r.Auth, r.ChatByUids)
	g.GET("/record/ids", r.Auth, r.ChatRecordIds)
	g.GET("/record/by-ids", r.Auth, r.ChatRecordByIds)
}

func Message(g *gin.RouterGroup) {
	g.GET("/sys", r.Auth, r.GetSysMsg)
	g.POST("/", r.Auth, r.Msg)
	g.POST("/revert", r.Auth, r.RevertSelfMsg)
	g.POST("/revert/by-manager", r.Auth, r.ManagerChatMsgRevert)
	g.POST("/exchange", r.Auth, r.Exchange)
}

// 文件相关
func File(g *gin.RouterGroup) {
	//获取所有头像
	g.GET("/sys/avatars", r.SysAvatars)
	g.GET("/locations", r.Locations)
}

func Group(g *gin.RouterGroup) {
	g.POST("/create", r.Auth, r.GroupCreate)
	g.GET("/groups", r.Auth, r.GroupsByUid)
	g.GET("/info", r.Auth, r.GroupInfoById)
	g.GET("/infos", r.Auth, r.GroupInfoByIds)

	g.POST("/appoint/manager", r.Auth, r.GroupAppointManager)
	g.GET("/managers", r.Auth, r.GroupManagerList)
	g.GET("/member/ids", r.Auth, r.GroupMemberIds)
	g.GET("/member/by-ids", r.Auth, r.GroupMemberByIds)
	g.POST("/members", r.Auth, r.GroupMembersByGroupId)

	g.POST("/member/add", r.Auth, r.GroupMemberAdd)
	g.POST("/member/del", r.Auth, r.GroupMemberDel)
	g.POST("/leave", r.Auth, r.GroupLeave)
	g.POST("/dismiss", r.Auth, r.GroupDismiss)

	g.POST("/join/apply", r.Auth, r.GroupJoin)
	g.POST("/join/agree", r.Auth, r.GroupJoinAgree)
	g.POST("/join/reject", r.Auth, r.GroupJoinReject)
	g.POST("/join/ignore", r.Auth, r.GroupJoinIgnore)

	g.POST("/transfer", r.Auth, r.GroupTransfer)
	g.POST("/mute", r.Auth, r.GroupMute)
	g.POST("/mute/user", r.Auth, r.GroupMuteUser)
	g.GET("/mute/list", r.Auth, r.GroupMuteList)
	g.GET("/block/list", r.Auth, r.GroupBlockList)
	g.POST("/block", r.Auth, r.GroupBlock)

	g.POST("/update/name", r.Auth, r.GroupUpdateName)
	g.POST("/update/notice", r.Auth, r.GroupUpdateNotice)
	g.POST("/update/avatar", r.Auth, r.GroupUpdateAvatar)
	g.POST("/update/jointype", r.Auth, r.GroupUpdateJoinType)

}

func ChatRoom(g *gin.RouterGroup) {

}
