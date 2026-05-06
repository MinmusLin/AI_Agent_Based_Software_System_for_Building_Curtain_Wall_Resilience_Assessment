package assets

import (
	"context"

	"icw_common/gen/core/biz"
	"icw_core_biz/repositories/mysql"
)

// GetProjectAssets 获取项目图像列表
func (s *Service) GetProjectAssets(ctx context.Context, req *bizpb.GetProjectAssetsRequest) (*bizpb.GetProjectAssetsResponse, error) {
	resp := &bizpb.GetProjectAssetsResponse{}
	err := s.CallRPC(ctx, req, resp, func() error {
		return s.getProjectAssets(req, resp)
	})
	return resp, err
}

func (s *Service) getProjectAssets(req *bizpb.GetProjectAssetsRequest, resp *bizpb.GetProjectAssetsResponse) error {
	groups, err := s.MySQL().ListProjectGroups(s.Ctx(), req.UserId, req.ProjectId)
	if err != nil {
		return err
	}
	images, err := s.MySQL().ListProjectImages(s.Ctx(), req.UserId, req.ProjectId)
	if err != nil {
		return err
	}

	imageMap := make(map[uint64][]*mysql.ProjectImageRecord, len(groups))
	for _, imageRecord := range images {
		if imageRecord == nil {
			continue
		}
		imageMap[imageRecord.GroupId] = append(imageMap[imageRecord.GroupId], imageRecord)
	}

	resp.Groups = make([]*bizpb.ProjectGroup, 0, len(groups))
	for _, groupRecord := range groups {
		group, err := mysql.ProjectGroupRecordToDTO(s.Ctx(), s.MinIO(), s.Redis(), groupRecord, imageMap[groupRecord.Id], s.Config().ProjectImageGetTTL)
		if err != nil {
			return err
		}
		resp.Groups = append(resp.Groups, group)
	}

	return nil
}
