import {
  AuditOutlined,
  BarChartOutlined,
  ClockCircleOutlined,
  ClusterOutlined,
  FileTextOutlined,
  InfoCircleOutlined,
  PictureOutlined,
  SafetyCertificateOutlined,
  ToolOutlined,
  WarningOutlined,
} from '@ant-design/icons';
import type { ReactElement, ReactNode } from 'react';

import { formatDateTime } from '@/utils/datetime';

const DISPLAY_INDEX_OFFSET = 1;
const EMPTY_ITEM_COUNT = 0;

export interface ProjectReportPayload {
  dataset_overview: string;
  group_count?: number;
  image_count?: number;
  limitations: string;
  maintenance_and_improvement_recommendations: string[];
  overall_conclusion: ProjectReportOverallConclusion;
  overall_resilience_assessment: string[];
  overall_risk_level: string;
  project_background: string;
  regional_assessments: ProjectReportRegionalAssessment[];
  resilience_assessment_objective: string;
}

interface ProjectReportOverallConclusion {
  curtain_wall_overall_status: string;
  defect_concentrated_areas: string;
  main_defect_types: string;
  risk_tendency: string;
}

interface ProjectReportRegionalAssessment {
  regional_conclusion: ProjectReportRegionalConclusion;
  regional_risk_level: string;
  section_title: string;
}

interface ProjectReportRegionalConclusion {
  curtain_wall_region_status: string;
  main_defect_types: string;
  risk_tendency: string;
}

interface ProjectReportViewerProps {
  fallbackText: string;
  payload?: ProjectReportPayload;
  updatedAt?: string;
}

interface InfoCardProps {
  icon: ReactElement;
  label: string;
  value: ReactNode;
  valueClassName?: string;
}

interface ReportSectionProps {
  children: ReactNode;
  icon: ReactElement;
  title: string;
}

interface ParagraphBlockProps {
  label: string;
  text: string;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null && !Array.isArray(value);
}

function stringField(record: Record<string, unknown>, key: string): string {
  const value = record[key];
  return typeof value === 'string' ? value.trim() : '';
}

function numberField(record: Record<string, unknown>, key: string): number | undefined {
  const value = record[key];
  return typeof value === 'number' ? value : undefined;
}

function stringArrayField(record: Record<string, unknown>, key: string): string[] {
  const value = record[key];
  if (!Array.isArray(value)) {
    return [];
  }
  return value.filter((item): item is string => typeof item === 'string').map((item) => item.trim());
}

function parseOverallConclusion(value: unknown): ProjectReportOverallConclusion {
  const record = isRecord(value) ? value : {};
  return {
    curtain_wall_overall_status: stringField(record, 'curtain_wall_overall_status'),
    defect_concentrated_areas: stringField(record, 'defect_concentrated_areas'),
    main_defect_types: stringField(record, 'main_defect_types'),
    risk_tendency: stringField(record, 'risk_tendency'),
  };
}

function parseRegionalConclusion(value: unknown): ProjectReportRegionalConclusion {
  const record = isRecord(value) ? value : {};
  return {
    curtain_wall_region_status: stringField(record, 'curtain_wall_region_status'),
    main_defect_types: stringField(record, 'main_defect_types'),
    risk_tendency: stringField(record, 'risk_tendency'),
  };
}

function parseRegionalAssessments(value: unknown): ProjectReportRegionalAssessment[] {
  if (!Array.isArray(value)) {
    return [];
  }
  return value.filter(isRecord).map((item) => ({
    regional_conclusion: parseRegionalConclusion(item.regional_conclusion),
    regional_risk_level: stringField(item, 'regional_risk_level'),
    section_title: stringField(item, 'section_title'),
  }));
}

