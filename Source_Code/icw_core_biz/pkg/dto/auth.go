package dto

type LoginRequest struct {
	Meta  *Meta
	Email string
	Scene string
	Code  string
}

type LoginResponse struct {
	AccessToken           string
	AccessTokenExpiresIn  int
	RefreshToken          string
	RefreshTokenExpiresIn int
	User                  *User
}

type LogoutRequest struct {
	Meta         *Meta
	AccessToken  string
	RefreshToken string
}

type LogoutResponse struct{}

type MeRequest struct {
	Meta        *Meta
	AccessToken string
}

type MeResponse struct {
	User *User
}

type RefreshRequest struct {
	Meta         *Meta
	RefreshToken string
}

type RefreshResponse struct {
	AccessToken           string
	AccessTokenExpiresIn  int
	RefreshToken          string
	RefreshTokenExpiresIn int
	User                  *User
}

type RegisterRequest struct {
	Meta      *Meta
	Email     string
	EmailCode string
	Password  string
	Name      string
}

type RegisterResponse struct{}

type ResetPasswordRequest struct {
	Meta        *Meta
	Email       string
	EmailCode   string
	NewPassword string
}

type ResetPasswordResponse struct{}

type SendEmailCodeRequest struct {
	Meta  *Meta
	Email string
	Scene string
}

type SendEmailCodeResponse struct {
	ExpiresInSeconds int
}
