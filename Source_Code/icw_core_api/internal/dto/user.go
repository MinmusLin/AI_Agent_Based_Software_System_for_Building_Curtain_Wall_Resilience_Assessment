package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

// NewGetAvatarResponse 将 BIZ 获取头像响应转换为 API 获取头像响应
func NewGetAvatarResponse(resp *bizpb.GetAvatarResponse) *apipb.GetAvatarResponse {
	if resp == nil {
		return nil
	}
	return &apipb.GetAvatarResponse{
		AvatarUrl:  resp.AvatarUrl,
		AvatarType: resp.AvatarType,
	}
}

// NewUploadAvatarResponse 将 BIZ 上传头像响应转换为 API 上传头像响应
func NewUploadAvatarResponse(resp *bizpb.UploadAvatarResponse) *apipb.UploadAvatarResponse {
	if resp == nil {
		return nil
	}
	return &apipb.UploadAvatarResponse{
		UploadUrl: resp.UploadUrl,
	}
}

// NewDeleteAvatarResponse 将 BIZ 删除头像响应转换为 API 删除头像响应
func NewDeleteAvatarResponse(_ *bizpb.DeleteAvatarResponse) *apipb.DeleteAvatarResponse {
	return &apipb.DeleteAvatarResponse{}
}
