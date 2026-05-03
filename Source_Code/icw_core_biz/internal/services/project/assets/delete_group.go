package assets

import (
	"errors"
	"log"

	"icw_core_biz/internal/services/common"
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/mysql"
)

// DeleteProjectGroup 删除图像组
func (s *Service) DeleteProjectGroup(req *project.DeleteProjectGroupRequest, resp *project.DeleteProjectGroupResponse) error {
	return s.CallRPC("ProjectAssetsService.DeleteProjectGroup", req, resp, func() error {
		return s.deleteProjectGroup(req, resp)
	})
}

func (s *Service) deleteProjectGroup(req *project.DeleteProjectGroupRequest, _ *project.DeleteProjectGroupResponse) error {
	images, err := s.MySQL().ListProjectImagesByGroupId(s.Ctx(), req.UserId, req.ProjectId, req.GroupId)
	if err != nil {
		return err
	}

	deleted, err := s.MySQL().DeleteProjectGroup(s.Ctx(), req.UserId, req.ProjectId, req.GroupId)
	if errors.Is(err, mysql.ErrProjectGroupCannotDeleteLast) {
		return rpc_err.BadRequest(rpc_err.DetailProjectAtLeastOneGroupRequired, "project must keep at least one group")
	}
	if err != nil {
		return err
	}
	if !deleted {
		return rpc_err.BadRequest(rpc_err.DetailProjectNotAccessible, "project group is not accessible")
	}

	for _, imageRecord := range images {
		if imageRecord == nil {
			continue
		}
		if err := utils.RemoveProjectImageObjects(s.Ctx(), s.MinIO(), s.Redis(), req.UserId, req.ProjectId, imageRecord.Uuid); err != nil {
			log.Printf("%s Remove project image objects failed, project_id: %d, image_uuid: %s, err: %v", common.WarnPrefix(), req.ProjectId, imageRecord.Uuid, err)
		}
	}

	return nil
}
