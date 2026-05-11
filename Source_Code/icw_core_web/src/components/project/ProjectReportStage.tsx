import { ReloadOutlined } from '@ant-design/icons';
import { Button, Empty, Spin } from 'antd';
import type { ReactElement } from 'react';
import { useCallback, useEffect, useMemo, useState } from 'react';

import { getErrorMessage } from '@/api/http';
import { getProjectReport } from '@/api/project/report';
import type { Project } from '@/gen/core/api/common';
import type { ProjectReport } from '@/gen/core/common';
import { ProjectReportStatus_Value } from '@/gen/core/common';
import { useProjectReportSocket } from '@/hooks/project/useProjectReportSocket';

const JSON_FORMAT_INDENT = 2;

interface ProjectReportStageProps {
  loading?: boolean;
  project: Project;
  projectId: string;
}

interface RefreshOptions {
  silent?: boolean;
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

export function ProjectReportStage({ loading = false, projectId }: ProjectReportStageProps): ReactElement {
  const [report, setReport] = useState<ProjectReport | undefined>();
  const [reportLoading, setReportLoading] = useState(true);
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

  const content = useMemo(() => reportContent(report), [report]);
  const pageLoading = loading || reportLoading || report?.status === ProjectReportStatus_Value.Pending || !report;

  useProjectReportSocket({
    enabled: !loading,
    onReportChanged: handleReportChanged,
    projectId,
  });

  return (
    <div className="flex min-h-0 flex-1 flex-col rounded-lg border border-slate-200 bg-white p-5">
      <div className="mb-4 flex items-center justify-between gap-4">
        <div>
          <h2 className="text-base font-semibold text-slate-900">评估报告生成</h2>
          <p className="mt-1 text-sm text-slate-500">
            汇总项目基础信息、图像检测结果与人工复核意见，生成项目级韧性评估报告。
          </p>
        </div>
        <Button disabled={loading || reportLoading} icon={<ReloadOutlined />} onClick={() => void loadReport()}>
          刷新
        </Button>
      </div>
      <div className="min-h-0 flex-1 overflow-hidden rounded-lg border border-slate-200 bg-slate-50 p-5">
        {pageLoading ? (
          <div className="flex h-full flex-col items-center justify-center gap-3 text-sm text-slate-500">
            <Spin />
            <span>报告正在生成中</span>
          </div>
        ) : null}
        {!pageLoading && report.status === ProjectReportStatus_Value.Failed ? (
          <div className="flex h-full items-center justify-center">
            <Empty description={errorText || '报告生成失败，请稍后重试'} />
          </div>
        ) : null}
        {!pageLoading && report.status === ProjectReportStatus_Value.Succeeded ? (
          <pre className="h-full overflow-auto overscroll-contain whitespace-pre-wrap break-all rounded bg-white p-4 text-sm leading-6 text-slate-700">
            {content || '报告已生成'}
          </pre>
        ) : null}
      </div>
    </div>
  );
}
