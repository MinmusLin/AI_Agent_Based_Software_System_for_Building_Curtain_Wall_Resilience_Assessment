import { LeftOutlined, RightOutlined } from '@ant-design/icons';
import { Button, Modal, Spin, Tabs } from 'antd';
import type { ReactElement } from 'react';
import { useMemo, useState } from 'react';

import type { ProjectDetectionSubStatus } from '@/types/common';
import {
  PROJECT_DETECTION_SUB_STATUS_FAILED,
  PROJECT_DETECTION_SUB_STATUS_PENDING,
  PROJECT_DETECTION_SUB_STATUS_SUCCEEDED,
} from '@/types/common';
import type { ProjectImage } from '@/types/project/assets';
import type {
  GetImageDetectionResultResponse,
  ProjectDetectionArtifacts,
  ProjectDetectionCorrosionResult,
  ProjectDetectionCrackResult,
  ProjectDetectionFlatnessResult,
  ProjectDetectionSpallingResult,
  ProjectDetectionStainResult,
  ProjectDetectionStatus,
  ProjectDetectionSummaryResult,
} from '@/types/project/detection';
import type { ViewerImage } from '@/utils/assetsStage';
import { FIRST_INDEX, formatProjectImageMetadata, KILOBYTE_SIZE_BYTES, NEXT_INDEX_OFFSET } from '@/utils/assetsStage';
import { formatDateTime } from '@/utils/datetime';

const DEFAULT_ACTIVE_TAB = 'original';
const SUMMARY_TAB = 'summary';
const DETECTION_DATA_PAGE_INDEX = 0;
const EMPTY_ITEMS_COUNT = 0;
const DEFAULT_DECIMAL_PLACES = 6;
const KILOBYTE_DECIMAL_PLACES = 0;
const JSON_FORMAT_INDENT = 2;
const MAX_PERCENT_DECIMAL_PLACES = 4;
const PERCENT_RATIO_BASE = 100;
const RUNTIME_DECIMAL_PLACES = 2;
const RATIO_FIELD_PATTERN = /(^|_)ratio$/u;
const BBOX_FIELD = 'bbox_xyxy';
const MAX_RATIO_VALUE = 1;

const REASONING_TASK_LABELS: Record<string, string> = {
  corrosion: '金属锈蚀',
  crack: '石材裂缝',
  flatness: '玻璃平整度',
  spalling: '玻璃爆裂',
  stain: '石材污渍',
};

const FIELD_LABELS: Record<string, string> = {
  angle_std: '角度标准差',
  average_confidence: '平均置信度',
  average_stain_ratio: '平均污渍占比',
  confidence: '置信度',
  bbox_x1: '边界框 x1 位置',
  bbox_x2: '边界框 x2 位置',
  bbox_y1: '边界框 y1 位置',
  bbox_y2: '边界框 y2 位置',
  corrosion_count: '锈蚀区域数量',
  corrosion_pixels: '锈蚀像素',
  corrosion_ratio: '锈蚀占比',
  crack_count: '裂缝区域数量',
  crack_pixels: '裂缝像素',
  crack_ratio: '裂缝占比',
  edge_count: '边缘数量',
  edge_uneven_detected: '边缘异常',
  frequency_max: '频谱最大值',
  frequency_min: '频谱最小值',
  frequency_uneven_detected: '频谱异常',
  gradient_mean: '梯度均值',
  gradient_std: '梯度标准差',
  gradient_uneven_detected: '梯度异常',
  laplacian_variance: '拉普拉斯方差',
  line_count: '直线数量',
  line_uneven_detected: '直线异常',
  mask_pixels: '裂缝像素',
  mask_ratio: '裂缝占比',
  max_confidence: '最高置信度',
  max_stain_ratio: '最高污渍占比',
  region_height: '区块高度',
  region_width: '区块宽度',
  result: '检测结论',
  runtime_seconds: '检测器执行耗时',
  stain_count: '污渍区域数量',
  stain_pixels: '污渍像素',
  stain_ratio: '污渍占比',
  uneven_count: '不平整区域数量',
};

const REPORT_EXCLUDED_FIELDS = new Set([
  'artifacts',
  'finished_at',
  'has_corrosion',
  'has_crack',
  'has_spalling',
  'has_stain',
  'regions',
  'result',
  'result_json',
  'runtime_seconds',
  'started_at',
  'status',
  'task_uuid',
]);

