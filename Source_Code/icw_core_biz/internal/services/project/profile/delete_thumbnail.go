package profile

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc_err"
	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/redis"
)

// DeleteProjectThumbnail 删除项目缩略图
func (s *Service) DeleteProjectThumbnail(ctx context.Context, req *bizpb.DeleteProjectThumbnailRequest) (*bizpb.DeleteProjectThumbnailResponse, error) {
	resp := &bizpb.DeleteProjectThumbnailResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.deleteProjectThumbnail(req, resp)
	})
	return resp, err
}

func (s *Service) deleteProjectThumbnail(req *bizpb.DeleteProjectThumbnailRequest, _ *bizpb.DeleteProjectThumbnailResponse) error {
	// 生成项目缩略图对象 Key
	thumbnailKey, err := minio.GenProjectThumbnailKey(req.ProjectId)
	if err != nil {
		return rpc_err.BadRequestDefault(err.Error())
	}

	if s.Redis() != nil {
		// 清除预签名 URL 缓存
		_ = s.Redis().ClearPresignURL(s.Ctx(), redis.GenProjectThumbnailPresignURLKey(req.UserId, req.ProjectId))
	}

	// 删除项目缩略图
	if err := s.MinIO().RemoveObject(s.Ctx(), thumbnailKey); err != nil {
		return err
	}

	return nil
}
