package minio

import (
	"bytes"
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"icw_core_biz/configs"
)

// Repository MinIO 对象存储服务
type Repository struct {
	client *minio.Client
	bucket string
}

func NewRepository(cfg configs.Config) (*Repository, error) {
	endpoint, useSSL := normalizeEndpoint(cfg.MinIOEndpoint)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOAccessSecret, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return &Repository{
		client: client,
		bucket: cfg.MinIOBucket,
	}, nil
}

// Ping 检查 Bucket 是否可访问
func (r *Repository) Ping(ctx context.Context) bool {
	exists, err := r.client.BucketExists(ctx, r.bucket)
	return exists && err == nil
}

// StatObject 判断对象是否存在
func (r *Repository) StatObject(ctx context.Context, key string) (bool, error) {
	if _, err := r.client.StatObject(ctx, r.bucket, key, minio.StatObjectOptions{}); err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// PutObject 上传对象
func (r *Repository) PutObject(ctx context.Context, key string, contentType string, data []byte) error {
	_, err := r.client.PutObject(ctx, r.bucket, key, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// RemoveObject 删除对象
func (r *Repository) RemoveObject(ctx context.Context, key string) error {
	return r.client.RemoveObject(ctx, r.bucket, key, minio.RemoveObjectOptions{})
}

// PresignGetObject 生成对象下载预签名 URL
func (r *Repository) PresignGetObject(ctx context.Context, key string, ttl time.Duration) (string, error) {
	presignedURL, err := r.client.PresignedGetObject(ctx, r.bucket, key, ttl, nil)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// PresignPutObject 生成对象上传预签名 URL
func (r *Repository) PresignPutObject(ctx context.Context, key string, ttl time.Duration) (string, error) {
	presignedURL, err := r.client.PresignedPutObject(ctx, r.bucket, key, ttl)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}
