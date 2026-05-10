package detectors

import (
	"icw_common/enum"
	"icw_common/gen/core/common"

	"icw_activity_reasoning/internal/detectors/common"
)

// RegisterDetectors 注册 Python 原子检测能力
func RegisterDetectors(runtimeRoot, modelRoot string) (*common.Registry, error) {
	return common.NewRegistry(runtimeRoot, modelRoot, registry())
}

// registry Python 原子检测能力注册表
func registry() []*common.DetectorMeta {
	return []*common.DetectorMeta{
		common.NewDetectorMeta(
			enum.DetectionTaskCodeString(commonpb.DetectionTaskCode_Corrosion),
			"金属锈蚀检测能力",
			"best_weights_model.pt",
		),
		common.NewDetectorMeta(
			enum.DetectionTaskCodeString(commonpb.DetectionTaskCode_Crack),
			"石材裂缝检测能力",
			"best_weights_model.pt",
		),
		common.NewDetectorMeta(
			enum.DetectionTaskCodeString(commonpb.DetectionTaskCode_Stain),
			"石材污渍检测能力",
			"best_weights_model.pt",
		),
		common.NewDetectorMeta(
			enum.DetectionTaskCodeString(commonpb.DetectionTaskCode_Flatness),
			"玻璃平整度检测能力",
			"best_weights_model.pt",
		),
		common.NewDetectorMeta(
			enum.DetectionTaskCodeString(commonpb.DetectionTaskCode_Spalling),
			"玻璃爆裂检测能力",
			"best_weights_model.pt",
		),
	}
}
