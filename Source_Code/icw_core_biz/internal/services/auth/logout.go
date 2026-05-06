package auth

import (
	"context"
	"time"

	"icw_common/gen/core/biz"
	"icw_common/utils"
	bizUtils "icw_core_biz/internal/services/auth/utils"
	"icw_core_biz/internal/services/common"
)

// Logout 登出
func (s *Service) Logout(ctx context.Context, req *bizpb.LogoutRequest) (*bizpb.LogoutResponse, error) {
	resp := &bizpb.LogoutResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.logout(req, resp)
	})
	return resp, err
}

func (s *Service) logout(req *bizpb.LogoutRequest, _ *bizpb.LogoutResponse) error {
	// 将未过期的 Access Token 加入 Redis 黑名单
	if claims, err := s.tokens.ParseAny(req.AccessToken); err == nil && claims.ID != "" {
		if claims.ExpiresAt != nil {
			ttl := time.Until(claims.ExpiresAt.Time)
			if ttl > 0 {
				if err := s.Redis().BlacklistAccessToken(s.Ctx(), claims.ID, ttl); err != nil {
					common.RpcWarn("Blacklist access token failed, claims: %s, err: %v", utils.JSONF(claims), err)
				}
			}
		}
	}

	// 吊销 Refresh Token
	if tokenId := bizUtils.ParseRefreshTokenId(req.RefreshToken); tokenId != "" {
		if err := s.MySQL().RevokeRefreshTokensByTokenId(s.Ctx(), tokenId); err != nil {
			common.RpcWarn("Revoke refresh tokens by token id failed, token_id: %s, err: %v", tokenId, err)
		}
	}

	return nil
}
