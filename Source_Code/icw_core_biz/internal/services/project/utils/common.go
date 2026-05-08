package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/redis"
)

// ProjectProfileFields 项目基础信息字段
type ProjectProfileFields struct {
	Name                string
	BuildingName        string
	BuildingLocation    string
	BuildingDescription string
	KnownIssues         string
	AssessmentGoal      string
}

// NormalizeProjectProfileFields 标准化并校验项目基础信息字段
func NormalizeProjectProfileFields(name, buildingName, buildingLocation, buildingDescription, knownIssues, assessmentGoal string) (*ProjectProfileFields, error) {
	fields := &ProjectProfileFields{
		Name:                strings.TrimSpace(name),
		BuildingName:        strings.TrimSpace(buildingName),
		BuildingLocation:    strings.TrimSpace(buildingLocation),
		BuildingDescription: strings.TrimSpace(buildingDescription),
		KnownIssues:         strings.TrimSpace(knownIssues),
		AssessmentGoal:      strings.TrimSpace(assessmentGoal),
	}
	if err := validateStringMaxLength(fields.Name, consts.ProjectNameMaxLength, rpc_error.DetailProjectNameTooLong, "project name is too long"); err != nil {
		return nil, err
	}
	if err := validateStringMaxLength(fields.BuildingName, consts.ProjectBuildingNameMaxLength, rpc_error.DetailProjectBuildingNameTooLong, "project building name is too long"); err != nil {
		return nil, err
	}
	if err := validateStringMaxLength(fields.BuildingLocation, consts.ProjectBuildingLocationMaxLength, rpc_error.DetailProjectBuildingLocationTooLong, "project building location is too long"); err != nil {
		return nil, err
	}
	if err := validateStringMaxLength(fields.BuildingDescription, consts.ProjectBuildingDescriptionMaxLength, rpc_error.DetailProjectBuildingDescriptionTooLong, "project building description is too long"); err != nil {
		return nil, err
	}
	if err := validateStringMaxLength(fields.KnownIssues, consts.ProjectKnownIssuesMaxLength, rpc_error.DetailProjectKnownIssuesTooLong, "project known issues is too long"); err != nil {
		return nil, err
	}
	if err := validateStringMaxLength(fields.AssessmentGoal, consts.ProjectAssessmentGoalMaxLength, rpc_error.DetailProjectAssessmentGoalTooLong, "project assessment goal is too long"); err != nil {
		return nil, err
	}
	return fields, nil
}

// NormalizeProjectGroupName 标准化并校验项目图像组名称
func NormalizeProjectGroupName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", rpc_error.BadRequest(rpc_error.DetailProjectGroupNameRequired, "project group name is required")
	}
	if err := validateStringMaxLength(name, consts.ProjectGroupNameMaxLength, rpc_error.DetailProjectGroupNameTooLong, "project group name is too long"); err != nil {
		return "", err
	}
	return name, nil
}

// NormalizeProjectImageMetadata 标准化并校验项目图像元数据，并返回压缩后的 JSON 字符串
func NormalizeProjectImageMetadata(metadata string) (string, error) {
	metadata = strings.TrimSpace(metadata)
	if metadata == "" {
		return "{}", nil
	}
	var compacted bytes.Buffer
	if err := json.Compact(&compacted, []byte(metadata)); err != nil {
		return "", rpc_error.BadRequestDefault("project image metadata must be json formatted string")
	}
	return compacted.String(), nil
}

// RemoveProjectImageObjects 删除项目图像原图和缩略图对象
func RemoveProjectImageObjects(ctx context.Context, minioRepo *minio.Repository, redisRepo *redis.Repository, userId, projectId uint64, imageUuid string) error {
	originalKey, err := minio.GenProjectImageOriginalKey(projectId, imageUuid)
	if err != nil {
		return rpc_error.BadRequestDefault(err.Error())
	}
	thumbnailKey, err := minio.GenProjectImageThumbnailKey(projectId, imageUuid)
	if err != nil {
		return rpc_error.BadRequestDefault(err.Error())
	}

	if redisRepo != nil {
		// 清除预签名 URL 缓存
		_ = redisRepo.ClearPresignURL(ctx, originalKey)
		_ = redisRepo.ClearPresignURL(ctx, thumbnailKey)
	}

	var removeErr error
	if err := minioRepo.RemoveObject(ctx, originalKey); err != nil {
		removeErr = errors.Join(removeErr, fmt.Errorf("remove project image original object failed: %w", err))
	}
	if err := minioRepo.RemoveObject(ctx, thumbnailKey); err != nil {
		removeErr = errors.Join(removeErr, fmt.Errorf("remove project image thumbnail object failed: %w", err))
	}

	return removeErr
}

// validateStringMaxLength 校验字符串最大字符数
func validateStringMaxLength(value string, maxLength int, detailCode rpc_error.DetailCode, message string) error {
	if utf8.RuneCountInString(value) > maxLength {
		return rpc_error.BadRequest(detailCode, message)
	}
	return nil
}

// ArtifactSha256MapJSON 将图像检测推理产物上传结果转换为 Sha256 Map JSON
func ArtifactSha256MapJSON(artifacts []*bizpb.ReasoningArtifactUploadResult) (string, error) {
	artifactSha256Map := make(map[string]string)
	for _, artifact := range artifacts {
		if artifact == nil || !artifact.Uploaded {
			continue
		}
		name := strings.TrimSpace(artifact.Name)
		sha256 := strings.TrimSpace(artifact.Sha256)
		if name == "" || sha256 == "" {
			continue
		}
		artifactSha256Map[name] = sha256
	}
	bytes, err := json.Marshal(artifactSha256Map)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
