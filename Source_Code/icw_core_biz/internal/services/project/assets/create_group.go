package assets

import (
	"fmt"
	"time"

	"icw_core_biz/internal/services/project/consts"
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/pkg/rpc_err"
	"icw_core_biz/repositories/mysql"
)

const (
	// CreateProjectGroupTimeout 创建图像组超时配置
	CreateProjectGroupTimeout = 3 * time.Second
)

// CreateProjectGroup 创建图像组
func (s *Service) CreateProjectGroup(req *project.CreateProjectGroupRequest, resp *project.CreateProjectGroupResponse) error {
	return s.CallRPC("ProjectAssetsService.CreateProjectGroup", req, resp, func() error {
		return s.createProjectGroup(req, resp)
	})
}

func (s *Service) createProjectGroup(req *project.CreateProjectGroupRequest, resp *project.CreateProjectGroupResponse) error {
	start := time.Now()

	for {
		if time.Since(start) > CreateProjectGroupTimeout {
			return rpc_err.InternalErrorDefault("create project group timeout")
		}

		// 分配新图像组的下一个序号
		seq, err := s.Redis().NextProjectGroupSequence(s.Ctx(), req.ProjectId)
		if err != nil {
			return err
		}
		if seq <= 0 {
			return rpc_err.InternalErrorDefault("create project group sequence error")
		}

		// 生成新图像组名称
		name := consts.DefaultNewProjectGroupName
		if seq > 1 {
			name = fmt.Sprintf("%s %d", name, seq)
		}

		groupRecord, err := s.MySQL().CreateProjectGroup(s.Ctx(), req.UserId, req.ProjectId, name)
		if mysql.IsDuplicateEntryError(err) {
			continue
		}
		if err != nil {
			return err
		}
		if groupRecord == nil {
			return rpc_err.InternalErrorDefault("create project group failed")
		}

		resp.Group, err = mysql.ProjectGroupRecordToDTO(s.Ctx(), s.MinIO(), groupRecord, nil, s.Config().ProjectThumbnailGetTTL)
		if err != nil {
			return err
		}

		return nil
	}
}
