package workers

import (
	"context"
	"encoding/json"
	"errors"

	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/gen/activity/classification"
	"icw_common/gen/activity/reasoning"
	"icw_common/gen/activity/summary"
	"icw_common/gen/core/biz"
	"icw_common/rpc"

	"icw_core_biz/configs"
	"icw_core_biz/internal/services/project/events"
	"icw_core_biz/repositories/minio"
	"icw_core_biz/repositories/mysql"
	"icw_core_biz/repositories/mysql/model"
	"icw_core_biz/repositories/redis"
	"icw_core_biz/repositories/rocketmq"
	"icw_core_biz/rpc/icw_activity_classification"
	"icw_core_biz/rpc/icw_activity_reasoning"
	"icw_core_biz/rpc/icw_activity_summary"
)

// DetectionWorker 负责项目图像检测任务编排推进
type DetectionWorker struct {
	cfg                  configs.Config
	mysql                *mysql.Repository
	minio                *minio.Repository
	redis                *redis.Repository
	rocketMQ             *rocketmq.Producer
	classificationClient *icw_activity_classification.Client
	reasoningClient      *icw_activity_reasoning.Client
	summaryClient        *icw_activity_summary.Client
	tasks                chan *detectionTask
}

type detectionTask struct {
	id        uint64
	requestId string
}

// NewDetectionWorker 创建项目图像检测任务 Worker
func NewDetectionWorker(
	cfg configs.Config,
	mysqlRepo *mysql.Repository,
	minioRepo *minio.Repository,
	redisRepo *redis.Repository,
	rocketMQProducer *rocketmq.Producer,
	classificationClient *icw_activity_classification.Client,
	reasoningClient *icw_activity_reasoning.Client,
	summaryClient *icw_activity_summary.Client,
) *DetectionWorker {
	return &DetectionWorker{
		cfg:                  cfg,
		mysql:                mysqlRepo,
		minio:                minioRepo,
		redis:                redisRepo,
		rocketMQ:             rocketMQProducer,
		classificationClient: classificationClient,
		reasoningClient:      reasoningClient,
		summaryClient:        summaryClient,
		tasks:                make(chan *detectionTask, 2048),
	}
}

// Start 启动项目图像检测任务 Worker
func (w *DetectionWorker) Start(ctx context.Context) {
	if w == nil {
		return
	}
	for i := 0; i < w.cfg.DetectionWorkerMaxConcurrency; i++ {
		go w.loop(ctx)
	}
}

// Enqueue 投递项目图像检测任务
func (w *DetectionWorker) Enqueue(ctx context.Context, taskId uint64) {
	if w == nil || taskId == 0 {
		return
	}
	w.tasks <- &detectionTask{
		id:        taskId,
		requestId: rpc.RequestIdFromGRPCContext(ctx),
	}
}

func (w *DetectionWorker) loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-w.tasks:
			w.processClassification(ctx, task)
		}
	}
}

func (w *DetectionWorker) processClassification(ctx context.Context, item *detectionTask) {
	if item == nil {
		return
	}
	ctx = rpc.WithRequestIdToOutgoingContext(ctx, item.requestId)

	task, err := w.mysql.StartProjectDetectionClassificationTask(ctx, item.id)
	if err != nil || task == nil {
		return
	}
	w.publishNode(ctx, task, events.DetectionNodeCodeClassification, "", events.DetectionNodeStatusPending)

	originalURL, err := minio.PresignProjectImageOriginalURL(ctx, w.minio, w.redis, task.ProjectId, task.ImageUuid, w.cfg.ProjectImageGetTTL)
	if err != nil || originalURL == "" {
		w.failTask(ctx, task)
		return
	}

	req := &classificationpb.StartRequest{
		TaskUuid:      task.Uuid,
		ImageUuid:     task.ImageUuid,
		PresignGetUrl: originalURL,
	}
	resp := &classificationpb.StartResponse{}
	if err := icw_activity_classification.Start(ctx, w.classificationClient, req, resp); err != nil {
		w.failTask(ctx, task)
	}
}

