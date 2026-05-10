import { http } from '@/api/http';
import type { ApiEnvelope } from '@/constants/common';
import type {
  GetProjectDetectionReviewRequest,
  GetProjectDetectionReviewResponse,
  UpdateProjectDetectionReviewRequest,
  UpdateProjectDetectionReviewResponse,
} from '@/gen/core/api/project_review';

// 获取图像检测人工复核信息
// @router /project/review/detail [GET]
export async function getProjectDetectionReview(
  request: GetProjectDetectionReviewRequest,
): Promise<GetProjectDetectionReviewResponse> {
  const { data } = await http.get<ApiEnvelope<GetProjectDetectionReviewResponse>>('/project/review/detail', {
    params: {
      project_id: request.project_id,
      task_uuid: request.task_uuid,
    },
  });
  return data.data;
}

// 更新图像检测人工复核信息
// @router /project/review/update [POST]
export async function updateProjectDetectionReview(
  request: UpdateProjectDetectionReviewRequest,
): Promise<UpdateProjectDetectionReviewResponse> {
  const { data } = await http.post<ApiEnvelope<UpdateProjectDetectionReviewResponse>>(
    '/project/review/update',
    request,
  );
  return data.data;
}
