package profile

import (
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/minio"
)

// DeleteProjectThumbnail 删除项目缩略图
func (s *Service) DeleteProjectThumbnail(req *project.DeleteProjectThumbnailRequest, resp *project.DeleteProjectThumbnailResponse) error {
	return s.CallRPC("ProjectProfileService.DeleteProjectThumbnail", req, resp, func() error {
		return s.deleteProjectThumbnail(req, resp)
	})
}

func (s *Service) deleteProjectThumbnail(req *project.DeleteProjectThumbnailRequest, _ *project.DeleteProjectThumbnailResponse) error {
	// 生成项目缩略图对象 Key
	thumbnailKey, err := minio.GenProjectThumbnailKey(req.ProjectId)
	if err != nil {
		return rpc_err.BadRequestDefault(err.Error())
	}

	// 删除项目缩略图
	if err := s.MinIO().RemoveObject(s.Ctx(), thumbnailKey); err != nil {
		return err
	}

	return nil
}
