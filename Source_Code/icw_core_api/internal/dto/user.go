package dto

import (
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

func NewDeleteAvatarResponse(_ *bizpb.DeleteAvatarResponse) *apipb.DeleteAvatarResponse {
	return &apipb.DeleteAvatarResponse{}
}

func NewGetAvatarResponse(resp *bizpb.GetAvatarResponse) *apipb.GetAvatarResponse {
	if resp == nil {
		return nil
	}
	return &apipb.GetAvatarResponse{
		AvatarUrl:  resp.AvatarUrl,
		AvatarType: resp.AvatarType,
	}
}

func NewUploadAvatarResponse(resp *bizpb.UploadAvatarResponse) *apipb.UploadAvatarResponse {
	if resp == nil {
		return nil
	}
	return &apipb.UploadAvatarResponse{
		UploadUrl: resp.UploadUrl,
	}
}
