import type { ProjectEventType } from '@/types/common';
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
