package auth

import (
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/utils"
)

// Me 获取用户信息
func (s *Service) Me(req *dto.MeRequest, resp *dto.MeResponse) error {
	return s.CallRPC("AuthService.Me", req, resp, func() error {
		return s.me(req, resp)
	})
}

func (s *Service) me(req *dto.MeRequest, resp *dto.MeResponse) (err error) {
	// 校验 Access Token 的签名、过期时间和签名算法
	claims, err := s.tokens.Verify(req.AccessToken)
	if err != nil {
		return rpc_err.UnauthorizedDefault(err.Error())
	}

	// 检查 Access Token 是否已被黑名单禁用
	if claims.ID != "" {
		blacklisted, err := s.Redis().AccessTokenBlacklisted(s.Ctx, claims.ID)
		if err != nil {
			return err
		}
		if blacklisted {
			return rpc_err.UnauthorizedDefault("token is blacklisted")
		}
	}

	// 按用户 ID 查询用户
	user, err := s.MySQL().FindUserById(s.Ctx, claims.UserId)
	if err != nil {
		return err
	}
	if user == nil || user.Id == 0 {
		return rpc_err.UnauthorizedDefault("user not found")
	}

	resp.User = utils.UserRecordToDTO(user)

	return nil
}