// StartReasoningTasks 启动项目图像检测推理子任务
func (w *DetectionWorker) StartReasoningTasks(ctx context.Context, subTasks map[string]*model.ProjectDetectionSubTaskRecord) {
	if w == nil {
		return
	}
	for taskCode, subTask := range subTasks {
		if subTask == nil {
			continue
		}
		startedTask, startedSubTask, err := w.mysql.StartProjectDetectionReasoningTask(ctx, taskCode, subTask.Uuid)
		if err != nil {
			continue
		}
		if startedTask == nil || startedSubTask == nil || startedSubTask.Status != bizpb.ProjectDetectionSubTaskStatus_Pending {
			continue
		}
		subTask = startedSubTask
		w.publishNode(ctx, startedTask, events.ReasoningNodeCode(taskCode), subTask.Uuid, events.DetectionNodeStatusPending)
		if err := w.startReasoningTask(ctx, startedTask, taskCode, subTask); err != nil {
			updatedTask, updatedSubTask, _, updateErr := w.mysql.UpdateProjectDetectionReasoningTaskResult(ctx, taskCode, subTask.Uuid, bizpb.ProjectDetectionSubTaskStatus_Failed, "", "")
			if updateErr == nil && updatedTask != nil && updatedSubTask != nil {
				w.publishNode(ctx, updatedTask, events.ReasoningNodeCode(taskCode), updatedSubTask.Uuid, events.DetectionNodeStatusFailed)
			}
		}
	}
}

// StartDetectionSummaryTask 启动项目图像检测总结任务
func (w *DetectionWorker) StartDetectionSummaryTask(ctx context.Context, task *model.ProjectDetectionTaskRecord, summaryTask *model.ProjectDetectionSummaryTaskRecord) {
	if w == nil || task == nil || summaryTask == nil {
		return
	}
	startedTask, startedSummaryTask, err := w.mysql.StartProjectDetectionSummaryTask(ctx, summaryTask.Uuid)
	if err != nil {
		return
	}
	if startedTask == nil || startedSummaryTask == nil || startedSummaryTask.Status != bizpb.ProjectDetectionSubTaskStatus_Pending {
		return
	}
	task = startedTask
	summaryTask = startedSummaryTask
	w.publishNode(ctx, task, events.DetectionNodeCodeSummary, summaryTask.Uuid, events.DetectionNodeStatusPending)

	req := &summarypb.StartDetectionSummaryRequest{
		TaskUuid:  summaryTask.Uuid,
		ImageUuid: task.ImageUuid,
	}
	req.ReportJson, _ = w.detectionSummaryReportJSON(ctx, task)
	resp := &summarypb.StartDetectionSummaryResponse{}
	if err := icw_activity_summary.StartDetectionSummary(ctx, w.summaryClient, req, resp); err != nil {
		updatedTask, updatedSummaryTask, updateErr := w.mysql.UpdateProjectDetectionSummaryResult(ctx, summaryTask.Uuid, bizpb.ProjectDetectionSubTaskStatus_Failed, "")
		if updateErr == nil && updatedTask != nil && updatedSummaryTask != nil {
			w.publishNode(ctx, updatedTask, events.DetectionNodeCodeSummary, updatedSummaryTask.Uuid, events.DetectionNodeStatusFailed)
		}
	}
}

