package consts

import (
	"icw_common/gen/activity/classification"
	"icw_common/gen/activity/reasoning"
	"icw_common/gen/activity/summary"
	"icw_common/gen/core/api"
	"icw_common/gen/core/biz"
)

var (
	// CoreApiPSM icw.core.api 服务标识
	CoreApiPSM = string(apipb.File_core_api_proto.Package())
	// CoreBizPSM icw.core.biz 服务标识
	CoreBizPSM = string(bizpb.File_core_biz_proto.Package())
	// ActivityClassificationPSM icw.activity.classification 服务标识
	ActivityClassificationPSM = string(classificationpb.File_activity_classification_proto.Package())
	// ActivityReasoningPSM icw.activity.reasoning 服务标识
	ActivityReasoningPSM = string(reasoningpb.File_activity_reasoning_proto.Package())
	// ActivitySummaryPSM icw.activity.summary 服务标识
	ActivitySummaryPSM = string(summarypb.File_activity_summary_proto.Package())
)
