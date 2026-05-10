import type { ProjectReportStatus } from '@/types/common';

export interface ProjectReportStatusChangedMessage {
  occurred_at?: string;
  project_id?: string;
  report_uuid?: string;
  status?: ProjectReportStatus;
  type?: string;
}

export interface ProjectReport {
  report_uuid?: string;
  result_json?: string;
  status?: ProjectReportStatus;
  updated_at?: string;
}

export interface GetProjectReportResponse {
  report?: ProjectReport;
}
