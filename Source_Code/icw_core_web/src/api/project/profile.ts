import { http } from '@/api/http';
import type { ApiEnvelope } from '@/types/common';
import type {
  GetProjectProfileResponse,
  UpdateProjectProfileRequest,
  UpdateProjectProfileResponse,
} from '@/types/project/profile';

// 获取项目基础信息
// @router /project/profile/detail [GET]
export async function getProjectProfile(projectId: string): Promise<GetProjectProfileResponse> {
  const { data } = await http.get<ApiEnvelope<GetProjectProfileResponse>>('/project/profile/detail', {
    params: {
      project_id: projectId,
    },
  });
  return data.data;
}

// 更新项目基础信息
// @router /project/profile/update [POST]
export async function updateProjectProfile(
  payload: UpdateProjectProfileRequest,
): Promise<UpdateProjectProfileResponse> {
  const { data } = await http.post<ApiEnvelope<UpdateProjectProfileResponse>>('/project/profile/update', payload);
  return data.data;
}