interface ProjectDetectionResultViewerProps {
  loading: boolean;
  onClose: () => void;
  onOpenImage: (imageUuid: string) => void;
  open: boolean;
  result: GetImageDetectionResultResponse | null;
  uploadedImages: ViewerImage[];
  viewerIndex: number;
}

interface OriginalInfoTabProps {
  image: ProjectImage | undefined;
  mainTaskUuid: string | undefined;
  status: ProjectDetectionStatus | undefined;
}

interface DetectionResultTabProps {
  groupIndex: number;
  onGroupIndexChange: (index: number) => void;
  result: DetectionResult | undefined;
  taskKey: string;
}

interface ResultNavigationProps {
  currentIndex: number;
  onChange: (index: number) => void;
  regionCount: number;
}

interface ReportDetailsProps {
  entries: [string, unknown][];
  title: string;
}

interface DetectionTabItem {
  key: string;
  label: string;
  result: DetectionResult | undefined;
}

interface DetectionPreviewProps {
  activeTab: string;
  groupIndex: number;
  image: ProjectImage | undefined;
  loading: boolean;
  originalUrl: string | undefined;
  result: DetectionResult | undefined;
  taskKey: string;
}

type DetectionResult =
  | ProjectDetectionCorrosionResult
  | ProjectDetectionCrackResult
  | ProjectDetectionFlatnessResult
  | ProjectDetectionSpallingResult
  | ProjectDetectionStainResult;

function imageDimensionText(image?: ProjectImage): string {
  if (!image) {
    return '-';
  }
  return `${String(image.width)} × ${String(image.height)}`;
}

function imageSizeText(image?: ProjectImage): string {
  if (!image) {
    return '-';
  }
  return `${(image.size_bytes / KILOBYTE_SIZE_BYTES).toFixed(KILOBYTE_DECIMAL_PLACES)} KB`;
}

function imageDateText(value?: string): string {
  if (!value) {
    return '-';
  }
  return formatDateTime(value, true);
}

function runtimeText(value?: number): string {
  if (value === undefined || value <= EMPTY_ITEMS_COUNT) {
    return '-';
  }
  return `${value.toFixed(RUNTIME_DECIMAL_PLACES)} 秒`;
}

function fieldLabel(key: string): string {
  return FIELD_LABELS[key] ?? key;
}

function formatFlatnessResult(value: unknown): string {
  if (value === 'uneven') {
    return '不平整';
  }
  if (value === 'flat') {
    return '平整';
  }
  if (value === 'notglass') {
    return '非玻璃';
  }
  return formatValue(value);
}

function formatValue(value: unknown): string {
  if (value === null || value === undefined || value === '') {
    return '-';
  }
  if (Array.isArray(value)) {
    return `[${value.map((item) => formatValue(item)).join(', ')}]`;
  }
  if (typeof value === 'boolean') {
    return value ? '是' : '否';
  }
  if (typeof value === 'number') {
    return Number.isInteger(value)
      ? String(value)
      : value.toFixed(DEFAULT_DECIMAL_PLACES).replace(/0+$/u, '').replace(/\.$/u, '');
  }
  if (typeof value === 'object') {
    return JSON.stringify(value, null, JSON_FORMAT_INDENT);
  }
  if (typeof value === 'string') {
    return value;
  }
  if (typeof value === 'bigint') {
    return value.toString();
  }
  return '-';
}

function formatPercent(value: number): string {
  const percentValue = Math.abs(value) > MAX_RATIO_VALUE ? value : value * PERCENT_RATIO_BASE;
  return `${percentValue.toFixed(MAX_PERCENT_DECIMAL_PLACES).replace(/0+$/u, '').replace(/\.$/u, '')}%`;
}

function formatReportValue(key: string, value: unknown): string {
  if (key === 'result') {
    return formatFlatnessResult(value);
  }
  if (RATIO_FIELD_PATTERN.test(key) && typeof value === 'number') {
    return formatPercent(value);
  }
  return formatValue(value);
}

