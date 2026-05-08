package assets

import (
	"context"
	"errors"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/internal/services/common"
	"icw_core_biz/internal/services/project/utils"
	"icw_core_biz/repositories/mysql/model"
)

// DeleteProjectGroup 删除图像组
func (s *Service) DeleteProjectGroup(ctx context.Context, req *bizpb.DeleteProjectGroupRequest) (*bizpb.DeleteProjectGroupResponse, error) {
	resp := &bizpb.DeleteProjectGroupResponse{}
	err := s.CallRPC(req, func() error {
		return s.deleteProjectGroup(req, resp)
	})
	return resp, err
}

func (s *Service) deleteProjectGroup(req *bizpb.DeleteProjectGroupRequest, _ *bizpb.DeleteProjectGroupResponse) error {
	images, err := s.MySQL().ListProjectImagesByGroupId(s.Ctx(), req.UserId, req.ProjectId, req.GroupId)
	if err != nil {
		return err
	}

	deleted, err := s.MySQL().DeleteProjectGroup(s.Ctx(), req.UserId, req.ProjectId, req.GroupId)
	if errors.Is(err, model.ErrProjectGroupCannotDeleteLast) {
		return rpc_error.BadRequest(rpc_error.DetailProjectAtLeastOneGroupRequired, "project must keep at least one group")
	}
	if err != nil {
		return err
	}
	if !deleted {
		return rpc_error.BadRequest(rpc_error.DetailProjectNotAccessible, "project group is not accessible")
	}

	for _, imageRecord := range images {
		if imageRecord == nil {
			continue
		}
		if err := utils.RemoveProjectImageObjects(s.Ctx(), s.MinIO(), s.Redis(), req.UserId, req.ProjectId, imageRecord.Uuid); err != nil {
			common.RpcWarn("Remove project image objects failed, project_id: %d, image_uuid: %s, err: %v", req.ProjectId, imageRecord.Uuid, err)
		}
	}

	return nil
}
