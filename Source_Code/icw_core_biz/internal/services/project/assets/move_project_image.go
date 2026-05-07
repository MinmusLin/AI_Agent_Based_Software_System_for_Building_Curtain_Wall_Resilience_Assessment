package assets

import (
	"context"
	"strings"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/repositories/mysql/model"
)

// MoveProjectImage 移动图像
func (s *Service) MoveProjectImage(ctx context.Context, req *bizpb.MoveProjectImageRequest) (*bizpb.MoveProjectImageResponse, error) {
	resp := &bizpb.MoveProjectImageResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.moveProjectImage(req, resp)
	})
	return resp, err
}

func (s *Service) moveProjectImage(req *bizpb.MoveProjectImageRequest, resp *bizpb.MoveProjectImageResponse) error {
	if len(req.ImageUuids) == 0 {
		return rpc_error.BadRequestDefault("image uuids are required")
	}
	if req.TargetGroupId == 0 {
		return rpc_error.BadRequestDefault("target group id is required")
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

	imageRecords, err := s.MySQL().MoveProjectImages(s.Ctx(), req.UserId, req.ProjectId, imageUuids, req.TargetGroupId)
	if err != nil {
		return err
	}
	if len(imageRecords) != len(imageUuids) {
		return rpc_error.BadRequest(rpc_error.DetailProjectNotAccessible, "project group is not accessible")
	}

	resp.Images = make([]*bizpb.ProjectImage, 0, len(imageRecords))
	for _, imageRecord := range imageRecords {
		projectImage, err := model.ProjectImageRecordToDTO(s.Ctx(), s.MinIO(), s.Redis(), imageRecord, s.Config().ProjectImageGetTTL)
		if err != nil {
			return err
		}
		resp.Images = append(resp.Images, projectImage)
	}

	return nil
}
