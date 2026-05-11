package workers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"

	"icw_common/enum"
	"icw_common/gen/core/biz"

	"icw_core_biz/repositories/mysql/model"
	"icw_core_biz/repositories/mysql/utils"
)

// projectReportSourceJSON 构建项目评估报告原始数据 JSON
func (w *DetectionWorker) projectReportSourceJSON(ctx context.Context, userId, projectId uint64) (string, error) {
	project, err := w.mysql.FindProjectByIdAndUserId(ctx, userId, projectId)
	if err != nil {
		return "", err
	}
	if project == nil {
		return "", errors.New("project is nil")
	}
	groups, err := w.mysql.ListProjectGroups(ctx, userId, projectId)
	if err != nil {
		return "", err
	}
	if len(groups) == 0 {
		return "", errors.New("project groups is nil")
	}
	images, err := w.mysql.ListProjectImages(ctx, userId, projectId)
	if err != nil {
		return "", err
	}
	if len(images) == 0 {
		return "", errors.New("project images is nil")
	}

	groupItems := make([]*projectReportSourceGroup, 0, len(groups))
	groupById := make(map[uint64]*projectReportSourceGroup, len(groups))
	for _, group := range groups {
		if group == nil {
			continue
		}
		item := &projectReportSourceGroup{
			GroupName: group.Name,
			Images:    make([]*projectReportSourceImage, 0),
		}
		groupItems = append(groupItems, item)
		groupById[group.Id] = item
	}

	for _, image := range images {
		if image == nil || image.Status != bizpb.ProjectImageStatus_Uploaded {
			continue
		}
		group := groupById[image.GroupId]
		if group == nil {
			group = &projectReportSourceGroup{
				Images: make([]*projectReportSourceImage, 0),
			}
			groupItems = append(groupItems, group)
			groupById[image.GroupId] = group
		}

		imageItem, err := w.projectReportSourceImage(ctx, userId, projectId, image)
		if err != nil {
			return "", err
		}
		group.Images = append(group.Images, imageItem)
	}

	data, err := json.Marshal(&projectReportSource{
		Project: projectReportSourceProject{
			ProjectName:         project.Name,
			BuildingName:        project.BuildingName,
			BuildingLocation:    project.BuildingLocation,
			BuiltYear:           utils.NullUint32(project.BuiltYear),
			BuildingDescription: utils.NullString(project.BuildingDescription),
			KnownIssues:         utils.NullString(project.KnownIssues),
			AssessmentGoal:      utils.NullString(project.AssessmentGoal),
		},
		Groups: groupItems,
	})
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// projectReportSourceImage 构建项目评估报告图像原始数据
func (w *DetectionWorker) projectReportSourceImage(ctx context.Context, userId, projectId uint64, image *model.ProjectImageRecord) (*projectReportSourceImage, error) {
	item := &projectReportSourceImage{
		FileName: image.FileName,
	}

	task, err := w.mysql.FindProjectDetectionTaskByImageUuid(ctx, userId, projectId, image.Uuid)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return item, nil
	}

	detection, err := w.projectReportSourceDetection(ctx, task)
	if err != nil {
		return nil, err
	}
	item.Detection = detection

	review, err := w.mysql.GetProjectDetectionReview(ctx, userId, projectId, task.Uuid)
	if errors.Is(err, sql.ErrNoRows) {
		return item, nil
	}
	if err != nil {
		return nil, err
	}
	item.Review = w.projectReportSourceReview(review)

	return item, nil
}

// projectReportSourceDetection 构建项目评估报告图像检测原始数据
func (w *DetectionWorker) projectReportSourceDetection(ctx context.Context, task *model.ProjectDetectionTaskRecord) (*projectReportSourceDetection, error) {
	detection := &projectReportSourceDetection{}

	var err error
	if task.CorrosionShouldExecute && task.CorrosionTaskId.Valid {
		detection.Corrosion, err = w.projectReportSourceCorrosion(ctx, uint64(task.CorrosionTaskId.Int64))
		if err != nil {
			return nil, err
		}
	}
	if task.CrackShouldExecute && task.CrackTaskId.Valid {
		detection.Crack, err = w.projectReportSourceCrack(ctx, uint64(task.CrackTaskId.Int64))
		if err != nil {
			return nil, err
		}
	}
	if task.StainShouldExecute && task.StainTaskId.Valid {
		detection.Stain, err = w.projectReportSourceStain(ctx, uint64(task.StainTaskId.Int64))
		if err != nil {
			return nil, err
		}
	}
	if task.FlatnessShouldExecute && task.FlatnessTaskId.Valid {
		detection.Flatness, err = w.projectReportSourceFlatness(ctx, uint64(task.FlatnessTaskId.Int64))
		if err != nil {
			return nil, err
		}
	}
	if task.SpallingShouldExecute && task.SpallingTaskId.Valid {
		detection.Spalling, err = w.projectReportSourceSpalling(ctx, uint64(task.SpallingTaskId.Int64))
		if err != nil {
			return nil, err
		}
	}
	if task.SummaryShouldExecute && task.SummaryTaskId.Valid {
		detection.Summary, err = w.projectReportSourceSummary(ctx, uint64(task.SummaryTaskId.Int64))
		if err != nil {
			return nil, err
		}
	}

	return detection, nil
}

