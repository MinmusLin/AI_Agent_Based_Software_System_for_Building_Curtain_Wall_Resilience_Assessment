import { DislikeOutlined, LeftOutlined, LikeOutlined, RightOutlined, RobotOutlined } from '@ant-design/icons';
import { Button, Input, Space, Spin } from 'antd';
import type { ReactElement } from 'react';

import type {
  ProjectDetectionCorrosionResult,
  ProjectDetectionCrackResult,
  ProjectDetectionFlatnessResult,
  ProjectDetectionSpallingResult,
  ProjectDetectionStainResult,
  ProjectDetectionStatus,
  ProjectImage,
} from '@/gen/core/common';
import { ProjectDetectionReviewVerdict_Value, ProjectDetectionSubTaskStatus_Value } from '@/gen/core/common';
import { FIRST_INDEX, formatProjectImageMetadata, KILOBYTE_SIZE_BYTES, NEXT_INDEX_OFFSET } from '@/utils/assetsStage';
import { formatDateTime } from '@/utils/datetime';

export const DEFAULT_ACTIVE_TAB = 'original';
export const SUMMARY_TAB = 'summary';
export const EMPTY_ITEMS_COUNT = 0;

const BUTTON_TYPE_DEFAULT = 'default';
const BUTTON_TYPE_PRIMARY = 'primary';
const DETECTION_DATA_PAGE_INDEX = 0;
const DEFAULT_DECIMAL_PLACES = 6;
const INACCURATE_BUTTON_HOVER_CLASS_NAME = 'hover:!border-red-500 hover:!text-red-500';
const KILOBYTE_DECIMAL_PLACES = 0;
const JSON_FORMAT_INDENT = 2;
const MAX_PERCENT_DECIMAL_PLACES = 4;
const MAX_REVIEW_COMMENT_LENGTH = 500;
const PERCENT_RATIO_BASE = 100;
const RUNTIME_DECIMAL_PLACES = 2;
const RATIO_FIELD_PATTERN = /(^|_)ratio$/u;
const BBOX_FIELD = 'bbox_xyxy';
const MAX_RATIO_VALUE = 1;

const FIELD_LABELS: Record<string, string> = {
  angle_std: '角度标准差',
  average_confidence: '平均置信度',
  average_stain_ratio: '平均污渍占比',
  bbox_x1: '边界框 x1 位置',
  bbox_x2: '边界框 x2 位置',
  bbox_y1: '边界框 y1 位置',
  bbox_y2: '边界框 y2 位置',
  confidence: '置信度',
  corrosion_count: '锈蚀区域数量',
  corrosion_pixels: '锈蚀像素数',
  corrosion_ratio: '锈蚀占比',
  crack_count: '裂缝区域数量',
  crack_pixels: '裂缝像素数',
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
  mask_pixels: '区域像素数',
  mask_ratio: '区域占比',
  max_confidence: '最高置信度',
  max_stain_ratio: '最高污渍占比',
  region_height: '校正区块高度',
  region_width: '校正区块宽度',
  result: '检测结论',
  runtime_seconds: '检测器执行耗时',
  stain_count: '污渍区域数量',
  stain_pixels: '污渍像素数',
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
  'runtime_seconds',
  'started_at',
  'status',
  'task_uuid',
]);

const DETECTION_REPORT_FIELD_ORDERS: Record<string, string[]> = {
  corrosion: ['corrosion_count', 'average_confidence', 'max_confidence', 'corrosion_pixels', 'corrosion_ratio'],
  crack: ['crack_count', 'crack_pixels', 'crack_ratio'],
  flatness: ['uneven_count'],
  spalling: ['confidence'],
  stain: ['stain_count', 'average_stain_ratio', 'max_stain_ratio'],
};

const DETECTION_REGION_FIELD_ORDERS: Record<string, string[]> = {
  corrosion: ['confidence', 'mask_pixels', 'mask_ratio', 'bbox_x1', 'bbox_y1', 'bbox_x2', 'bbox_y2'],
  crack: ['mask_pixels', 'mask_ratio', 'bbox_x1', 'bbox_y1', 'bbox_x2', 'bbox_y2'],
  flatness: [
    'edge_uneven_detected',
    'gradient_uneven_detected',
    'line_uneven_detected',
    'frequency_uneven_detected',
    'edge_count',
    'line_count',
    'angle_std',
    'gradient_mean',
    'gradient_std',
    'laplacian_variance',
    'frequency_max',
    'frequency_min',
    'bbox_x1',
    'bbox_y1',
    'bbox_x2',
    'bbox_y2',
  ],
  stain: [
    'confidence',
    'stain_pixels',
    'stain_ratio',
    'region_height',
    'region_width',
    'bbox_x1',
    'bbox_y1',
    'bbox_x2',
    'bbox_y2',
  ],
};

