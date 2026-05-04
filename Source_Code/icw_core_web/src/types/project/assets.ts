export const PROJECT_IMAGE_STATUS_PENDING = 'pending';
export const PROJECT_IMAGE_STATUS_UPLOADED = 'uploaded';
export const PROJECT_IMAGE_STATUS_FAILED = 'failed';

export type ProjectImageStatus =
  | typeof PROJECT_IMAGE_STATUS_PENDING
  | typeof PROJECT_IMAGE_STATUS_UPLOADED
  | typeof PROJECT_IMAGE_STATUS_FAILED;

export interface ProjectGroup {
  id: string;
  name: string;
  sort_order: string;
  images: ProjectImage[];
}

export interface ProjectImage {
  uuid: string;
  file_name: string;
  content_type: string;
  size_bytes: number;
  width: number;
  height: number;
  metadata: string;
  status: ProjectImageStatus;
  thumbnail_url: string;
  uploaded_at: string;
  created_at: string;
}

export interface UploadProjectImageItem {
  file_name: string;
  content_type: string;
  size_bytes: number;
  width: number;
  height: number;
  metadata: string;
}

export interface UploadProjectImageResult {
  image: ProjectImage;
  original_upload_url: string;
  thumbnail_upload_url: string;
}

export interface GetProjectAssetsResponse {
  groups: ProjectGroup[];
}

export interface CreateProjectGroupResponse {
  group: ProjectGroup;
}

export interface DeleteProjectGroupRequest {
  project_id: string;
  group_id: string;
}

export type DeleteProjectGroupResponse = Record<string, never>;

export interface MoveProjectGroupRequest {
  project_id: string;
  group_id: string;
  previous_group_id: string;
  next_group_id: string;
  move_to_first: boolean;
  move_to_last: boolean;
}

export interface MoveProjectGroupResponse {
  group: ProjectGroup;
}

export interface UpdateProjectGroupRequest {
  project_id: string;
  group_id: string;
  name: string;
}

export interface UpdateProjectGroupResponse {
  group: ProjectGroup;
}

export interface DeleteProjectImageRequest {
  project_id: string;
  image_uuids: string[];
}

export type DeleteProjectImageResponse = Record<string, never>;

export interface GetProjectImageOriginalResponse {
  original_url: string;
}

export interface MoveProjectImageRequest {
  project_id: string;
  image_uuids: string[];
  target_group_id: string;
}

export interface MoveProjectImageResponse {
  images: ProjectImage[];
}

export interface ReportProjectImageRequest {
  project_id: string;
  image_uuid: string;
  status: ProjectImageStatus;
}

export interface ReportProjectImageResponse {
  image: ProjectImage;
}

export interface UploadProjectImageRequest {
  project_id: string;
  group_id: string;
  images: UploadProjectImageItem[];
}

export interface UploadProjectImageResponse {
  images: UploadProjectImageResult[];
}
