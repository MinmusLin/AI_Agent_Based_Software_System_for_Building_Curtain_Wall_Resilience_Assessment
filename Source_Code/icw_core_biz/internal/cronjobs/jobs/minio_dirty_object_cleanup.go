package jobs

import (
	"errors"
	"fmt"
	"strings"

	"icw_common/utils"

	"icw_core_biz/internal/cronjobs/common"
	"icw_core_biz/repositories/mysql/model"
)

// MinIODirtyObjectCleanupJob MinIO 脏对象清理任务
type MinIODirtyObjectCleanupJob struct {
	*common.BaseCronJob
}

// NewMinIODirtyObjectCleanupJob 创建 MinIO 脏对象清理任务
func NewMinIODirtyObjectCleanupJob(baseJob *common.BaseCronJob) common.CronJob {
	return &MinIODirtyObjectCleanupJob{
		BaseCronJob: baseJob,
	}
}

// MinIODirtyObjectCleanupJobResult MinIO 脏对象清理任务执行结果
type MinIODirtyObjectCleanupJobResult struct {
	CandidatePrefixCount     int `json:"candidate_prefix_count"`
	DeletedPrefixCount       int `json:"deleted_prefix_count"`
	DeletedObjectCount       int `json:"deleted_object_count"`
	SkippedInvalidObjectKeys int `json:"skipped_invalid_object_keys"`
	FailedPrefixCount        int `json:"failed_prefix_count"`
}

// Start 执行 MinIO 脏对象清理任务
func (j *MinIODirtyObjectCleanupJob) Start() (interface{}, error) {
	imageRefs, err := j.MySQL().ListProjectImageObjectReferences(j.Ctx())
	if err != nil {
		return nil, err
	}
	detectionRefs, err := j.MySQL().ListProjectDetectionObjectReferences(j.Ctx())
	if err != nil {
		return nil, err
	}
	reportRefs, err := j.MySQL().ListProjectReportObjectReferences(j.Ctx())
	if err != nil {
		return nil, err
	}

	prefixes := make(map[string]struct{})
	if err := j.addFailedObjectPrefixes(prefixes); err != nil {
		return nil, err
	}

	skippedInvalidObjectKeys, err := j.addOrphanObjectPrefixes(prefixes, imageRefs, detectionRefs, reportRefs)
	if err != nil {
		return nil, err
	}

	result := &MinIODirtyObjectCleanupJobResult{
		CandidatePrefixCount:     len(prefixes),
		SkippedInvalidObjectKeys: skippedInvalidObjectKeys,
	}

	var removeErr error
	for prefix := range prefixes {
		deletedObjectCount, err := j.removeObjectsByPrefix(prefix)
		if err != nil {
			result.FailedPrefixCount++
			removeErr = errors.Join(removeErr, fmt.Errorf("remove minio dirty object prefix failed, prefix: %s, err: %w", prefix, err))
			continue
		}
		result.DeletedPrefixCount++
		result.DeletedObjectCount += deletedObjectCount
	}

	return result, removeErr
}

func (j *MinIODirtyObjectCleanupJob) addFailedObjectPrefixes(prefixes map[string]struct{}) error {
	failedImages, err := j.MySQL().ListFailedProjectImageObjectReferences(j.Ctx())
	if err != nil {
		return err
	}
	for _, item := range failedImages {
		prefix, err := projectImageAssetPrefix(item.ProjectId, item.ImageUuid)
		if err != nil {
			return err
		}
		addCleanupPrefix(prefixes, prefix)
	}

	failedDetections, err := j.MySQL().ListFailedProjectDetectionObjectReferences(j.Ctx())
	if err != nil {
		return err
	}
	for _, item := range failedDetections {
		prefix, err := projectDetectionTaskPrefix(item.ProjectId, item.ImageUuid, item.TaskCode)
		if err != nil {
			return err
		}
		addCleanupPrefix(prefixes, prefix)
	}

	failedReports, err := j.MySQL().ListFailedProjectReportObjectReferences(j.Ctx())
	if err != nil {
		return err
	}
	for _, item := range failedReports {
		prefix, err := projectReportPrefix(item.ProjectId)
		if err != nil {
			return err
		}
		addCleanupPrefix(prefixes, prefix)
	}

	return nil
}

func (j *MinIODirtyObjectCleanupJob) addOrphanObjectPrefixes(
	prefixes map[string]struct{},
	imageRefs []model.ProjectImageObjectReference,
	detectionRefs []model.ProjectDetectionObjectReference,
	reportRefs []model.ProjectReportObjectReference,
) (int, error) {
	imageRefSet := projectImageObjectReferenceSet(imageRefs)
	detectionRefSet := projectDetectionObjectReferenceSet(detectionRefs)
	reportRefSet := projectReportObjectReferenceSet(reportRefs)

	objectKeys, err := j.MinIO().ListObjectKeysByPrefix(j.Ctx(), "projects/")
	if err != nil {
		return 0, err
	}

	skippedInvalidObjectKeys := 0
	for _, objectKey := range objectKeys {
		segments := strings.Split(objectKey, "/")
		if len(segments) < 3 || segments[0] != strings.TrimSuffix("projects/", "/") {
			continue
		}

		projectId, err := utils.Decode(segments[1])
		if err != nil {
			skippedInvalidObjectKeys++
			continue
		}

		switch segments[2] {
		case "assets":
			if len(segments) < 5 {
				continue
			}
			imageUuid := strings.TrimSpace(segments[3])
			if _, ok := imageRefSet[projectImageObjectReferenceKey(projectId, imageUuid)]; !ok {
				addCleanupPrefix(prefixes, projectImageAssetPrefixFromCode(segments[1], imageUuid))
			}
		case "detections":
			if len(segments) < 5 {
				continue
			}
			imageUuid := strings.TrimSpace(segments[3])
			if _, ok := imageRefSet[projectImageObjectReferenceKey(projectId, imageUuid)]; !ok {
				addCleanupPrefix(prefixes, projectDetectionImagePrefixFromCode(segments[1], imageUuid))
				continue
			}
			if len(segments) < 6 {
				continue
			}
			taskCode := strings.TrimSpace(segments[4])
			if _, ok := detectionRefSet[projectDetectionObjectReferenceKey(projectId, imageUuid, taskCode)]; !ok {
				addCleanupPrefix(prefixes, projectDetectionTaskPrefixFromCode(segments[1], imageUuid, taskCode))
			}
		case "reports":
			if len(segments) < 4 {
				continue
			}
			if _, ok := reportRefSet[projectReportObjectReferenceKey(projectId)]; !ok {
				addCleanupPrefix(prefixes, projectReportPrefixFromCode(segments[1]))
			}
		}
	}

	return skippedInvalidObjectKeys, nil
}

