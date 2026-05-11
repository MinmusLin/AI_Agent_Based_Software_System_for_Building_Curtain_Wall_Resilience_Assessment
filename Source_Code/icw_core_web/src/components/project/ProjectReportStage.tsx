import { ReloadOutlined, RetweetOutlined } from '@ant-design/icons';
import { Button, Empty, message, Spin } from 'antd';
import type { ReactElement } from 'react';
import { useCallback, useEffect, useMemo, useState } from 'react';

import { getErrorMessage } from '@/api/http';
import { getProjectReport, retryProjectReport } from '@/api/project/report';
import type { Project } from '@/gen/core/api/common';
import type { ProjectReport } from '@/gen/core/common';
import { ProjectProgress_Value, ProjectReportStatus_Value } from '@/gen/core/common';
import { useProjectReportSocket } from '@/hooks/project/useProjectReportSocket';

const JSON_FORMAT_INDENT = 2;

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
  loading: boolean;
  onRefresh: () => void;
  onRetry: () => void;
  reportFailed: boolean;
  reportGenerating: boolean;
  retrying: boolean;
}

interface ReportBodyProps {
  content: string;
  errorText: string;
  reportFailed: boolean;
  reportLoading: boolean;
  reportSucceeded: boolean;
}

function reportContent(report?: ProjectReport): string {
  const resultJson = report?.result_json.trim() ?? '';
  if (resultJson === '') {
    return '';
  }
  try {
    const parsed = JSON.parse(resultJson) as unknown;
    if (typeof parsed === 'object' && parsed !== null && 'summary' in parsed && typeof parsed.summary === 'string') {
      return parsed.summary;
    }
    return JSON.stringify(parsed, null, JSON_FORMAT_INDENT);
  } catch {
    return resultJson;
  }
}

function ReportActions({
  loading,
  onRefresh,
  onRetry,
  reportFailed,
  reportGenerating,
  retrying,
}: ReportActionsProps): ReactElement | null {
  if (loading) {
    return null;
  }
  if (reportGenerating) {
    return (
      <Button icon={<ReloadOutlined />} onClick={onRefresh}>
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
  content,
  errorText,
  reportFailed,
  reportLoading,
  reportSucceeded,
}: ReportBodyProps): ReactElement | null {
  if (reportLoading) {
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
    return (
      <pre className="h-full overflow-auto overscroll-contain whitespace-pre-wrap break-all rounded bg-white p-4 text-sm leading-6 text-slate-700">
        {content || '报告生成成功'}
      </pre>
    );
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
  const [retrying, setRetrying] = useState(false);
  const [errorText, setErrorText] = useState('');

  const loadReport = useCallback(
    async (options: RefreshOptions = {}): Promise<void> => {
      if (projectId === '') {
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

  const content = useMemo(() => reportContent(report), [report]);
  const pageDataLoading = loading || reportLoading;
  const reportGenerating = !pageDataLoading && (!report || report.status === ProjectReportStatus_Value.Pending);
  const reportFailed = !pageDataLoading && report?.status === ProjectReportStatus_Value.Failed;
  const reportSucceeded = !pageDataLoading && report?.status === ProjectReportStatus_Value.Succeeded;
  const pageLoading = pageDataLoading || reportGenerating;

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
    <div className="flex min-h-0 flex-1 flex-col rounded-lg border border-slate-200 bg-white p-5">
      {contextHolder}
      <div className="mb-4 flex items-center justify-between gap-4">
        <div>
          <h2 className="text-base font-semibold text-slate-900">评估报告生成</h2>
          <p className="mt-1 text-sm text-slate-500">
            汇总建筑背景信息、Agent
            智能检测结果与人工复核意见，生成项目级建筑幕墙韧性评估报告，构建完整的现代智慧幕墙数字运维闭环
          </p>
        </div>
        <ReportActions
          loading={pageDataLoading}
          onRefresh={() => void loadReport()}
          onRetry={() => void handleRetry()}
          reportFailed={reportFailed}
          reportGenerating={reportGenerating}
          retrying={retrying}
        />
      </div>
      <div className="min-h-0 flex-1 overflow-hidden rounded-lg border border-slate-200 bg-slate-50 p-5">
        <ReportBody
          content={content}
          errorText={errorText}
          reportFailed={reportFailed}
          reportLoading={pageLoading}
          reportSucceeded={reportSucceeded}
        />
      </div>
    </div>
  );
}
