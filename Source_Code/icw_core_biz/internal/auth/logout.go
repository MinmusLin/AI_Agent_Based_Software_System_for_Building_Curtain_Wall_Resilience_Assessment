package auth

import (
	"context"
	"time"

	"icw_core_biz/internal/auth/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
)

// Logout 登出
func (s *Service) Logout(req *dto.LogoutRequest, _ *dto.LogoutResponse) error {
	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	// 将未过期的 Access Token 加入 Redis 黑名单
	if claims, err := s.tokens.ParseAny(req.AccessToken); err == nil && claims.ID != "" {
		if claims.ExpiresAt != nil {
			ttl := time.Until(claims.ExpiresAt.Time)
			if ttl > 0 {
				_ = s.redis.BlacklistAccessToken(ctx, claims.ID, ttl)
			}
		}
	}

	// 吊销 Refresh Token
	if tokenId := utils.ParseRefreshTokenId(req.RefreshToken); tokenId != "" {
		_ = s.mysql.RevokeRefreshTokenByTokenId(ctx, tokenId)
	}

	return nil
}
