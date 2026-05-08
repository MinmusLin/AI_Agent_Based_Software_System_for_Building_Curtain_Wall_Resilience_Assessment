package core

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"
)

// CheckProjectAccess 校验项目访问权限
func (s *Service) CheckProjectAccess(ctx context.Context, req *bizpb.CheckProjectAccessRequest) (*bizpb.CheckProjectAccessResponse, error) {
	resp := &bizpb.CheckProjectAccessResponse{}
	err := s.CallRPC(req, func() error {
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
		return rpc_error.BadRequest(rpc_error.DetailProjectNotAccessible, "project is not accessible")
	}

	resp.ProjectId = projectRecord.Id
	resp.Progress = projectRecord.Progress
	resp.Status = projectRecord.Status

	return nil
}