function detailEntries(source: Record<string, unknown>, excludedFields: Set<string>): [string, unknown][] {
  const entries: [string, unknown][] = [];
  Object.entries(source).forEach(([key, value]) => {
    if (excludedFields.has(key)) {
      return;
    }
    if (key === BBOX_FIELD && Array.isArray(value)) {
      entries.push(['bbox_x1', value[0] ?? '-']);
      entries.push(['bbox_y1', value[1] ?? '-']);
      entries.push(['bbox_x2', value[2] ?? '-']);
      entries.push(['bbox_y2', value[3] ?? '-']);
      return;
    }
    entries.push([key, value]);
  });
  return entries;
}

function reportEntries(report: Record<string, unknown>): [string, unknown][] {
  return detailEntries(report, REPORT_EXCLUDED_FIELDS);
}

function resultReport(result?: DetectionResult): Record<string, unknown> {
  if (!result) {
    return {};
  }
  return Object.fromEntries(
    Object.entries(result as Record<string, unknown>).filter(([key]) => !REPORT_EXCLUDED_FIELDS.has(key)),
  );
}

function resultRegions(result?: DetectionResult): Record<string, unknown>[] {
  if (!result || !('regions' in result) || !Array.isArray(result.regions)) {
    return [];
  }
  return result.regions.map((item) => ({ ...item }));
}

function resultPageIndex(result: DetectionResult | undefined, pageIndex: number): number {
  const pageCount = resultRegions(result).length + NEXT_INDEX_OFFSET;
  return Math.min(pageIndex, Math.max(pageCount - NEXT_INDEX_OFFSET, DETECTION_DATA_PAGE_INDEX));
}

function resultArtifacts(result?: DetectionResult): ProjectDetectionArtifacts {
  if (!result || !('artifacts' in result)) {
    return {};
  }
  return result.artifacts ?? {};
}

function artifactUrl(artifacts: ProjectDetectionArtifacts, name: string): string | undefined {
  return artifacts[name];
}

function detectionResultTabs(result: GetImageDetectionResultResponse | null): DetectionTabItem[] {
  if (!result) {
    return [];
  }
  return [
    {
      key: 'corrosion',
      label: REASONING_TASK_LABELS.corrosion,
      result: result.corrosion_result,
    },
    {
      key: 'crack',
      label: REASONING_TASK_LABELS.crack,
      result: result.crack_result,
    },
    {
      key: 'stain',
      label: REASONING_TASK_LABELS.stain,
      result: result.stain_result,
    },
    {
      key: 'flatness',
      label: REASONING_TASK_LABELS.flatness,
      result: result.flatness_result,
    },
    {
      key: 'spalling',
      label: REASONING_TASK_LABELS.spalling,
      result: result.spalling_result,
    },
  ].filter((item): item is DetectionTabItem => Boolean(item.result));
}

function statusText(status?: ProjectDetectionSubStatus): string {
  if (status === PROJECT_DETECTION_SUB_STATUS_PENDING) {
    return '检测中';
  }
  if (status === PROJECT_DETECTION_SUB_STATUS_SUCCEEDED) {
    return '成功';
  }
  if (status === PROJECT_DETECTION_SUB_STATUS_FAILED) {
    return '失败';
  }
  return '-';
}

function detectionConclusion(taskKey: string, result?: DetectionResult): string {
  if (!result) {
    return '-';
  }
  if (taskKey === 'corrosion' && 'has_corrosion' in result) {
    return result.has_corrosion ? '存在锈蚀' : '不存在锈蚀';
  }
  if (taskKey === 'crack' && 'has_crack' in result) {
    return result.has_crack ? '存在裂缝' : '不存在裂缝';
  }
  if (taskKey === 'stain' && 'has_stain' in result) {
    return result.has_stain ? '存在污渍' : '不存在污渍';
  }
  if (taskKey === 'flatness' && 'result' in result) {
    const resultText = formatFlatnessResult(result.result);
    return resultText === '非玻璃' ? '非玻璃材质' : `玻璃${resultText}`;
  }
  if (taskKey === 'spalling' && 'has_spalling' in result) {
    return result.has_spalling ? '存在爆裂' : '不存在爆裂';
  }
  return statusText(result.status);
}

function parseSummaryJSON(reportJson?: string): string {
  if (!reportJson || reportJson.trim() === '') {
    return '-';
  }
  try {
    const parsed = JSON.parse(reportJson) as unknown;
    if (typeof parsed === 'object' && parsed !== null && 'summary' in parsed && typeof parsed.summary === 'string') {
      return parsed.summary;
    }
    return JSON.stringify(parsed, null, JSON_FORMAT_INDENT);
  } catch {
    return reportJson;
  }
}

