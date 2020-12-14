package gclient

import (
	"context"
	"fmt"
	"github.com/qsock/qf/qlog"
	"github.com/qsock/qf/service/qgrpc"
	"github.com/qsock/qim/lib/method"
	"github.com/qsock/qim/lib/proto/user"
)

func FullUserNameLabel(ctx context.Context, id int64) (string, error) {
	nameLabel := `<im user="%d">%s</im>`
	info, err := Info(ctx, id)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(nameLabel, id, info.Name), nil
}

func Info(ctx context.Context, userId int64) (*user.UserInfo, error) {
	creq := &user.InfoReq{UserId: userId}
	cresp := &user.InfoResp{}
	if err := qgrpc.CallIn500ms(ctx, method.UserInfo, creq, cresp); err != nil {
		qlog.Get().Ctx(ctx).Error(err, userId)
		return nil, err
	}
	return cresp.User, nil
}

func Infos(ctx context.Context, userIds []int64) ([]*user.UserInfo, error) {
	if len(userIds) == 0 {
		return nil, nil
	}
	creq := &user.InfosReq{UserIds: userIds}
	cresp := &user.InfosResp{}
	if err := qgrpc.CallIn500ms(ctx, method.UserInfos, creq, cresp); err != nil {
		qlog.Get().Ctx(ctx).Error(err, userIds)
		return nil, err
	}
	return cresp.Users, nil
}

func Groups(ctx context.Context, ids []int64) ([]*user.GroupInfo, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	creq := &user.GroupInfosReq{GroupIds: ids}
	cresp := &user.GroupInfosResp{}
	if err := qgrpc.CallIn500ms(ctx, method.UserGroupInfoByIds, creq, cresp); err != nil {
		qlog.Get().Ctx(ctx).Error(err, ids)
		return nil, err
	}
	return cresp.Infos, nil
}

func FriendByIds(ctx context.Context, ids []int64, userId int64) ([]*user.FriendItem, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	creq := &user.FriendByIdsReq{Ids: ids, UserId: userId}
	cresp := &user.FriendByIdsResp{}
	if err := qgrpc.CallIn500ms(ctx, method.UserFriendByIds, creq, cresp); err != nil {
		qlog.Get().Ctx(ctx).Error(err, ids)
		return nil, err
	}
	return cresp.Items, nil
}