export function parseProjectReportPayload(resultJson: string): ProjectReportPayload | undefined {
  const text = resultJson.trim();
  if (text === '') {
    return undefined;
  }
  try {
    const parsed = JSON.parse(text) as unknown;
    if (!isRecord(parsed)) {
      return undefined;
    }
    return {
      dataset_overview: stringField(parsed, 'dataset_overview'),
      group_count: numberField(parsed, 'group_count'),
      image_count: numberField(parsed, 'image_count'),
      limitations: stringField(parsed, 'limitations'),
      maintenance_and_improvement_recommendations: stringArrayField(
        parsed,
        'maintenance_and_improvement_recommendations',
      ),
      overall_conclusion: parseOverallConclusion(parsed.overall_conclusion),
      overall_resilience_assessment: stringArrayField(parsed, 'overall_resilience_assessment'),
      overall_risk_level: stringField(parsed, 'overall_risk_level'),
      project_background: stringField(parsed, 'project_background'),
      regional_assessments: parseRegionalAssessments(parsed.regional_assessments),
      resilience_assessment_objective: stringField(parsed, 'resilience_assessment_objective'),
    };
  } catch {
    return undefined;
  }
}

function formatReportTime(value?: string): string {
  if (!value) {
    return '暂无';
  }
  return formatDateTime(value, true) || value;
}

function riskLevelClassName(value: string): string {
  const baseClassName = 'inline-flex h-7 min-w-12 items-center justify-center rounded border px-3 text-sm font-medium';
  switch (value) {
    case '高':
      return `${baseClassName} border-red-200 bg-red-50 text-red-600`;
    case '中':
      return `${baseClassName} border-amber-200 bg-amber-50 text-amber-600`;
    case '低':
      return `${baseClassName} border-green-200 bg-green-50 text-green-600`;
    default:
      return `${baseClassName} border-slate-200 bg-slate-50 text-slate-600`;
  }
}

function ReportSection({ children, icon, title }: ReportSectionProps): ReactElement {
  return (
    <section className="rounded-lg border border-slate-200 bg-white p-5">
      <h3 className="mb-4 flex items-center gap-2 text-base font-semibold text-slate-900">
        <span className="inline-flex h-7 w-7 items-center justify-center rounded bg-blue-50 text-[#1677FF]">
          {icon}
        </span>
        <span>{title}</span>
      </h3>
      {children}
    </section>
  );
}

function ParagraphBlock({ label, text }: ParagraphBlockProps): ReactElement {
  return (
    <div>
      <div className="mb-1 text-sm font-medium text-slate-900">{label}</div>
      <p className="text-sm leading-7 text-slate-600">{text || '-'}</p>
    </div>
  );
}

function InfoCard({
  icon,
  label,
  value,
  valueClassName = 'text-2xl font-semibold text-slate-950',
}: InfoCardProps): ReactElement {
  return (
    <div className="rounded-lg border border-slate-200 bg-slate-50 px-4 py-3">
      <div className="flex items-center gap-2 text-sm text-slate-500">
        <span className="inline-flex h-7 w-7 items-center justify-center rounded bg-white text-[#1677FF]">{icon}</span>
        <span>{label}</span>
      </div>
      <div className={`mt-2 ${valueClassName}`}>{value}</div>
    </div>
  );
}

function TextList({ items }: { items: string[] }): ReactElement {
  if (items.length === EMPTY_ITEM_COUNT) {
    return <div className="text-sm text-slate-500">-</div>;
  }
  return (
    <ol className="space-y-3">
      {items.map((item, index) => (
        <li className="flex gap-3 text-sm leading-7 text-slate-600" key={`${item}-${index.toString()}`}>
          <span className="mt-1 inline-flex h-5 w-5 shrink-0 items-center justify-center rounded-full bg-blue-50 text-xs font-semibold text-[#1677FF]">
            {(index + DISPLAY_INDEX_OFFSET).toString()}
          </span>
          <span>{item}</span>
        </li>
      ))}
    </ol>
  );
}

function ReportFallback({ text }: { text: string }): ReactElement {
  return (
    <pre className="h-full overflow-auto overscroll-contain whitespace-pre-wrap break-all rounded border border-slate-200 bg-white p-4 text-sm leading-6 text-slate-700">
      {text || '报告生成成功'}
    </pre>
  );
}