func (w *DetectionWorker) detectionSummaryReportJSON(ctx context.Context, task *model.ProjectDetectionTaskRecord) (string, error) {
	reports := map[string]string{}
	checks := []struct {
		taskCode      string
		shouldExecute bool
		taskId        int64
	}{
		{taskCode: enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Corrosion), shouldExecute: task.CorrosionShouldExecute && task.CorrosionTaskId.Valid, taskId: task.CorrosionTaskId.Int64},
		{taskCode: enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Crack), shouldExecute: task.CrackShouldExecute && task.CrackTaskId.Valid, taskId: task.CrackTaskId.Int64},
		{taskCode: enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Stain), shouldExecute: task.StainShouldExecute && task.StainTaskId.Valid, taskId: task.StainTaskId.Int64},
		{taskCode: enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Flatness), shouldExecute: task.FlatnessShouldExecute && task.FlatnessTaskId.Valid, taskId: task.FlatnessTaskId.Int64},
		{taskCode: enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Spalling), shouldExecute: task.SpallingShouldExecute && task.SpallingTaskId.Valid, taskId: task.SpallingTaskId.Int64},
	}
	for _, check := range checks {
		if !check.shouldExecute {
			continue
		}
		reportJSON, err := w.mysql.GetProjectDetectionSubReportJSON(ctx, check.taskCode, uint64(check.taskId))
		if err != nil {
			return "", err
		}
		if reportJSON != "" {
			reports[check.taskCode] = reportJSON
		}
	}
	data, err := json.Marshal(reports)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (w *DetectionWorker) startReasoningTask(ctx context.Context, task *model.ProjectDetectionTaskRecord, taskCode string, subTask *model.ProjectDetectionSubTaskRecord) error {
	originalURL, err := minio.PresignProjectImageOriginalURL(ctx, w.minio, w.redis, task.ProjectId, task.ImageUuid, w.cfg.ProjectImageGetTTL)
	if err != nil {
		return err
	}
	if originalURL == "" {
		return errors.New("project image original object is not found")
	}
	artifactPolicy, err := w.reasoningArtifactPolicy(ctx, task, taskCode)
	if err != nil {
		return err
	}
	req := &reasoningpb.StartRequest{
		TaskUuid:       subTask.Uuid,
		TaskCode:       enum.ParseDetectionTaskCode(taskCode),
		ImageUuid:      task.ImageUuid,
		PresignGetUrl:  originalURL,
		ArtifactPolicy: artifactPolicy,
	}
	resp := &reasoningpb.StartResponse{}
	return icw_activity_reasoning.Start(ctx, w.reasoningClient, req, resp)
}

func (w *DetectionWorker) reasoningArtifactPolicy(ctx context.Context, task *model.ProjectDetectionTaskRecord, taskCode string) (*reasoningpb.ReasoningArtifactUploadPolicy, error) {
	keyPrefix, err := minio.GenProjectDetectionArtifactPrefixByTask(task.ProjectId, task.ImageUuid, taskCode)
	if err != nil {
		return nil, err
	}
	policyURL, formData, err := w.minio.PresignPostPolicy(ctx, keyPrefix, w.cfg.ProjectImageUploadTTL)
	if err != nil {
		return nil, err
	}
	return &reasoningpb.ReasoningArtifactUploadPolicy{
		Url:       policyURL,
		KeyPrefix: keyPrefix,
		FormData:  formData,
	}, nil
}

func (w *DetectionWorker) failTask(ctx context.Context, task *model.ProjectDetectionTaskRecord) {
	if task == nil {
		return
	}
	updatedTask, err := w.mysql.UpdateProjectDetectionTaskStatus(ctx, task.Id, bizpb.ProjectDetectionTaskStatus_Failed)
	if err != nil || updatedTask == nil {
		return
	}
	w.publishNode(ctx, updatedTask, events.DetectionNodeCodeClassification, "", events.DetectionNodeStatusFailed)
}

func (w *DetectionWorker) publishNode(ctx context.Context, task *model.ProjectDetectionTaskRecord, nodeCode, subTaskUuid string, subStatus bizpb.ProjectDetectionSubTaskStatus_Value) {
	if task == nil {
		return
	}
	events.PublishProjectDetectionNodeStatusChangedEvent(
		ctx,
		w.rocketMQ,
		task.UserId,
		task.ProjectId,
		task.ImageUuid,
		nodeCode,
		task.Uuid,
		task.Status,
		subTaskUuid,
		subStatus,
	)
}
