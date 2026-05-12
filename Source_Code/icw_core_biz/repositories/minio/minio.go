package minio

import (
	"bytes"
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"icw_core_biz/configs"
)

// Repository MinIO 对象存储服务
type Repository struct {
	client      *minio.Client
	adminClient *madmin.AdminClient
	bucket      string
}

// NewRepository 创建 MinIO 对象存储服务
func NewRepository(client *minio.Client, adminClient *madmin.AdminClient, bucket string) *Repository {
	return &Repository{
		client:      client,
		adminClient: adminClient,
		bucket:      bucket,
	}
}

// BucketStats MinIO Bucket 统计数据
type BucketStats struct {
	ObjectCount    uint64
	UsedBytes      uint64
	QuotaBytes     uint64
	RemainingBytes uint64
}

// ObjectMetadata MinIO 对象元数据
type ObjectMetadata struct {
	Key          string
	LastModified time.Time
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

// NewAdminClient 创建 MinIO Admin SDK Client
func NewAdminClient(cfg configs.Config) (*madmin.AdminClient, error) {
	endpoint, useSSL := normalizeEndpoint(cfg.MinIOEndpoint)
	return madmin.NewWithOptions(endpoint, &madmin.Options{
		Creds:  credentials.NewStaticV4(cfg.MinIOAccessKey, cfg.MinIOAccessSecret, ""),
		Secure: useSSL,
	})
}

// BucketStats 查询当前 Bucket 的对象数量和容量统计
func (r *Repository) BucketStats(ctx context.Context) (BucketStats, error) {
	if r == nil || r.adminClient == nil {
		return BucketStats{}, errors.New("minio admin client is nil")
	}
	usage, err := r.adminClient.DataUsageInfo(ctx)
	if err != nil {
		return BucketStats{}, err
	}
	stats := BucketStats{
		RemainingBytes: usage.TotalFreeCapacity,
	}
	if bucketUsage, ok := usage.BucketsUsage[r.bucket]; ok {
		stats.ObjectCount = bucketUsage.ObjectsCount
		stats.UsedBytes = bucketUsage.Size
	}
	if quota, err := r.adminClient.GetBucketQuota(ctx, r.bucket); err == nil {
		stats.QuotaBytes = quota.Size
		if stats.QuotaBytes == 0 {
			stats.QuotaBytes = quota.Quota
		}
		if stats.QuotaBytes > stats.UsedBytes {
			stats.RemainingBytes = stats.QuotaBytes - stats.UsedBytes
		} else if stats.QuotaBytes > 0 {
			stats.RemainingBytes = 0
		}
	}
	return stats, nil
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

// ListObjectKeysByPrefix 按对象 Key 前缀查询对象 Key 列表
func (r *Repository) ListObjectKeysByPrefix(ctx context.Context, prefix string) ([]string, error) {
	objects, err := r.ListObjectsByPrefix(ctx, prefix)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(objects))
	for _, object := range objects {
		keys = append(keys, object.Key)
	}
	return keys, nil
}

// ListObjectsByPrefix 按对象 Key 前缀查询对象元数据列表
func (r *Repository) ListObjectsByPrefix(ctx context.Context, prefix string) ([]ObjectMetadata, error) {
	objects := r.client.ListObjects(ctx, r.bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})
	metadata := make([]ObjectMetadata, 0)
	for object := range objects {
		if object.Err != nil {
			return nil, object.Err
		}
		metadata = append(metadata, ObjectMetadata{
			Key:          object.Key,
			LastModified: object.LastModified,
		})
	}
	return metadata, nil
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

// PresignPostPolicy 生成受限前缀的 POST Policy 上传授权
func (r *Repository) PresignPostPolicy(ctx context.Context, keyPrefix string, ttl time.Duration) (string, map[string]string, error) {
	policy := minio.NewPostPolicy()
	if err := policy.SetBucket(r.bucket); err != nil {
		return "", nil, err
	}
	if err := policy.SetKeyStartsWith(keyPrefix); err != nil {
		return "", nil, err
	}
	if err := policy.SetContentType("image/png"); err != nil {
		return "", nil, err
	}
	if err := policy.SetExpires(time.Now().Add(ttl)); err != nil {
		return "", nil, err
	}
	presignedURL, formData, err := r.client.PresignedPostPolicy(ctx, policy)
	if err != nil {
		return "", nil, err
	}
	return presignedURL.String(), formData, nil
}

// normalizeEndpoint 将 Endpoint 配置标准化为 MinIO SDK 需要的 <host>:<port> 格式
func normalizeEndpoint(endpoint string) (string, bool) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		return "", false
	}
	parsedURL, err := url.Parse(endpoint)
	if err != nil || parsedURL.Host == "" {
		return "", false
	}
	return parsedURL.Host, parsedURL.Scheme == "https"
}
