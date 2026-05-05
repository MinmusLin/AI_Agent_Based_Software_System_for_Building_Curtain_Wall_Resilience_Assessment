package detectors

import (
	"icw_activity_reasoning/consts"
	"icw_activity_reasoning/internal/detectors/common"
)

// RegisterDetectors 注册 Python 原子检测能力
func RegisterDetectors(pythonBin, runtimeRoot string) *common.Registry {
	return common.NewRegistry(pythonBin, runtimeRoot, registry())
}

// registry Python 原子检测能力注册表
func registry() []common.DetectorMeta {
	return []common.DetectorMeta{
		common.NewDetectorMeta(consts.CorrosionDetectionTaskCode, "金属锈蚀检测能力"),
		common.NewDetectorMeta(consts.CrackDetectionTaskCode, "石材裂缝检测能力"),
		common.NewDetectorMeta(consts.StainDetectionTaskCode, "石材污渍检测能力"),
		common.NewDetectorMeta(consts.FlatnessDetectionTaskCode, "玻璃平整度检测能力"),
		common.NewDetectorMeta(consts.SpallingDetectionTaskCode, "玻璃爆裂检测能力"),
	}
}
