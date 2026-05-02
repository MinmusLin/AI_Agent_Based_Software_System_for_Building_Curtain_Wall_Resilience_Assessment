package core

import (
	"context"

	"icw_core_biz/internal/rpc_log"
	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
)

// AdvanceProject 项目进度流转
func (s *Service) AdvanceProject(req *project.AdvanceProjectRequest, resp *project.AdvanceProjectResponse) (err error) {
	start := rpc_log.Start("ProjectCoreService.AdvanceProject", req)
	defer func() {
		rpc_log.Finish("ProjectCoreService.AdvanceProject", req, resp, start, err)
	}()

	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	// 校验项目进度推进合法
	maxProgress := consts.ProjectProgressReportFinished.Uint8()
	if req.FromProgress >= maxProgress || req.ToProgress > maxProgress {
		return rpc_err.BadRequestDefault("project progress is out of range")
	}
	if req.ToProgress != req.FromProgress+1 {
		return rpc_err.BadRequestDefault("project progress can only advance one step")
	}
	fromProgress := consts.ParseProjectProgressFromUint8(req.FromProgress)
	toProgress := consts.ParseProjectProgressFromUint8(req.ToProgress)
	nextStatus := consts.ProjectStatusActive
	if toProgress == consts.ProjectProgressReportFinished {
		nextStatus = consts.ProjectStatusCompleted
	}

	// 校验用户是否拥有项目访问权限
	projectRecord, err := utils.ValidateProjectOwnership(ctx, s.MySQL(), req.UserId, req.ProjectId)
	if err != nil {
		return err
	}

	// 如果项目已经推进到目标阶段，允许重试项目进度流转后置扩展点
	if projectRecord.Progress == toProgress && projectRecord.Status == nextStatus {
		return utils.PostAdvanceProject(ctx, s.MySQL(), req.UserId, req.ProjectId, fromProgress, toProgress)
	}

	// 校验项目状态和进度是否已经发生变化
	if projectRecord.Status != consts.ProjectStatusActive || projectRecord.Progress != fromProgress {
		return rpc_err.BadRequestDefault("project status or progress has changed")
	}

	// 执行项目进度流转前置扩展点
	if err := utils.PreAdvanceProject(ctx, s.MySQL(), req.UserId, req.ProjectId, fromProgress, toProgress); err != nil {
		return err
	}

	advanced, err := s.MySQL().AdvanceProject(ctx, req.UserId, req.ProjectId, fromProgress, toProgress, nextStatus)
	if err != nil {
		return err
	}
	if !advanced {
		return rpc_err.BadRequestDefault("project status or progress has changed")
	}

	// 执行项目进度流转后置扩展点
	if err := utils.PostAdvanceProject(ctx, s.MySQL(), req.UserId, req.ProjectId, fromProgress, toProgress); err != nil {
		return err
	}

	return nil
}
