import type { ProjectDetectionMainStatus, ProjectDetectionSubStatus, ProjectEventType } from '@/types/common';
import type { ProjectImage } from '@/types/project/assets';

export interface CreateSocketTicketResponse {
  ticket: string;
  expires_in: number;
}

export interface ProjectImageStatusChangedMessage {
  type: ProjectEventType;
  project_id: string;
  image: ProjectImage;
}

export interface ProjectDetectionTaskStatusChangedMessage {
  type: ProjectEventType;
  project_id: string;
  image_uuid: string;
  node_code: string;
  main_task_id: string;
  main_status: ProjectDetectionMainStatus;
  sub_task_id: string;
  sub_status: ProjectDetectionSubStatus;
  occurred_at: string;
}
