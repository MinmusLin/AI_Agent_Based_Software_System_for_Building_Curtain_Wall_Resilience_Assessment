package dto

import (
	"icw_core_biz/pkg/dto"
)

type GetAvatarResponse struct {
	AvatarURL  string `json:"avatar_url"`
	AvatarType string `json:"avatar_type"`
}

func NewGetAvatarResponse(resp *dto.GetAvatarResponse) *GetAvatarResponse {
	if resp == nil {
		return &GetAvatarResponse{}
	}
	return &GetAvatarResponse{
		AvatarURL:  resp.AvatarURL,
		AvatarType: resp.AvatarType,
	}
}

type UploadAvatarResponse struct {
	UploadURL string `json:"upload_url"`
}

func NewUploadAvatarResponse(resp *dto.UploadAvatarResponse) *UploadAvatarResponse {
	if resp == nil {
		return &UploadAvatarResponse{}
	}
	return &UploadAvatarResponse{
		UploadURL: resp.UploadURL,
	}
}

type DeleteAvatarResponse struct{}

func NewDeleteAvatarResponse(resp *dto.DeleteAvatarResponse) *DeleteAvatarResponse {
	if resp == nil {
		return &DeleteAvatarResponse{}
	}
	return &DeleteAvatarResponse{}
}
