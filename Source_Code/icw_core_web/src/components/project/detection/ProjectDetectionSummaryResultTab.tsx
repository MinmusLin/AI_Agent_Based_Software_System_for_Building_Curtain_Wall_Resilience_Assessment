import type { ReactElement } from 'react';

import type { ProjectDetectionSummaryResult } from '@/gen/core/common';
import { formatDateTime } from '@/utils/datetime';

const EMPTY_SECTION_COUNT = 0;
const SUMMARY_TEXT_FIELD_PATTERN = /\s+/gu;

type AgentSummaryPayload = Partial<
  Record<'corrosion' | 'crack' | 'flatness' | 'recommendation' | 'spalling' | 'stain' | 'summary', string>
>;

interface AgentSummarySection {
  color: string;
  label: string;
  text: string;
}

const AGENT_SUMMARY_FIELDS: Array<{ color: string; key: keyof AgentSummaryPayload; label: string }> = [
  { color: 'blue', key: 'summary', label: '总体概述' },
  { color: 'red', key: 'corrosion', label: '金属锈蚀' },
  { color: 'red', key: 'crack', label: '石材裂缝' },
  { color: 'red', key: 'stain', label: '石材污渍' },
  { color: 'red', key: 'flatness', label: '玻璃平整度' },
  { color: 'red', key: 'spalling', label: '玻璃爆裂' },
  { color: 'green', key: 'recommendation', label: '后续建议' },
];

function imageDateText(value?: string): string {
  if (!value) {
    return '-';
  }
  return formatDateTime(value, true);
}

function compactSummaryText(value: unknown): string {
  if (typeof value !== 'string') {
    return '';
  }
  return value.replace(SUMMARY_TEXT_FIELD_PATTERN, '').trim();
}

function summarySections(result?: string): AgentSummarySection[] {
  const text = result?.trim() ?? '';
  if (text === '') {
    return [];
  }
  try {
    const payload = JSON.parse(text) as AgentSummaryPayload;
    return AGENT_SUMMARY_FIELDS.flatMap((field) => {
      const fieldText = compactSummaryText(payload[field.key]);
      if (fieldText === '') {
        return [];
      }
      return [{ color: field.color, label: field.label, text: fieldText }];
    });
  } catch {
    const fallbackText = compactSummaryText(text);
    return fallbackText === '' ? [] : [{ color: 'blue', label: '总体概述', text: fallbackText }];
  }
}

function summaryTagClassName(color: string): string {
  const baseClassName =
    'inline-flex h-6 w-20 shrink-0 items-center justify-center rounded border px-2 text-xs font-medium leading-none';
  switch (color) {
    case 'red':
      return `${baseClassName} border-red-200 bg-red-50 text-red-600`;
    case 'green':
      return `${baseClassName} border-green-200 bg-green-50 text-green-600`;
    case 'blue':
    default:
      return `${baseClassName} border-blue-200 bg-blue-50 text-[#1677FF]`;
  }
}

export function SummaryResultTab({ result }: { result: ProjectDetectionSummaryResult | undefined }): ReactElement {
  const sections = summarySections(result?.result);

  return (
    <div className="flex h-full min-h-0 flex-col gap-4 text-sm text-slate-600">
      <div>
        <div className="mb-1 font-medium text-slate-900">任务 ID</div>
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
        <div className="min-h-0 flex-1 space-y-3 overflow-auto overscroll-contain text-sm leading-6 text-slate-700">
          {sections.length > EMPTY_SECTION_COUNT ? (
            sections.map((section) => (
              <div className="flex items-start gap-2" key={section.label}>
                <span className={summaryTagClassName(section.color)}>{section.label}</span>
                <span className="break-all">{section.text}</span>
              </div>
            ))
          ) : (
            <span>-</span>
          )}
        </div>
      </div>
    </div>
  );
}
