package core

import (
	"context"

	"icw_common/enum"
	"icw_common/gen/core/biz"
	"icw_common/rpc_err"
)

// CheckProjectAccess 校验项目访问权限
func (s *Service) CheckProjectAccess(ctx context.Context, req *bizpb.CheckProjectAccessRequest) (*bizpb.CheckProjectAccessResponse, error) {
	resp := &bizpb.CheckProjectAccessResponse{}
	err := s.CallRPC(ctx, req, resp, func() error {
		return s.checkProjectAccess(req, resp)
	})
	return resp, err
}

func (s *Service) checkProjectAccess(req *bizpb.CheckProjectAccessRequest, resp *bizpb.CheckProjectAccessResponse) error {
	projectRecord, err := s.MySQL().FindProjectByIdAndUserId(s.Ctx(), req.UserId, req.ProjectId)
	if err != nil {
		return err
	}
	if projectRecord == nil || (projectRecord.Status != bizpb.ProjectStatus_Active && projectRecord.Status != bizpb.ProjectStatus_Completed) {
		return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project is not accessible")
	}

	resp.ProjectId = projectRecord.Id
	resp.Progress = uint32(projectRecord.Progress.Uint8())
	resp.Status = enum.ProjectStatusString(projectRecord.Status)

	return nil
}
