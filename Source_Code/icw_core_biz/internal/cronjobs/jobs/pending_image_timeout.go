package jobs

import (
	"icw_core_biz/internal/cronjobs/common"
	"icw_core_biz/internal/services/project/events"
	"icw_core_biz/repositories/mysql"
)

// PendingImageTimeoutJob 上传中图像超时失败任务
type PendingImageTimeoutJob struct {
	*common.BaseCronJob
}

// NewPendingImageTimeoutJob 创建上传中图像超时失败任务
func NewPendingImageTimeoutJob(baseJob *common.BaseCronJob) common.CronJob {
	return &PendingImageTimeoutJob{
		BaseCronJob: baseJob,
	}
}

// PendingImageTimeoutJobResult 上传中图像超时失败任务执行结果
type PendingImageTimeoutJobResult struct {
	PublishImageUuids []string `json:"publish_image_uuids"`
}

// Start 执行上传中图像超时失败任务
func (j *PendingImageTimeoutJob) Start() (interface{}, error) {
	// 将超时的上传中项目图像状态更新为上传失败
	imageRecords, err := j.MySQL().FailTimeoutPendingProjectImages(j.Ctx(), j.Config().ProjectImagePendingTimeout)
	if err != nil {
		return nil, err
	}

	result := &PendingImageTimeoutJobResult{
		PublishImageUuids: make([]string, 0),
	}

	for _, imageRecord := range imageRecords {
		image, err := mysql.ProjectImageRecordToDTO(j.Ctx(), j.MinIO(), j.Redis(), imageRecord, j.Config().ProjectImageGetTTL)
		if err != nil {
			common.CronWarn("Convert timeout pending image failed, image_uuid: %s, err: %v", imageRecord.Uuid, err)
			continue
		}

		// 发布项目图像状态变化事件
		events.PublishProjectImageStatusChangedEvent(j.Ctx(), j.RocketMQ(), imageRecord.UserId, imageRecord.ProjectId, image)
		result.PublishImageUuids = append(result.PublishImageUuids, image.Uuid)
	}

	return result, nil
}
