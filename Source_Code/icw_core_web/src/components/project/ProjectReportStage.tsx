import { ReloadOutlined, RetweetOutlined } from '@ant-design/icons';
import { Button, Empty, message, Spin } from 'antd';
import type { ReactElement } from 'react';
import { useCallback, useEffect, useMemo, useState } from 'react';

import { getErrorMessage } from '@/api/http';
import { getProjectReport, retryProjectReport } from '@/api/project/report';
import type { ProjectReportPayload } from '@/components/project/report/ProjectReportViewer';
import { parseProjectReportPayload, ProjectReportViewer } from '@/components/project/report/ProjectReportViewer';
import type { Project } from '@/gen/core/api/common';
import type { ProjectReport } from '@/gen/core/common';
import { ProjectProgress_Value, ProjectReportStatus_Value } from '@/gen/core/common';
import { useProjectReportSocket } from '@/hooks/project/useProjectReportSocket';

interface ProjectReportStageProps {
  loading?: boolean;
  onProgressChange: (progress: ProjectProgress_Value) => void;
  onProjectChange: (project: Project) => void;
  project: Project;
  projectId: string;
}

interface RefreshOptions {
  silent?: boolean;
}

interface ReportActionsProps {
  pageLoading: boolean;
  onRefresh: () => void;
  onRetry: () => void;
  reportFailed: boolean;
  reportGenerating: boolean;
  refreshing: boolean;
  retrying: boolean;
}

interface ReportBodyProps {
  errorText: string;
  fallbackText: string;
  payload?: ProjectReportPayload;
  reportGenerating: boolean;
  reportFailed: boolean;
  reportDataLoading: boolean;
  reportSucceeded: boolean;
  updatedAt?: string;
}

function reportText(report?: ProjectReport): string {
  return report?.result_json.trim() ?? '';
}

function ReportActions({
  onRefresh,
  onRetry,
  pageLoading,
  reportFailed,
  reportGenerating,
  refreshing,
  retrying,
}: ReportActionsProps): ReactElement | null {
  if (pageLoading) {
    return null;
  }
  if (reportGenerating) {
    return (
      <Button disabled={refreshing} icon={<ReloadOutlined />} loading={refreshing} onClick={onRefresh}>
        刷新
      </Button>
    );
  }
  if (reportFailed) {
    return (
      <Button danger icon={<RetweetOutlined />} loading={retrying} onClick={onRetry}>
        失败重试
      </Button>
    );
  }
  return null;
}

function ReportBody({
  errorText,
  fallbackText,
  payload,
  reportGenerating,
  reportFailed,
  reportDataLoading,
  reportSucceeded,
  updatedAt,
}: ReportBodyProps): ReactElement | null {
  if (reportDataLoading) {
    return (
      <div className="flex h-full flex-col items-center justify-center gap-3 text-sm text-slate-500">
        <Spin />
        <span>报告加载中</span>
      </div>
    );
  }
  if (reportGenerating) {
    return (
      <div className="flex h-full flex-col items-center justify-center gap-3 text-sm text-slate-500">
        <Spin />
        <span>报告正在生成中</span>
      </div>
    );
  }
  if (reportFailed) {
    return (
      <div className="flex h-full items-center justify-center">
        <Empty description={errorText || '报告生成失败，请稍后重试'} />
      </div>
    );
  }
  if (reportSucceeded) {
    return <ProjectReportViewer fallbackText={fallbackText} payload={payload} updatedAt={updatedAt} />;
  }
  return null;
}

