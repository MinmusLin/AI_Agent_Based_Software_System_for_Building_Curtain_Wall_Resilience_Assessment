package dto

import (
	bizDto "icw_core_biz/pkg/dto"
)

type GetAvatarResponse struct {
	AvatarURL string `json:"avatar_url"`
}

func NewGetAvatarResponse(resp *bizDto.GetAvatarResponse) *GetAvatarResponse {
	if resp == nil {
		return &GetAvatarResponse{}
	}
	return &GetAvatarResponse{
		AvatarURL: resp.AvatarURL,
	}
}

type UploadAvatarResponse struct {
	UploadURL string `json:"upload_url"`
}

func NewUploadAvatarResponse(resp *bizDto.UploadAvatarResponse) *UploadAvatarResponse {
	if resp == nil {
		return &UploadAvatarResponse{}
	}
	return &UploadAvatarResponse{
		UploadURL: resp.UploadURL,
	}
}

type DeleteAvatarResponse struct {
	AvatarURL string `json:"avatar_url"`
}

func NewDeleteAvatarResponse(resp *bizDto.DeleteAvatarResponse) *DeleteAvatarResponse {
	if resp == nil {
		return &DeleteAvatarResponse{}
	}
	return &DeleteAvatarResponse{
		AvatarURL: resp.AvatarURL,
	}
}
