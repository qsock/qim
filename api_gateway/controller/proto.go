package controller

const (
	fileProto = `syntax = "proto3";

package file;
import "ret.proto";
import "errmsg.proto";
import "model.proto";

option go_package="github.com/qsock/qim/lib/proto/file";

service File {
  rpc Ping (ret.NoArgs) returns (ret.NoArgs) {}
  rpc GetUploadToken(GetUploadTokenReq) returns (GetUploadTokenResp) {}
  rpc GetUserFile(GetUserFileReq) returns (GetUserFileResp) {}
  rpc UploadFileByUrl(UploadFileByUrlReq) returns (UploadFileByUrlResp) {}
  rpc UserUploadSucceed(UserUploadSucceedReq) returns (ret.EmptyResp) {}
  rpc GetSysAvatars(GetSysAvatarsReq) returns (GetSysAvatarsResp) {}

  // 得到省和市
  rpc GetProvinceAndCity(GetProvinceAndCityReq) returns (GetProvinceAndCityResp) {}
}


message GetProvinceAndCityReq {
}

message GetProvinceAndCityResp {
  // 省
  repeated model.Cnarea2019 provinces=1;
  // 城市
  repeated model.Cnarea2019 cities=2;
}

message UserUploadSucceedReq {
  int64 user_id=1;
  string url=2;
  string path=3;
}

message GetUserFileReq {
  int64 user_id=1;
  string path=2;
  int32 page=3;
  int32 page_size=4;
}

message GetUserFileResp {
  errmsg.ErrMsg err=1;
  repeated UserFile files=2;
}

message UserFile {
  int64 id=1;
  int64 user_id=2;
  string url=3;
  string path=4;
  int64 created_on=5;
}

enum UploadType {
  // 上传到本地
  UploadLocal=0;
  // 上传到阿里云oss
  UploadOss=1;
  // 上传到qiniu
  UploadQiniu=2;
  // 上传到腾讯云cos
  UploadCos = 3;
}

// 获取upload token
message GetUploadTokenReq {
  int64 user_id=1;
  // url path
  string path=2;
}

// 用户上传的回调
message GetUploadTokenResp {
  errmsg.ErrMsg err=1;
  map<int32,string> tokens=2;
  string path=3;
}

message UploadFileByUrlReq {
  int64 user_id=1;
  // url path
  string path=2;
  // url
  string url=3;
}

message UploadFileByUrlResp {
  string url=1;
}

message GetSysAvatarsReq {
}

message GetSysAvatarsResp {
   repeated string avatars=1;
}`
	errmsgProto = `syntax = "proto3";

package errmsg;
option go_package="github.com/qsock/qim/lib/proto/errmsg";

message ErrMsg {
  int32 code = 1;
  string message = 2;
}`
	wsProto = `syntax = "proto3";

package ws;
import "ret.proto";
//import "stream.proto";

option go_package="github.com/qsock/qim/lib/proto/ws";

// websocket server
service Ws {
  rpc Ping (ret.NoArgs) returns (ret.NoArgs) {}

  // 用户是否在线
  rpc IsSessOnline(IsSessOnlineReq) returns (ret.BoolResp) {}
  // 关闭用户
  rpc CloseUser(CloseUserReq) returns (ret.EmptyResp) {}
  // 交换令牌
  rpc Exchange(ExchangeReq) returns(ret.EmptyResp) {}

  rpc Msg(MsgReq) returns (ret.EmptyResp) {}
  rpc AllMsg(AllMsgReq) returns (ret.EmptyResp) {}
}

message ExchangeReq {
  string uuid=1;
  string sess_id=2;
}

// 发送消息
message MsgReq {
  string sess_id=1;
  bytes content=2;
}

// 发送消息
message AllMsgReq {
  bytes content =1;
}

// 关闭用户
message CloseUserReq {
  string sess_id=1;
  bytes content=2;
}

// 得到在线用户
message IsSessOnlineReq {
  string sess_id=1;
}`
	retProto = `syntax = "proto3";

import "errmsg.proto";

package ret;
option go_package="github.com/qsock/qim/lib/proto/ret";

message NoArgs {}

message EmptyResp {
  errmsg.ErrMsg err = 1;
}

message IntResp {
  errmsg.ErrMsg err = 1;
  int64 val = 2;
}

message BoolResp {
  errmsg.ErrMsg err = 1;
  bool flag = 2;
}

message StringResp {
  errmsg.ErrMsg err = 1;
  string str = 2;
}

message BytesResp {
  errmsg.ErrMsg err = 1;
  bytes val = 2;
}`
	idProto = `syntax = "proto3";

package id;
import "ret.proto";
import "errmsg.proto";

option go_package="github.com/qsock/qim/lib/proto/id";

service Id {
  rpc Ping (ret.NoArgs) returns (ret.NoArgs) {}

  rpc RegistKey (RegistKeyReq) returns (ret.EmptyResp) {}
  rpc SnowflakeIdToTime (SnowflakeIdToTimeReq) returns (SnowflakeIdToTimeResp) {}
  rpc GenSnowflakeId(GenSnowflakeIdReq) returns (GenSnowflakeIdResp) {}
  rpc GenDbId(GenDbIdReq) returns (GenDbIdResp) {}
}

message KeyItem {
  string key=1;
  int64 offset=2;
  int32 size=3;
}

// 注册一个key
message RegistKeyReq {
  repeated KeyItem keys = 1;
}

message GenSnowflakeIdReq {
}

message GenSnowflakeIdResp {
  errmsg.ErrMsg err=1;
  int64 id=2;
}

message GenDbIdReq {
  string key=1;
}

message GenDbIdResp {
  errmsg.ErrMsg err=1;
  int64 id=2;
}

message SnowflakeIdToTimeReq {
  int64 id=1;
}

message SnowflakeIdToTimeResp {
  errmsg.ErrMsg err=1;
  int64 unix_time=2;
}
`
	eventProto = `syntax = "proto3";

import "ret.proto";
import "model.proto";

package event;
option go_package="github.com/qsock/qim/lib/proto/event";

service Event {
  rpc Ping (ret.NoArgs) returns (ret.NoArgs) {}
}

message HttpReport {
  model.RequestMeta meta=1;
  string method = 2;
  int32 status_code = 3;
  int64 end_on=4;
  string path=5;
  string alias = 6;
  string headers=7;
  string req=8;
  string resp=9;
}

message RpcReport {
  string err=1;
  string method=2;
  int64 created_on=3;
  int64 end_on=4;
  string req=5;
  string resp=6;
}

message Fresher{
  int64 user_id=1;
}

message NewMsg {
  string chat_id=1;
  int64 msg_id=2;
  repeated int64 ids=3;
}

message Friend {
  // 发起操作的人
  int64 operator_id=1;
  int64 recver_id=2;
  // 附带的信息s
  string content=3;
  int64 id=4;
}

message GroupNewApply {
  int64 group_id=1;
  int64 operator_id=2;
  int64 recver_id=3;
  string content=4;
  int64 id=5;
}

message GroupManager {
  int64 user_id=1;
  int64 group_id=2;
  int64 manager_id=3;
  bool is_manager=4;
}

message GroupMember {
  int64 operator_id=1;
  int64 user_id=2;
  int64 group_id=3;
  bool flag=4;
  int64 time=5;
}

message GroupUpdate {
  int64 group_id=1;
  string str=2;
  int64 time=3;
  int32 type=4;
  int64 user_id=5;
}

message GroupDismiss {
  int64 operator_id=1;
  int64 group_id=2;
  repeated int64 ids=3;
}`
	modelProto = `syntax = "proto3";

package model;
option go_package="github.com/qsock/qim/lib/proto/model";

enum Device {
  DeviceFalse=0;
  // mac
  DeviceMac=1;
  // win
  DeviceWin=2;
  // 安卓
  DeviceAndroid=3;
  // ios
  DeviceIos=4;
  // web
  DeviceWeb=5;
  // 小程序
  DeviceMiniProgram=6;
}

enum Gender {
  GenderUnknown = 0;
  GenderMale=1;
  GenderFemale=2;
  GenderSecret=3;
}

message RequestMeta {
  string trace_id = 1;
  Device device = 2;
  string app_name = 3;
  string app_version = 4;
  string device_id = 5;
  string user_ip = 6;
  string lat= 7;
  string lng = 8;
  int64 user_id=9;
  int64 created_on=10;
}

// 用户auth信息
message UserAuth {
  // 用户的id
  int64 user_id=1;
  // 用户鉴权用的token,jwt的token,24小时过期，用refresh-token去刷新使用
  // 每次需要放在header:x-token中传过来
  string token=2;
  // 刷新的token
  string refresh_token=3;
}

message Cnarea2019 {
  // 区域的id
  int32 id=1;
  // 等级
  int32 level=2;
  // 父级行政代码
  int64 parent_code=3;
  // 行政代码
  int64 area_code=4;
  // 邮政编码
  int32 zip_code=5;
  // 区号
  string city_code=6;
  // 名称
  string name=7;
  // 简称
  string short_name=8;
  // 组合名
  string merger_name=9;
  // 拼音
  string pinyin=10;
  // 经度
  double lng=11;
  // 纬度
  double lat=12;
}

message Gift {
  // 礼物id
  int64 id=1;
  // 礼物封面图
  string cover_img=2;
  // 礼物的序列帧图
  repeated string imgs=3;
  // 礼物价格
  int64 gift_cost=4;
}`
	streamProto = `syntax = "proto3";

//import "errmsg.proto";
import "model.proto";

package stream;
option go_package="github.com/qsock/qim/lib/proto/stream";

// 新消息类型
message NewMsgModel {
  int64 msg_id=1;
  string chat_id=2;
}

// 消息枚举
// 总结了很多枚举，大枚举里面包小枚举，最后发现处理并不容易
// 所以整体处理成一个大枚举
enum StreamType {
  StreamFalse=0;
  Ping=1;
  Pong=2;
  // 连接反馈
  ConnectS2C = 3;
  // 消息部分
  NewMsgC2S=101;
  NewMsgS2C=102;
  DirectMsgS2C = 103;
  // 系统消息
  SysMsgS2C=201;
  DirectSysMsgS2C=202;
}

enum MsgType {
  MsgTypeFalse =0;
  // 普通文本消息
  MsgTypeText=1;
  // 发送音频消息
  MsgTypeAudio=2;
  // 发送视频消息
  MsgTypeVideo=3;
  // 发送图片消息
  MsgTypePic=4;
  // 发送语音电话消息
  MsgTypeAudioCall=5;
  // 发送视频电话消息
  MsgTypeVideoCall=6;
  // 定位消息
  MsgTypeLocation =7;
  // 提示类型消息
  MsgTypeHint=8;
  // 撤回消息
  MsgTypeRevertMsg=9;
  // 只是指令，并不显示
  MsgTypeCommand=10;
  // 发送礼物消息
  MsgTypeGift = 11;
}

message MsgModel {
  // C端忽略
  int64 msg_id=1;
  // C端忽略
  int64 sender_id=2;
  // 接受者id
  int64 recv_id=3;
  // 创建时间
  int64 created_on=4;
  // 消息状态
  MsgStatus status=5;
  // 会话id
  string chat_id=6;
  // 消息类型
  MsgType msg_type=7;
  // 会话类型,1单聊，2群聊，3聊天室
  ChatType chat_type=8;
  // 发送的设备
  model.Device device=9;
  //oneof data {
    // 文本消息
    TextMsg text =10;
    // 音频消息
    AudioMsg audio=11;
    // 图片消息
    PicMsg pic = 12;
    // 视频消息
    VideoMsg video=13;
    // 地理位置消息
    LocationMsg location=14;
    // 其他的额外附带消息
    RevertMsg revert=15;
    // 显示在中间的提示类型消息
    HintMsg hint=16;
    // 命令类型的消息
    CommandMsg command=17;
    // 礼物消息，送礼物的消息
    GiftMsg gift=18;
  //}
}

message GiftMsg {
  // 礼物对象
  model.Gift gift=1;
  // 是否显示跑马灯,即弹窗效果，很酷炫的这种
  bool is_horse_lamp=2;
}

message RevertMsg {
  // 操作者id
  int64 operator_id=1;
  int64 msg_id=2;
  string content=3;
}

// 命令类型的消息
message CommandMsg {
  // 命令的枚举
  CommandType t=1;
  // 是否需要震动/声音
  bool noise = 2;
  // 请求附带的内容信息
  string content=3;

  // 操作者,是谁操作的这个事儿
  int64 operator_id=4;
  // 接受者,谁来接受的这个事儿
  int64 recver_id=5;

  // 额外附带的信息，有可能提供前端，或者自己使用
  string extra=6;
}

enum CommandType {
  CommandTypeFalse=0;
  // 单聊1-1000
  // 会话被清理了
  ChatClear=1;
  // 消息被删除
  MsgDelete=2;
}

message HintMsg {
  // 通知消息的枚举
  HintType t=1;
  // 消息内容
  string content=2;
}

enum HintType {
  HintTypeFalse=0;
  // 变成好友了,可以愉快交谈了等等
  BecomeFriends=1;
  // 创建群组
  GroupCreate = 101;
  // 加入群组
  GroupJoin = 102;
  // 离开群组
  GroupLeave = 103;
  // 群组解散
  GroupDismiss = 104;
  // 群组任命管理员
  GroupCharger = 105;
  // 群组管理员任命取消
  GroupChargerCancel = 106;
  // 修改群组名称
  GroupUpdateName = 107;
  // 修改群组通知
  GroupUpdateNotice = 108;
  // 群组禁言某个人
  GroupMuteSomeOne = 109;
  // 群组取消禁言某个人
  GroupMuteSomeOneCancel = 110;
  // 群组禁言全部人
  GroupMuteAll = 111;
  // 群组取消禁言全部人
  GroupMuteAllCancel=112;
  // 群组头像更新
  GroupUpdateAvatar = 113;
  // 群组修改用户的备注
  GroupUpdateUserRemark = 114;
  // 群组转让
  GroupTransfer=115;
}

message ExtraMsg {
  int64 int_val=1;
  string str_val=2;
  bool bool_val=3;
  double float_val=4;
}

message LocationMsg {
  // 封面图
  string cover_url=1;
  // 经度
  string lng=2;
  // 纬度
  string lat=3;
  // 描述信息
  string desc=4;
}

message TextMsg {
  // 消息内容
  string content=1;
  // @的用户
  repeated int64 at_uids=2;
}

message AudioMsg {
  // 音频时长，毫秒为单位
  int32 duration=1;
  // 音频地址
  string src_url=2;
}

message PicMsg {
  // 原图地址
  string src_url=1;
  // 缩略图地址
  string cover_url=2;
  // 图片宽高
  int32 height=3;
  int32 width=4;
}

message VideoMsg {
  // 视频地址
  string src_url=1;
  // 封面图地址
  string cover_url=2;
  // 图片宽高
  int32 height=3;
  int32 width=4;
  // 视频时长，毫秒为单位
  int32 duration=5;
}

// 会话类型
enum ChatType {
  ChatTypeFalse = 0;
  // 默认单聊
  ChatTypeSingle=1;
  // 群聊
  ChatTypeGroup=2;
  // 聊天室
  ChatTypeRoom=3;
}

enum SysMsgType {
  SysMsgTypeFalse =0;
  // 纯文本信息
  Text=1;
  // 指令类型的信息，不需要显示，只是app，作出相应的操作
  Command=2;
  // 弹窗类型的信息
  Pop=3;
  // 点击跳入链接的信息
  Link=4;
  // 关闭的消息
  Close = 5;
}

// 系统消息
// 比如点赞等，都可以用这个系统消息
message SysMsgModel {
  // 消息id
  int64 msg_id=1;
  // 发送者id,0就是系统发送的
  int64 sender_id=2;
  // 接收条件,0就是所有人都接收
  int64 recver_id=3;
  // 创建时间
  int64 created_on=4;
  // 消息类型
  SysMsgType msg_type=5;
  // 是否需要发送push消息
  bool need_push=6;
  // 系统消息的发送时间
  int64 send_on=7;
  // 撤回消息
  MsgStatus status=8;
  // 是否需要存储
  bool need_save=9;
  // 消息体
  // 不用one-of的原因是json无法解析了
  //oneof data {
    // 文本消息
    SysTextMsg text=10;
    // 指令消息
    SysCommandMsg command=11;
    // 弹窗消息
    SysPopMsg pop=12;
    SysLinkMsg link=13;
    SysCloseMsg close=14;
  //}
}

message SysTextMsg {
  // 富文本消息
  string content=1;
  // 封面图
  string cover_url=2;
}

enum SysCommandType {
  SysCommandTypeFalse = 0;
  // 好友申请
  SysCommandFriendApply=1;
  // 好友申请被同意,直接发送消息了，没必要，再发这个消息
  // SysCommandFriendApplyAgree = 2;
  // 好友申请被拒绝
  // SysCommandFriendApplyReject = 3;
  // 好友被删除
  SysCommandFriendDel = 4;
  // 会话已读
  SysChatRead = 5;
  // 会话置顶
  SysChatAhead=6;
  // 会话取消
  SysChatAheadCancel=7;
  // 删除会话
  SysChatDeleted=8;
  // 会话被创建
  SysChatTouch=9;

  // 会话被设置静音
  SysChatMute=10;
  SysChatMuteCancel=11;
  // 备注姓名
  SysCommandFriendMarkname=12;
  SysCommandGroupApply = 13;
}

// 完全静音的操作消息
message SysCommandMsg {
  SysCommandType t=1;
  // 是否触发震动/声音
  bool noise=2;
  // 请求附带的内容信息
  string content=3;
  // 操作者,是谁操作的这个事儿
  int64 operator_id=4;
  // 接受者,谁来接受的这个事儿
  int64 recver_id=5;
  // 额外附带的信息，有可能提供前端，或者自己使用
  string extra=6;
}

message SysPopMsg {
   // 弹窗出现多长时间
  int64 duration=1;
  string content =2;
  // 弹窗开始时间
  int64 start_at=3;
  // 是否可以手动关闭
  bool can_close=4;
}

message SysLinkMsg {
  string url=1;
  //  简单描述
  string desc=2;
  // 封面图
  string cover_url=3;
}

message SysCloseMsg {
  // 操作者,是谁操作的这个事儿
  int64 operator_id=1;
  // 接受者,谁来接受的这个事儿
  int64 recver_id=2;
  // 附带信息
  string content=3;
  // need_pop,是否需要弹窗
  bool pop=4;
}

enum MsgStatus {
  MsgNormal=0;
  MsgRevert=1;
  MsgDeleted=2;
}

//message Packet {
//  StreamType t=1;
//  bytes content=2;
//  errmsg.ErrMsg err=3;
//}

message ConnectMsgModel {
  string uuid=1;
  string key=2;
}`
	passportProto = `syntax = "proto3";

package passport;

import "ret.proto";
import "model.proto";
import "errmsg.proto";

option go_package="github.com/qsock/qim/lib/proto/passport";

service Passport {
  rpc Ping (ret.NoArgs) returns (ret.NoArgs) {}

  // 发送手机短信
  rpc Sms(SmsReq) returns (ret.EmptyResp) {}
  rpc TelLogin(TelLoginReq) returns (LoginResp) {}
  rpc Qqlogin(QqloginReq) returns (LoginResp) {}
  rpc WxLogin(WxLoginReq) returns (LoginResp) {}

  // 鉴权
  rpc Auth (AuthReq) returns(ret.IntResp) {}
  rpc Refresh(RefreshReq) returns (ret.StringResp) {}
  rpc Logout (LogoutReq) returns (ret.EmptyResp) {}

  rpc Ban(BanReq) returns (ret.EmptyResp) {}
  rpc UnBan(UnBanReq) returns (ret.EmptyResp) {}
}

message LoginResp {
  errmsg.ErrMsg err=1;
  // 用户auth信息
  model.UserAuth auth=2;
}

message SmsReq {
  // 手机号
  string tel=1;
}

message SmsModel {
  // 验证码
  string code=1;
  // 创建时间
  int64 created_on=2;
}

// 手机短信登陆
message TelLoginReq {
  //  忽略
  model.RequestMeta meta=1;
  // 手机号
  string tel=2;
  // 验证码
  string code=3;
}

// qq登陆
message QqloginReq {
  //  忽略
  model.RequestMeta meta=1;
  // qq返回的accesstoken
  string token=2;
  // qq返回的openid
  string openid=3;
}

// 微信登陆
message WxLoginReq {
  //  忽略
  model.RequestMeta meta=1;
  // 微信登陆的code
  string code=2;
}

message RefreshReq {
  //  忽略
  model.RequestMeta meta=1;
  string token=2;
  // refresh_token
  string refresh_token=3;
}

message LogoutReq {
  //  忽略
  model.RequestMeta meta=1;
}

message AuthReq {
  model.RequestMeta meta=1;
  string token=2;
}

message JwtClaims {
  model.Device device=1;
  int64 user_id=2;
  string user_ip=3;
  int64 seq_id=4;
}

message BanReq {
  int64 user_id=1;
  int64 end_on=2;
}

message UnBanReq {
  int64 user_id=1;
}`
	msgProto = `syntax = "proto3";
package msg;

import "ret.proto";
import "stream.proto";
import "errmsg.proto";
import "model.proto";

option go_package="github.com/qsock/qim/lib/proto/msg";

service Msg {
  rpc Ping (ret.NoArgs) returns (ret.NoArgs) {}
  rpc Msg (MsgReq) returns (MsgResp) {}
  rpc SysMsg (SysMsgReq) returns (ret.IntResp) {}
  rpc RevertSelfMsg (RevertSelfMsgReq) returns (ret.EmptyResp){}
  rpc ManagerChatMsgRevert (ManagerChatMsgRevertReq) returns (ret.EmptyResp){}
  rpc GetSysMsg (GetSysMsgReq) returns (GetSysMsgResp){}
  rpc GetMemberIdByChatId (GetMemberIdByChatIdReq) returns (GetMemberIdByChatIdResp){}
  rpc UserClosed (UserClosedReq) returns (ret.EmptyResp){}
  rpc SessConnect(SessConnectReq) returns (ret.BytesResp) {}
  rpc Exchange (ExchangeReq) returns (ret.EmptyResp){}
  rpc CloseWithMsg (CloseWithMsgReq) returns (ret.EmptyResp){}

  rpc MarkChatRead (MarkChatReadReq) returns (ret.EmptyResp){}
  rpc ChatAhead (ChatAheadReq) returns (ret.EmptyResp){}
  // 创建会话
  rpc ChatTouch (ChatTouchReq) returns (ret.EmptyResp){}
  // 删除单边会话,下一个消息，还会回来
  rpc ChatRemove (ChatRemoveReq) returns (ret.EmptyResp){}
  // 清空会话
  rpc ChatClear (ChatClearReq) returns (ret.EmptyResp){}
  rpc ChatByUids (ChatByUidsReq) returns (ChatByUidsResp){}
  rpc ChatIds (ChatIdsReq) returns (ChatIdsResp){}
  rpc ChatByIds (ChatByIdsReq) returns (ChatByIdsResp){}
  rpc ChatRecordIds (ChatRecordIdsReq) returns (ChatRecordIdsResp){}
  rpc ChatRecordByIds (ChatRecordByIdsReq) returns (ChatRecordByIdsResp){}
  rpc ChatMute (ChatMuteReq) returns (ret.EmptyResp){}
}

message SessConnectReq {
  string sess_id=1;
  string server_key=2;
}

message SysMsgReq {
  stream.SysMsgModel m=1;
}

message MsgReq {
  stream.MsgModel m=1;
}

message MsgResp {
  errmsg.ErrMsg err=1;
  int64 msg_id=2;
  string chat_id=3;
}

// 标记会话已读
message MarkChatReadReq {
  // 哪个用户去标记这个会话已读
  int64 user_id=1;
  // 标记这个会话已读
  string chat_id=2;
}

// 置顶会话
message ChatAheadReq {
  int64 user_id=1;
  // 会话id
  string chat_id=2;
  bool is_ahead=3;
}

// 清理会话，清理所有人的
message ChatClearReq {
  string chat_id=1;
}

// 删除单边会话
message ChatRemoveReq {
  int64 user_id=1;
  string chat_id=2;
}

// 新建会话
message ChatTouchReq {
  int64 user_id=1;
  int64 recver_id=2;
  stream.ChatType type=3;
}

// 得到所有会话id
message ChatIdsReq {
  int64 user_id=1;
}

message ChatIdsResp {
  errmsg.ErrMsg err=1;
  repeated string ids=2;
}

message ChatByUidsReq {
  int64 user_id=1;
  int32 page=2;
  int32 page_size=3;
}

message ChatByUidsResp {
  errmsg.ErrMsg err=1;
  repeated ChatItem items=2;
  int32 total=3;
}

// 获取所有会话
message ChatByIdsReq {
  repeated string ids=1;
  // 只能查自己的
  int64 user_id=2;
}

message ChatItem {
  string chat_id=1;
  // 1单聊，2群聊，3聊天室
  stream.ChatType t= 2;
  // 置顶时间
  int64 ahead_on=3;
  // 更新时间
  int64 updated_on=4;
  int32 unread_ct=5;
  // 最后一条消息
  int64 last_msg_id=6;
  // 已读的最后一条消息
  int64 read_last_msg_id=7;
  // 是否被静音
  bool is_mute=8;
  // 头像
  repeated string avatars=9;
  // 群名称
  string name=10;
  int64 recver_id=11;
}

message ChatByIdsResp {
  errmsg.ErrMsg err=1;
  repeated ChatItem items=2;
}

// 获取聊天记录id
message ChatRecordIdsReq {
  int64 user_id=1;
  string chat_id=2;
  int32 page=3;
  int32 page_size=4;
}

message ChatRecordIdsResp {
  errmsg.ErrMsg err=1;
  repeated int64 ids=2;
  int64 total=3;
}

// 获取聊天记录
message ChatRecordByIdsReq {
  int64 user_id=1;
  // 会话id
  string chat_id=2;
  // 聊天id
  repeated int64 ids=3;
}

// 获取聊天记录
message ChatRecordByIdsResp {
  errmsg.ErrMsg err=1;
  // 消息
  repeated stream.MsgModel items=2;
}

// 将会话静音
message ChatMuteReq {
  int64 user_id=1;
  string chat_id=2;
  bool is_mute=3;
}

// 撤回自己的消息
message RevertSelfMsgReq {
  model.RequestMeta meta=1;
  int64 user_id=2;
  int64 msg_id=3;
  string chat_id=4;
}

// 管理员撤回消息
message ManagerChatMsgRevertReq {
  model.RequestMeta meta=1;
  int64 user_id=2;
  int64 msg_id=3;
  string chat_id=4;
}

// 得到系统消息id
message GetSysMsgReq {
  int64 user_id=1;
  int32 page=2;
  int32 page_size=3;
}

message GetSysMsgResp {
  errmsg.ErrMsg err=1;
  repeated stream.SysMsgModel items=2;
  int32 total=3;
}

// 通过会话id，得到会话成员id, 单聊/群聊/聊天室id
message GetMemberIdByChatIdReq {
  string chat_id=1;
}

message GetMemberIdByChatIdResp {
  errmsg.ErrMsg err=1;
  repeated int64 ids=2;
}

// 用户关闭的请求
message UserClosedReq {
  string sess_id=1;
}

// 交换令牌
message ExchangeReq {
  int64 user_id=1;
  string uuid=2;
  string key=3;
}

// 关闭用户
message CloseWithMsgReq {
  int64 user_id=1;
  string msg = 2;
}

`
	userProto = `syntax = "proto3";

package user;
import "ret.proto";
import "model.proto";
import "errmsg.proto";
import "stream.proto";

option go_package="github.com/qsock/qim/lib/proto/user";

service User {
  rpc Ping (ret.NoArgs) returns (ret.NoArgs) {}
  rpc Create (CreateReq) returns (ret.EmptyResp) {}
  // 用户上一次的更新
  rpc Lastactive(LastactiveReq) returns (ret.EmptyResp) {}
  // 用户的更新状态
  rpc Update (UpdateReq) returns (ret.EmptyResp) {}

  rpc Info(InfoReq) returns (InfoResp) {}
  rpc Infos(InfosReq) returns (InfosResp) {}

  // 更新备注名称
  rpc FriendMarknameUpdate(FriendMarknameUpdateReq) returns (ret.EmptyResp) {}
  // 好友列表
  rpc FriendIds(FriendIdsReq) returns (FriendIdsResp) {}
  rpc FriendByIds(FriendByIdsReq) returns (FriendByIdsResp) {}
  rpc FriendDel(FriendDelReq) returns (ret.EmptyResp) {}
  // 是否是好友
  rpc IsFriend(IsFriendReq) returns (ret.BoolResp) {}

  // 新的朋友列表
  rpc NewApplyUserList (NewApplyUserListReq) returns (NewApplyListResp) {}

  // 添加新好友
  rpc FriendNewApply(FriendNewApplyReq) returns (ret.EmptyResp) {}
  // 同意拒绝添加好友
  rpc FriendNewAgree(FriendNewAgreeReq) returns (ret.EmptyResp) {}
  rpc FriendNewReject(FriendNewRejectReq) returns (ret.EmptyResp) {}
  rpc FriendNewDel(FriendNewDelReq) returns (ret.EmptyResp) {}
  // 忽略这个好友
  rpc FriendNewIgnore(FriendNewIgnoreReq) returns (ret.EmptyResp) {}
  rpc FriendsByUid(FriendsByUidReq) returns (FriendsByUidResp) {}

  // 添加黑名单
  rpc BlacklistAdd(BlacklistAddReq) returns (ret.EmptyResp) {}
  // 添加黑名单
  rpc BlacklistDel(BlacklistDelReq) returns (ret.EmptyResp) {}
  // 黑名单列表
  rpc Blacklist(BlacklistReq) returns (BlacklistResp) {}

  rpc GroupCreate(GroupCreateReq) returns (GroupCreateResp) {}
  rpc GroupInfoById(GroupInfoReq) returns (GroupInfoResp) {}
  rpc GroupInfoByIds(GroupInfosReq) returns (GroupInfosResp) {}
  rpc GroupsByUid(GroupsByUidReq) returns (GroupsByUidResp) {}

  // 任命管理员
  rpc GroupAppointManager(GroupAppointManagerReq) returns (ret.EmptyResp) {}
  rpc GroupManagerList(GroupManagerListReq) returns (GroupManagerListResp) {}
  rpc IsGroupManager(IsGroupManagerReq) returns (ret.BoolResp) {}

  // 群组成员id
  rpc GroupMemberIds(GroupMemberIdsReq) returns (GroupMemberIdsResp) {}
  rpc GroupMemberByIds(GroupMemberByIdsReq) returns (GroupMemberByIdsResp) {}
  rpc GroupMembersByGroupId(GroupMembersByGroupIdReq) returns (GroupMembersByGroupIdResp) {}
  // 是否是群成员
  rpc IsGroupMember(IsGroupMemberReq) returns (ret.BoolResp) {}

  // 管理员直接添加
  rpc GroupMemberAdd(GroupMemberAddReq) returns (ret.EmptyResp) {}
  rpc GroupMemberDel(GroupMemberDelReq) returns (ret.EmptyResp) {}
  rpc GroupLeave(GroupLeaveReq) returns (ret.EmptyResp) {}
  rpc GroupDismiss(GroupDismissReq) returns (ret.EmptyResp) {}

  // 申请加入群组
  rpc GroupJoin(GroupJoinReq) returns (ret.EmptyResp) {}
  rpc GroupJoinAgree(GroupJoinAgreeReq) returns (ret.EmptyResp) {}
  rpc GroupJoinReject(GroupJoinRejectReq) returns (ret.EmptyResp) {}
  rpc GroupJoinIgnore(GroupJoinIgnoreReq) returns (ret.EmptyResp) {}

  // 转让群组
  rpc GroupTransfer(GroupTransferReq) returns (ret.EmptyResp) {}

  // 是否被禁言
  rpc IsGroupMemberBeenMute(IsGroupMemberBeenMuteReq) returns (ret.BoolResp) {}
  rpc GroupMute(GroupMuteReq) returns (ret.EmptyResp) {}
  rpc GroupMuteUser(GroupMuteUserReq) returns (ret.EmptyResp) {}
  rpc GroupMuteList(GroupMuteListReq) returns (GroupMuteListResp) {}
  rpc GroupBlock(GroupBlockReq) returns (ret.EmptyResp) {}
  rpc GroupBlockList(GroupBlockListReq) returns (GroupBlockListResp){}
  rpc GroupUpdateName(GroupUpdateNameReq) returns (ret.EmptyResp) {}
  rpc GroupUpdateNotice(GroupUpdateNoticeReq) returns (ret.EmptyResp) {}
  rpc GroupUpdateAvatar(GroupUpdateAvatarReq) returns (ret.EmptyResp) {}
  rpc GroupUpdateJoinType(GroupUpdateJoinTypeReq) returns (ret.EmptyResp) {}
}

enum NewApplyType {
  NewApplyUser=0;
  NewApplyGroup=1;
}

message GroupUpdateJoinTypeReq {
  int64 user_id=1;
  int64 group_id=2;
  GroupJoinType join_type=3;
}

message GroupUpdateAvatarReq {
  int64 user_id=1;
  int64 group_id=2;
  string avatar=3;
}

message GroupUpdateNoticeReq {
  int64 user_id=1;
  int64 group_id=2;
  string notice=3;
}

message GroupUpdateNameReq {
  int64 user_id=1;
  int64 group_id=2;
  string name=3;
}

message GroupMuteListReq {
  int64 user_id=1;
  int64 group_id=2;
}

message GroupMuteListResp {
  errmsg.ErrMsg err=1;
  repeated GroupMember members = 2;
}

message GroupMuteUserReq{
  int64 user_id=1;
  int64 group_id=2;
  int64 mute_until=3;
  repeated int64 member_id=4;
}

message GroupMuteReq{
  int64 user_id=1;
  int64 group_id=2;
  int64 mute_until=3;
}

message GroupBlockListReq {
  int64 group_id=1;
}

message GroupBlockListResp {
  errmsg.ErrMsg err=1;
  repeated GroupMember members = 2;
}

message GroupBlockReq {
  int64 user_id=1;
  int64 group_id=2;

  // 拉黑的用户
  int64 member_id=3;
  bool is_block=4;
}

message GroupTransferReq{
  int64 user_id=1;
  int64 group_id=2;
  int64 member_id=3;
}

message GroupJoinIgnoreReq {
  int64 user_id=1;
  int64 group_id=2;
  int64 member_id=3;
}

message GroupJoinRejectReq {
  int64 user_id=1;
  int64 group_id=2;
  int64 member_id=3;
  string reason=4;
}

message GroupJoinAgreeReq {
  int64 user_id=1;
  int64 group_id=2;
  int64 member_id=3;
}

message GroupJoinReq {
  int64 user_id=1;
  int64 group_id=2;
  string reason=3;
}

message GroupDismissReq {
  int64 user_id=1;
  int64 group_id=2;
}

message GroupLeaveReq {
  int64 user_id=1;
  int64 group_id=2;
}

message GroupMemberAddReq {
  int64 user_id=1;
  int64 group_id=2;
  repeated int64 member_ids=3;
}

message GroupMemberDelReq {
  int64 user_id=1;
  int64 group_id=2;
  repeated int64 member_ids=3;
}

message GroupMemberByIdsReq {
  repeated int64 user_ids=1;
  int64 group_id=2;
}

message GroupMemberByIdsResp {
  errmsg.ErrMsg err=1;
  repeated GroupMember members = 2;
}

message GroupMembersByGroupIdReq {
  int64 group_id=1;
  int32 page=2;
  int32 page_size=3;
  int64 user_id=4;
}

message GroupMembersByGroupIdResp {
  errmsg.ErrMsg err=1;
  repeated GroupMember members = 2;
  int32 total=3;
}

message IsGroupManagerReq {
  int64 user_id=1;
  int64 group_id=2;
}

message GroupManagerListReq {
  int64 group_id=1;
}

message GroupMember {
  int64 user_id =1;
  GroupRoleType role_type =2;
  string mark_name =3;
  int64 mute_until =4;
  bool not_disturb =5;
  int64 created_on =6;
  bool is_blocked =7;
  UserInfo user=8;
}

enum GroupRoleType {
  GroupRoleNormal=0;
  GroupRoleOwner=1;
  GroupRoleManager=2;
}
message GroupManagerListResp {
  errmsg.ErrMsg err=1;
  repeated GroupMember members = 2;
}

message GroupAppointManagerReq {
  // 忽略，不用传
  int64 user_id=1;
  // 群组id
  int64 group_id=2;
  // 所有的用户id
  repeated int64 manager_ids=3;
  // true指定，false取消
  bool is_appoint=4;
}

message GroupsByUidReq {
  // 忽略，不用传
  int64 user_id=1;
  // 页号，第几页，从0开始
  int32 page=2;
  // 每页数量
  int32 page_size=3;
}

message GroupsByUidResp {
  errmsg.ErrMsg err=1;
  repeated GroupInfo infos = 2;
  int32 total=3;
}

message GroupInfosReq {
  repeated int64 group_ids=1;
}

message GroupInfosResp {
  errmsg.ErrMsg err=1;
  repeated GroupInfo infos = 2;
}

message GroupInfoReq{
  int64 group_id=1;
}

message GroupInfoResp {
  errmsg.ErrMsg err=1;
  GroupInfo info=2;
}

enum GroupJoinType {
  // 需要验证
  GroupJoinAnyone=0;
  GroupJoinVerify=1;
  // 无人可以主动加入
  GroupJoinNone=2;
}

message GroupInfo {
  int64 id=1;
  string name=2;
  int64 mute_util=3;
  string notice=4;
  int64 created_on=5;
  int32 max_member_ct=6;
  repeated string avatars=7;
  GroupJoinType join_type=8;
  int32 current_ct=9;
  int64 deleted_on=10;
}

message GroupCreateReq {
  int64 user_id=1;
  // 名称
  string name=2;
  // 头像
  string avatar=3;
  // 添加的群id
  repeated int64 member_ids=5;
  // 总的人数
  int64 max_member_ct=6;
  // 会话类型,1单聊，2群聊，3聊天室
  stream.ChatType t=7;
}

message GroupCreateResp {
  errmsg.ErrMsg err=1;
  int64 group_id=2;
}

// 群组部分
message GroupMemberIdsReq {
  int64 group_id=1;
}

message GroupMemberIdsResp {
  errmsg.ErrMsg err=1;
  repeated int64 ids=2;
}

message IsGroupMemberReq {
  int64 user_id=1;
  int64 group_id=2;
}

message IsGroupMemberBeenMuteReq {
  int64 user_id=1;
  int64 group_id=2;
}

message BlacklistAddReq {
  int64 user_id=1;
  int64 black_user_id=2;
}

message BlacklistDelReq {
  int64 user_id=1;
  int64 black_user_id=2;
}

message BlacklistReq {
  int64 user_id=1;
  int32 page=2;
  int32 page_size=3;
}

message BlacklistResp {
  errmsg.ErrMsg err=1;
  repeated UserInfo users=2;
  int32 total=3;
}


message InfosReq {
  repeated int64 user_ids=1;
}

message InfosResp {
  errmsg.ErrMsg err=1;
  repeated UserInfo users=2;
}

message InfoReq {
  int64 user_id=1;
}

message InfoResp {
  errmsg.ErrMsg err=1;
  UserInfo user=2;
}

message IsFriendReq {
  int64 user_id=1;
  int64 friend_id=2;
}

message FriendNewIgnoreReq {
  // 忽略，不用传
  int64 user_id=1;
  // 好友id
  int64 friend_id=2;
}

message FriendNewDelReq {
  // 忽略，不用传
  int64 user_id=1;
  // 返回的列表id
  int64 id=2;
}

message FriendNewAgreeReq {
  // 忽略，不用传
  int64 user_id=1;
  // 好友id
  int64 friend_id=2;
}

message FriendNewRejectReq {
  // 忽略，不用传
  int64 user_id=1;
  // 好友id
  int64 friend_id=2;
  // 拒绝理由
  string reason=3;
}

// 用户查询自己相关的
message NewApplyUserListReq {
  int64 user_id=1;
  int32 page=2;
  int32 page_size=3;
}

message NewItem {
  // id
  int64 id=1;
  // 申请人的id
  int64 apply_id=2;
  // 接收群组/人的id
  int64 recver_id=3;
  // 申请类型，群组还是用户，0用户，1群组
  NewApplyType apply_type=4;
  // 申请的次数
  int32 ct=5;
  // 申请是否被忽略
  bool ignore=6;
  // 创建时间
  int64 created_on=7;
  // 申请状态，1成功，2失败
  NewApplyStatus status=8;
  // 更新时间
  int64 updated_on=9;
  // 申请原因
  string reason=10;
  // 操作用户id
  int64 operator_id=11;
  // 操作原因
  string operator_reason=12;
  // 接收的人，如果apply_type是用户，才会有
  UserInfo recv_user=13;
  // 接收的人，如果apply_type是群组，才会有
  GroupInfo recv_group=14;
  // 操作者
  UserInfo operator_user=15;
  // 申请用户
  UserInfo apply_user=16;
}

message NewApplyListResp {
  errmsg.ErrMsg err=1;
  int32 total=2;
  repeated NewItem items=3;
}

message FriendNewApplyReq {
  // 忽略不用传
  int64 user_id=1;
  // 好友id
  int64 friend_id=2;
  // 备注
  string reason=3;
}

// 修改好友备注
message FriendMarknameUpdateReq {
  // 用户id
  int64 user_id=1;
  // 好友id
  int64 friend_id=2;
  // 备注姓名
  string name=3;
}

message FriendDelReq {
  // 忽略
  int64 user_id=1;
  // 好友id
  int64 friend_id=2;
}

message FriendByIdsReq {
  repeated int64 ids=1;
  int64 user_id=2;
}

message FriendItem {
  UserInfo user=1;
  int64 friend_time=2;
  string mark_name=3;
  int64 user_id=4;
}

message FriendByIdsResp {
  errmsg.ErrMsg err=1;
  repeated FriendItem items=2;
}

message FriendIdsReq {
  int64 user_id=1;
}

message FriendIdsResp {
  errmsg.ErrMsg err=1;
  repeated int64 ids=2;
}

message FriendsByUidReq {
  int64 user_id=1;
  int32 page=2;
  int32 page_size=3;
}

message FriendsByUidResp {
  errmsg.ErrMsg err=1;
  // 好友 列表
  repeated FriendItem items=2;
  // 总好友数
  int32 total=3;
}

enum FriendType {
  // 陌生人
  FriendTypeAnnonymous = 0;
  // 好友
  FriendTypeFriend = 1;
  // 黑名单
  FriendTypeBlacklist = 2;
}

message UpdateReq {
  // 忽略，不用传 用户id
  int64 user_id=1;
  // 用户名称
  string name=2;
  // 用户头像
  string avatar=3;
  // 性别
  model.Gender gender=4;
  // 生日
  int64 birthday=5;
  // 简介
  string brief=6;
  // 添加好友的类型
  FriendAddType add_friend_type = 7;
}

message LastactiveReq {
  model.RequestMeta meta=1;
}

message CreateReq{
   int64 user_id=1;
   string name=2;
   string avatar=3;
   // 性别
   model.Gender gender=4;
}

message UserInfo {
  // 用户id
  int64 user_id=1;
  // 用户名称
  string name=2;
  // 用户头像
  string avatar=3;
  // 性别
  model.Gender gender=4;
  // 生日
  int64 birthday=5;
  // 简介
  string brief=6;
  // 添加好友的方式
  FriendAddType add_friend_type = 7;
}

enum FriendAddType {
  FriendAddTypeFalse = 0;
  // 默认是需要询问
  FriendAddNeedAsk=1;
  // 所有人都可以
  FriendAddAnyone=2;
  // 任何人都不行
  FriendAddNone=3;
}

enum NewApplyStatus {
   NewApply=0;
   NewApplySucceed=1;
   NewApplyRejected=2;
}`
)
