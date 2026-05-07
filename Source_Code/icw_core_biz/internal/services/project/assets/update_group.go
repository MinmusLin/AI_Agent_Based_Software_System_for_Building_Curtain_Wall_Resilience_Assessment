package assets

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/repositories/mysql"
)

// UpdateProjectGroup 更新图像组
func (s *Service) UpdateProjectGroup(ctx context.Context, req *bizpb.UpdateProjectGroupRequest) (*bizpb.UpdateProjectGroupResponse, error) {
	resp := &bizpb.UpdateProjectGroupResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.updateProjectGroup(req, resp)
	})
	return resp, err
}

func (s *Service) updateProjectGroup(req *bizpb.UpdateProjectGroupRequest, resp *bizpb.UpdateProjectGroupResponse) error {
	// 标准化项目图像组名称
	name, err := utils.NormalizeProjectGroupName(req.Name)
	if err != nil {
		return err
	}

	groupRecord, err := s.MySQL().UpdateProjectGroupName(s.Ctx(), req.UserId, req.ProjectId, req.GroupId, name)
	if mysql.IsDuplicateEntryError(err) {
		return rpc_error.BadRequest(rpc_error.DetailProjectGroupNameDuplicated, "project group name already exists")
	}
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