// projectReportSourceCorrosion 构建项目评估报告图像金属锈蚀检测原始数据
func (w *DetectionWorker) projectReportSourceCorrosion(ctx context.Context, taskId uint64) (*projectReportSourceCorrosion, error) {
	result, err := w.mysql.GetProjectDetectionCorrosionResult(ctx, taskId)
	if err != nil || result == nil {
		return nil, err
	}
	regions := make([]*projectReportSourceCorrosionRegion, 0, len(result.Regions))
	for _, region := range result.Regions {
		if region == nil {
			continue
		}
		regions = append(regions, &projectReportSourceCorrosionRegion{
			Confidence: region.Confidence,
			MaskPixels: region.MaskPixels,
			MaskRatio:  region.MaskRatio,
		})
	}
	return &projectReportSourceCorrosion{
		HasCorrosion:      result.HasCorrosion,
		CorrosionCount:    result.CorrosionCount,
		MaxConfidence:     result.MaxConfidence,
		AverageConfidence: result.AverageConfidence,
		CorrosionPixels:   result.CorrosionPixels,
		CorrosionRatio:    result.CorrosionRatio,
		Regions:           regions,
	}, nil
}

// projectReportSourceCrack 构建项目评估报告图像石材裂缝检测原始数据
func (w *DetectionWorker) projectReportSourceCrack(ctx context.Context, taskId uint64) (*projectReportSourceCrack, error) {
	result, err := w.mysql.GetProjectDetectionCrackResult(ctx, taskId)
	if err != nil || result == nil {
		return nil, err
	}
	regions := make([]*projectReportSourceCrackRegion, 0, len(result.Regions))
	for _, region := range result.Regions {
		if region == nil {
			continue
		}
		regions = append(regions, &projectReportSourceCrackRegion{
			MaskPixels: region.MaskPixels,
			MaskRatio:  region.MaskRatio,
		})
	}
	return &projectReportSourceCrack{
		HasCrack:    result.HasCrack,
		CrackCount:  result.CrackCount,
		CrackPixels: result.CrackPixels,
		CrackRatio:  result.CrackRatio,
		Regions:     regions,
	}, nil
}

// projectReportSourceStain 构建项目评估报告图像石材污渍检测原始数据
func (w *DetectionWorker) projectReportSourceStain(ctx context.Context, taskId uint64) (*projectReportSourceStain, error) {
	result, err := w.mysql.GetProjectDetectionStainResult(ctx, taskId)
	if err != nil || result == nil {
		return nil, err
	}
	regions := make([]*projectReportSourceStainRegion, 0, len(result.Regions))
	for _, region := range result.Regions {
		if region == nil {
			continue
		}
		regions = append(regions, &projectReportSourceStainRegion{
			Confidence:   region.Confidence,
			RegionWidth:  region.RegionWidth,
			RegionHeight: region.RegionHeight,
			StainPixels:  region.StainPixels,
			StainRatio:   region.StainRatio,
		})
	}
	return &projectReportSourceStain{
		HasStain:          result.HasStain,
		StainCount:        result.StainCount,
		AverageStainRatio: result.AverageStainRatio,
		MaxStainRatio:     result.MaxStainRatio,
		Regions:           regions,
	}, nil
}

// projectReportSourceFlatness 构建项目评估报告图像玻璃平整度检测原始数据
func (w *DetectionWorker) projectReportSourceFlatness(ctx context.Context, taskId uint64) (*projectReportSourceFlatness, error) {
	result, err := w.mysql.GetProjectDetectionFlatnessResult(ctx, taskId)
	if err != nil || result == nil {
		return nil, err
	}
	regions := make([]*projectReportSourceFlatnessRegion, 0, len(result.Regions))
	for _, region := range result.Regions {
		if region == nil {
			continue
		}
		regions = append(regions, &projectReportSourceFlatnessRegion{
			EdgeUnevenDetected:      region.EdgeUnevenDetected,
			LineUnevenDetected:      region.LineUnevenDetected,
			GradientUnevenDetected:  region.GradientUnevenDetected,
			FrequencyUnevenDetected: region.FrequencyUnevenDetected,
			EdgeCount:               region.EdgeCount,
			LaplacianVariance:       region.LaplacianVariance,
			LineCount:               region.LineCount,
			AngleStd:                region.AngleStd,
			GradientMean:            region.GradientMean,
			GradientStd:             region.GradientStd,
			FrequencyMin:            region.FrequencyMin,
			FrequencyMax:            region.FrequencyMax,
		})
	}
	return &projectReportSourceFlatness{
		Result:      result.Result,
		UnevenCount: result.UnevenCount,
		Regions:     regions,
	}, nil
}

// projectReportSourceSpalling 构建项目评估报告图像玻璃爆裂检测原始数据
func (w *DetectionWorker) projectReportSourceSpalling(ctx context.Context, taskId uint64) (*projectReportSourceSpalling, error) {
	result, err := w.mysql.GetProjectDetectionSpallingResult(ctx, taskId)
	if err != nil || result == nil {
		return nil, err
	}
	return &projectReportSourceSpalling{
		HasSpalling: result.HasSpalling,
		Confidence:  result.Confidence,
	}, nil
}

// projectReportSourceSummary 构建项目评估报告图像总结原始数据
func (w *DetectionWorker) projectReportSourceSummary(ctx context.Context, taskId uint64) (*projectReportSourceSummary, error) {
	result, err := w.mysql.GetProjectDetectionSummaryTypedResult(ctx, taskId)
	if err != nil || result == nil {
		return nil, err
	}
	return &projectReportSourceSummary{
		Result: result.Result,
	}, nil
}

// projectReportSourceReview 构建项目评估报告图像人工复核原始数据
func (w *DetectionWorker) projectReportSourceReview(review *bizpb.ProjectDetectionReview) *projectReportSourceReview {
	if review == nil {
		return nil
	}
	item := &projectReportSourceReview{
		Verdict: enum.ProjectDetectionReviewVerdictString(review.Verdict),
		Comment: strings.TrimSpace(review.Comment),
	}
	if item.Verdict == "" && item.Comment == "" {
		return nil
	}
	return item
}
