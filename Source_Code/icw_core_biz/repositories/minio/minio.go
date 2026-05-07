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

// NewRepository 创建 MinIO 对象存储服务
func NewRepository(client *minio.Client, bucket string) *Repository {
	return &Repository{
		client: client,
		bucket: bucket,
	}
}

// NewClient 创建 MinIO SDK Client
func NewClient(cfg configs.Config) (*minio.Client, error) {
	endpoint, useSSL := normalizeEndpoint(cfg.MinIOEndpoint)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOAccessSecret, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
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
func (r *Repository) PutObject(ctx context.Context, key, contentType string, data []byte) error {
	_, err := r.client.PutObject(ctx, r.bucket, key, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// RemoveObject 删除对象
func (r *Repository) RemoveObject(ctx context.Context, key string) error {
	return r.client.RemoveObject(ctx, r.bucket, key, minio.RemoveObjectOptions{})
}

// RemoveObjectsByPrefix 按对象 Key 前缀删除对象
func (r *Repository) RemoveObjectsByPrefix(ctx context.Context, prefix string) error {
	objects := r.client.ListObjects(ctx, r.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})
	for object := range objects {
		if object.Err != nil {
			return object.Err
		}
		if err := r.RemoveObject(ctx, object.Key); err != nil {
			return err
		}
	}
	return nil
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
