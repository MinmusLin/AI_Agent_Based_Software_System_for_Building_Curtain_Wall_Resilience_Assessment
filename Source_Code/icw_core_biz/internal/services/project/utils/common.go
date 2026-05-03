package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/minio"
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
	if err := validateStringMaxLength(fields.Name, consts.ProjectNameMaxLength, rpc_err.DetailProjectNameTooLong, "project name is too long"); err != nil {
		return nil, err
	}
	if err := validateStringMaxLength(fields.BuildingName, consts.ProjectBuildingNameMaxLength, rpc_err.DetailProjectBuildingNameTooLong, "project building name is too long"); err != nil {
		return nil, err
	}
	if err := validateStringMaxLength(fields.BuildingLocation, consts.ProjectBuildingLocationMaxLength, rpc_err.DetailProjectBuildingLocationTooLong, "project building location is too long"); err != nil {
		return nil, err
	}
	if err := validateStringMaxLength(fields.BuildingDescription, consts.ProjectBuildingDescriptionMaxLength, rpc_err.DetailProjectBuildingDescriptionTooLong, "project building description is too long"); err != nil {
		return nil, err
	}
	if err := validateStringMaxLength(fields.KnownIssues, consts.ProjectKnownIssuesMaxLength, rpc_err.DetailProjectKnownIssuesTooLong, "project known issues is too long"); err != nil {
		return nil, err
	}
	if err := validateStringMaxLength(fields.AssessmentGoal, consts.ProjectAssessmentGoalMaxLength, rpc_err.DetailProjectAssessmentGoalTooLong, "project assessment goal is too long"); err != nil {
		return nil, err
	}
	return fields, nil
}

// NormalizeProjectGroupName 标准化并校验项目图像组名称
func NormalizeProjectGroupName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", rpc_err.BadRequest(rpc_err.DetailProjectGroupNameRequired, "project group name is required")
	}
	if err := validateStringMaxLength(name, consts.ProjectGroupNameMaxLength, rpc_err.DetailProjectGroupNameTooLong, "project group name is too long"); err != nil {
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
		return "", rpc_err.BadRequestDefault("project image metadata must be json formatted string")
	}
	return compacted.String(), nil
}

// validateStringMaxLength 校验字符串最大字符数
func validateStringMaxLength(value string, maxLength int, detailCode rpc_err.DetailCode, message string) error {
	if utf8.RuneCountInString(value) > maxLength {
		return rpc_err.BadRequest(detailCode, message)
	}
	return nil
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
	var removeErr error
	if err := repo.RemoveObject(ctx, originalKey); err != nil {
		removeErr = errors.Join(removeErr, fmt.Errorf("remove project image original object failed: %w", err))
	}
	if err := repo.RemoveObject(ctx, thumbnailKey); err != nil {
		removeErr = errors.Join(removeErr, fmt.Errorf("remove project image thumbnail object failed: %w", err))
	}
	return removeErr
}
