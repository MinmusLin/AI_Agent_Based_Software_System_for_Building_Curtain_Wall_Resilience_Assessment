package user

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_core_api/rpc/common"
	"icw_core_api/rpc/icw_core_biz"
)

// DeleteAvatar 删除用户自定义头像
func DeleteAvatar(ctx context.Context, client *icw_core_biz.Client, req *bizpb.DeleteAvatarRequest, resp *bizpb.DeleteAvatarResponse) error {
	return common.CallGRPC[bizpb.DeleteAvatarRequest, bizpb.DeleteAvatarResponse](ctx, client, req, resp, client.User().DeleteAvatar)
}

// GetAvatar 获取用户头像
func GetAvatar(ctx context.Context, client *icw_core_biz.Client, req *bizpb.GetAvatarRequest, resp *bizpb.GetAvatarResponse) error {
	return common.CallGRPC[bizpb.GetAvatarRequest, bizpb.GetAvatarResponse](ctx, client, req, resp, client.User().GetAvatar)
}

// UploadAvatar 上传用户自定义头像
func UploadAvatar(ctx context.Context, client *icw_core_biz.Client, req *bizpb.UploadAvatarRequest, resp *bizpb.UploadAvatarResponse) error {
	return common.CallGRPC[bizpb.UploadAvatarRequest, bizpb.UploadAvatarResponse](ctx, client, req, resp, client.User().UploadAvatar)
}
