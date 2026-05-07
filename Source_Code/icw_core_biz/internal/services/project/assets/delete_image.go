package assets

import (
	"context"
	"strings"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"
	"icw_core_biz/internal/services/common"
	"icw_core_biz/internal/services/project/utils"
)

// DeleteProjectImage 删除图像
func (s *Service) DeleteProjectImage(ctx context.Context, req *bizpb.DeleteProjectImageRequest) (*bizpb.DeleteProjectImageResponse, error) {
	resp := &bizpb.DeleteProjectImageResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.deleteProjectImage(req, resp)
	})
	return resp, err
}

func (s *Service) deleteProjectImage(req *bizpb.DeleteProjectImageRequest, _ *bizpb.DeleteProjectImageResponse) error {
	if len(req.ImageUuids) == 0 {
		return rpc_error.BadRequestDefault("image uuids are required")
	}

	imageUuids := make([]string, 0, len(req.ImageUuids))
	imageUuidSet := make(map[string]struct{}, len(req.ImageUuids))
	for _, rawImageUuid := range req.ImageUuids {
		imageUuid := strings.TrimSpace(rawImageUuid)
		if imageUuid == "" {
			return rpc_error.BadRequestDefault("image uuid is required")
		}
		if _, ok := imageUuidSet[imageUuid]; ok {
			return rpc_error.BadRequestDefault("image uuid is duplicated")
		}
		imageUuidSet[imageUuid] = struct{}{}
		imageUuids = append(imageUuids, imageUuid)
	}

	deleted, err := s.MySQL().DeleteProjectImages(s.Ctx(), req.UserId, req.ProjectId, imageUuids)
	if err != nil {
		return err
	}
	if !deleted {
		return rpc_error.BadRequest(rpc_error.DetailProjectNotAccessible, "project group is not accessible")
	}

	for _, imageUuid := range imageUuids {
		if err := utils.RemoveProjectImageObjects(s.Ctx(), s.MinIO(), s.Redis(), req.UserId, req.ProjectId, imageUuid); err != nil {
			common.RpcWarn("Remove project image objects failed, project_id: %d, image_uuid: %s, err: %v", req.ProjectId, imageUuid, err)
		}
	}

	return nil
}
