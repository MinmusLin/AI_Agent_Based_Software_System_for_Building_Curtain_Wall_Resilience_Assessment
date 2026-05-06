package detectors

import (
	"icw_activity_reasoning/internal/detectors/common"
	"icw_common/enum"
	"icw_common/gen/activity"
)

// RegisterDetectors 注册 Python 原子检测能力
func RegisterDetectors(pythonBin, runtimeRoot string) *common.Registry {
	return common.NewRegistry(pythonBin, runtimeRoot, registry())
}

// registry Python 原子检测能力注册表
func registry() []*common.DetectorMeta {
	return []*common.DetectorMeta{
		common.NewDetectorMeta(
			enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Corrosion),
			"金属锈蚀检测能力",
		),
		common.NewDetectorMeta(
			enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Crack),
			"石材裂缝检测能力",
		),
		common.NewDetectorMeta(
			enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Stain),
			"石材污渍检测能力",
		),
		common.NewDetectorMeta(
			enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Flatness),
			"玻璃平整度检测能力",
		),
		common.NewDetectorMeta(
			enum.DetectionTaskCodeString(activitypb.DetectionTaskCode_Spalling),
			"玻璃爆裂检测能力",
		),
	}
}
