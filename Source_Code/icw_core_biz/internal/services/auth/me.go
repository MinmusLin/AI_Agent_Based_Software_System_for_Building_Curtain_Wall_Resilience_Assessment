package auth

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/repositories/mysql/model"
)

// Me 获取用户信息
func (s *Service) Me(ctx context.Context, req *bizpb.MeRequest) (*bizpb.MeResponse, error) {
	resp := &bizpb.MeResponse{}
	err := s.CallRPC(req, func() error {
		return s.me(req, resp)
	})
	return resp, err
}

func (s *Service) me(req *bizpb.MeRequest, resp *bizpb.MeResponse) error {
	// 校验 Access Token 的签名、过期时间和签名算法
	claims, err := s.tokens.Verify(req.AccessToken)
	if err != nil {
		return rpc_error.UnauthorizedDefault(err.Error())
	}

	// 检查 Access Token 是否已被黑名单禁用
	if claims.ID != "" {
		blacklisted, err := s.Redis().AccessTokenBlacklisted(s.Ctx(), claims.ID)
		if err != nil {
			return err
		}
		if blacklisted {
			return rpc_error.UnauthorizedDefault("token is blacklisted")
		}
	}

	// 按用户 ID 查询用户
	user, err := s.MySQL().FindUserById(s.Ctx(), claims.UserId)
	if err != nil {
		return err
	}
	if user == nil || user.Id == 0 {
		return rpc_error.UnauthorizedDefault("user not found")
	}

	resp.User = model.UserRecordToDTO(user)

	return nil
}
