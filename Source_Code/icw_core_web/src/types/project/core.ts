export interface Project {
  id: string;
  name: string;
  building_name: string;
  building_location: string;
  built_year: number;
  building_description: string;
  known_issues: string;
  assessment_goal: string;
  thumbnail_url: string;
  progress: number;
  created_at: string;
  updated_at: string;
}

export interface ProjectListItem {
  id: string;
  name: string;
  building_name: string;
  building_location: string;
  thumbnail_url: string;
  progress: number;
  created_at: string;
}

export interface AdvanceProjectProgressRequest {
  project_id: string;
  from_progress: number;
  to_progress: number;
}

export interface CreateProjectResponse {
  project: Project;
}

export interface DeleteProjectRequest {
  project_id: string;
}

export interface DeleteProjectResponse {
  active_projects: ProjectListItem[];
  completed_projects: ProjectListItem[];
}

export interface ListProjectsResponse {
  active_projects: ProjectListItem[];
  completed_projects: ProjectListItem[];
}
