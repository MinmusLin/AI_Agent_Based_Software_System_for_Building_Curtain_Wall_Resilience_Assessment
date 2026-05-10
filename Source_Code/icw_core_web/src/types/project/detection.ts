import type {
  DetectionTaskCode,
  ProjectDetectionMainStatus,
  ProjectDetectionNodeCode,
  ProjectDetectionSubStatus,
  ProjectEventType,
} from '@/types/common';
import type { ProjectImage } from '@/types/project/assets';

export interface ProjectDetectionNodeStatus {
  node_code: ProjectDetectionNodeCode;
  sub_task_uuid: string;
  sub_status: ProjectDetectionSubStatus;
}

export interface ProjectDetectionStatus {
  classification_status?: ProjectDetectionNodeStatus;
  detection_status?: ProjectDetectionNodeStatus[];
  finished_at?: string;
  image_uuid: string;
  main_task_uuid: string;
  main_status: ProjectDetectionMainStatus;
  started_at?: string;
  summary_status?: ProjectDetectionNodeStatus;
}

export type ProjectDetectionArtifacts = Record<string, string>;

export interface ProjectDetectionCorrosionRegion {
  bbox_xyxy?: number[];
  confidence?: number;
  id?: number;
  mask_pixels?: number;
  mask_ratio?: number;
}

export interface ProjectDetectionCrackRegion {
  bbox_xyxy?: number[];
  id?: number;
  mask_pixels?: number;
  mask_ratio?: number;
}

export interface ProjectDetectionStainRegion {
  bbox_xyxy?: number[];
  confidence?: number;
  id?: number;
  region_height?: number;
  region_width?: number;
  stain_pixels?: number;
  stain_ratio?: number;
}

export interface ProjectDetectionFlatnessRegion {
  angle_std?: number;
  bbox_xyxy?: number[];
  edge_count?: number;
  edge_uneven_detected?: boolean;
  frequency_max?: number;
  frequency_min?: number;
  frequency_uneven_detected?: boolean;
  gradient_mean?: number;
  gradient_std?: number;
  gradient_uneven_detected?: boolean;
  id?: number;
  laplacian_variance?: number;
  line_count?: number;
  line_uneven_detected?: boolean;
}

export interface ProjectDetectionBaseResult {
  artifacts?: ProjectDetectionArtifacts;
  finished_at?: string;
  status?: ProjectDetectionSubStatus;
  started_at?: string;
  task_uuid?: string;
}

export interface ProjectDetectionCorrosionResult extends ProjectDetectionBaseResult {
  average_confidence?: number;
  corrosion_count?: number;
  corrosion_pixels?: number;
  corrosion_ratio?: number;
  has_corrosion?: boolean;
  max_confidence?: number;
  regions?: ProjectDetectionCorrosionRegion[];
  runtime_seconds?: number;
}

export interface ProjectDetectionCrackResult extends ProjectDetectionBaseResult {
  crack_count?: number;
  crack_pixels?: number;
  crack_ratio?: number;
  has_crack?: boolean;
  regions?: ProjectDetectionCrackRegion[];
  runtime_seconds?: number;
}

export interface ProjectDetectionStainResult extends ProjectDetectionBaseResult {
  average_stain_ratio?: number;
  has_stain?: boolean;
  max_stain_ratio?: number;
  regions?: ProjectDetectionStainRegion[];
  runtime_seconds?: number;
  stain_count?: number;
}

export interface ProjectDetectionFlatnessResult extends ProjectDetectionBaseResult {
  regions?: ProjectDetectionFlatnessRegion[];
  result?: string;
  runtime_seconds?: number;
  uneven_count?: number;
}

export interface ProjectDetectionSpallingResult extends ProjectDetectionBaseResult {
  confidence?: number;
  has_spalling?: boolean;
  runtime_seconds?: number;
}

export interface ProjectDetectionSummaryResult {
  finished_at?: string;
  result?: string;
  status?: ProjectDetectionSubStatus;
  started_at?: string;
  task_uuid?: string;
}

export interface GetImageDetectionResultResponse {
  corrosion_result?: ProjectDetectionCorrosionResult;
  crack_result?: ProjectDetectionCrackResult;
  flatness_result?: ProjectDetectionFlatnessResult;
  image?: ProjectImage;
  original_url?: string;
  spalling_result?: ProjectDetectionSpallingResult;
  status?: ProjectDetectionStatus;
  stain_result?: ProjectDetectionStainResult;
  summary_result?: ProjectDetectionSummaryResult;
  task_codes?: DetectionTaskCode[];
}

export interface GetImageDetectionResultRequest {
  image_uuid: string;
  project_id: string;
}

export interface GetProjectDetectionTasksResponse {
  tasks?: ProjectDetectionStatus[];
}

export interface StartProjectDetectionResponse {
  task_count?: number;
}

export interface RetryProjectDetectionResponse {
  task_count?: number;
}

export interface ProjectDetectionTaskStatusChangedMessage {
  type: ProjectEventType;
  project_id: string;
  image_uuid: string;
  node_code: ProjectDetectionNodeCode;
  main_task_uuid: string;
  main_status: ProjectDetectionMainStatus;
  sub_task_uuid: string;
  sub_status: ProjectDetectionSubStatus;
  occurred_at: string;
}
