package user

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_core_biz/internal/services/user/utils"
	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/redis"
)

// DeleteAvatar 删除用户自定义头像
func (s *Service) DeleteAvatar(ctx context.Context, req *bizpb.DeleteAvatarRequest) (*bizpb.DeleteAvatarResponse, error) {
	resp := &bizpb.DeleteAvatarResponse{}
	err := s.CallRPC(ctx, req, resp, func() error {
		return s.deleteAvatar(req, resp)
	})
	return resp, err
}

func (s *Service) deleteAvatar(req *bizpb.DeleteAvatarRequest, _ *bizpb.DeleteAvatarResponse) error {
	// 对标准化邮箱地址做 SHA-256 哈希
	emailHash, err := utils.NormalizeEmailHash(req.Email)
	if err != nil {
		return err
	}

	if s.Redis() != nil {
		// 清除预签名 URL 缓存
		_ = s.Redis().ClearPresignURL(s.Ctx(), redis.GenCustomAvatarPresignURLKey(emailHash))
	}

	// 删除用户自定义头像
	if err := s.MinIO().RemoveObject(s.Ctx(), minio.GenCustomAvatarKey(emailHash)); err != nil {
		return err
	}

	return nil
}
