package auth

import (
	"context"

	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositoies"
)

// Me 获取用户信息
func (s *Service) Me(req *dto.MeRequest, resp *dto.MeResponse) error {
	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	// 校验 Access Token：
	// 1. JWT 签名是否正确
	// 2. Token 是否过期
	// 3. 签名算法是否是允许的 HMAC
	claims, err := s.tokens.Verify(req.AccessToken)
	if err != nil {
		return rpc_err.UnauthorizedDefault(err.Error())
	}

	// 检查 Access Token 是否已被黑名单禁用
	if claims.ID != "" && s.redis.AccessTokenBlacklisted(ctx, claims.ID) {
		return rpc_err.UnauthorizedDefault("token is blacklisted")
	}

	// 按用户 ID 查询用户
	user, err := s.mysql.FindUserById(ctx, claims.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return rpc_err.UnauthorizedDefault("user not found")
	}

	resp.User = repositoies.UserRecordToDTO(user)
	return nil
}
