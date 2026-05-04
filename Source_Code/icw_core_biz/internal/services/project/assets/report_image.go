package assets

import (
	"errors"
	"strings"

	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/internal/services/project/events"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/mysql"
)

// ReportProjectImage 上报图像
func (s *Service) ReportProjectImage(req *project.ReportProjectImageRequest, resp *project.ReportProjectImageResponse) error {
	return s.CallRPC("ProjectAssetsService.ReportProjectImage", req, resp, func() error {
		return s.reportProjectImage(req, resp)
	})
}

func (s *Service) reportProjectImage(req *project.ReportProjectImageRequest, _ *project.ReportProjectImageResponse) error {
	imageUuid := strings.TrimSpace(req.ImageUuid)
	if imageUuid == "" {
		return rpc_err.BadRequestDefault("image uuid is required")
	}
	status := consts.ParseProjectImageStatus(req.Status)
	if status != consts.ProjectImageStatusUploaded && status != consts.ProjectImageStatusFailed {
		return rpc_err.BadRequestDefault("project image status is invalid")
	}

	imageRecord, err := s.MySQL().FindProjectImageByUuid(s.Ctx(), req.UserId, req.ProjectId, imageUuid)
	if err != nil {
		return err
	}
	if imageRecord == nil {
		return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project group is not accessible")
	}

	if status == consts.ProjectImageStatusUploaded {
		originalKey, err := minio.GenProjectImageOriginalKey(req.ProjectId, imageUuid)
		if err != nil {
			return rpc_err.BadRequestDefault(err.Error())
		}
		thumbnailKey, err := minio.GenProjectImageThumbnailKey(req.ProjectId, imageUuid)
		if err != nil {
			return rpc_err.BadRequestDefault(err.Error())
		}
		originalExists, err := s.MinIO().StatObject(s.Ctx(), originalKey)
		if err != nil {
			return err
		}
		thumbnailExists, err := s.MinIO().StatObject(s.Ctx(), thumbnailKey)
		if err != nil {
			return err
		}
		if !originalExists || !thumbnailExists {
			return rpc_err.BadRequestDefault("project image objects are not uploaded")
		}
	}

	imageRecord, err = s.MySQL().UpdateProjectImageStatus(s.Ctx(), req.UserId, req.ProjectId, imageUuid, status)
	if errors.Is(err, mysql.ErrProjectImageStatusTransitionInvalid) {
		return rpc_err.BadRequestDefault("project image status transition is invalid")
	}
	if err != nil {
		return err
	}
	if imageRecord == nil {
		return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project group is not accessible")
	}

	image, err := mysql.ProjectImageRecordToDTO(s.Ctx(), s.MinIO(), s.Redis(), imageRecord, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}

	// 发布项目图像状态变化事件
	events.PublishProjectImageStatusChangedEvent(s.Ctx(), s.RocketMQ(), req.UserId, req.ProjectId, image)

	return nil
}
