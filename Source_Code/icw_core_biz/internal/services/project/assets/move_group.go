package assets

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/repositories/mysql"
)

// MoveProjectGroup 移动图像组
func (s *Service) MoveProjectGroup(ctx context.Context, req *bizpb.MoveProjectGroupRequest) (*bizpb.MoveProjectGroupResponse, error) {
	resp := &bizpb.MoveProjectGroupResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.moveProjectGroup(req, resp)
	})
	return resp, err
}

func (s *Service) moveProjectGroup(req *bizpb.MoveProjectGroupRequest, resp *bizpb.MoveProjectGroupResponse) error {
	if req.PreviousGroupId == req.NextGroupId || req.GroupId == req.PreviousGroupId || req.GroupId == req.NextGroupId {
		return rpc_error.BadRequestDefault("project group move target is invalid")
	}
	if req.MoveToFirst && req.MoveToLast {
		return rpc_error.BadRequestDefault("move to first and move to last cannot both be true")
	}
	if !req.MoveToFirst && !req.MoveToLast {
		if req.PreviousGroupId == 0 || req.NextGroupId == 0 {
			return rpc_error.BadRequestDefault("previous group id and next group id are required")
		}
	}
	if req.MoveToFirst && !req.MoveToLast {
		if req.NextGroupId == 0 {
			return rpc_error.BadRequestDefault("next group id is required")
		}
	}
	if !req.MoveToFirst && req.MoveToLast {
		if req.PreviousGroupId == 0 {
			return rpc_error.BadRequestDefault("previous group id is required")
		}
	}

	groupRecord, err := s.MySQL().FindProjectGroupById(s.Ctx(), req.UserId, req.ProjectId, req.GroupId)
	if err != nil {
		return err
	}
	if groupRecord == nil {
		return rpc_error.BadRequest(rpc_error.DetailProjectNotAccessible, "project group is not accessible")
	}

	groupRecord, err = s.MySQL().MoveProjectGroup(s.Ctx(), req.UserId, req.ProjectId, req.GroupId, req.PreviousGroupId, req.NextGroupId, req.MoveToFirst, req.MoveToLast)
	if err != nil {
		return err
	}
	if groupRecord == nil {
		return rpc_error.BadRequest(rpc_error.DetailProjectNotAccessible, "project group is not accessible")
	}

	resp.Group, err = mysql.ProjectGroupRecordToDTO(s.Ctx(), s.MinIO(), s.Redis(), groupRecord, nil, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}

	return nil
}
