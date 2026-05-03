package assets

import (
	"strings"

	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/mysql"
)

// MoveProjectImage 移动图像
func (s *Service) MoveProjectImage(req *project.MoveProjectImageRequest, resp *project.MoveProjectImageResponse) error {
	return s.CallRPC("ProjectAssetsService.MoveProjectImage", req, resp, func() error {
		return s.moveProjectImage(req, resp)
	})
}

func (s *Service) moveProjectImage(req *project.MoveProjectImageRequest, resp *project.MoveProjectImageResponse) error {
	if len(req.ImageUuids) == 0 {
		return rpc_err.BadRequestDefault("image uuids are required")
	}
	if req.TargetGroupId == 0 {
		return rpc_err.BadRequestDefault("target group id is required")
	}

	imageUuids := make([]string, 0, len(req.ImageUuids))
	imageUuidSet := make(map[string]struct{}, len(req.ImageUuids))
	for _, rawImageUuid := range req.ImageUuids {
		imageUuid := strings.TrimSpace(rawImageUuid)
		if imageUuid == "" {
			return rpc_err.BadRequestDefault("image uuid is required")
		}
		if _, ok := imageUuidSet[imageUuid]; ok {
			return rpc_err.BadRequestDefault("image uuid is duplicated")
		}
		imageUuidSet[imageUuid] = struct{}{}
		imageUuids = append(imageUuids, imageUuid)
	}

	imageRecords, err := s.MySQL().MoveProjectImages(s.Ctx(), req.UserId, req.ProjectId, imageUuids, req.TargetGroupId)
	if err != nil {
		return err
	}
	if len(imageRecords) != len(imageUuids) {
		return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project group is not accessible")
	}

	resp.Images = make([]*project.ProjectImage, 0, len(imageRecords))
	for _, imageRecord := range imageRecords {
		projectImage, err := mysql.ProjectImageRecordToDTO(s.Ctx(), s.MinIO(), imageRecord, s.Config().ProjectImageGetTTL)
		if err != nil {
			return err
		}
		resp.Images = append(resp.Images, projectImage)
	}

	return nil
}
