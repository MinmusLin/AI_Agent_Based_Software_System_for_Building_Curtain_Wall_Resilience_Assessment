package jobs

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"icw_common/utils"

	"icw_core_biz/internal/cronjobs/common"
	"icw_core_biz/repositories/mysql/model"
)

const (
	// MinIODirtyObjectCleanupGracePeriod MinIO 脏对象清理保护窗口
	MinIODirtyObjectCleanupGracePeriod = 30 * time.Minute
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
	DeletedFolderPaths []string `json:"deleted_folder_paths"`
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
	if err := j.addOrphanObjectPrefixes(prefixes, imageRefs, detectionRefs, reportRefs); err != nil {
		return nil, err
	}

	result := &MinIODirtyObjectCleanupJobResult{
		DeletedFolderPaths: make([]string, 0, len(prefixes)),
	}

	var removeErr error
	for prefix := range prefixes {
		removed, err := j.removeObjectsByPrefix(prefix)
		if err != nil {
			removeErr = errors.Join(removeErr, fmt.Errorf("remove minio dirty object folder failed, folder_path: %s, err: %w", prefix, err))
			continue
		}
		if removed {
			result.DeletedFolderPaths = append(result.DeletedFolderPaths, prefix)
		}
	}
	sort.Strings(result.DeletedFolderPaths)

	return result, removeErr
}

// addFailedObjectPrefixes 添加数据库中明确失败状态对应的待清理文件夹路径
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

// addOrphanObjectPrefixes 添加 MinIO 中存在但数据库已无有效引用的待清理文件夹路径
func (j *MinIODirtyObjectCleanupJob) addOrphanObjectPrefixes(prefixes map[string]struct{}, imageRefs []model.ProjectImageObjectReference, detectionRefs []model.ProjectDetectionObjectReference, reportRefs []model.ProjectReportObjectReference) error {
	imageRefSet := projectImageObjectReferenceSet(imageRefs)
	detectionRefSet := projectDetectionObjectReferenceSet(detectionRefs)
	reportRefSet := projectReportObjectReferenceSet(reportRefs)

	objectKeys, err := j.MinIO().ListObjectKeysByPrefix(j.Ctx(), "projects/")
	if err != nil {
		return err
	}

	for _, objectKey := range objectKeys {
		segments := strings.Split(objectKey, "/")
		if len(segments) < 3 || segments[0] != strings.TrimSuffix("projects/", "/") {
			continue
		}
		projectId, err := utils.Decode(segments[1])
		if err != nil {
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

	return nil
}

// removeObjectsByPrefix 删除指定文件夹路径下超过保护窗口的全部对象并清理预签名 URL 缓存
func (j *MinIODirtyObjectCleanupJob) removeObjectsByPrefix(prefix string) (bool, error) {
	objects, err := j.MinIO().ListObjectsByPrefix(j.Ctx(), prefix)
	if err != nil {
		return false, err
	}
	if len(objects) == 0 {
		return false, nil
	}

	cutoffTime := time.Now().Add(-MinIODirtyObjectCleanupGracePeriod)
	for _, object := range objects {
		if object.LastModified.IsZero() || object.LastModified.After(cutoffTime) {
			return false, nil
		}
	}

	var removeErr error
	for _, object := range objects {
		if j.Redis() != nil {
			_ = j.Redis().ClearPresignURL(j.Ctx(), object.Key)
		}
		if err := j.MinIO().RemoveObject(j.Ctx(), object.Key); err != nil {
			removeErr = errors.Join(removeErr, err)
			continue
		}
	}

	return removeErr == nil, removeErr
}

// projectImageObjectReferenceSet 将有效项目图像对象引用列表转换为集合
func projectImageObjectReferenceSet(items []model.ProjectImageObjectReference) map[string]struct{} {
	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		set[projectImageObjectReferenceKey(item.ProjectId, item.ImageUuid)] = struct{}{}
	}
	return set
}

// projectDetectionObjectReferenceSet 将有效检测子任务对象引用列表转换为集合
func projectDetectionObjectReferenceSet(items []model.ProjectDetectionObjectReference) map[string]struct{} {
	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		set[projectDetectionObjectReferenceKey(item.ProjectId, item.ImageUuid, item.TaskCode)] = struct{}{}
	}
	return set
}