func (j *MinIODirtyObjectCleanupJob) removeObjectsByPrefix(prefix string) (int, error) {
	objectKeys, err := j.MinIO().ListObjectKeysByPrefix(j.Ctx(), prefix)
	if err != nil {
		return 0, err
	}

	var removeErr error
	deletedObjectCount := 0
	for _, objectKey := range objectKeys {
		if j.Redis() != nil {
			_ = j.Redis().ClearPresignURL(j.Ctx(), objectKey)
		}
		if err := j.MinIO().RemoveObject(j.Ctx(), objectKey); err != nil {
			removeErr = errors.Join(removeErr, err)
			continue
		}
		deletedObjectCount++
	}

	return deletedObjectCount, removeErr
}

func projectImageObjectReferenceSet(items []model.ProjectImageObjectReference) map[string]struct{} {
	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		set[projectImageObjectReferenceKey(item.ProjectId, item.ImageUuid)] = struct{}{}
	}
	return set
}

func projectDetectionObjectReferenceSet(items []model.ProjectDetectionObjectReference) map[string]struct{} {
	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		set[projectDetectionObjectReferenceKey(item.ProjectId, item.ImageUuid, item.TaskCode)] = struct{}{}
	}
	return set
}

func projectReportObjectReferenceSet(items []model.ProjectReportObjectReference) map[string]struct{} {
	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		set[projectReportObjectReferenceKey(item.ProjectId)] = struct{}{}
	}
	return set
}

func projectImageObjectReferenceKey(projectId uint64, imageUuid string) string {
	return fmt.Sprintf("%d:%s", projectId, strings.TrimSpace(imageUuid))
}

func projectDetectionObjectReferenceKey(projectId uint64, imageUuid, taskCode string) string {
	return fmt.Sprintf("%d:%s:%s", projectId, strings.TrimSpace(imageUuid), strings.TrimSpace(taskCode))
}

func projectReportObjectReferenceKey(projectId uint64) string {
	return fmt.Sprintf("%d", projectId)
}

func projectImageAssetPrefix(projectId uint64, imageUuid string) (string, error) {
	projectCode, err := projectCode(projectId)
	if err != nil {
		return "", err
	}
	return projectImageAssetPrefixFromCode(projectCode, imageUuid), nil
}

func projectDetectionTaskPrefix(projectId uint64, imageUuid, taskCode string) (string, error) {
	projectCode, err := projectCode(projectId)
	if err != nil {
		return "", err
	}
	return projectDetectionTaskPrefixFromCode(projectCode, imageUuid, taskCode), nil
}

func projectReportPrefix(projectId uint64) (string, error) {
	projectCode, err := projectCode(projectId)
	if err != nil {
		return "", err
	}
	return projectReportPrefixFromCode(projectCode), nil
}

func projectCode(projectId uint64) (string, error) {
	projectCode := utils.Encode(projectId)
	if projectCode == "" {
		return "", fmt.Errorf("project id is invalid: %d", projectId)
	}
	return projectCode, nil
}

func projectImageAssetPrefixFromCode(projectCode, imageUuid string) string {
	return fmt.Sprintf("projects/%s/assets/%s/", strings.TrimSpace(projectCode), strings.TrimSpace(imageUuid))
}

func projectDetectionImagePrefixFromCode(projectCode, imageUuid string) string {
	return fmt.Sprintf("projects/%s/detections/%s/", strings.TrimSpace(projectCode), strings.TrimSpace(imageUuid))
}

func projectDetectionTaskPrefixFromCode(projectCode, imageUuid, taskCode string) string {
	return fmt.Sprintf("projects/%s/detections/%s/%s/", strings.TrimSpace(projectCode), strings.TrimSpace(imageUuid), strings.TrimSpace(taskCode))
}

func projectReportPrefixFromCode(projectCode string) string {
	return fmt.Sprintf("projects/%s/reports/", strings.TrimSpace(projectCode))
}

func addCleanupPrefix(prefixes map[string]struct{}, prefix string) {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return
	}
	for existingPrefix := range prefixes {
		if strings.HasPrefix(prefix, existingPrefix) {
			return
		}
		if strings.HasPrefix(existingPrefix, prefix) {
			delete(prefixes, existingPrefix)
		}
	}
	prefixes[prefix] = struct{}{}
}
