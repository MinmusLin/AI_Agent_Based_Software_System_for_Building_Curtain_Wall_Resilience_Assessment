package classification

import (
	"context"
	"time"

	"google.golang.org/protobuf/proto"

	"icw_common/enum"
	"icw_common/gen/activity"
	"icw_common/gen/activity/classification"
	"icw_common/gen/core/biz"
	"icw_common/rpc"
	"icw_common/utils"

	classificationUtils "icw_activity_classification/internal/services/classification/utils"
	"icw_activity_classification/internal/services/common"
	"icw_activity_classification/rpc/icw_core_biz"
)

// Start 启动分类任务
func (s *Service) Start(ctx context.Context, req *classificationpb.StartRequest) (*classificationpb.StartResponse, error) {
	resp := &classificationpb.StartResponse{}
	err := s.CallRPC(req, func() error {
		if err := classificationUtils.ValidateRequest(req); err != nil {
			return err
		}
		requestId := rpc.RequestIdFromIncomingContext(ctx)
		taskReq := proto.Clone(req).(*classificationpb.StartRequest)
		go s.asyncExecuteClassification(requestId, taskReq)
		return nil
	})
	return resp, err
}

// asyncExecuteClassification 异步执行分类任务并回调 icw.core.biz
func (s *Service) asyncExecuteClassification(requestId string, req *classificationpb.StartRequest) {
	s.Acquire()
	defer s.Release()

	ctx, cancel := context.WithTimeout(s.Ctx(), time.Duration(s.Config().AgentRequestTimeoutSeconds)*time.Second)
	defer cancel()

	callbackReq := &bizpb.ReportClassificationResultRequest{
		TaskUuid:  req.TaskUuid,
		ImageUuid: req.ImageUuid,
	}

	// 执行分类任务
	classificationStart := time.Now()
	taskCodes, modelOutput, err := classificationUtils.ExecuteClassification(ctx, req, s.AgentClient(), s.Config().AgentImageSize)
	classificationCost := time.Since(classificationStart)

	if utils.IsEmptyError(err) {
		callbackReq.Status = activitypb.DetectionStatus_Succeeded
		callbackReq.TaskCodes = make([]activitypb.DetectionTaskCode_Value, 0, len(taskCodes))
		for _, taskCode := range taskCodes {
			callbackReq.TaskCodes = append(callbackReq.TaskCodes, enum.ParseDetectionTaskCode(taskCode))
		}
		callbackReq.ErrorMessage = ""
		common.ClassificationInfo(requestId, req.TaskUuid, req.ImageUuid, taskCodes, modelOutput, classificationCost)
	} else {
		callbackReq.Status = activitypb.DetectionStatus_Failed
		callbackReq.TaskCodes = nil
		callbackReq.ErrorMessage = err.Error()
		common.ClassificationError(requestId, req.TaskUuid, req.ImageUuid, taskCodes, modelOutput, classificationCost, err)
	}

	// 上报图像检测分类结果
	callbackCtx := rpc.WithRequestIdToOutgoingContext(context.Background(), requestId)
	callbackResp := &bizpb.ReportClassificationResultResponse{}
	callbackStart := time.Now()
	if err := icw_core_biz.ReportClassificationResult(callbackCtx, s.CoreBizClient(), callbackReq, callbackResp); utils.IsEmptyError(err) {
		common.CallbackInfo(requestId, req.TaskUuid, req.ImageUuid, enum.DetectionStatusString(callbackReq.Status), callbackStart)
		return
	} else {
		common.CallbackError(requestId, req.TaskUuid, req.ImageUuid, enum.DetectionStatusString(callbackReq.Status), callbackStart, err)
	}
}