// projectReportObjectReferenceSet 将有效项目报告对象引用列表转换为集合
func projectReportObjectReferenceSet(items []model.ProjectReportObjectReference) map[string]struct{} {
	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		set[projectReportObjectReferenceKey(item.ProjectId)] = struct{}{}
	}
	return set
}

// projectImageObjectReferenceKey 生成项目图像对象引用集合键
func projectImageObjectReferenceKey(projectId uint64, imageUuid string) string {
	return fmt.Sprintf("%d:%s", projectId, strings.TrimSpace(imageUuid))
}

// projectDetectionObjectReferenceKey 生成检测子任务对象引用集合键
func projectDetectionObjectReferenceKey(projectId uint64, imageUuid, taskCode string) string {
	return fmt.Sprintf("%d:%s:%s", projectId, strings.TrimSpace(imageUuid), strings.TrimSpace(taskCode))
}

// projectReportObjectReferenceKey 生成项目报告对象引用集合键
func projectReportObjectReferenceKey(projectId uint64) string {
	return fmt.Sprintf("%d", projectId)
}

// projectImageAssetPrefix 根据项目 ID 和图像 UUID 生成图像资产文件夹路径
func projectImageAssetPrefix(projectId uint64, imageUuid string) (string, error) {
	projectCode, err := projectCode(projectId)
	if err != nil {
		return "", err
	}
	return projectImageAssetPrefixFromCode(projectCode, imageUuid), nil
}

// projectDetectionTaskPrefix 根据项目 ID、图像 UUID 和任务编码生成检测子任务文件夹路径
func projectDetectionTaskPrefix(projectId uint64, imageUuid, taskCode string) (string, error) {
	projectCode, err := projectCode(projectId)
	if err != nil {
		return "", err
	}
	return projectDetectionTaskPrefixFromCode(projectCode, imageUuid, taskCode), nil
}

// projectReportPrefix 根据项目 ID 生成项目报告文件夹路径
func projectReportPrefix(projectId uint64) (string, error) {
	projectCode, err := projectCode(projectId)
	if err != nil {
		return "", err
	}
	return projectReportPrefixFromCode(projectCode), nil
}

// projectCode 将项目 ID 编码为对象存储路径中的项目短 ID
func projectCode(projectId uint64) (string, error) {
	projectCode := utils.Encode(projectId)
	if projectCode == "" {
		return "", fmt.Errorf("project id is invalid: %d", projectId)
	}
	return projectCode, nil
}

// projectImageAssetPrefixFromCode 根据项目短 ID 和图像 UUID 生成图像资产文件夹路径
func projectImageAssetPrefixFromCode(projectCode, imageUuid string) string {
	return fmt.Sprintf("projects/%s/assets/%s/", strings.TrimSpace(projectCode), strings.TrimSpace(imageUuid))
}

// projectDetectionImagePrefixFromCode 根据项目短 ID 和图像 UUID 生成整张图像检测文件夹路径
func projectDetectionImagePrefixFromCode(projectCode, imageUuid string) string {
	return fmt.Sprintf("projects/%s/detections/%s/", strings.TrimSpace(projectCode), strings.TrimSpace(imageUuid))
}

// projectDetectionTaskPrefixFromCode 根据项目短 ID、图像 UUID 和任务编码生成检测子任务文件夹路径
func projectDetectionTaskPrefixFromCode(projectCode, imageUuid, taskCode string) string {
	return fmt.Sprintf("projects/%s/detections/%s/%s/", strings.TrimSpace(projectCode), strings.TrimSpace(imageUuid), strings.TrimSpace(taskCode))
}

// projectReportPrefixFromCode 根据项目短 ID 生成项目报告文件夹路径
func projectReportPrefixFromCode(projectCode string) string {
	return fmt.Sprintf("projects/%s/reports/", strings.TrimSpace(projectCode))
}

// addCleanupPrefix 添加清理文件夹路径并合并父子路径，避免重复删除
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
