import {
  BarChartOutlined,
  CheckCircleOutlined,
  CloudServerOutlined,
  DatabaseOutlined,
  FileImageOutlined,
  FolderOpenOutlined,
  ProjectOutlined,
  RocketOutlined,
} from '@ant-design/icons';
import { Card, Empty, message } from 'antd';
import type { CSSProperties, ReactElement, ReactNode } from 'react';
import { useCallback, useEffect, useMemo, useState } from 'react';

import { getErrorMessage } from '@/api/http';
import { getProjectDashboard } from '@/api/project/core';
import type { GetProjectDashboardResponse } from '@/gen/core/api/project_core';

const BAR_MIN_PERCENT = 4;
const BYTE_UNIT_BASE = 1024;
const BYTE_UNIT_START_INDEX = 0;
const BYTE_UNIT_LAST_INDEX_OFFSET = 1;
const DISPLAY_FRACTION_DIGITS = 2;
const MINIO_REMAINING_COLOR = '#22C55E';
const MINIO_USED_COLOR = '#2563EB';
const MIN_VISIBLE_PIE_PERCENT = 0.8;
const PERCENT_FULL = 100;
const PERCENT_TINY_THRESHOLD = 0.01;
const PIE_CHART_SIZE = 'clamp(10rem, min(30vh, 22vw), 20rem)';
const TOP_CARD_SKELETON_COUNT = 6;
const CHART_SKELETON_COUNT = 8;
const UNIT_INDEX_INCREMENT = 1;
const ZERO_VALUE = 0;
const BYTE_UNITS = ['B', 'KB', 'MB', 'GB', 'TB'];
const NUMBER_FORMATTER = new Intl.NumberFormat('zh-CN');

interface TopMetric {
  icon: ReactNode;
  label: string;
  value: number;
}

interface BarMetric {
  label: string;
  value: number;
}

interface StorageSegment {
  color: string;
  label: string;
  value: number;
}

interface ChartTitleProps {
  icon: ReactNode;
  title: string;
}

function formatNumber(value: number): string {
  return NUMBER_FORMATTER.format(value);
}

function formatBytes(value: number): string {
  if (value <= BYTE_UNIT_START_INDEX) {
    return `0 ${BYTE_UNITS[BYTE_UNIT_START_INDEX]}`;
  }
  let size = value;
  let unitIndex = BYTE_UNIT_START_INDEX;
  while (size >= BYTE_UNIT_BASE && unitIndex < BYTE_UNITS.length - BYTE_UNIT_LAST_INDEX_OFFSET) {
    size /= BYTE_UNIT_BASE;
    unitIndex += UNIT_INDEX_INCREMENT;
  }
  return `${size.toFixed(DISPLAY_FRACTION_DIGITS)} ${BYTE_UNITS[unitIndex]}`;
}

function percentText(value: number): string {
  return value.toFixed(DISPLAY_FRACTION_DIGITS);
}

function displayPercentText(value: number): string {
  if (value > ZERO_VALUE && value < PERCENT_TINY_THRESHOLD) {
    return `<${percentText(PERCENT_TINY_THRESHOLD)}`;
  }
  return percentText(value);
}

function visualPiePercent(value: number, remainingValue: number): number {
  if (value <= ZERO_VALUE) {
    return ZERO_VALUE;
  }
  if (remainingValue <= ZERO_VALUE) {
    return PERCENT_FULL;
  }
  if (value < MIN_VISIBLE_PIE_PERCENT) {
    return MIN_VISIBLE_PIE_PERCENT;
  }
  if (value > PERCENT_FULL - MIN_VISIBLE_PIE_PERCENT) {
    return PERCENT_FULL - MIN_VISIBLE_PIE_PERCENT;
  }
  return value;
}

