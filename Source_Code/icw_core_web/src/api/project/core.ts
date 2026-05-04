import { http } from '@/api/http';
import type { ApiEnvelope } from '@/types/common';
import type {
  AdvanceProjectProgressRequest,
  CreateProjectResponse,
  DeleteProjectResponse,
  ListProjectsResponse,
} from '@/types/project/core';

// 项目进度流转
// @router /project/core/advance [POST]
export async function advanceProject(payload: AdvanceProjectProgressRequest): Promise<void> {
  await http.post<ApiEnvelope<Record<string, never>>>('/project/core/advance', payload);
}

// 创建项目
// @router /project/core/create [POST]
export async function createProject(): Promise<CreateProjectResponse> {
  const { data } = await http.post<ApiEnvelope<CreateProjectResponse>>('/project/core/create');
  return data.data;
}

// 删除项目
// @router /project/core/delete [POST]
export async function deleteProject(projectId: string): Promise<DeleteProjectResponse> {
  const { data } = await http.post<ApiEnvelope<DeleteProjectResponse>>('/project/core/delete', {
    project_id: projectId,
  });
  return data.data;
}

// 获取项目列表
// @router /project/core/list [GET]
export async function listProjects(): Promise<ListProjectsResponse> {
  const { data } = await http.get<ApiEnvelope<ListProjectsResponse>>('/project/core/list');
  return data.data;
}
