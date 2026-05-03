package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"

	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/minio"
)

// NormalizeProjectGroupName 标准化项目图像组名称
func NormalizeProjectGroupName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", rpc_err.BadRequest(rpc_err.DetailProjectGroupNameRequired, "project group name is required")
	}
	return name, nil
}

// NormalizeProjectImageMetadata 标准化项目图像元数据，并返回压缩后的 JSON 字符串
func NormalizeProjectImageMetadata(metadata string) (string, error) {
	metadata = strings.TrimSpace(metadata)
	if metadata == "" {
		return "{}", nil
	}
	var compacted bytes.Buffer
	if err := json.Compact(&compacted, []byte(metadata)); err != nil {
		return "", rpc_err.BadRequestDefault("project image metadata must be json formatted string")
	}
	return compacted.String(), nil
}

// RemoveProjectImageObjects 删除项目图像原图和缩略图对象
func RemoveProjectImageObjects(ctx context.Context, repo *minio.Repository, projectId uint64, imageUuid string) error {
	originalKey, err := minio.GenProjectImageOriginalKey(projectId, imageUuid)
	if err != nil {
		return rpc_err.BadRequestDefault(err.Error())
	}
	thumbnailKey, err := minio.GenProjectImageThumbnailKey(projectId, imageUuid)
	if err != nil {
		return rpc_err.BadRequestDefault(err.Error())
	}
	if err := repo.RemoveObject(ctx, originalKey); err != nil {
		return err
	}
	if err := repo.RemoveObject(ctx, thumbnailKey); err != nil {
		return err
	}
	return nil
}