function topMetrics(data: GetProjectDashboardResponse): TopMetric[] {
  return [
    {
      icon: <RocketOutlined />,
      label: '进行中项目数量',
      value: data.active_project_count,
    },
    {
      icon: <CheckCircleOutlined />,
      label: '已完成项目数量',
      value: data.completed_project_count,
    },
    {
      icon: <ProjectOutlined />,
      label: '总创建项目数量',
      value: data.total_project_count,
    },
    {
      icon: <FileImageOutlined />,
      label: '项目图像上传总数',
      value: data.uploaded_image_count,
    },
    {
      icon: <FolderOpenOutlined />,
      label: '项目图像组创建总数',
      value: data.project_group_count,
    },
    {
      icon: <CloudServerOutlined />,
      label: '项目中间产物总数',
      value: data.minio_object_count,
    },
  ];
}

function barMetrics(data: GetProjectDashboardResponse): BarMetric[] {
  return [
    {
      label: '检测任务执行总次数',
      value: data.detection_task_count,
    },
    {
      label: '金属锈蚀检测执行总次数',
      value: data.corrosion_detection_task_count,
    },
    {
      label: '石材裂缝检测执行总次数',
      value: data.crack_detection_task_count,
    },
    {
      label: '石材污渍检测执行总次数',
      value: data.stain_detection_task_count,
    },
    {
      label: '玻璃平整度检测执行总次数',
      value: data.flatness_detection_task_count,
    },
    {
      label: '玻璃爆裂检测执行总次数',
      value: data.spalling_detection_task_count,
    },
    {
      label: 'Agent 图像检测结果总结次数',
      value: data.detection_summary_task_count,
    },
    {
      label: 'Agent 韧性评估报告生成次数',
      value: data.report_task_count,
    },
  ];
}

function storageSegments(data: GetProjectDashboardResponse): StorageSegment[] {
  return [
    {
      color: MINIO_USED_COLOR,
      label: '对象存储已使用空间',
      value: data.minio_bucket_used_bytes,
    },
    {
      color: MINIO_REMAINING_COLOR,
      label: '对象存储剩余可用空间',
      value: data.minio_bucket_remaining_bytes,
    },
  ];
}

function SkeletonBar({ className, style }: { className: string; style?: CSSProperties }): ReactElement {
  return <span className={`block animate-pulse rounded-lg bg-slate-200 ${className}`} style={style} />;
}

function ChartTitle({ icon, title }: ChartTitleProps): ReactElement {
  return (
    <div className="relative flex h-10 w-full shrink-0 items-center justify-center text-base font-semibold text-slate-900">
      <span className="absolute inset-0 flex items-center justify-center gap-2 whitespace-nowrap">
        {icon}
        <span>{title}</span>
      </span>
    </div>
  );
}

function DashboardSkeleton(): ReactElement {
  return (
    <div className="flex min-h-0 flex-1 flex-col gap-4">
      <div className="grid grid-cols-1 gap-4 lg:grid-cols-3">
        {Array.from({ length: TOP_CARD_SKELETON_COUNT }, (_, index) => (
          <Card className="h-32 border-slate-200 shadow-none" key={index}>
            <div className="flex h-full items-center justify-between">
              <div className="space-y-3">
                <SkeletonBar className="h-4 w-28" />
                <SkeletonBar className="h-8 w-20" />
              </div>
              <SkeletonBar className="size-12 rounded-xl" />
            </div>
          </Card>
        ))}
      </div>
      <div className="grid min-h-[28rem] flex-1 grid-cols-1 gap-4 lg:min-h-0 lg:grid-cols-2">
        <Card className="h-full min-h-[28rem] border-slate-200 shadow-none lg:min-h-0 [&_.ant-card-body]:flex [&_.ant-card-body]:h-full [&_.ant-card-body]:min-h-0 [&_.ant-card-body]:flex-col [&_.ant-card-body]:px-8 [&_.ant-card-body]:py-5">
          <div className="flex h-10 shrink-0 justify-center">
            <SkeletonBar className="h-5 w-32" />
          </div>
          <div className="flex min-h-0 flex-1 flex-col justify-evenly">
            {Array.from({ length: CHART_SKELETON_COUNT }, (_, index) => (
              <SkeletonBar className="h-6 w-full" key={index} />
            ))}
          </div>
        </Card>
        <Card className="h-full min-h-[28rem] border-slate-200 shadow-none lg:min-h-0 [&_.ant-card-body]:flex [&_.ant-card-body]:h-full [&_.ant-card-body]:min-h-0 [&_.ant-card-body]:flex-col [&_.ant-card-body]:px-8 [&_.ant-card-body]:py-5">
          <div className="flex h-10 shrink-0 justify-center">
            <SkeletonBar className="h-5 w-36" />
          </div>
          <div className="flex min-h-0 flex-1 items-center justify-center">
            <SkeletonBar className="rounded-full" style={{ height: PIE_CHART_SIZE, width: PIE_CHART_SIZE }} />
          </div>
          <div className="shrink-0">
            <SkeletonBar className="h-6 w-64" />
          </div>
        </Card>
      </div>
    </div>
  );
}

