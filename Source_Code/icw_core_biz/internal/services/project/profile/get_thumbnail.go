package profile

import (
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/repositories/minio"
)

// GetProjectThumbnail 获取项目缩略图
func (s *Service) GetProjectThumbnail(req *project.GetProjectThumbnailRequest, resp *project.GetProjectThumbnailResponse) error {
	return s.CallRPC("ProjectProfileService.GetProjectThumbnail", req, resp, func() error {
		return s.getProjectThumbnail(req, resp)
	})
}

func (s *Service) getProjectThumbnail(req *project.GetProjectThumbnailRequest, resp *project.GetProjectThumbnailResponse) error {
	// 获取项目缩略图下载预签名 URL
	thumbnailURL, err := minio.PresignProjectThumbnailURL(s.Ctx(), s.MinIO(), req.ProjectId, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}
	resp.ThumbnailURL = thumbnailURL

	return nil
}
