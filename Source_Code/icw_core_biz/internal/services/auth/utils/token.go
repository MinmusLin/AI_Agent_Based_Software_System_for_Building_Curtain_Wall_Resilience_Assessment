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
	"icw_core_biz/repositories/mysql"
	"icw_core_biz/utils"
)

// TokenMetadata Token 元数据
type TokenMetadata struct {
	AccessToken  string
	RefreshToken string
	TokenId      string
	TokenHash    string
	ExpiresAt    time.Time
	User         *dto.User
}

// NewTokenMetadata 创建 Token 元数据
func NewTokenMetadata(cfg configs.Config, tokens *TokenManager, user *mysql.UserRecord) (*TokenMetadata, error) {
	userDTO := utils.UserRecordToDTO(user)
	if userDTO == nil {
		return nil, rpc_err.InternalErrorDefault("user is nil")
	}

	// 签发短期 Access Token
	// jti 已包含在 Access Token 中，当前函数不需要额外返回
	accessToken, _, err := tokens.Sign(userDTO)
	if err != nil {
		return nil, err
	}

	// 生成长期 Refresh Token
	refreshToken, tokenId, err := NewRefreshToken()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(cfg.RefreshTokenTTL)
	tokenHash := HashRefreshToken(refreshToken)

	return &TokenMetadata{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenId:      tokenId,
		TokenHash:    tokenHash,
		ExpiresAt:    expiresAt,
		User:         userDTO,
	}, nil
}

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
	return &TokenManager{
		secret: []byte(secret),
		ttl:    ttl,
	}
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

// IssueTokens 登录场景下，签发 Access Token 和 Refresh Token，并更新用户最近登录时间
func IssueTokens(ctx context.Context, cfg configs.Config, repo *mysql.Repository, tokens *TokenManager, user *mysql.UserRecord, resp *dto.LoginResponse) error {
	pair, err := NewTokenMetadata(cfg, tokens, user)
	if err != nil {
		return err
	}

	// 登录时保存登录态 Refresh Token，并更新用户最近登录时间
	if err := repo.CreateLoginSession(ctx, pair.TokenId, user.Id, pair.TokenHash, pair.ExpiresAt); err != nil {
		return err
	}

	resp.AccessToken = pair.AccessToken
	resp.AccessTokenExpiresIn = int(cfg.AccessTokenTTL.Seconds())
	resp.RefreshToken = pair.RefreshToken
	resp.RefreshTokenExpiresIn = int(cfg.RefreshTokenTTL.Seconds())
	resp.User = pair.User

	return nil
}

// IssueRotatedTokens 刷新场景下，签发新的 Access Token 和 Refresh Token，并吊销旧 Refresh Token
func IssueRotatedTokens(ctx context.Context, cfg configs.Config, repo *mysql.Repository, tokens *TokenManager, oldTokenId string, user *mysql.UserRecord, resp *dto.RefreshResponse) error {
	if resp == nil {
		return rpc_err.InternalErrorDefault("response is nil")
	}
	pair, err := NewTokenMetadata(cfg, tokens, user)
	if err != nil {
		return err
	}

	// 刷新时签发新的 Access Token 和 Refresh Token，并吊销旧 Refresh Token
	if err := repo.RotateRefreshToken(ctx, oldTokenId, pair.TokenId, user.Id, pair.TokenHash, pair.ExpiresAt); err != nil {
		if errors.Is(err, mysql.ErrRefreshTokenNotReplaceable) {
			return rpc_err.Unauthorized(rpc_err.DetailUnauthorized, err.Error())
		}
		return err
	}

	resp.AccessToken = pair.AccessToken
	resp.AccessTokenExpiresIn = int(cfg.AccessTokenTTL.Seconds())
	resp.RefreshToken = pair.RefreshToken
	resp.RefreshTokenExpiresIn = int(cfg.RefreshTokenTTL.Seconds())
	resp.User = pair.User

	return nil
}
