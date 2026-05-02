export const PROJECT_STAGE_KEYS = ['profile', 'assets', 'detection', 'review', 'report'] as const;

export type ProjectStageKey = (typeof PROJECT_STAGE_KEYS)[number];

export interface ProjectStageMeta {
  description: string;
  title: string;
}

export const PROJECT_STAGES: ProjectStageMeta[] = [
  {
    title: '项目基础信息',
    description: '完善项目、建筑与评估目标等背景信息，为后续 AI 模型推理和韧性评估报告生成提供稳定上下文',
  },
  {
    title: '图像资产构建',
    description: '按建筑立面或区域创建图像组并上传建筑幕墙原始图像，形成可追踪、可管理的项目幕墙图像资产库',
  },
  {
    title: 'Agent 智能检测',
    description: 'AI Agent 批量调用建筑材料分类与多维度缺陷检测能力，持续生成图像级指标、掩码和结构化检测结果',
  },
  {
    title: '人工复核确认',
    description: '对关键区域幕墙图像和多维度缺陷检测结果进行人工复核确认、补充评论与修正，为最终评估引入专家判断',
  },
  {
    title: '评估报告生成',
    description: '汇总建筑背景、分类结果、多维度缺陷检测指标、专家意见与 AI 知识库，生成结构化的建筑幕墙韧性评估报告',
  },
];

export const LAST_VISIBLE_PROGRESS = PROJECT_STAGE_KEYS.length - 1;

export function stageKeyFromProgress(progress: number): ProjectStageKey {
  const visibleProgress = Math.min(Math.max(progress, 0), LAST_VISIBLE_PROGRESS);
  return PROJECT_STAGE_KEYS[visibleProgress];
}

export function progressFromStageKey(stageKey: string | undefined): number | null {
  if (!stageKey) {
    return null;
  }

  const progress = PROJECT_STAGE_KEYS.findIndex((key) => key === stageKey);
  if (progress < 0) {
    return null;
  }
  return progress;
}
