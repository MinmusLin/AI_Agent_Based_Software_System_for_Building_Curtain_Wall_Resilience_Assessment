package core

import (
	"context"

	"icw_core_biz/internal/rpc_log"
	projectUtils "icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/mysql"
)

// AdvanceProject 项目进度流转
func (s *Service) AdvanceProject(req *project.AdvanceProjectRequest, resp *project.AdvanceProjectResponse) (err error) {
	start := rpc_log.Start("ProjectCoreService.AdvanceProjectProgress", req)
	defer func() {
		rpc_log.Finish("ProjectCoreService.AdvanceProjectProgress", req, resp, start, err)
	}()

	if req == nil {
		return rpc_err.BadRequestDefault("request is nil")
	}
	ctx := context.Background()

	// 校验项目进度推进合法
	maxProgress := uint8(mysql.ProjectProgressReportFinished)
	if req.FromProgress >= maxProgress || req.ToProgress > maxProgress {
		return rpc_err.BadRequestDefault("project progress is out of range")
	}
	if req.ToProgress != req.FromProgress+1 {
		return rpc_err.BadRequestDefault("project progress can only advance one step")
	}
	fromProgress := mysql.ProjectProgress(req.FromProgress)
	toProgress := mysql.ProjectProgress(req.ToProgress)

	// 校验用户是否拥有项目访问权限
	projectRecord, err := projectUtils.ValidateProjectOwnership(ctx, s.MySQL(), req.UserId, req.ProjectId)
	if err != nil {
		return err
	}

	// 校验项目状态和进度是否已经发生变化
	if projectRecord.Status != string(mysql.ProjectStatusActive) || projectRecord.Progress != fromProgress {
		return rpc_err.BadRequestDefault("project status or progress has changed")
	}

	// 如果项目进度推进结束，则更新项目状态为已完成
	nextStatus := mysql.ProjectStatusActive
	if toProgress == mysql.ProjectProgressReportFinished {
		nextStatus = mysql.ProjectStatusCompleted
	}

	advanced, err := s.MySQL().AdvanceProject(ctx, req.UserId, req.ProjectId, fromProgress, toProgress, nextStatus)
	if err != nil {
		return err
	}
	if !advanced {
		return rpc_err.BadRequestDefault("project status or progress has changed")
	}

	return nil
}