export function ProjectReportStage({
  loading = false,
  onProgressChange,
  onProjectChange,
  project,
  projectId,
}: ProjectReportStageProps): ReactElement {
  const [messageApi, contextHolder] = message.useMessage();
  const [report, setReport] = useState<ProjectReport | undefined>();
  const [reportLoading, setReportLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [retrying, setRetrying] = useState(false);
  const [errorText, setErrorText] = useState('');

  const loadReport = useCallback(
    async (options: RefreshOptions = {}): Promise<void> => {
      if (projectId === '') {
        if (!options.silent) {
          setReportLoading(false);
        }
        return;
      }
      if (!options.silent) {
        setReportLoading(true);
        setErrorText('');
      }
      try {
        const data = await getProjectReport(projectId);
        setReport(data.report);
      } catch (error: unknown) {
        if (!options.silent) {
          setErrorText(getErrorMessage(error));
        }
      } finally {
        if (!options.silent) {
          setReportLoading(false);
        }
      }
    },
    [projectId],
  );

  const handleRefresh = useCallback(async (): Promise<void> => {
    if (projectId === '' || refreshing) {
      return;
    }
    setRefreshing(true);
    try {
      const data = await getProjectReport(projectId);
      setReport(data.report);
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
    } finally {
      setRefreshing(false);
    }
  }, [messageApi, projectId, refreshing]);

  useEffect(() => {
    if (loading) {
      return;
    }
    void loadReport();
  }, [loadReport, loading]);

  const handleReportChanged = useCallback((): void => {
    void loadReport({ silent: true });
  }, [loadReport]);

  const handleRetry = useCallback(async (): Promise<void> => {
    if (projectId === '') {
      return;
    }
    setRetrying(true);
    setErrorText('');
    try {
      await retryProjectReport({
        project_id: projectId,
      });
      await loadReport({ silent: true });
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
    } finally {
      setRetrying(false);
    }
  }, [loadReport, messageApi, projectId]);

  const fallbackText = useMemo(() => reportText(report), [report]);
  const payload = useMemo(() => parseProjectReportPayload(fallbackText), [fallbackText]);
  const pageDataLoading = loading || reportLoading;
  const reportGenerating = !pageDataLoading && (!report || report.status === ProjectReportStatus_Value.Pending);
  const reportFailed = !pageDataLoading && report?.status === ProjectReportStatus_Value.Failed;
  const reportSucceeded = !pageDataLoading && report?.status === ProjectReportStatus_Value.Succeeded;

  useEffect(() => {
    if (!reportSucceeded || project.progress >= ProjectProgress_Value.ReportFinished) {
      return;
    }
    onProjectChange({
      ...project,
      progress: ProjectProgress_Value.ReportFinished,
    });
    onProgressChange(ProjectProgress_Value.ReportFinished);
  }, [onProgressChange, onProjectChange, project, reportSucceeded]);

  useProjectReportSocket({
    enabled: !loading,
    onReportChanged: handleReportChanged,
    projectId,
  });

  return (
    <div className="flex min-h-0 flex-1 flex-col overflow-hidden rounded-lg border border-slate-200 bg-white">
      {contextHolder}
      <div className="flex items-center justify-between gap-4 px-5 py-5">
        <div>
          <h2 className="text-base font-semibold text-slate-900">评估报告生成</h2>
          <p className="mt-1 text-sm text-slate-500">
            汇总建筑背景信息、Agent
            智能检测结果与人工复核意见，生成项目级建筑幕墙韧性评估报告，构建完整的现代智慧幕墙数字运维闭环
          </p>
        </div>
        <ReportActions
          onRefresh={() => void handleRefresh()}
          onRetry={() => void handleRetry()}
          pageLoading={pageDataLoading}
          refreshing={refreshing}
          reportFailed={reportFailed}
          reportGenerating={reportGenerating}
          retrying={retrying}
        />
      </div>
      <div className="min-h-0 flex-1 overflow-hidden px-5 pb-5">
        <ReportBody
          errorText={errorText}
          fallbackText={fallbackText}
          payload={payload}
          reportDataLoading={pageDataLoading}
          reportFailed={reportFailed}
          reportGenerating={reportGenerating}
          reportSucceeded={reportSucceeded}
          updatedAt={report?.updated_at}
        />
      </div>
    </div>
  );
}
