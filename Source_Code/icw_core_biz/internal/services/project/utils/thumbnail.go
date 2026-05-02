package utils

import (
	"context"
	"time"

	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/mysql"
	"icw_core_biz/utils"
)

// ProjectRecordToDTOWithThumbnail 将 MySQL 数据模型转换为 RPC 数据模型，并添加项目缩略图下载预签名 URL
func ProjectRecordToDTOWithThumbnail(ctx context.Context, repo *minio.Repository, record *mysql.ProjectRecord, ttl time.Duration) (*project.Project, error) {
	item := utils.ProjectRecordToDTO(record)
	if item == nil {
		return nil, nil
	}
	thumbnailURL, err := presignProjectThumbnailURL(ctx, repo, record.Id, ttl)
	if err != nil {
		return nil, err
	}
	item.ThumbnailURL = thumbnailURL
	return item, nil
}

// ProjectRecordsToListItemsDTOWithThumbnail 将 MySQL 数据模型转换为 RPC 数据模型，并添加项目缩略图下载预签名 URL
func ProjectRecordsToListItemsDTOWithThumbnail(ctx context.Context, repo *minio.Repository, records []*mysql.ProjectRecord, ttl time.Duration) ([]*project.ProjectListItem, error) {
	items := utils.ProjectRecordsToListItemsDTO(records)
	if items == nil {
		return nil, nil
	}
	for _, item := range items {
		if item == nil {
			continue
		}
		thumbnailURL, err := presignProjectThumbnailURL(ctx, repo, item.Id, ttl)
		if err != nil {
			return nil, err
		}
		item.ThumbnailURL = thumbnailURL
	}
	return items, nil
}

// presignProjectThumbnailURL 获取项目缩略图下载预签名 URL
func presignProjectThumbnailURL(ctx context.Context, repo *minio.Repository, projectId uint64, ttl time.Duration) (string, error) {
	thumbnailKey, err := minio.GenProjectThumbnailKey(projectId)
	if err != nil {
		return "", err
	}

	exists, err := repo.StatObject(ctx, thumbnailKey)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", nil
	}

	return repo.PresignGetObject(ctx, thumbnailKey, ttl)
}
