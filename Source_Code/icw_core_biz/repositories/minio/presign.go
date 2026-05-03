package minio

import (
	"context"
	"time"
)

// PresignProjectThumbnailURL 获取项目缩略图下载预签名 URL
func PresignProjectThumbnailURL(ctx context.Context, repo *Repository, projectId uint64, ttl time.Duration) (string, error) {
	thumbnailKey, err := GenProjectThumbnailKey(projectId)
	if err != nil {
		return "", err
	}
	return presignExistingObjectURL(ctx, repo, thumbnailKey, ttl)
}

// PresignProjectImageOriginalURL 获取项目原图下载预签名 URL
func PresignProjectImageOriginalURL(ctx context.Context, repo *Repository, projectId uint64, imageUuid string, ttl time.Duration) (string, error) {
	originalKey, err := GenProjectImageOriginalKey(projectId, imageUuid)
	if err != nil {
		return "", err
	}
	return presignExistingObjectURL(ctx, repo, originalKey, ttl)
}

// PresignDefaultAvatarURL 获取用户默认头像下载预签名 URL
func PresignDefaultAvatarURL(ctx context.Context, repo *Repository, emailHash string, ttl time.Duration) (string, error) {
	return presignExistingObjectURL(ctx, repo, GenDefaultAvatarKey(emailHash), ttl)
}

// PresignCustomAvatarURL 获取用户自定义头像下载预签名 URL
func PresignCustomAvatarURL(ctx context.Context, repo *Repository, emailHash string, ttl time.Duration) (string, error) {
	return presignExistingObjectURL(ctx, repo, GenCustomAvatarKey(emailHash), ttl)
}

// presignExistingObjectURL 获取下载预签名 URL
func presignExistingObjectURL(ctx context.Context, repo *Repository, key string, ttl time.Duration) (string, error) {
	exists, err := repo.StatObject(ctx, key)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", nil
	}
	return repo.PresignGetObject(ctx, key, ttl)
}
