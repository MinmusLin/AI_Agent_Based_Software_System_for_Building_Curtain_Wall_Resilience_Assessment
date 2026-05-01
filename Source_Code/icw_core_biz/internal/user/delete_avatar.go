package user

import (
	"context"

	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/internal/user/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/rpc_err"
)

// DeleteAvatar 删除用户自定义头像
func (s *Service) DeleteAvatar(req *dto.DeleteAvatarRequest, resp *dto.DeleteAvatarResponse) (err error) {
	start := rpc_log.Start("UserService.DeleteAvatar", req)
	defer func() {
		rpc_log.Finish("UserService.DeleteAvatar", req, resp, start, err)
	}()

	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	// 对标准化邮箱地址做 SHA-256 哈希
	emailHash, err := utils.NormalizeEmailHash(req.Email)
	if err != nil {
		return err
	}

	// 删除用户自定义头像
	if err := s.minio.RemoveObject(ctx, utils.GenCustomAvatarKey(emailHash)); err != nil {
		return err
	}

	return nil
}
