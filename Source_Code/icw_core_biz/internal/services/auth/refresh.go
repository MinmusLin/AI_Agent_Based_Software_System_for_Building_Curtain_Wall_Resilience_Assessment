package auth

import (
	"context"
	"strings"
	"time"

	"icw_common/gen/core/biz"
	"icw_common/rpc_err"
	"icw_core_biz/internal/services/auth/utils"
	"icw_core_biz/internal/services/common"
)

// Refresh 刷新 Token
func (s *Service) Refresh(ctx context.Context, req *bizpb.RefreshRequest) (*bizpb.RefreshResponse, error) {
	resp := &bizpb.RefreshResponse{}
	err := s.CallRPC(ctx, req, resp, func() error {
		return s.refresh(req, resp)
	})
	return resp, err
}

func (s *Service) refresh(req *bizpb.RefreshRequest, resp *bizpb.RefreshResponse) error {
	// 解析 Refresh Token Id
	refreshToken := strings.TrimSpace(req.RefreshToken)
	tokenId := utils.ParseRefreshTokenId(refreshToken)
	if tokenId == "" {
		return rpc_err.UnauthorizedDefault("invalid refresh token")
	}

	// 防止同一个 Refresh Token 被多个请求同时刷新
	ok, err := s.Redis().SetRefreshReuseLock(s.Ctx(), tokenId, 30*time.Second)
	if err != nil {
		return err
	}
	if !ok {
		return rpc_err.UnauthorizedDefault("refresh in progress")
	}
	defer func() {
		if err := s.Redis().ClearRefreshReuseLock(s.Ctx(), tokenId); err != nil {
			common.RpcWarn("Clear refresh reuse lock failed, token_id: %s, err: %v", tokenId, err)
		}
	}()

	// 按 Token Id 和 Token Hash 查询 Refresh Token 及所属用户
	token, user, err := s.MySQL().FindRefreshToken(s.Ctx(), tokenId, utils.HashRefreshToken(refreshToken))
	if err != nil {
		return err
	}
	if token == nil || user == nil {
		return rpc_err.UnauthorizedDefault("refresh token not found")
	}

	// 检查 Refresh Token 是否已吊销或已过期
	if token.RevokedAt.Valid {
		return rpc_err.UnauthorizedDefault("refresh token revoked")
	}
	if time.Now().After(token.ExpiresAt) {
		return rpc_err.UnauthorizedDefault("refresh token expired")
	}

	// 签发新的 Access Token 和 Refresh Token，并吊销旧 Refresh Token
	if err := utils.IssueRotatedTokens(s.Ctx(), s.Config(), s.MySQL(), s.tokens, tokenId, user, resp); err != nil {
		return err
	}
	newTokenId := utils.ParseRefreshTokenId(resp.RefreshToken)
	if newTokenId == "" {
		return rpc_err.UnauthorizedDefault("invalid new refresh token")
	}

	return nil
}