function TopMetricCard({ metric }: { metric: TopMetric }): ReactElement {
  return (
    <Card className="h-32 border-slate-200 shadow-none">
      <div className="flex h-full items-center justify-between gap-4">
        <div className="min-w-0">
          <div className="truncate text-sm text-slate-500">{metric.label}</div>
          <div className="mt-3 truncate text-3xl font-semibold text-slate-950">{formatNumber(metric.value)}</div>
        </div>
        <div className="flex size-12 shrink-0 items-center justify-center rounded-xl bg-blue-50 text-xl text-blue-600">
          {metric.icon}
        </div>
      </div>
    </Card>
  );
}

function DetectionBarChart({ metrics }: { metrics: BarMetric[] }): ReactElement {
  const maxValue = Math.max(...metrics.map((metric) => metric.value));
  return (
    <Card className="h-full min-h-[28rem] overflow-hidden border-slate-200 shadow-none lg:min-h-0 [&_.ant-card-body]:flex [&_.ant-card-body]:h-full [&_.ant-card-body]:min-h-0 [&_.ant-card-body]:w-full [&_.ant-card-body]:flex-col [&_.ant-card-body]:px-8 [&_.ant-card-body]:py-5">
      <ChartTitle icon={<BarChartOutlined className="text-blue-600" />} title="检测任务执行概览" />
      <div className="flex min-h-0 w-full flex-1 flex-col justify-evenly py-2">
        {metrics.map((metric) => {
          const percent =
            maxValue > ZERO_VALUE ? Math.max((metric.value / maxValue) * PERCENT_FULL, BAR_MIN_PERCENT) : ZERO_VALUE;
          return (
            <div className="space-y-2" key={metric.label}>
              <div className="flex items-center justify-between gap-3 text-sm">
                <span className="truncate text-slate-600">{metric.label}</span>
                <span className="shrink-0 font-medium text-slate-900">{formatNumber(metric.value)}</span>
              </div>
              <div className="h-3 overflow-hidden rounded-full bg-slate-100">
                <div
                  className="h-full rounded-full bg-blue-500 transition-all duration-300"
                  style={{
                    width: `${percentText(percent)}%`,
                  }}
                />
              </div>
            </div>
          );
        })}
      </div>
    </Card>
  );
}

