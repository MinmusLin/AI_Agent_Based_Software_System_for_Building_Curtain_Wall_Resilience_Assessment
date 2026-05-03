package dto

type GetAvatarRequest struct {
	Meta   *Meta
	UserId uint64
	Email  string
}

type GetAvatarResponse struct {
	AvatarURL  string
	AvatarType string
}

type UploadAvatarRequest struct {
	Meta        *Meta
	UserId      uint64
	Email       string
	ContentType string
}

type UploadAvatarResponse struct {
	UploadURL string
}

type DeleteAvatarRequest struct {
	Meta   *Meta
	UserId uint64
	Email  string
}

type DeleteAvatarResponse struct{}
