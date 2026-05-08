package profile

import (
	"context"

	"icw_common/gen/core/biz"

	"icw_core_biz/repositories/minio"
)

// GetProjectThumbnail 获取项目缩略图
func (s *Service) GetProjectThumbnail(ctx context.Context, req *bizpb.GetProjectThumbnailRequest) (*bizpb.GetProjectThumbnailResponse, error) {
	resp := &bizpb.GetProjectThumbnailResponse{}
	err := s.CallRPC(req, func() error {
		return s.getProjectThumbnail(req, resp)
	})
	return resp, err
}

func (s *Service) getProjectThumbnail(req *bizpb.GetProjectThumbnailRequest, resp *bizpb.GetProjectThumbnailResponse) error {
	// 获取项目缩略图下载预签名 URL
	thumbnailURL, err := minio.PresignProjectThumbnailURL(s.Ctx(), s.MinIO(), s.Redis(), req.ProjectId, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}
	resp.ThumbnailUrl = thumbnailURL

	return nil
}