function MinioStorageChart({ data }: { data: GetProjectDashboardResponse }): ReactElement {
  const segments = storageSegments(data);
  const total = segments.reduce((sum, segment) => sum + segment.value, ZERO_VALUE);
  const usedPercent = total > ZERO_VALUE ? (data.minio_bucket_used_bytes / total) * PERCENT_FULL : ZERO_VALUE;
  const usedPercentText = displayPercentText(usedPercent);
  const chartUsedPercent = visualPiePercent(usedPercent, data.minio_bucket_remaining_bytes);
  const chartUsedPercentText = percentText(chartUsedPercent);
  const pieBackground = `conic-gradient(${MINIO_USED_COLOR} 0 ${chartUsedPercentText}%, ${MINIO_REMAINING_COLOR} ${chartUsedPercentText}% ${String(PERCENT_FULL)}%)`;
  const quotaText =
    data.minio_bucket_quota_bytes > ZERO_VALUE
      ? `Bucket 配额：${formatBytes(data.minio_bucket_quota_bytes)}`
      : 'Bucket 未配置固定配额，剩余空间按对象存储当前可用容量展示';

  return (
    <Card className="h-full min-h-[28rem] overflow-hidden border-slate-200 shadow-none lg:min-h-0 [&_.ant-card-body]:flex [&_.ant-card-body]:h-full [&_.ant-card-body]:min-h-0 [&_.ant-card-body]:w-full [&_.ant-card-body]:flex-col [&_.ant-card-body]:px-8 [&_.ant-card-body]:py-5">
      <ChartTitle icon={<DatabaseOutlined className="text-blue-600" />} title="对象存储容量概览" />
      {data.minio_storage_available && total > ZERO_VALUE ? (
        <>
          <div className="flex min-h-0 w-full flex-1 items-center justify-center">
            <div
              className="relative rounded-full"
              style={{
                background: pieBackground,
                height: PIE_CHART_SIZE,
                width: PIE_CHART_SIZE,
              }}
            >
              <div className="absolute inset-[14%] flex flex-col items-center justify-center rounded-full bg-white text-center">
                <span className="text-xs text-slate-500">已使用</span>
                <span className="mt-1 text-2xl font-semibold text-slate-950">{usedPercentText}%</span>
              </div>
            </div>
          </div>
          <div className="w-full shrink-0">
            <div className="space-y-3">
              {segments.map((segment) => (
                <div className="flex items-center justify-between gap-3 text-sm" key={segment.label}>
                  <span className="flex min-w-0 items-center gap-2 text-slate-600">
                    <span className="size-3 shrink-0 rounded-full" style={{ backgroundColor: segment.color }} />
                    <span className="truncate">{segment.label}</span>
                  </span>
                  <span className="shrink-0 font-medium text-slate-900">{formatBytes(segment.value)}</span>
                </div>
              ))}
              <div className="pt-2 text-xs text-slate-400">{quotaText}</div>
            </div>
          </div>
        </>
      ) : (
        <div className="flex flex-1 items-center justify-center">
          <Empty description="对象存储统计暂不可用" />
        </div>
      )}
    </Card>
  );
}

export default function DashboardPage(): ReactElement {
  const [messageApi, contextHolder] = message.useMessage();
  const [dashboard, setDashboard] = useState<GetProjectDashboardResponse | null>(null);
  const [loading, setLoading] = useState(true);

  const loadDashboard = useCallback(async (): Promise<void> => {
    setLoading(true);
    try {
      setDashboard(await getProjectDashboard());
    } catch (error: unknown) {
      void messageApi.error(getErrorMessage(error));
    } finally {
      setLoading(false);
    }
  }, [messageApi]);

  useEffect(() => {
    void loadDashboard();
  }, [loadDashboard]);

  const topMetricItems = useMemo(() => (dashboard ? topMetrics(dashboard) : []), [dashboard]);
  const barMetricItems = useMemo(() => (dashboard ? barMetrics(dashboard) : []), [dashboard]);

  return (
    <div className="flex min-h-[calc(100vh-112px)] flex-col lg:h-[calc(100vh-112px)] lg:overflow-hidden">
      {contextHolder}
      <div className="mb-6 shrink-0">
        <h1 className="text-xl font-semibold text-slate-900">工作台</h1>
        <p className="mt-1 text-sm text-slate-500">统一查看账号下建筑项目、图像资产、智能检测与对象存储资源概览</p>
      </div>
      {loading ? (
        <DashboardSkeleton />
      ) : dashboard ? (
        <div className="flex min-h-0 flex-1 flex-col gap-4">
          <div className="grid shrink-0 grid-cols-1 gap-4 lg:grid-cols-3">
            {topMetricItems.map((metric) => (
              <TopMetricCard key={metric.label} metric={metric} />
            ))}
          </div>
          <div className="grid min-h-[28rem] flex-1 grid-cols-1 items-stretch gap-4 lg:min-h-0 lg:grid-cols-2">
            <DetectionBarChart metrics={barMetricItems} />
            <MinioStorageChart data={dashboard} />
          </div>
        </div>
      ) : (
        <Card className="border-slate-200 shadow-none">
          <Empty description="暂无工作台数据" />
        </Card>
      )}
    </div>
  );
}
