package core

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/internal/services/project/utils"
)

// AdvanceProject 项目进度流转
func (s *Service) AdvanceProject(ctx context.Context, req *bizpb.AdvanceProjectRequest) (*bizpb.AdvanceProjectResponse, error) {
	resp := &bizpb.AdvanceProjectResponse{}
	err := s.CallRPC(req, func() error {
		return s.advanceProject(req, resp)
	})
	return resp, err
}

func (s *Service) advanceProject(req *bizpb.AdvanceProjectRequest, _ *bizpb.AdvanceProjectResponse) error {
	// 校验项目进度推进合法
	maxProgress := bizpb.ProjectProgress_ReportFinished
	if req.FromProgress >= maxProgress || req.ToProgress > maxProgress {
		return rpc_error.BadRequestDefault("project progress is out of range")
	}
	if req.ToProgress != req.FromProgress+1 {
		return rpc_error.BadRequestDefault("project progress can only advance one step")
	}
	fromProgress := req.FromProgress
	toProgress := req.ToProgress
	nextStatus := bizpb.ProjectStatus_Active
	if toProgress == bizpb.ProjectProgress_ReportFinished {
		nextStatus = bizpb.ProjectStatus_Completed
	}

	// 如果项目进度已被并发请求流转，请求视为成功
	alreadyAdvanced, err := utils.ProjectAlreadyAdvanced(s.Ctx(), s.MySQL(), req.UserId, req.ProjectId, toProgress, nextStatus)
	if err != nil {
		return err
	}
	if alreadyAdvanced {
		return nil
	}

	// 执行项目进度流转前置扩展点
	if err := utils.BeforeAdvanceProject(s.Ctx(), s.MySQL(), req.UserId, req.ProjectId, fromProgress, toProgress); err != nil {
		return err
	}

	// 执行项目进度流转扩展点
	advanced, err := utils.AdvanceProject(s.Ctx(), s.MySQL(), req.UserId, req.ProjectId, fromProgress, toProgress, nextStatus)
	if err != nil {
		return err
	}
	if !advanced {
		// 如果项目进度已被并发请求流转，请求视为成功
		alreadyAdvanced, err := utils.ProjectAlreadyAdvanced(s.Ctx(), s.MySQL(), req.UserId, req.ProjectId, toProgress, nextStatus)
		if err != nil {
			return err
		}
		if alreadyAdvanced {
			return nil
		}
		return rpc_error.BadRequestDefault("project status or progress has changed")
	}

	if fromProgress == bizpb.ProjectProgress_DetectionFinished && toProgress == bizpb.ProjectProgress_ReviewFinished {
		go s.DetectionWorker().StartProjectReportTask(s.Ctx(), req.UserId, req.ProjectId)
	}

	return nil
}
