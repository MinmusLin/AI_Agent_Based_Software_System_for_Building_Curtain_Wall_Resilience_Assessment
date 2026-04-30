package auth

import (
	"context"
	"strings"
	"time"

	"icw_core_biz/internal/auth/utils"
	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories"
)

// Refresh 刷新 Token
func (s *Service) Refresh(req *dto.RefreshRequest, resp *dto.RefreshResponse) (err error) {
	start := rpc_log.Start("AuthService.Refresh", req)
	defer func() {
		rpc_log.Finish("AuthService.Refresh", req, resp, start, err)
	}()

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
		return rpc_err.UnauthorizedDefault("refresh in progress")
	}
	defer func(redis *repositories.RedisRepository, ctx context.Context, tokenId string) {
		_ = redis.ClearRefreshReuseLock(ctx, tokenId)
	}(s.redis, ctx, tokenId)

	// 按 Token Id 和 Token Hash 查询 Refresh Token 及所属用户
	token, user, err := s.mysql.FindRefreshToken(ctx, tokenId, utils.HashRefreshToken(refreshToken))
	if err != nil {
		return err
	}
	if token == nil || user == nil {
		return rpc_err.Unauthorized(rpc_err.DetailUnauthorized, "refresh token not found")
	}

	// 检查 Refresh Token 是否已吊销或已过期
	if token.RevokedAt.Valid {
		return rpc_err.Unauthorized(rpc_err.DetailUnauthorized, "refresh token revoked")
	}
	if time.Now().After(token.ExpiresAt) {
		return rpc_err.Unauthorized(rpc_err.DetailUnauthorized, "refresh token expired")
	}

	// 签发新的 Access Token 和 Refresh Token，并吊销旧 Refresh Token
	if err := utils.IssueRotatedTokens(ctx, s.cfg, s.mysql, s.tokens, tokenId, user, resp); err != nil {
		return err
	}
	newTokenId := utils.ParseRefreshTokenId(resp.RefreshToken)
	if newTokenId == "" {
		return rpc_err.Unauthorized(rpc_err.DetailUnauthorized, "invalid new refresh token")
	}

	return nil
}
