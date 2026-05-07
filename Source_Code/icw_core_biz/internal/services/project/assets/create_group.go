package assets

import (
	"context"
	"errors"
	"fmt"
	"time"

	"icw_common/gen/core/biz"
	"icw_common/rpc/error"

	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/repositories/mysql/model"
	mysqlUtils "icw_core_biz/repositories/mysql/utils"
)

const (
	// CreateProjectGroupTimeout 创建图像组超时配置
	CreateProjectGroupTimeout = 3 * time.Second
)

// CreateProjectGroup 创建图像组
func (s *Service) CreateProjectGroup(ctx context.Context, req *bizpb.CreateProjectGroupRequest) (*bizpb.CreateProjectGroupResponse, error) {
	resp := &bizpb.CreateProjectGroupResponse{}
	err := s.CallRPC(ctx, req, func() error {
		return s.createProjectGroup(req, resp)
	})
	return resp, err
}

func (s *Service) createProjectGroup(req *bizpb.CreateProjectGroupRequest, resp *bizpb.CreateProjectGroupResponse) error {
	start := time.Now()

	for {
		if time.Since(start) > CreateProjectGroupTimeout {
			return errors.New("create project group timeout")
		}

		// 分配新图像组的下一个序号
		seq, err := s.Redis().NextProjectGroupSequence(s.Ctx(), req.ProjectId)
		if err != nil {
			return err
		}
		if seq <= 0 {
			return errors.New("create project group sequence error")
		}

		// 生成新图像组名称
		name := consts.DefaultNewProjectGroupName
		if seq > 1 {
			name = fmt.Sprintf("%s %d", name, seq)
		}

		groupRecord, err := s.MySQL().CreateProjectGroup(s.Ctx(), req.UserId, req.ProjectId, name)
		if mysqlUtils.IsDuplicateEntryError(err) {
			continue
		}
		if err != nil {
			return err
		}
		if groupRecord == nil {
			return rpc_error.BadRequest(rpc_error.DetailProjectNotAccessible, "project group is not accessible")
		}

		resp.Group, err = model.ProjectGroupRecordToDTO(s.Ctx(), s.MinIO(), s.Redis(), groupRecord, nil, s.Config().ProjectImageGetTTL)
		if err != nil {
			return err
		}

		return nil
	}
}
