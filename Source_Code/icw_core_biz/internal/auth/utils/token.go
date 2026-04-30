package utils

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"icw_core_biz/configs"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories"
)

// TokenManager JWT 管理器，负责 Access Token 的签发、校验和解析
type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

// AccessClaims Access Token 中携带的 JWT claims
type AccessClaims struct {
	UserId uint64 `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// NewTokenManager 创建 JWT 管理器
func NewTokenManager(secret string, ttl time.Duration) *TokenManager {
	return &TokenManager{secret: []byte(secret), ttl: ttl}
}

// Sign 签发 Access Token
func (m *TokenManager) Sign(user *dto.User) (string, string, error) {
	now := time.Now()
	jti := uuid.NewString()
	claims := AccessClaims{
		UserId: user.Id,
		Email:  user.Email,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			Issuer:    "icw.core.biz",
			Subject:   user.Email,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	return signed, jti, err
}

// Verify 校验 Access Token
// 1. JWT 签名是否正确
// 2. Token 是否过期
// 3. 签名算法是否是允许的 HMAC
func (m *TokenManager) Verify(raw string) (*AccessClaims, error) {
	claims := &AccessClaims{}
	token, err := jwt.ParseWithClaims(raw, claims, func(token *jwt.Token) (interface{}, error) {
		// 只接受 HMAC 签名算法，避免算法混淆类问题
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return m.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// ParseAny 在不校验签名和过期时间的情况下解析 claims
// 用于将未过期的 Access Token 加入 Redis 黑名单
func (m *TokenManager) ParseAny(raw string) (*AccessClaims, error) {
	claims := &AccessClaims{}
	_, _, err := jwt.NewParser().ParseUnverified(raw, claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// IssueTokens 签发 Access Token 和 Refresh Token
func IssueTokens(ctx context.Context, cfg configs.Config, mysql *repositories.MySQLRepository, tokens *TokenManager, user *repositories.UserRecord, res interface{}) error {
	if user == nil {
		return rpc_err.InternalErrorDefault("user is nil")
	}

	// 签发短期 Access Token
	// jti 已包含在 Access Token 中，当前函数不需要额外返回
	accessToken, _, err := tokens.Sign(repositories.UserRecordToDTO(user))
	if err != nil {
		return err
	}

	// 生成长期 Refresh Token 并落库
	refreshToken, tokenId, err := NewRefreshToken()
	if err != nil {
		return err
	}
	expiresAt := time.Now().Add(cfg.RefreshTokenTTL)
	tokenHash := HashRefreshToken(refreshToken)

	// Login 和 Refresh 返回同样的 Token 响应字段，因此共用该函数填充响应体
	switch out := res.(type) {
	case *dto.LoginResponse:
		if err := mysql.CreateLoginSession(ctx, tokenId, user.Id, tokenHash, expiresAt); err != nil {
			return err
		}
		out.AccessToken = accessToken
		out.AccessTokenExpiresIn = int(cfg.AccessTokenTTL.Seconds())
		out.RefreshToken = refreshToken
		out.RefreshTokenExpiresIn = int(cfg.RefreshTokenTTL.Seconds())
		out.User = repositories.UserRecordToDTO(user)
	case *dto.RefreshResponse:
		if err := mysql.InsertRefreshToken(ctx, tokenId, user.Id, tokenHash, expiresAt); err != nil {
			return err
		}
		out.AccessToken = accessToken
		out.AccessTokenExpiresIn = int(cfg.AccessTokenTTL.Seconds())
		out.RefreshToken = refreshToken
		out.RefreshTokenExpiresIn = int(cfg.RefreshTokenTTL.Seconds())
		out.User = repositories.UserRecordToDTO(user)
	default:
		return rpc_err.InternalErrorDefault("invalid response type")
	}

	return nil
}
