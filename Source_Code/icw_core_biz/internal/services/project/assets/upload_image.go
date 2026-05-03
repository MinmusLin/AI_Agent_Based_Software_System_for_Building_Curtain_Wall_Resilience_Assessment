package assets

import (
	"strings"

	"github.com/google/uuid"

	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/mysql"
)

// UploadProjectImage 上传图像
func (s *Service) UploadProjectImage(req *project.UploadProjectImageRequest, resp *project.UploadProjectImageResponse) error {
	return s.CallRPC("ProjectAssetsService.UploadProjectImage", req, resp, func() error {
		return s.uploadProjectImage(req, resp)
	})
}

func (s *Service) uploadProjectImage(req *project.UploadProjectImageRequest, resp *project.UploadProjectImageResponse) error {
	if len(req.Images) == 0 {
		return rpc_err.BadRequestDefault("project images are required")
	}

	// 校验图像组访问权限
	groupRecord, err := s.MySQL().FindProjectGroupById(s.Ctx(), req.UserId, req.ProjectId, req.GroupId)
	if err != nil {
		return err
	}
	if groupRecord == nil {
		return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project group is not accessible")
	}

	resp.Images = make([]*project.UploadProjectImageResult, 0, len(req.Images))
	for _, image := range req.Images {
		if image == nil {
			continue
		}

		contentType := image.ContentType
		if image.ContentType != consts.ProjectImageContentType {
			return rpc_err.BadRequest(rpc_err.DetailInvalidImageContentType, "image content type must be image/png")
		}
		fileName := strings.TrimSpace(image.FileName)
		if fileName == "" {
			return rpc_err.BadRequestDefault("project image file name is required")
		}

		// 标准化项目图像元数据，并返回压缩后的 JSON 字符串
		metadata, err := utils.NormalizeProjectImageMetadata(image.Metadata)
		if err != nil {
			return err
		}

		imageUuid := uuid.NewString()

		// 生成项目图像原图对象 Key
		originalKey, err := minio.GenProjectImageOriginalKey(req.ProjectId, imageUuid)
		if err != nil {
			return rpc_err.BadRequestDefault(err.Error())
		}
		// 生成项目图像缩略图对象 Key
		thumbnailKey, err := minio.GenProjectImageThumbnailKey(req.ProjectId, imageUuid)
		if err != nil {
			return rpc_err.BadRequestDefault(err.Error())
		}

		// 生成对象上传预签名 URL
		originalUploadURL, err := s.MinIO().PresignPutObject(s.Ctx(), originalKey, s.Config().ProjectImageUploadTTL)
		if err != nil {
			return err
		}
		// 生成对象下载预签名 URL
		thumbnailUploadURL, err := s.MinIO().PresignPutObject(s.Ctx(), thumbnailKey, s.Config().ProjectImageUploadTTL)
		if err != nil {
			return err
		}

		imageRecord, err := s.MySQL().CreateProjectImage(s.Ctx(), req.UserId, req.ProjectId, req.GroupId, imageUuid, fileName, contentType, image.SizeBytes, image.Width, image.Height, metadata)
		if err != nil {
			return err
		}
		if imageRecord == nil {
			return rpc_err.BadRequestDefault("project image is not accessible")
		}

		projectImage, err := mysql.ProjectImageRecordToDTO(s.Ctx(), s.MinIO(), imageRecord, s.Config().ProjectImageGetTTL)
		if err != nil {
			return err
		}

		resp.Images = append(resp.Images, &project.UploadProjectImageResult{
			Image:              projectImage,
			OriginalUploadURL:  originalUploadURL,
			ThumbnailUploadURL: thumbnailUploadURL,
		})
	}

	if len(resp.Images) == 0 {
		return rpc_err.BadRequestDefault("project images are required")
	}

	return nil
}
