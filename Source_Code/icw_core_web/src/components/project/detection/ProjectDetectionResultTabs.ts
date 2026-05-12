import type { GetImageDetectionResultResponse } from '@/gen/core/api/project_detection';
import type { ViewerImage } from '@/utils/assetsStage';
import { FIRST_INDEX } from '@/utils/assetsStage';

import type { DetectionResult } from './ProjectDetectionResultViewerParts';

const REASONING_TASK_LABELS: Record<string, string> = {
  corrosion: '金属锈蚀',
  crack: '石材裂缝',
  flatness: '玻璃平整度',
  spalling: '玻璃爆裂',
  stain: '石材污渍',
};

export interface DetectionTabItem {
  key: string;
  label: string;
  result: DetectionResult | undefined;
}

export function detectionResultTabs(result: GetImageDetectionResultResponse | null): DetectionTabItem[] {
  if (!result) {
    return [];
  }
  const tabs: DetectionTabItem[] = [
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
  ];
  return tabs.filter((item) => Boolean(item.result));
}

export function viewerTitle(
  image: GetImageDetectionResultResponse['image'],
  uploadedImages: ViewerImage[],
  viewerIndex: number,
): string {
  if (image?.file_name) {
    return image.file_name;
  }
  if (viewerIndex >= FIRST_INDEX && uploadedImages[viewerIndex]) {
    return uploadedImages[viewerIndex].image.file_name;
  }
  return '检测结果';
}
