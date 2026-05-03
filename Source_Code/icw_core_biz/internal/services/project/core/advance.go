package core

import (
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
)

// AdvanceProject 项目进度流转
func (s *Service) AdvanceProject(req *project.AdvanceProjectRequest, resp *project.AdvanceProjectResponse) error {
	return s.CallRPC("ProjectCoreService.AdvanceProject", req, resp, func() error {
		return s.advanceProject(req, resp)
	})
}

func (s *Service) advanceProject(req *project.AdvanceProjectRequest, resp *project.AdvanceProjectResponse) (err error) {
	// 校验项目进度推进合法
	maxProgress := dto.ProjectProgressReportFinished.Uint8()
	if req.FromProgress >= maxProgress || req.ToProgress > maxProgress {
		return rpc_err.BadRequestDefault("project progress is out of range")
	}
	if req.ToProgress != req.FromProgress+1 {
		return rpc_err.BadRequestDefault("project progress can only advance one step")
	}
	fromProgress := dto.ParseProjectProgress(req.FromProgress)
	toProgress := dto.ParseProjectProgress(req.ToProgress)
	nextStatus := dto.ProjectStatusActive
	if toProgress == dto.ProjectProgressReportFinished {
		nextStatus = dto.ProjectStatusCompleted
	}

	// 执行项目进度流转前置扩展点
	if err := utils.PreAdvanceProject(s.Ctx, s.MySQL(), req.UserId, req.ProjectId, fromProgress, toProgress); err != nil {
		return err
	}

	advanced, err := utils.AdvanceProject(s.Ctx, s.MySQL(), req.UserId, req.ProjectId, fromProgress, toProgress, nextStatus)
	if err != nil {
		return err
	}
	if !advanced {
		return rpc_err.BadRequestDefault("project status or progress has changed")
	}

	return nil
}
