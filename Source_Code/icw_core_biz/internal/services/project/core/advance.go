package core

import (
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
)

// AdvanceProject 项目进度流转
func (s *Service) AdvanceProject(req *project.AdvanceProjectRequest, resp *project.AdvanceProjectResponse) error {
	return s.CallRPC(req, resp, func() error {
		return s.advanceProject(req, resp)
	})
}

func (s *Service) advanceProject(req *project.AdvanceProjectRequest, _ *project.AdvanceProjectResponse) error {
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
		return rpc_err.BadRequestDefault("project status or progress has changed")
	}

	return nil
}
