package auth

import (
	"context"
	"strings"
	"time"

	"icw_core_biz/internal/auth/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositoies"
)

// Refresh 刷新 Token
func (s *Service) Refresh(req *dto.RefreshRequest, resp *dto.RefreshResponse) error {
	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	// 解析 Refresh Token Id
	refreshToken := strings.TrimSpace(req.RefreshToken)
	tokenId := utils.ParseRefreshTokenId(refreshToken)
	if tokenId == "" {
		return rpc_err.UnauthorizedDefault("invalid refresh token")
	}

	// 防止同一个 Refresh Token 被多个请求同时刷新
	ok, err := s.redis.SetRefreshReuseLock(ctx, tokenId, 30*time.Second)
	if err != nil {
		return err
	}
	if !ok {
		return rpc_err.UnauthorizedDefault("")
	}
	defer func(redis *repositoies.RedisRepository, ctx context.Context, tokenId string) {
		_ = redis.ClearRefreshReuseLock(ctx, tokenId)
	}(s.redis, ctx, tokenId)

	// 按 Token Id 和 Token Hash 查询 Refresh Token 及所属用户
	token, user, err := s.mysql.FindRefreshToken(ctx, tokenId, utils.HashRefreshToken(refreshToken))
	if err != nil {
		return err
	}
	if token == nil || user == nil {
		return rpc_err.Unauthorized(rpc_err.DetailUnauthorized, "")
	}

	// 检查 Token 是否已吊销或过期
	if token.RevokedAt.Valid || time.Now().After(token.ExpiresAt) {
		return rpc_err.Unauthorized(rpc_err.DetailUnauthorized, "")
	}

	// 签发新的 Access Token 和 Refresh Token
	if err := utils.IssueTokens(ctx, s.cfg, s.mysql, s.tokens, user, resp); err != nil {
		return err
	}
	newTokenId := utils.ParseRefreshTokenId(resp.RefreshToken)
	if newTokenId == "" {
		return rpc_err.Unauthorized(rpc_err.DetailUnauthorized, "")
	}

	// 吊销旧 Refresh Token
	_ = s.mysql.ReplaceRefreshToken(ctx, tokenId, newTokenId)
	return nil
}