function OriginalInfoTab({ image, mainTaskUuid, status }: OriginalInfoTabProps): ReactElement {
  return (
    <div className="flex h-full min-h-0 flex-col">
      <div className="shrink-0 space-y-4 text-sm text-slate-600">
        <div>
          <div className="mb-1 font-medium text-slate-900">任务 ID</div>
          <div className="line-clamp-1 break-all">{mainTaskUuid ?? '-'}</div>
        </div>
        <div>
          <div className="mb-1 font-medium text-slate-900">文件名称</div>
          <div className="overflow-hidden break-all [display:-webkit-box] [-webkit-box-orient:vertical] [-webkit-line-clamp:2]">
            {image?.file_name ?? '-'}
          </div>
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <div className="mb-1 font-medium text-slate-900">图像尺寸（宽 × 高）</div>
            <div>{imageDimensionText(image)}</div>
          </div>
          <div>
            <div className="mb-1 font-medium text-slate-900">文件大小</div>
            <div>{imageSizeText(image)}</div>
          </div>
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <div className="mb-1 font-medium text-slate-900">任务开始时间</div>
            <div>{imageDateText(status?.started_at)}</div>
          </div>
          <div>
            <div className="mb-1 font-medium text-slate-900">任务完成时间</div>
            <div>{imageDateText(status?.finished_at)}</div>
          </div>
        </div>
      </div>
      <div className="mt-4 flex min-h-0 flex-1 flex-col text-sm text-slate-600">
        <div className="mb-1 font-medium text-slate-900">元数据</div>
        <div className="min-h-0 flex-1 overflow-hidden rounded bg-slate-50 p-3">
          <pre className="h-full overflow-auto overscroll-contain whitespace-pre-wrap break-all text-xs leading-5">
            {formatProjectImageMetadata(image?.metadata ?? '{}')}
          </pre>
        </div>
      </div>
    </div>
  );
}

function ResultNavigation({ currentIndex, onChange, regionCount }: ResultNavigationProps): ReactElement {
  const total = regionCount + NEXT_INDEX_OFFSET;
  const label =
    currentIndex === DETECTION_DATA_PAGE_INDEX ? '区域数据' : `区域 ${String(currentIndex)} / ${String(regionCount)}`;

  return (
    <div className="flex items-center justify-between rounded-t bg-slate-50 px-3 py-2">
      <Button
        disabled={currentIndex <= FIRST_INDEX}
        icon={<LeftOutlined />}
        onClick={() => {
          onChange(currentIndex - NEXT_INDEX_OFFSET);
        }}
        size="small"
      >
        上一个结果
      </Button>
      <span className="text-xs text-slate-500">{label}</span>
      <Button
        disabled={currentIndex >= total - NEXT_INDEX_OFFSET}
        onClick={() => {
          onChange(currentIndex + NEXT_INDEX_OFFSET);
        }}
        size="small"
      >
        下一个结果
        <RightOutlined />
      </Button>
    </div>
  );
}

function ReportDetails({ entries, title }: ReportDetailsProps): ReactElement {
  return (
    <div className="flex min-h-0 flex-1 flex-col rounded-b bg-slate-50 p-3">
      <div className="mb-2 font-medium text-slate-900">{title}</div>
      <div className="min-h-0 flex-1 space-y-2 overflow-auto overscroll-contain">
        {entries.map(([key, value]) => (
          <div className="grid grid-cols-[118px_minmax(0,1fr)] gap-3" key={key}>
            <span className="text-slate-500">{fieldLabel(key)}</span>
            <span className="break-all text-slate-700">{formatReportValue(key, value)}</span>
          </div>
        ))}
      </div>
    </div>
  );
}

