package dto

type GetAvatarRequest struct {
	UserId uint64
	Email  string
}

type GetAvatarResponse struct {
	AvatarURL string
}

type UploadAvatarRequest struct {
	UserId      uint64
	Email       string
	ContentType string
}

type UploadAvatarResponse struct {
	UploadURL string
}

type DeleteAvatarRequest struct {
	UserId uint64
	Email  string
}

type DeleteAvatarResponse struct {
	AvatarURL string
}
