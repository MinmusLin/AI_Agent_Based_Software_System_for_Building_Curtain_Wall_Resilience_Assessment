package assets

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"

	"icw_common/gen/core/biz"
	"icw_common/rpc_err"
	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/mysql"
)

// UploadProjectImage 上传图像
func (s *Service) UploadProjectImage(ctx context.Context, req *bizpb.UploadProjectImageRequest) (*bizpb.UploadProjectImageResponse, error) {
	resp := &bizpb.UploadProjectImageResponse{}
	err := s.CallRPC(ctx, req, resp, func() error {
		return s.uploadProjectImage(req, resp)
	})
	return resp, err
}

func (s *Service) uploadProjectImage(req *bizpb.UploadProjectImageRequest, resp *bizpb.UploadProjectImageResponse) error {
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

	uploadItems := make([]*bizpb.UploadProjectImageItem, 0, len(req.Images))
	for _, image := range req.Images {
		if image == nil {
			continue
		}

		contentType := strings.TrimSpace(image.ContentType)
		if contentType != consts.ProjectImageContentType {
			return rpc_err.BadRequest(rpc_err.DetailInvalidImageContentType, "image content type must be image/png")
		}
		fileName := strings.TrimSpace(image.FileName)
		if fileName == "" {
			return rpc_err.BadRequestDefault("project image file name is required")
		}
		if utf8.RuneCountInString(fileName) > consts.ProjectImageFileNameMaxLength {
			return rpc_err.BadRequest(rpc_err.DetailProjectImageFileNameTooLong, "project image file name is too long")
		}
		if image.SizeBytes == 0 || image.Width == 0 || image.Height == 0 {
			return rpc_err.BadRequest(rpc_err.DetailProjectImageFormatInvalid, "project image format is invalid")
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

		uploadItems = append(uploadItems, &bizpb.UploadProjectImageItem{
			FileName:     fileName,
			ContentType:  contentType,
			SizeBytes:    image.SizeBytes,
			Width:        image.Width,
			Height:       image.Height,
			Metadata:     metadata,
			ImageUuid:    imageUuid,
			OriginalKey:  originalKey,
			ThumbnailKey: thumbnailKey,
		})
	}

	if len(uploadItems) == 0 {
		return rpc_err.BadRequestDefault("project images are required")
	}

	createRecords := make([]*mysql.ProjectImageCreateRecord, 0, len(uploadItems))
	for _, item := range uploadItems {
		// 生成项目图像原图上传预签名 URL
		originalUploadURL, err := s.MinIO().PresignPutObject(s.Ctx(), item.OriginalKey, s.Config().ProjectImageUploadTTL)
		if err != nil {
			return err
		}
		// 生成项目图像缩略图上传预签名 URL
		thumbnailUploadURL, err := s.MinIO().PresignPutObject(s.Ctx(), item.ThumbnailKey, s.Config().ProjectImageUploadTTL)
		if err != nil {
			return err
		}

		item.OriginalUploadUrl = originalUploadURL
		item.ThumbnailUploadUrl = thumbnailUploadURL

		createRecords = append(createRecords, &mysql.ProjectImageCreateRecord{
			ImageUuid:   item.ImageUuid,
			FileName:    item.FileName,
			ContentType: item.ContentType,
			SizeBytes:   item.SizeBytes,
			Width:       item.Width,
			Height:      item.Height,
			Metadata:    item.Metadata,
		})
	}

	imageRecords, err := s.MySQL().CreateProjectImages(s.Ctx(), req.UserId, req.ProjectId, req.GroupId, createRecords)
	if err != nil {
		return err
	}
	if len(imageRecords) != len(uploadItems) {
		return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project group is not accessible")
	}

	resp.Images = make([]*bizpb.UploadProjectImageResult, 0, len(uploadItems))
	for index, item := range uploadItems {
		projectImage, err := mysql.ProjectImageRecordToDTO(s.Ctx(), s.MinIO(), s.Redis(), imageRecords[index], s.Config().ProjectImageGetTTL)
		if err != nil {
			return err
		}

		resp.Images = append(resp.Images, &bizpb.UploadProjectImageResult{
			Image:              projectImage,
			OriginalUploadUrl:  item.OriginalUploadUrl,
			ThumbnailUploadUrl: item.ThumbnailUploadUrl,
		})
	}

	return nil
}
