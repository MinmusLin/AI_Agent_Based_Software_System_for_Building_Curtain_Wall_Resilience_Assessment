package core

import (
	"icw_core_biz/pkg/dto"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
)

// CheckProjectAccess 校验用户是否拥有项目访问权限，并返回项目当前进度与状态
func (s *Service) CheckProjectAccess(req *project.CheckProjectAccessRequest, resp *project.CheckProjectAccessResponse) error {
	return s.CallRPC("ProjectCoreService.CheckProjectAccess", req, resp, func() error {
		return s.checkProjectAccess(req, resp)
	})
}

func (s *Service) checkProjectAccess(req *project.CheckProjectAccessRequest, resp *project.CheckProjectAccessResponse) (err error) {
	projectRecord, err := s.MySQL().FindProjectByIdAndUserId(s.Ctx(), req.UserId, req.ProjectId)
	if err != nil {
		return err
	}
	if projectRecord == nil || projectRecord.Status == dto.ProjectStatusDeleted {
		return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project is not accessible")
	}

	resp.ProjectId = projectRecord.Id
	resp.Progress = projectRecord.Progress.Uint8()
	resp.Status = projectRecord.Status.String()

	return nil
}
