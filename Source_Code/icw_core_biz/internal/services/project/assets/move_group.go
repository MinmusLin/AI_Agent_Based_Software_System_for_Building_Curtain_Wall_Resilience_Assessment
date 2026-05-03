package assets

import (
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/mysql"
)

// MoveProjectGroup 移动图像组
func (s *Service) MoveProjectGroup(req *project.MoveProjectGroupRequest, resp *project.MoveProjectGroupResponse) error {
	return s.CallRPC("ProjectAssetsService.MoveProjectGroup", req, resp, func() error {
		return s.moveProjectGroup(req, resp)
	})
}

func (s *Service) moveProjectGroup(req *project.MoveProjectGroupRequest, resp *project.MoveProjectGroupResponse) error {
	if req.PreviousGroupId == req.NextGroupId || req.GroupId == req.PreviousGroupId || req.GroupId == req.NextGroupId {
		return rpc_err.BadRequestDefault("project group move target is invalid")
	}
	if req.MoveToFirst && req.MoveToLast {
		return rpc_err.BadRequestDefault("move to first and move to last cannot both be true")
	}
	if !req.MoveToFirst && !req.MoveToLast {
		if req.PreviousGroupId == 0 || req.NextGroupId == 0 {
			return rpc_err.BadRequestDefault("previous group id and next group id are required")
		}
	}
	if req.MoveToFirst && !req.MoveToLast {
		if req.NextGroupId == 0 {
			return rpc_err.BadRequestDefault("next group id is required")
		}
	}
	if !req.MoveToFirst && req.MoveToLast {
		if req.PreviousGroupId == 0 {
			return rpc_err.BadRequestDefault("previous group id is required")
		}
	}

	groupRecord, err := s.MySQL().FindProjectGroupById(s.Ctx(), req.UserId, req.ProjectId, req.GroupId)
	if err != nil {
		return err
	}
	if groupRecord == nil {
		return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project group is not accessible")
	}

	groupRecord, err = s.MySQL().MoveProjectGroup(s.Ctx(), req.UserId, req.ProjectId, req.GroupId, req.PreviousGroupId, req.NextGroupId, req.MoveToFirst, req.MoveToLast)
	if err != nil {
		return err
	}
	if groupRecord == nil {
		return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project group is not accessible")
	}

	resp.Group, err = mysql.ProjectGroupRecordToDTO(s.Ctx(), s.MinIO(), groupRecord, nil, s.Config().ProjectImageGetTTL)
	if err != nil {
		return err
	}

	return nil
}
