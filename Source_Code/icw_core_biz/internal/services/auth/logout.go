package auth

import (
	"log"
	"time"

	authUtils "icw_core_biz/internal/services/auth/utils"
	"icw_core_biz/internal/services/common"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/utils"
)

// Logout 登出
func (s *Service) Logout(req *dto.LogoutRequest, resp *dto.LogoutResponse) error {
	return s.CallRPC("AuthService.Logout", req, resp, func() error {
		return s.logout(req, resp)
	})
}

func (s *Service) logout(req *dto.LogoutRequest, _ *dto.LogoutResponse) error {
	// 将未过期的 Access Token 加入 Redis 黑名单
	if claims, err := s.tokens.ParseAny(req.AccessToken); err == nil && claims.ID != "" {
		if claims.ExpiresAt != nil {
			ttl := time.Until(claims.ExpiresAt.Time)
			if ttl > 0 {
				if err := s.Redis().BlacklistAccessToken(s.Ctx(), claims.ID, ttl); err != nil {
					log.Printf("%s Blacklist access token failed, claims: %s, err: %v", common.RpcWarnPrefix(), utils.JSONF(claims), err)
				}
			}
		}
	}

	// 吊销 Refresh Token
	if tokenId := authUtils.ParseRefreshTokenId(req.RefreshToken); tokenId != "" {
		if err := s.MySQL().RevokeRefreshTokensByTokenId(s.Ctx(), tokenId); err != nil {
			log.Printf("%s Revoke refresh tokens by token id failed, token_id: %s, err: %v", common.RpcWarnPrefix(), tokenId, err)
		}
	}

	return nil
}
