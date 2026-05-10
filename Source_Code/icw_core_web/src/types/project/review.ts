import type { ProjectDetectionReviewVerdict } from '@/types/common';

export interface ProjectDetectionReview {
  comment?: string;
  image_uuid?: string;
  task_uuid?: string;
  updated_at?: string;
  verdict?: ProjectDetectionReviewVerdict;
}

export interface GetProjectDetectionReviewRequest {
  project_id: string;
  task_uuid: string;
}

export interface GetProjectDetectionReviewResponse {
  review?: ProjectDetectionReview;
}

export interface UpdateProjectDetectionReviewRequest {
  comment?: string;
  project_id: string;
  task_uuid: string;
  verdict?: ProjectDetectionReviewVerdict;
}

export interface UpdateProjectDetectionReviewResponse {
  review?: ProjectDetectionReview;
}
