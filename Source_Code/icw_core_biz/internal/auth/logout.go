package auth

import (
	"context"
	"log"
	"time"

	authUtils "icw_core_biz/internal/auth/utils"
	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/utils"
)

// Logout 登出
func (s *Service) Logout(req *dto.LogoutRequest, resp *dto.LogoutResponse) (err error) {
	start := rpc_log.Start("AuthService.Logout", req)
	defer func() {
		rpc_log.Finish("AuthService.Logout", req, resp, start, err)
	}()

	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	// 将未过期的 Access Token 加入 Redis 黑名单
	if claims, err := s.tokens.ParseAny(req.AccessToken); err == nil && claims.ID != "" {
		if claims.ExpiresAt != nil {
			ttl := time.Until(claims.ExpiresAt.Time)
			if ttl > 0 {
				if err := s.redis.BlacklistAccessToken(ctx, claims.ID, ttl); err != nil {
					log.Printf("[WARN] Blacklist access token failed, claims: %s, err: %v", utils.JSONF(claims), err)
				}
			}
		}
	}

	// 吊销 Refresh Token
	if tokenId := authUtils.ParseRefreshTokenId(req.RefreshToken); tokenId != "" {
		if err := s.mysql.RevokeRefreshTokensByTokenId(ctx, tokenId); err != nil {
			log.Printf("[WARN] Revoke refresh tokens by token id failed, token_id: %s, err: %v", tokenId, err)
		}
	}

	return nil
}
