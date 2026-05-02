import type { Project } from './core';

export interface GetProjectProfileResponse {
  project: Project;
}

export interface UploadProjectThumbnailResponse {
  upload_url: string;
}

export interface UpdateProjectProfileRequest {
  project_id: string;
  name: string;
  building_name: string;
  building_location: string;
  built_year: number;
  building_description: string;
  known_issues: string;
  assessment_goal: string;
}

export interface UpdateProjectProfileResponse {
  project: Project;
}