function DetectionResultTab({
  groupIndex,
  onGroupIndexChange,
  result,
  taskKey,
}: DetectionResultTabProps): ReactElement {
  const report = resultReport(result);
  const regions = resultRegions(result);
  const currentPageIndex = resultPageIndex(result, groupIndex);
  const currentRegion =
    currentPageIndex > DETECTION_DATA_PAGE_INDEX ? (regions.at(currentPageIndex - NEXT_INDEX_OFFSET) ?? null) : null;
  const detailRows =
    currentRegion === null
      ? reportEntries(report)
      : detailEntries(currentRegion, new Set([...REPORT_EXCLUDED_FIELDS, 'id']));
  const showNavigation = regions.length > EMPTY_ITEMS_COUNT;

  return (
    <div className="flex h-full min-h-0 flex-col gap-4 text-sm text-slate-600">
      <div>
        <div className="mb-1 font-medium text-slate-900">子任务 ID</div>
        <div className="break-all">{result?.task_uuid ?? '-'}</div>
      </div>
      <div className="grid grid-cols-2 gap-4">
        <div>
          <div className="mb-1 font-medium text-slate-900">执行耗时</div>
          <div>{runtimeText(result?.runtime_seconds)}</div>
        </div>
        <div>
          <div className="mb-1 font-medium text-slate-900">检测结果</div>
          <div>{detectionConclusion(taskKey, result)}</div>
        </div>
      </div>
      <div className="grid grid-cols-2 gap-4">
        <div>
          <div className="mb-1 font-medium text-slate-900">任务开始时间</div>
          <div>{imageDateText(result?.started_at)}</div>
        </div>
        <div>
          <div className="mb-1 font-medium text-slate-900">任务完成时间</div>
          <div>{imageDateText(result?.finished_at)}</div>
        </div>
      </div>
      <div className="flex min-h-0 flex-1 flex-col">
        {showNavigation ? (
          <ResultNavigation
            currentIndex={currentPageIndex}
            onChange={onGroupIndexChange}
            regionCount={regions.length}
          />
        ) : null}
        <ReportDetails entries={detailRows} title={currentRegion === null ? '检测数据' : '区域数据'} />
      </div>
    </div>
  );
}

function OriginalImagePreview({
  image,
  loading,
  originalUrl,
}: {
  image: ProjectImage | undefined;
  loading: boolean;
  originalUrl: string | undefined;
}): ReactElement {
  if (loading) {
    return <Spin />;
  }
  return (
    <img
      alt={image?.file_name ?? '图像原图'}
      className="max-h-[520px] max-w-full object-contain"
      draggable={false}
      src={originalUrl}
    />
  );
}

function ArtifactImage({
  cellClassName = 'items-center justify-center',
  name,
  url,
}: {
  cellClassName?: string;
  name: string;
  url: string;
}): ReactElement {
  return (
    <div className={`flex size-full min-h-0 overflow-hidden ${cellClassName}`}>
      <img alt={name} className="max-h-full max-w-full object-contain" draggable={false} src={url} />
    </div>
  );
}

function SingleArtifactPreview({ name, url }: { name: string; url: string | undefined }): ReactElement | null {
  if (!url) {
    return null;
  }
  return <ArtifactImage name={name} url={url} />;
}

function HorizontalArtifactsPreview({
  left,
  right,
}: {
  left: { name: string; url: string | undefined };
  right: { name: string; url: string | undefined };
}): ReactElement | null {
  if (!left.url || !right.url) {
    return null;
  }
  return (
    <div className="grid size-full grid-cols-2 gap-0 overflow-hidden bg-slate-100">
      <ArtifactImage cellClassName="items-center justify-end" name={left.name} url={left.url} />
      <ArtifactImage cellClassName="items-center justify-start" name={right.name} url={right.url} />
    </div>
  );
}

function FlatnessRegionPreview({
  artifacts,
  regionId,
}: {
  artifacts: ProjectDetectionArtifacts;
  regionId: string;
}): ReactElement | null {
  const items = [
    { name: `region_${regionId}.png`, url: artifactUrl(artifacts, `region_${regionId}.png`) },
    { name: `gradient_${regionId}.png`, url: artifactUrl(artifacts, `gradient_${regionId}.png`) },
    { name: `lines_${regionId}.png`, url: artifactUrl(artifacts, `lines_${regionId}.png`) },
    { name: `frequency_${regionId}.png`, url: artifactUrl(artifacts, `frequency_${regionId}.png`) },
  ];
  if (items.some((item) => !item.url)) {
    return null;
  }
  return (
    <div className="grid size-full grid-cols-2 grid-rows-2 gap-0 overflow-hidden bg-slate-100">
      <ArtifactImage cellClassName="items-end justify-end" name={items[0].name} url={items[0].url ?? ''} />
      <ArtifactImage cellClassName="items-end justify-start" name={items[1].name} url={items[1].url ?? ''} />
      <ArtifactImage cellClassName="items-start justify-end" name={items[2].name} url={items[2].url ?? ''} />
      <ArtifactImage cellClassName="items-start justify-start" name={items[3].name} url={items[3].url ?? ''} />
    </div>
  );
}

