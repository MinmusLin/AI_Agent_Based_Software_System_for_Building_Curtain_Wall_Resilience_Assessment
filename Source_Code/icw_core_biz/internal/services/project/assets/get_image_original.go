package assets

import (
	"strings"

	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/minio"
)

// GetProjectImageOriginal 获取原图
func (s *Service) GetProjectImageOriginal(req *project.GetProjectImageOriginalRequest, resp *project.GetProjectImageOriginalResponse) error {
	return s.CallRPC("ProjectAssetsService.GetProjectImageOriginal", req, resp, func() error {
		return s.getProjectImageOriginal(req, resp)
	})
}

func (s *Service) getProjectImageOriginal(req *project.GetProjectImageOriginalRequest, resp *project.GetProjectImageOriginalResponse) error {
	imageUuid := strings.TrimSpace(req.ImageUuid)
	if imageUuid == "" {
		return rpc_err.BadRequestDefault("image uuid is required")
	}

	imageRecord, err := s.MySQL().FindProjectImageByUuid(s.Ctx(), req.UserId, req.ProjectId, imageUuid)
	if err != nil {
		return err
	}
	if imageRecord == nil {
		return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project group is not accessible")
	}
	if imageRecord.Status != consts.ProjectImageStatusUploaded {
		return rpc_err.BadRequestDefault("project image status is not uploaded")
	}

	// 获取项目图像原图下载预签名 URL
	originalURL, err := minio.PresignProjectImageOriginalURL(s.Ctx(), s.MinIO(), s.Redis(), req.UserId, req.ProjectId, imageUuid, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}
	if originalURL == "" {
		return rpc_err.BadRequest(rpc_err.DetailProjectImageExpired, "project image original object is not found")
	}
	resp.OriginalURL = originalURL

	return nil
}
