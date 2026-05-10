import { http } from '@/api/http';
import type { ApiEnvelope } from '@/constants/common';
import type {
  CreateProjectGroupResponse,
  DeleteProjectGroupRequest,
  DeleteProjectImageRequest,
  GetProjectAssetsResponse,
  GetProjectImageOriginalRequest,
  GetProjectImageOriginalResponse,
  MoveProjectGroupRequest,
  MoveProjectGroupResponse,
  MoveProjectImageRequest,
  MoveProjectImageResponse,
  ReportProjectImageRequest,
  UpdateProjectGroupRequest,
  UpdateProjectGroupResponse,
  UploadProjectImageRequest,
  UploadProjectImageResponse,
} from '@/gen/core/api/project_assets';

// 获取项目图像列表
// @router /project/assets/list [GET]
export async function getProjectAssets(projectId: string): Promise<GetProjectAssetsResponse> {
  const { data } = await http.get<ApiEnvelope<GetProjectAssetsResponse>>('/project/assets/list', {
    params: {
      project_id: projectId,
    },
  });
  return data.data;
}

// 创建图像组
// @router /project/assets/group/create [POST]
export async function createProjectGroup(projectId: string): Promise<CreateProjectGroupResponse> {
  const { data } = await http.post<ApiEnvelope<CreateProjectGroupResponse>>('/project/assets/group/create', {
    project_id: projectId,
  });
  return data.data;
}

// 删除图像组
// @router /project/assets/group/delete [POST]
export async function deleteProjectGroup(payload: DeleteProjectGroupRequest): Promise<Record<string, never>> {
  const { data } = await http.post<ApiEnvelope<Record<string, never>>>('/project/assets/group/delete', payload);
  return data.data;
}

// 移动图像组
// @router /project/assets/group/move [POST]
export async function moveProjectGroup(payload: MoveProjectGroupRequest): Promise<MoveProjectGroupResponse> {
  const { data } = await http.post<ApiEnvelope<MoveProjectGroupResponse>>('/project/assets/group/move', payload);
  return data.data;
}

// 更新图像组
// @router /project/assets/group/update [POST]
export async function updateProjectGroup(payload: UpdateProjectGroupRequest): Promise<UpdateProjectGroupResponse> {
  const { data } = await http.post<ApiEnvelope<UpdateProjectGroupResponse>>('/project/assets/group/update', payload);
  return data.data;
}

// 删除图像
// @router /project/assets/image/delete [POST]
export async function deleteProjectImage(payload: DeleteProjectImageRequest): Promise<Record<string, never>> {
  const { data } = await http.post<ApiEnvelope<Record<string, never>>>('/project/assets/image/delete', payload);
  return data.data;
}

// 获取原图
// @router /project/assets/image/original [GET]
export async function getProjectImageOriginal(
  payload: GetProjectImageOriginalRequest,
): Promise<GetProjectImageOriginalResponse> {
  const { data } = await http.get<ApiEnvelope<GetProjectImageOriginalResponse>>('/project/assets/image/original', {
    params: {
      image_uuid: payload.image_uuid,
      project_id: payload.project_id,
    },
  });
  return data.data;
}

// 移动图像
// @router /project/assets/image/move [POST]
export async function moveProjectImage(payload: MoveProjectImageRequest): Promise<MoveProjectImageResponse> {
  const { data } = await http.post<ApiEnvelope<MoveProjectImageResponse>>('/project/assets/image/move', payload);
  return data.data;
}

// 上报图像
// @router /project/assets/image/report [POST]
export async function reportProjectImage(payload: ReportProjectImageRequest): Promise<void> {
  await http.post<ApiEnvelope<Record<string, never>>>('/project/assets/image/report', payload);
}

// 上传图像
// @router /project/assets/image/upload [POST]
export async function uploadProjectImage(payload: UploadProjectImageRequest): Promise<UploadProjectImageResponse> {
  const { data } = await http.post<ApiEnvelope<UploadProjectImageResponse>>('/project/assets/image/upload', payload);
  return data.data;
}
