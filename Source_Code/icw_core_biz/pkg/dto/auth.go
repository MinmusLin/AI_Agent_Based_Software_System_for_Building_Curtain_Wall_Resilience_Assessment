package dto

type LoginRequest struct {
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
	AccessToken  string
	RefreshToken string
}

type LogoutResponse struct{}

type MeRequest struct {
	AccessToken string
}

type MeResponse struct {
	User *User
}

type RefreshRequest struct {
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
	Email     string
	EmailCode string
	Password  string
	Name      string
}

type RegisterResponse struct{}

type ResetPasswordRequest struct {
	Email       string
	EmailCode   string
	NewPassword string
}

type ResetPasswordResponse struct{}

type SendEmailCodeRequest struct {
	Email string
	Scene string
}

type SendEmailCodeResponse struct {
	ExpiresInSeconds int
}