function RegionalAssessments({ items }: { items: ProjectReportRegionalAssessment[] }): ReactElement {
  if (items.length === EMPTY_ITEM_COUNT) {
    return <div className="text-sm text-slate-500">-</div>;
  }
  return (
    <div className="space-y-4">
      {items.map((item) => (
        <div className="rounded-lg border border-slate-200 bg-slate-50 p-4" key={item.section_title}>
          <div className="mb-3 flex items-center justify-between gap-3">
            <h4 className="flex items-center gap-2 text-sm font-semibold text-slate-900">
              <ClusterOutlined className="text-[#1677FF]" />
              <span>{item.section_title || '未命名区域'}</span>
            </h4>
            <span className={riskLevelClassName(item.regional_risk_level)}>{item.regional_risk_level || '未知'}</span>
          </div>
          <div className="space-y-3">
            <ParagraphBlock label="区域状态" text={item.regional_conclusion.curtain_wall_region_status} />
            <ParagraphBlock label="主要缺陷类型" text={item.regional_conclusion.main_defect_types} />
            <ParagraphBlock label="风险倾向" text={item.regional_conclusion.risk_tendency} />
          </div>
        </div>
      ))}
    </div>
  );
}

function ProjectReportContent({
  payload,
  updatedAt,
}: {
  payload: ProjectReportPayload;
  updatedAt?: string;
}): ReactElement {
  return (
    <div className="h-full overflow-auto overscroll-contain">
      <div className="space-y-5">
        <ReportSection icon={<BarChartOutlined />} title="报告概览">
          <div className="grid gap-3 md:grid-cols-2 xl:grid-cols-4">
            <InfoCard icon={<ClusterOutlined />} label="图像组数量" value={payload.group_count ?? '-'} />
            <InfoCard icon={<PictureOutlined />} label="图像数量" value={payload.image_count ?? '-'} />
            <InfoCard
              icon={<WarningOutlined />}
              label="总体风险等级"
              value={
                <span className={riskLevelClassName(payload.overall_risk_level)}>
                  {payload.overall_risk_level || '未知'}
                </span>
              }
              valueClassName="text-sm"
            />
            <InfoCard
              icon={<ClockCircleOutlined />}
              label="报告生成时间"
              value={formatReportTime(updatedAt)}
              valueClassName="text-base font-semibold text-slate-950"
            />
          </div>
        </ReportSection>
        <ReportSection icon={<FileTextOutlined />} title="项目背景">
          <div className="space-y-4">
            <ParagraphBlock label="建筑背景" text={payload.project_background} />
            <ParagraphBlock label="韧性评估目标" text={payload.resilience_assessment_objective} />
            <ParagraphBlock label="数据集概述" text={payload.dataset_overview} />
          </div>
        </ReportSection>
        <ReportSection icon={<SafetyCertificateOutlined />} title="总体评估结论">
          <div className="grid gap-4 lg:grid-cols-2">
            <ParagraphBlock label="幕墙整体状态" text={payload.overall_conclusion.curtain_wall_overall_status} />
            <ParagraphBlock label="主要缺陷类型" text={payload.overall_conclusion.main_defect_types} />
            <ParagraphBlock label="缺陷集中区域" text={payload.overall_conclusion.defect_concentrated_areas} />
            <ParagraphBlock label="风险倾向" text={payload.overall_conclusion.risk_tendency} />
          </div>
        </ReportSection>
        <ReportSection icon={<ClusterOutlined />} title="区域评估结论">
          <RegionalAssessments items={payload.regional_assessments} />
        </ReportSection>
        <ReportSection icon={<AuditOutlined />} title="总体韧性评估">
          <TextList items={payload.overall_resilience_assessment} />
        </ReportSection>
        <ReportSection icon={<ToolOutlined />} title="维护和改进建议">
          <TextList items={payload.maintenance_and_improvement_recommendations} />
        </ReportSection>
        <ReportSection icon={<InfoCircleOutlined />} title="报告限制说明">
          <p className="text-sm leading-7 text-slate-600">{payload.limitations || '-'}</p>
        </ReportSection>
      </div>
    </div>
  );
}

export function ProjectReportViewer({ fallbackText, payload, updatedAt }: ProjectReportViewerProps): ReactElement {
  if (!payload) {
    return <ReportFallback text={fallbackText} />;
  }
  return <ProjectReportContent payload={payload} updatedAt={updatedAt} />;
}
