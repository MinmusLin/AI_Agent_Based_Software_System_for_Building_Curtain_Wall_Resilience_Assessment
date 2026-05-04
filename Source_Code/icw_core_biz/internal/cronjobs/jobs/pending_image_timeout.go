package jobs

import (
	"context"

	"icw_core_biz/internal/cronjobs/common"
	"icw_core_biz/internal/services/project/events"
	"icw_core_biz/repositories/mysql"
)

// PendingImageTimeoutJob 上传中图像超时失败任务
type PendingImageTimeoutJob struct {
	*common.CronJob
}

// StartPendingImageTimeoutJob 执行上传中图像超时失败任务
func StartPendingImageTimeoutJob(ctx context.Context, deps *common.Deps, cronExpression string) error {
	if deps == nil || deps.MySQL == nil || deps.Redis == nil || deps.MinIO == nil || deps.RocketMQ == nil {
		common.CronWarn("Skip project assets pending image timeout job: dependencies are nil")
		return nil
	}

	job := &PendingImageTimeoutJob{
		CronJob: common.NewCronJob(ctx, deps),
	}

	return job.Schedule("project_assets.pending_image_timeout", cronExpression, job.pendingImageTimeoutJob)
}

func (j *PendingImageTimeoutJob) pendingImageTimeoutJob() error {
	imageRecords, err := j.MySQL().FailTimeoutPendingProjectImages(j.Ctx(), j.Config().ProjectImagePendingTimeout)
	if err != nil {
		return err
	}
	for _, imageRecord := range imageRecords {
		image, err := mysql.ProjectImageRecordToDTO(j.Ctx(), j.MinIO(), j.Redis(), imageRecord, j.Config().ProjectImageGetTTL)
		if err != nil {
			common.CronWarn("Convert timeout pending project image failed, project_id: %d, image_uuid: %s, err: %v", imageRecord.ProjectId, imageRecord.Uuid, err)
			continue
		}
		// 发布项目图像状态变化事件
		events.PublishProjectImageStatusChangedEvent(j.Ctx(), j.RocketMQ(), imageRecord.UserId, imageRecord.ProjectId, image)
	}
	return nil
}
