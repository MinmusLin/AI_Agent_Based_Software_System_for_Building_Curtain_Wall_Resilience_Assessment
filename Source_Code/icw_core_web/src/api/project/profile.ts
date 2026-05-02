import { http } from '@/api/http';
import type { ApiEnvelope } from '@/types/common';
import type {
  GetProjectProfileResponse,
  GetProjectThumbnailResponse,
  UpdateProjectProfileRequest,
  UpdateProjectProfileResponse,
  UploadProjectThumbnailResponse,
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

// 获取项目缩略图
// @router /project/profile/thumbnail [GET]
export async function getProjectThumbnail(projectId: string): Promise<GetProjectThumbnailResponse> {
  const { data } = await http.get<ApiEnvelope<GetProjectThumbnailResponse>>('/project/profile/thumbnail', {
    params: {
      project_id: projectId,
    },
  });
  return data.data;
}

// 上传项目缩略图
// @router /project/profile/thumbnail [POST]
export async function uploadProjectThumbnail(projectId: string): Promise<UploadProjectThumbnailResponse> {
  const { data } = await http.post<ApiEnvelope<UploadProjectThumbnailResponse>>('/project/profile/thumbnail', {
    project_id: projectId,
  });
  return data.data;
}

// 删除项目缩略图
// @router /project/profile/thumbnail [DELETE]
export async function deleteProjectThumbnail(projectId: string): Promise<void> {
  await http.delete<ApiEnvelope<Record<string, never>>>('/project/profile/thumbnail', {
    params: {
      project_id: projectId,
    },
  });
}

// 更新项目基础信息
// @router /project/profile/update [POST]
export async function updateProjectProfile(
  payload: UpdateProjectProfileRequest,
): Promise<UpdateProjectProfileResponse> {
  const { data } = await http.post<ApiEnvelope<UpdateProjectProfileResponse>>('/project/profile/update', payload);
  return data.data;
}
