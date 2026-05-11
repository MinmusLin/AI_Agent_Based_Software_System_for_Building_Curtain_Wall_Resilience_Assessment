import { http } from '@/api/http';
import type { ApiEnvelope } from '@/constants/common';
import type {
  AdvanceProjectRequest,
  CreateProjectResponse,
  DeleteProjectResponse,
  GetProjectDashboardResponse,
  ListProjectsResponse,
} from '@/gen/core/api/project_core';

// 项目进度流转
// @router /project/core/advance [POST]
export async function advanceProject(payload: AdvanceProjectRequest): Promise<void> {
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

// 获取项目工作台统计
// @router /project/core/dashboard [GET]
export async function getProjectDashboard(): Promise<GetProjectDashboardResponse> {
  const { data } = await http.get<ApiEnvelope<GetProjectDashboardResponse>>('/project/core/dashboard');
  return data.data;
}

// 获取项目列表
// @router /project/core/list [GET]
export async function listProjects(): Promise<ListProjectsResponse> {
  const { data } = await http.get<ApiEnvelope<ListProjectsResponse>>('/project/core/list');
  return data.data;
}