function taskArtifactPreview(
  taskKey: string,
  artifacts: ProjectDetectionArtifacts,
  pageIndex: number,
): ReactElement | null {
  if (taskKey === 'corrosion') {
    return <SingleArtifactPreview name="annotated.png" url={artifactUrl(artifacts, 'annotated.png')} />;
  }
  if (taskKey === 'crack') {
    if (pageIndex === DETECTION_DATA_PAGE_INDEX) {
      return null;
    }
    return (
      <HorizontalArtifactsPreview
        left={{ name: 'overlay.png', url: artifactUrl(artifacts, 'overlay.png') }}
        right={{ name: 'mask.png', url: artifactUrl(artifacts, 'mask.png') }}
      />
    );
  }
  if (taskKey === 'stain') {
    if (pageIndex === DETECTION_DATA_PAGE_INDEX) {
      return <SingleArtifactPreview name="annotated.png" url={artifactUrl(artifacts, 'annotated.png')} />;
    }
    return (
      <HorizontalArtifactsPreview
        left={{
          name: `region_${String(pageIndex)}.png`,
          url: artifactUrl(artifacts, `region_${String(pageIndex)}.png`),
        }}
        right={{
          name: `overlay_${String(pageIndex)}.png`,
          url: artifactUrl(artifacts, `overlay_${String(pageIndex)}.png`),
        }}
      />
    );
  }
  if (taskKey === 'flatness') {
    if (pageIndex === DETECTION_DATA_PAGE_INDEX) {
      return (
        <HorizontalArtifactsPreview
          left={{ name: 'mask.png', url: artifactUrl(artifacts, 'mask.png') }}
          right={{ name: 'overlay.png', url: artifactUrl(artifacts, 'overlay.png') }}
        />
      );
    }
    return <FlatnessRegionPreview artifacts={artifacts} regionId={String(pageIndex)} />;
  }
  return null;
}

function DetectionPreview({
  activeTab,
  groupIndex,
  image,
  loading,
  originalUrl,
  result,
  taskKey,
}: DetectionPreviewProps): ReactElement {
  if (loading || activeTab === DEFAULT_ACTIVE_TAB || activeTab === SUMMARY_TAB || !result) {
    return <OriginalImagePreview image={image} loading={loading} originalUrl={originalUrl} />;
  }

  const preview = taskArtifactPreview(taskKey, resultArtifacts(result), resultPageIndex(result, groupIndex));
  if (!preview) {
    return <OriginalImagePreview image={image} loading={false} originalUrl={originalUrl} />;
  }

  return preview;
}

function viewerTitle(image: ProjectImage | undefined, uploadedImages: ViewerImage[], viewerIndex: number): string {
  if (image?.file_name) {
    return image.file_name;
  }
  if (viewerIndex >= FIRST_INDEX && uploadedImages[viewerIndex]) {
    return uploadedImages[viewerIndex].image.file_name;
  }
  return '检测结果';
}

function SummaryResultTab({ result }: { result: ProjectDetectionSummaryResult | undefined }): ReactElement {
  return (
    <div className="flex h-full min-h-0 flex-col gap-4 text-sm text-slate-600">
      <div>
        <div className="mb-1 font-medium text-slate-900">子任务 ID</div>
        <div className="break-all">{result?.task_uuid ?? '-'}</div>
      </div>
      <div className="grid grid-cols-2 gap-4">
        <div>
          <div className="mb-1 font-medium text-slate-900">任务开始时间</div>
          <div>{imageDateText(result?.started_at)}</div>
        </div>
        <div>
          <div className="mb-1 font-medium text-slate-900">任务完成时间</div>
          <div>{imageDateText(result?.finished_at)}</div>
        </div>
      </div>
      <div className="flex min-h-0 flex-1 flex-col rounded bg-slate-50 p-3">
        <div className="mb-2 font-medium text-slate-900">Agent 总结</div>
        <pre className="min-h-0 flex-1 overflow-auto overscroll-contain whitespace-pre-wrap break-all text-xs leading-5 text-slate-700">
          {parseSummaryJSON(result?.result_json)}
        </pre>
      </div>
    </div>
  );
}

