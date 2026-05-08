package assets

import (
	"context"
	"strings"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/repositories/minio"
)

// GetProjectImageOriginal 获取原图
func (s *Service) GetProjectImageOriginal(ctx context.Context, req *bizpb.GetProjectImageOriginalRequest) (*bizpb.GetProjectImageOriginalResponse, error) {
	resp := &bizpb.GetProjectImageOriginalResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.getProjectImageOriginal(req, resp)
	})
	return resp, err
}

func (s *Service) getProjectImageOriginal(req *bizpb.GetProjectImageOriginalRequest, resp *bizpb.GetProjectImageOriginalResponse) error {
	imageUuid := strings.TrimSpace(req.ImageUuid)
	if imageUuid == "" {
		return rpc_error.BadRequestDefault("image uuid is required")
	}

	imageRecord, err := s.MySQL().FindProjectImageByUuid(s.Ctx(), req.UserId, req.ProjectId, imageUuid)
	if err != nil {
		return err
	}
	if imageRecord == nil {
		return rpc_error.BadRequest(rpc_error.DetailProjectNotAccessible, "project group is not accessible")
	}
	if imageRecord.Status != bizpb.ProjectImageStatus_Uploaded {
		return rpc_error.BadRequestDefault("project image status is not uploaded")
	}

	// 获取项目图像原图下载预签名 URL
	originalURL, err := minio.PresignProjectImageOriginalURL(s.Ctx(), s.MinIO(), s.Redis(), req.ProjectId, imageUuid, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}
	if originalURL == "" {
		return rpc_error.BadRequest(rpc_error.DetailProjectImageExpired, "project image original object not found")
	}
	resp.OriginalUrl = originalURL

	return nil
}
