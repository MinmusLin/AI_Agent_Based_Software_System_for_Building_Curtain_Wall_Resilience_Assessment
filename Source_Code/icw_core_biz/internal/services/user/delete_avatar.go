package user

import (
	"icw_core_biz/internal/services/user/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/repositories/minio"
)

// DeleteAvatar 删除用户自定义头像
func (s *Service) DeleteAvatar(req *dto.DeleteAvatarRequest, resp *dto.DeleteAvatarResponse) error {
	return s.CallRPC("UserService.DeleteAvatar", req, resp, func() error {
		return s.deleteAvatar(req, resp)
	})
}

func (s *Service) deleteAvatar(req *dto.DeleteAvatarRequest, _ *dto.DeleteAvatarResponse) (err error) {
	// 对标准化邮箱地址做 SHA-256 哈希
	emailHash, err := utils.NormalizeEmailHash(req.Email)
	if err != nil {
		return err
	}

	// 删除用户自定义头像
	if err := s.MinIO().RemoveObject(s.Ctx, minio.GenCustomAvatarKey(emailHash)); err != nil {
		return err
	}

	return nil
}