export interface ProjectDetectionReviewProps {
  comment: string;
  enabled: boolean;
  onCommentChange: (comment: string) => void;
  onSave: () => void;
  onVerdictChange: (verdict: ProjectDetectionReviewVerdict_Value) => void;
  readOnly: boolean;
  saving: boolean;
  verdict: ProjectDetectionReviewVerdict_Value;
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
  taskKey: string;
  title: string;
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

export type DetectionResult =
  | ProjectDetectionCorrosionResult
  | ProjectDetectionCrackResult
  | ProjectDetectionFlatnessResult
  | ProjectDetectionSpallingResult
  | ProjectDetectionStainResult;

type ProjectDetectionArtifacts = Record<string, string>;

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

function taskFieldLabel(taskKey: string, key: string): string {
  if (taskKey === 'corrosion' && key === 'mask_pixels') {
    return '锈蚀像素数';
  }
  if (taskKey === 'corrosion' && key === 'mask_ratio') {
    return '锈蚀占比';
  }
  if (taskKey === 'crack' && key === 'mask_pixels') {
    return '裂缝像素数';
  }
  if (taskKey === 'crack' && key === 'mask_ratio') {
    return '裂缝占比';
  }
  return fieldLabel(key);
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

function sortedEntries(entries: [string, unknown][], fieldOrder: string[] | undefined): [string, unknown][] {
  if (!fieldOrder) {
    return entries;
  }
  const fieldRank = new Map(fieldOrder.map((field, index) => [field, index]));
  return [...entries].sort(([leftKey], [rightKey]) => {
    const leftRank = fieldRank.get(leftKey) ?? Number.MAX_SAFE_INTEGER;
    const rightRank = fieldRank.get(rightKey) ?? Number.MAX_SAFE_INTEGER;
    if (leftRank !== rightRank) {
      return leftRank - rightRank;
    }
    return leftKey.localeCompare(rightKey);
  });
}

function reportEntries(taskKey: string, report: Record<string, unknown>): [string, unknown][] {
  return sortedEntries(detailEntries(report, REPORT_EXCLUDED_FIELDS), DETECTION_REPORT_FIELD_ORDERS[taskKey]);
}

function regionEntries(taskKey: string, region: Record<string, unknown>): [string, unknown][] {
  return sortedEntries(
    detailEntries(region, new Set([...REPORT_EXCLUDED_FIELDS, 'id'])),
    DETECTION_REGION_FIELD_ORDERS[taskKey],
  );
}

function isFlatnessUneven(result?: DetectionResult): boolean {
  return Boolean(result && 'result' in result && result.result === 'uneven');
}

function showDetectionRegions(taskKey: string, result?: DetectionResult): boolean {
  if (!result) {
    return false;
  }
  if (taskKey === 'corrosion' && 'has_corrosion' in result) {
    return result.has_corrosion;
  }
  if (taskKey === 'crack' && 'has_crack' in result) {
    return result.has_crack;
  }
  if (taskKey === 'stain' && 'has_stain' in result) {
    return result.has_stain;
  }
  if (taskKey === 'flatness') {
    return isFlatnessUneven(result);
  }
  return false;
}

function visibleRegions(taskKey: string, result?: DetectionResult): Record<string, unknown>[] {
  return showDetectionRegions(taskKey, result) ? resultRegions(result) : [];
}

function detectionDataEntries(
  taskKey: string,
  report: Record<string, unknown>,
  result?: DetectionResult,
): [string, unknown][] {
  if (taskKey === 'corrosion' && result && 'has_corrosion' in result && !result.has_corrosion) {
    return [];
  }
  if (taskKey === 'crack' && result && 'has_crack' in result && !result.has_crack) {
    return [];
  }
  if (taskKey === 'stain' && result && 'has_stain' in result && !result.has_stain) {
    return [];
  }
  if (taskKey === 'flatness' && !isFlatnessUneven(result)) {
    return [];
  }
  return reportEntries(taskKey, report);
}

function resultReport(result?: DetectionResult): Record<string, unknown> {
  if (!result) {
    return {};
  }
  return Object.fromEntries(
    Object.entries(result as unknown as Record<string, unknown>).filter(([key]) => !REPORT_EXCLUDED_FIELDS.has(key)),
  );
}

function resultRegions(result?: DetectionResult): Record<string, unknown>[] {
  if (!result || !('regions' in result) || !Array.isArray(result.regions)) {
    return [];
  }
  return result.regions.map((item) => ({ ...item }));
}

function resultPageIndex(taskKey: string, result: DetectionResult | undefined, pageIndex: number): number {
  const pageCount = visibleRegions(taskKey, result).length + NEXT_INDEX_OFFSET;
  return Math.min(pageIndex, Math.max(pageCount - NEXT_INDEX_OFFSET, DETECTION_DATA_PAGE_INDEX));
}

function resultArtifacts(result?: DetectionResult): ProjectDetectionArtifacts {
  if (!result || !('artifacts' in result)) {
    return {};
  }
  return result.artifacts;
}

function artifactUrl(artifacts: ProjectDetectionArtifacts, name: string): string | undefined {
  return artifacts[name];
}

function statusText(status?: ProjectDetectionSubTaskStatus_Value): string {
  if (status === ProjectDetectionSubTaskStatus_Value.Pending) {
    return '检测中';
  }
  if (status === ProjectDetectionSubTaskStatus_Value.Succeeded) {
    return '成功';
  }
  if (status === ProjectDetectionSubTaskStatus_Value.Failed) {
    return '失败';
  }
  return '-';
}

function detectionConclusion(taskKey: string, result?: DetectionResult): string {
  if (!result) {
    return '-';
  }
  if (taskKey === 'flatness' && 'result' in result) {
    const resultText = formatFlatnessResult(result.result);
    return resultText === '非玻璃' ? '非玻璃材质' : `玻璃${resultText}`;
  }
  const conclusions: Record<string, string> = {
    corrosion: '锈蚀',
    crack: '裂缝',
    spalling: '爆裂',
    stain: '污渍',
  };
  const field = `has_${taskKey}`;
  if (field in result && taskKey in conclusions) {
    return Boolean((result as unknown as Record<string, unknown>)[field])
      ? `存在${conclusions[taskKey]}`
      : `不存在${conclusions[taskKey]}`;
  }
  return statusText(result.status);
}

export function OriginalInfoTab({ image, mainTaskUuid, status }: OriginalInfoTabProps): ReactElement {
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

function ReportDetails({ entries, taskKey, title }: ReportDetailsProps): ReactElement {
  return (
    <div className="flex min-h-0 flex-1 flex-col rounded-b bg-slate-50 p-3">
      <div className="mb-2 font-medium text-slate-900">{title}</div>
      <div className="min-h-0 flex-1 space-y-2 overflow-auto overscroll-contain">
        {entries.map(([key, value]) => (
          <div className="grid grid-cols-[118px_minmax(0,1fr)] gap-3" key={key}>
            <span className="text-slate-500">{taskFieldLabel(taskKey, key)}</span>
            <span className="break-all text-slate-700">{formatReportValue(key, value)}</span>
          </div>
        ))}
      </div>
    </div>
  );
}

export function DetectionResultTab({
  groupIndex,
  onGroupIndexChange,
  result,
  taskKey,
}: DetectionResultTabProps): ReactElement {
  const report = resultReport(result);
  const regions = visibleRegions(taskKey, result);
  const currentPageIndex = resultPageIndex(taskKey, result, groupIndex);
  const currentRegion =
    currentPageIndex > DETECTION_DATA_PAGE_INDEX ? (regions.at(currentPageIndex - NEXT_INDEX_OFFSET) ?? null) : null;
  const detailRows =
    currentRegion === null ? detectionDataEntries(taskKey, report, result) : regionEntries(taskKey, currentRegion);
  const showNavigation = regions.length > EMPTY_ITEMS_COUNT;
  const showReportDetails = showNavigation || detailRows.length > EMPTY_ITEMS_COUNT;

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
      {showReportDetails ? (
        <div className="flex min-h-0 flex-1 flex-col">
          {showNavigation ? (
            <ResultNavigation
              currentIndex={currentPageIndex}
              onChange={onGroupIndexChange}
              regionCount={regions.length}
            />
          ) : null}
          <ReportDetails
            entries={detailRows}
            taskKey={taskKey}
            title={currentRegion === null ? '检测数据' : '区域数据'}
          />
        </div>
      ) : null}
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
  result?: DetectionResult,
): ReactElement | null {
  if (taskKey === 'corrosion') {
    if (!showDetectionRegions(taskKey, result)) {
      return null;
    }
    return <SingleArtifactPreview name="annotated.png" url={artifactUrl(artifacts, 'annotated.png')} />;
  }
  if (taskKey === 'crack') {
    if (!showDetectionRegions(taskKey, result)) {
      return null;
    }
    return (
      <HorizontalArtifactsPreview
        left={{ name: 'overlay.png', url: artifactUrl(artifacts, 'overlay.png') }}
        right={{ name: '缺陷区域图像', url: artifactUrl(artifacts, 'mask.png') }}
      />
    );
  }
  if (taskKey === 'stain') {
    if (!showDetectionRegions(taskKey, result)) {
      return null;
    }
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
      return null;
    }
    if (!showDetectionRegions(taskKey, result)) {
      return null;
    }
    return <FlatnessRegionPreview artifacts={artifacts} regionId={String(pageIndex)} />;
  }
  return null;
}

export function DetectionPreview({
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

  const preview = taskArtifactPreview(
    taskKey,
    resultArtifacts(result),
    resultPageIndex(taskKey, result, groupIndex),
    result,
  );
  if (!preview) {
    return <OriginalImagePreview image={image} loading={false} originalUrl={originalUrl} />;
  }

  return preview;
}

export function ReviewPanel({ review }: { review: ProjectDetectionReviewProps }): ReactElement {
  const reviewOptions: Array<{ label: string; value: ProjectDetectionReviewVerdict_Value }> = [
    { label: '准确', value: ProjectDetectionReviewVerdict_Value.Accurate },
    { label: '不准确', value: ProjectDetectionReviewVerdict_Value.Inaccurate },
  ];
  const accurateSelected = review.verdict === ProjectDetectionReviewVerdict_Value.Accurate;
  const inaccurateSelected = review.verdict === ProjectDetectionReviewVerdict_Value.Inaccurate;
  const inputPlaceholder = review.readOnly ? '无补充评论与修正' : '在此补充评论与修正，为最终评估引入专家判断';

  return (
    <div className="mt-4 flex items-center gap-3">
      <div className="flex shrink-0 items-center gap-2">
        <RobotOutlined className="text-[#1677FF]" />
        <span className="text-sm font-medium text-slate-900">你认为 Agent 智能检测结果</span>
      </div>
      <Space.Compact>
        {reviewOptions.map((option) => {
          const selected = review.verdict === option.value;
          const accurateReadOnlyClassName = 'disabled:!border-blue-200 disabled:!bg-blue-50 disabled:!text-[#1677FF]';
          const inaccurateReadOnlyClassName = 'disabled:!border-red-200 disabled:!bg-red-50 disabled:!text-red-500';
          const readOnlySelectedClassName =
            option.value === ProjectDetectionReviewVerdict_Value.Accurate
              ? accurateReadOnlyClassName
              : inaccurateReadOnlyClassName;
          return (
            <Button
              className={
                review.readOnly && selected
                  ? readOnlySelectedClassName
                  : option.value === ProjectDetectionReviewVerdict_Value.Inaccurate
                    ? INACCURATE_BUTTON_HOVER_CLASS_NAME
                    : undefined
              }
              danger={option.value === ProjectDetectionReviewVerdict_Value.Inaccurate && inaccurateSelected}
              disabled={review.saving || review.readOnly}
              icon={
                option.value === ProjectDetectionReviewVerdict_Value.Accurate ? <LikeOutlined /> : <DislikeOutlined />
              }
              key={option.value}
              onClick={() => {
                if (review.readOnly) {
                  return;
                }
                review.onVerdictChange(selected ? ProjectDetectionReviewVerdict_Value.Unknown : option.value);
              }}
              type={
                selected || (option.value === ProjectDetectionReviewVerdict_Value.Accurate && accurateSelected)
                  ? BUTTON_TYPE_PRIMARY
                  : BUTTON_TYPE_DEFAULT
              }
            >
              {option.label}
            </Button>
          );
        })}
      </Space.Compact>
      <Input
        className="min-w-0 flex-1 read-only:!bg-white"
        disabled={review.saving}
        maxLength={MAX_REVIEW_COMMENT_LENGTH}
        onBlur={() => {
          if (!review.readOnly) {
            review.onSave();
          }
        }}
        onChange={(event) => {
          review.onCommentChange(event.target.value);
        }}
        placeholder={inputPlaceholder}
        readOnly={review.readOnly}
        value={review.comment}
      />
    </div>
  );
}

export function ReviewPanelPlaceholder(): ReactElement {
  return <div aria-hidden="true" className="mt-4 min-h-8" />;
}
