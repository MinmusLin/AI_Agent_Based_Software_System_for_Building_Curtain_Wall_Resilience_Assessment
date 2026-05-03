package assets

import (
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/mysql"
)

// UpdateProjectGroup 更新图像组
func (s *Service) UpdateProjectGroup(req *project.UpdateProjectGroupRequest, resp *project.UpdateProjectGroupResponse) error {
	return s.CallRPC("ProjectAssetsService.UpdateProjectGroup", req, resp, func() error {
		return s.updateProjectGroup(req, resp)
	})
}

func (s *Service) updateProjectGroup(req *project.UpdateProjectGroupRequest, resp *project.UpdateProjectGroupResponse) error {
	// 标准化项目图像组名称
	name, err := utils.NormalizeProjectGroupName(req.Name)
	if err != nil {
		return err
	}

	groupRecord, err := s.MySQL().UpdateProjectGroupName(s.Ctx(), req.UserId, req.ProjectId, req.GroupId, name)
	if mysql.IsDuplicateEntryError(err) {
		return rpc_err.BadRequest(rpc_err.DetailProjectGroupNameDuplicated, "project group name already exists")
	}
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
