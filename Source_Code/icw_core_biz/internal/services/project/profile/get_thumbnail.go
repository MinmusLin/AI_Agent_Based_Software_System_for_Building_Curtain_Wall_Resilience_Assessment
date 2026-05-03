package profile

import (
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/repositories/mysql"
)

// GetProjectThumbnail 获取项目缩略图
func (s *Service) GetProjectThumbnail(req *project.GetProjectThumbnailRequest, resp *project.GetProjectThumbnailResponse) error {
	return s.CallRPC("ProjectProfileService.GetProjectThumbnail", req, resp, func() error {
		return s.getProjectThumbnail(req, resp)
	})
}

func (s *Service) getProjectThumbnail(req *project.GetProjectThumbnailRequest, resp *project.GetProjectThumbnailResponse) error {
	// 获取项目缩略图下载预签名 URL
	thumbnailURL, err := mysql.PresignProjectThumbnailURL(s.Ctx(), s.MinIO(), req.ProjectId, s.Config().ProjectThumbnailGetTTL)
	if err != nil {
		return err
	}
	resp.ThumbnailURL = thumbnailURL

	return nil
}
