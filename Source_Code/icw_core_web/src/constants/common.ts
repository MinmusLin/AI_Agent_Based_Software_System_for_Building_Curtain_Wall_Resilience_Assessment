import {
  AvatarType_Value,
  DetectionNodeCode_Value,
  DetectionTaskCode_Value,
  ProjectEventType_Value,
  ProjectProgress_Value,
  SocketScope_Value,
} from '@/gen/core/common';

export const AVATAR_TYPE_CUSTOM = AvatarType_Value.Custom;
export const AVATAR_TYPE_DEFAULT = AvatarType_Value.Default;
export const AVATAR_TYPE_NONE = AvatarType_Value.None;

export const PROJECT_EVENT_TYPE_IMAGE_STATUS_CHANGED = ProjectEventType_Value.ImageStatusChanged;
export const PROJECT_EVENT_TYPE_DETECTION_TASK_STATUS_CHANGED = ProjectEventType_Value.DetectionTaskStatusChanged;
export const PROJECT_EVENT_TYPE_REPORT_STATUS_CHANGED = ProjectEventType_Value.ReportStatusChanged;

export const SOCKET_SCOPE_PROJECT_ASSETS = SocketScope_Value.ProjectAssets;
export const SOCKET_SCOPE_PROJECT_DETECTION = SocketScope_Value.ProjectDetection;
export const SOCKET_SCOPE_PROJECT_REPORT = SocketScope_Value.ProjectReport;

export const PROJECT_DETECTION_NODE_CLASSIFICATION = detectionNodeCodeValueToProtocol(
  DetectionNodeCode_Value.Classification,
);
export const PROJECT_DETECTION_NODE_SUMMARY = detectionNodeCodeValueToProtocol(DetectionNodeCode_Value.Summary);
export const PROJECT_DETECTION_NODE_REASONING_PREFIX = `${detectionNodeCodeValueToProtocol(
  DetectionNodeCode_Value.Reasoning,
)}:`;

export const DETECTION_PROGRESS_NODE_STATUS_FAILED = 'failed';
export const DETECTION_PROGRESS_NODE_STATUS_PENDING = 'pending';
export const DETECTION_PROGRESS_NODE_STATUS_RUNNING = 'running';
export const DETECTION_PROGRESS_NODE_STATUS_SKIPPED = 'skipped';
export const DETECTION_PROGRESS_NODE_STATUS_SUCCEEDED = 'succeeded';

export type DetectionProgressNodeStatus =
  | typeof DETECTION_PROGRESS_NODE_STATUS_FAILED
  | typeof DETECTION_PROGRESS_NODE_STATUS_PENDING
  | typeof DETECTION_PROGRESS_NODE_STATUS_RUNNING
  | typeof DETECTION_PROGRESS_NODE_STATUS_SKIPPED
  | typeof DETECTION_PROGRESS_NODE_STATUS_SUCCEEDED;

export const AUTH_STATUS_INITIALIZING = 'initializing';
export const AUTH_STATUS_AUTHENTICATED = 'authenticated';
export const AUTH_STATUS_ANONYMOUS = 'anonymous';

export type AuthStatus =
  | typeof AUTH_STATUS_INITIALIZING
  | typeof AUTH_STATUS_AUTHENTICATED
  | typeof AUTH_STATUS_ANONYMOUS;

export const PROJECT_STAGE_KEY_PROFILE = 'profile';
export const PROJECT_STAGE_KEY_ASSETS = 'assets';
export const PROJECT_STAGE_KEY_DETECTION = 'detection';
export const PROJECT_STAGE_KEY_REVIEW = 'review';
export const PROJECT_STAGE_KEY_REPORT = 'report';

export const PROJECT_STAGE_KEYS = [
  PROJECT_STAGE_KEY_PROFILE,
  PROJECT_STAGE_KEY_ASSETS,
  PROJECT_STAGE_KEY_DETECTION,
  PROJECT_STAGE_KEY_REVIEW,
  PROJECT_STAGE_KEY_REPORT,
] as const;

export type ProjectStageKey = (typeof PROJECT_STAGE_KEYS)[number];

export const PROJECT_PROGRESS_VALUES = [
  ProjectProgress_Value.InitializationFinished,
  ProjectProgress_Value.ProfileFinished,
  ProjectProgress_Value.AssetsFinished,
  ProjectProgress_Value.DetectionFinished,
  ProjectProgress_Value.ReviewFinished,
  ProjectProgress_Value.ReportFinished,
] as const;

export const LAST_VISIBLE_PROJECT_PROGRESS = ProjectProgress_Value.ReviewFinished;
export const LAST_VISIBLE_PROGRESS = LAST_VISIBLE_PROJECT_PROGRESS;

export const HTTP_STATUS_UNAUTHORIZED = 401;
export const HTTP_STATUS_INTERNAL_SERVER_ERROR = 500;

const NOT_FOUND_STAGE_INDEX = -1;

export interface ApiEnvelope<T> {
  code: string;
  message: string;
  data: T;
}

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

export function normalizeProjectProgressValue(progress: unknown): ProjectProgress_Value {
  const numericProgress = typeof progress === 'number' ? progress : Number(progress);
  if (PROJECT_PROGRESS_VALUES.includes(numericProgress)) {
    return numericProgress;
  }
  return ProjectProgress_Value.InitializationFinished;
}

export function detectionNodeCodeValueToProtocol(code: DetectionNodeCode_Value): string {
  switch (code) {
    case DetectionNodeCode_Value.Classification:
      return 'classification';
    case DetectionNodeCode_Value.Reasoning:
      return 'reasoning';
    case DetectionNodeCode_Value.Summary:
      return 'summary';
    case DetectionNodeCode_Value.Unknown:
    case DetectionNodeCode_Value.UNRECOGNIZED:
      return '';
  }
}

export function detectionTaskCodeValueToProtocol(code: DetectionTaskCode_Value): string {
  switch (code) {
    case DetectionTaskCode_Value.Corrosion:
      return 'corrosion';
    case DetectionTaskCode_Value.Crack:
      return 'crack';
    case DetectionTaskCode_Value.Stain:
      return 'stain';
    case DetectionTaskCode_Value.Flatness:
      return 'flatness';
    case DetectionTaskCode_Value.Spalling:
      return 'spalling';
    case DetectionTaskCode_Value.Unknown:
    case DetectionTaskCode_Value.UNRECOGNIZED:
      return '';
  }
}

export function reasoningDetectionNodeCode(taskCode: DetectionTaskCode_Value): string {
  return `${PROJECT_DETECTION_NODE_REASONING_PREFIX}${detectionTaskCodeValueToProtocol(taskCode)}`;
}

export function stageKeyFromProgress(progress: unknown): ProjectStageKey {
  const normalizedProgress = normalizeProjectProgressValue(progress);
  const visibleProgress = Math.min(
    Math.max(normalizedProgress, ProjectProgress_Value.InitializationFinished),
    LAST_VISIBLE_PROGRESS,
  );
  return PROJECT_STAGE_KEYS[visibleProgress] ?? PROJECT_STAGE_KEY_PROFILE;
}

export function progressFromStageKey(stageKey: string | undefined): ProjectProgress_Value | null {
  if (!stageKey) {
    return null;
  }

  const progress = PROJECT_STAGE_KEYS.findIndex((key) => key === stageKey);
  if (progress === NOT_FOUND_STAGE_INDEX) {
    return null;
  }
  return progress;
}