export function ProjectDetectionResultViewer({
  loading,
  onClose,
  onOpenImage,
  open,
  result,
  uploadedImages,
  viewerIndex,
}: ProjectDetectionResultViewerProps): ReactElement {
  const [activeTab, setActiveTab] = useState(DEFAULT_ACTIVE_TAB);
  const [groupIndex, setGroupIndex] = useState(FIRST_INDEX);
  const image = result?.image;
  const title = viewerTitle(image, uploadedImages, viewerIndex);
  const reasoningTabs = useMemo(() => (loading ? [] : detectionResultTabs(result)), [loading, result]);
  const activeDetectionResult = loading ? undefined : reasoningTabs.find((item) => item.key === activeTab)?.result;
  const tabs = useMemo(
    () => [
      {
        children: (
          <OriginalInfoTab image={image} mainTaskUuid={result?.status?.main_task_uuid} status={result?.status} />
        ),
        key: DEFAULT_ACTIVE_TAB,
        label: '图像数据',
      },
      ...reasoningTabs.map((item) => ({
        children: (
          <DetectionResultTab
            groupIndex={groupIndex}
            onGroupIndexChange={setGroupIndex}
            result={item.result}
            taskKey={item.key}
          />
        ),
        key: item.key,
        label: item.label,
      })),
      ...(reasoningTabs.length > EMPTY_ITEMS_COUNT && result?.summary_result
        ? [
            {
              children: <SummaryResultTab result={result.summary_result} />,
              key: SUMMARY_TAB,
              label: 'Agent 总结',
            },
          ]
        : []),
    ],
    [groupIndex, image, reasoningTabs, result],
  );
  const currentActiveTab = !loading && tabs.some((tab) => tab.key === activeTab) ? activeTab : DEFAULT_ACTIVE_TAB;
  if (loading) {
    return (
      <Modal centered footer={null} onCancel={onClose} open={open} title={title} width={1040}>
        <div className="flex h-[520px] flex-col items-center justify-center gap-3 text-sm text-slate-500">
          <Spin />
          <div>正在加载检测结果</div>
        </div>
      </Modal>
    );
  }

  return (
    <Modal centered footer={null} onCancel={onClose} open={open} title={title} width={1040}>
      <div className="grid h-[520px] grid-cols-[minmax(0,1fr)_320px] gap-5">
        <div className="relative flex min-h-0 items-center justify-center overflow-hidden rounded-lg bg-slate-100">
          <DetectionPreview
            activeTab={currentActiveTab}
            groupIndex={groupIndex}
            image={image}
            loading={loading}
            originalUrl={result?.original_url}
            result={activeDetectionResult}
            taskKey={currentActiveTab}
          />
        </div>
        <div className="flex min-h-0 flex-col">
          <Tabs
            activeKey={currentActiveTab}
            className="min-h-0 flex-1 [&_.ant-tabs-content-holder]:min-h-0 [&_.ant-tabs-content]:h-full [&_.ant-tabs-tabpane]:h-full"
            items={tabs}
            onChange={(key) => {
              setActiveTab(key);
              setGroupIndex(FIRST_INDEX);
            }}
            size="small"
          />
          <div className="flex justify-between pt-4">
            <Button
              disabled={viewerIndex <= FIRST_INDEX}
              icon={<LeftOutlined />}
              onClick={() => {
                if (viewerIndex > FIRST_INDEX) {
                  onOpenImage(uploadedImages[viewerIndex - NEXT_INDEX_OFFSET].image.uuid);
                }
              }}
            >
              上一张
            </Button>
            <Button
              disabled={viewerIndex < FIRST_INDEX || viewerIndex >= uploadedImages.length - NEXT_INDEX_OFFSET}
              onClick={() => {
                if (viewerIndex >= FIRST_INDEX && viewerIndex < uploadedImages.length - NEXT_INDEX_OFFSET) {
                  onOpenImage(uploadedImages[viewerIndex + NEXT_INDEX_OFFSET].image.uuid);
                }
              }}
            >
              下一张
              <RightOutlined />
            </Button>
          </div>
        </div>
      </div>
    </Modal>
  );
}
