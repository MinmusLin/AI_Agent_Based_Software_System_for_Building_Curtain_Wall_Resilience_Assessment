package assets

import (
	"log"
	"strings"

	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
)

// DeleteProjectImage 删除图像
func (s *Service) DeleteProjectImage(req *project.DeleteProjectImageRequest, resp *project.DeleteProjectImageResponse) error {
	return s.CallRPC("ProjectAssetsService.DeleteProjectImage", req, resp, func() error {
		return s.deleteProjectImage(req, resp)
	})
}

func (s *Service) deleteProjectImage(req *project.DeleteProjectImageRequest, _ *project.DeleteProjectImageResponse) error {
	if len(req.ImageUuids) == 0 {
		return rpc_err.BadRequestDefault("image uuids are required")
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

	deleted, err := s.MySQL().DeleteProjectImages(s.Ctx(), req.UserId, req.ProjectId, imageUuids)
	if err != nil {
		return err
	}
	if !deleted {
		return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project group is not accessible")
	}

	for _, imageUuid := range imageUuids {
		if err := utils.RemoveProjectImageObjects(s.Ctx(), s.MinIO(), req.ProjectId, imageUuid); err != nil {
			log.Printf("[WARN] Remove project image objects failed, project_id: %d, image_uuid: %s, err: %v", req.ProjectId, imageUuid, err)
		}
	}

	return nil
}
