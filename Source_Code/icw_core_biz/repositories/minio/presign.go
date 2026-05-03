package minio

import (
	"context"
	"time"

	"icw_core_biz/repositories/redis"
)

// PresignDefaultAvatarURL 获取用户默认头像下载预签名 URL
func PresignDefaultAvatarURL(ctx context.Context, repo *Repository, redisRepo *redis.Repository, userId uint64, emailHash string, ttl time.Duration) (string, error) {
	return presignExistingObjectURL(ctx, repo, redisRepo, redis.GenDefaultAvatarPresignURLKey(userId, emailHash), GenDefaultAvatarKey(emailHash), ttl)
}

// PresignCustomAvatarURL 获取用户自定义头像下载预签名 URL
func PresignCustomAvatarURL(ctx context.Context, repo *Repository, redisRepo *redis.Repository, userId uint64, emailHash string, ttl time.Duration) (string, error) {
	return presignExistingObjectURL(ctx, repo, redisRepo, redis.GenCustomAvatarPresignURLKey(userId, emailHash), GenCustomAvatarKey(emailHash), ttl)
}

// PresignProjectThumbnailURL 获取项目缩略图下载预签名 URL
func PresignProjectThumbnailURL(ctx context.Context, repo *Repository, redisRepo *redis.Repository, userId, projectId uint64, ttl time.Duration) (string, error) {
	thumbnailKey, err := GenProjectThumbnailKey(projectId)
	if err != nil {
		return "", err
	}
	return presignExistingObjectURL(ctx, repo, redisRepo, redis.GenProjectThumbnailPresignURLKey(userId, projectId), thumbnailKey, ttl)
}

// PresignProjectImageOriginalURL 获取项目图像原图下载预签名 URL
func PresignProjectImageOriginalURL(ctx context.Context, repo *Repository, redisRepo *redis.Repository, userId, projectId uint64, imageUuid string, ttl time.Duration) (string, error) {
	originalKey, err := GenProjectImageOriginalKey(projectId, imageUuid)
	if err != nil {
		return "", err
	}
	return presignExistingObjectURL(ctx, repo, redisRepo, redis.GenProjectImageOriginalPresignURLKey(userId, projectId, imageUuid), originalKey, ttl)
}

// PresignProjectImageThumbnailURL 获取项目图像缩略图下载预签名 URL
func PresignProjectImageThumbnailURL(ctx context.Context, repo *Repository, redisRepo *redis.Repository, userId, projectId uint64, imageUuid string, ttl time.Duration) (string, error) {
	thumbnailKey, err := GenProjectImageThumbnailKey(projectId, imageUuid)
	if err != nil {
		return "", err
	}
	return presignExistingObjectURL(ctx, repo, redisRepo, redis.GenProjectImageThumbnailPresignURLKey(userId, projectId, imageUuid), thumbnailKey, ttl)
}

// presignExistingObjectURL 获取下载预签名 URL
func presignExistingObjectURL(ctx context.Context, repo *Repository, redisRepo *redis.Repository, cacheKey, objectKey string, ttl time.Duration) (string, error) {
	if redisRepo != nil {
		// 获取预签名 URL 缓存
		cachedURL, err := redisRepo.GetPresignURL(ctx, cacheKey)
		if err == nil && cachedURL != "" {
			return cachedURL, nil
		}
	}

	// 判断对象是否存在
	exists, err := repo.StatObject(ctx, objectKey)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", nil
	}

	// 生成对象下载预签名 URL
	presignURL, err := repo.PresignGetObject(ctx, objectKey, ttl)
	if err != nil || presignURL == "" {
		return "", err
	}

	if redisRepo != nil {
		// 设置预签名 URL 缓存
		_ = redisRepo.SetPresignURL(ctx, cacheKey, presignURL, ttl)
	}

	return presignURL, nil
}
