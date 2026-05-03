package assets

import (
	"icw_core_biz/pkg/dto/project"
	"icw_core_biz/repositories/mysql"
)

// GetProjectAssets 获取项目图像列表
func (s *Service) GetProjectAssets(req *project.GetProjectAssetsRequest, resp *project.GetProjectAssetsResponse) error {
	return s.CallRPC("ProjectAssetsService.GetProjectAssets", req, resp, func() error {
		return s.getProjectAssets(req, resp)
	})
}

func (s *Service) getProjectAssets(req *project.GetProjectAssetsRequest, resp *project.GetProjectAssetsResponse) error {
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

	resp.Groups = make([]*project.ProjectGroup, 0, len(groups))
	for _, groupRecord := range groups {
		group, err := mysql.ProjectGroupRecordToDTO(s.Ctx(), s.MinIO(), groupRecord, imageMap[groupRecord.Id], s.Config().ProjectImageGetTTL)
		if err != nil {
			return err
		}
		resp.Groups = append(resp.Groups, group)
	}

	return nil
}
