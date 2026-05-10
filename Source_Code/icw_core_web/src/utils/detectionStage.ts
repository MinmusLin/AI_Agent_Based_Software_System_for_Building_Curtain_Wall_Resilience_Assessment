import {
  PROJECT_DETECTION_NODE_CLASSIFICATION,
  PROJECT_DETECTION_NODE_REASONING_PREFIX,
  PROJECT_DETECTION_NODE_SUMMARY,
  PROJECT_EVENT_TYPE_DETECTION_TASK_STATUS_CHANGED,
} from '@/constants/common';
import type { ProjectDetectionTaskStatusChangedMessage } from '@/gen/core/api/messages';
import {
  type ProjectDetectionNodeStatus,
  type ProjectDetectionStatus,
  ProjectDetectionSubTaskStatus_Value,
  ProjectDetectionTaskStatus_Value,
} from '@/gen/core/common';

export const DETECTION_PROGRESS_TOTAL_STEPS = 7;
export const DETECTION_CLASSIFICATION_COMPLETED_STEPS = 1;
export const DETECTION_REASONING_COMPLETED_STEPS = 6;
export const DETECTION_SUMMARY_COMPLETED_STEPS = 7;
export const PERCENT_FULL = 100;
export const EMPTY_DETECTION_STEPS = 0;
export const DETECTION_TASK_COUNT_INCREMENT = 1;

const EMPTY_STRING = '';
const STATUS_RANK_INITIAL = 0;
const STATUS_RANK_PENDING = 1;
const STATUS_RANK_CLASSIFYING = 2;
const STATUS_RANK_DETECTING = 3;
const STATUS_RANK_SUMMARIZING = 4;
const STATUS_RANK_TERMINAL = 5;

export type ProjectDetectionStatusMap = Partial<Record<string, ProjectDetectionStatus>>;

interface DetectionActionStateParams {
  detectionCompleted: boolean;
  detectionFailed: boolean;
  detectionStarted: boolean;
  pageLoading: boolean;
  readOnly: boolean;
  uploadedImageCount: number;
}

interface DetectionActionState {
  canComplete: boolean;
  canRetry: boolean;
  canStart: boolean;
  showComplete: boolean;
  showRetry: boolean;
  showStart: boolean;
}

interface ProjectDetectionTaskStatusChangedRecord extends Record<string, unknown> {
  image_uuid: string;
  main_status: ProjectDetectionTaskStatus_Value;
  main_task_uuid: string;
  node_code: string;
  occurred_at: string;
  project_id: string;
  sub_status: ProjectDetectionSubTaskStatus_Value;
  sub_task_uuid: string;
}

export function detectionTasksMap(tasks: ProjectDetectionStatus[] | undefined): ProjectDetectionStatusMap {
  return (tasks ?? []).reduce<ProjectDetectionStatusMap>((taskMap, task) => {
    if (task.image_uuid !== EMPTY_STRING) {
      taskMap[task.image_uuid] = normalizeDetectionTask(task);
    }
    return taskMap;
  }, {});
}

export function mergeDetectionTaskMessage(
  currentTasks: ProjectDetectionStatusMap,
  message: ProjectDetectionTaskStatusChangedMessage,
): ProjectDetectionStatusMap {
  const currentTask = currentTasks[message.image_uuid];
  if (isStaleDetectionTaskMessage(currentTask, message)) {
    return currentTasks;
  }

  const shouldReplaceTask = shouldReplaceDetectionTask(currentTask, message);
  const nextTask = normalizeDetectionTask(
    shouldReplaceTask || !currentTask
      ? {
          classification_status: undefined,
          detection_status: [],
          finished_at: '',
          image_uuid: message.image_uuid,
          main_status: message.main_status,
          main_task_uuid: message.main_task_uuid,
          started_at: '',
          summary_status: undefined,
        }
      : currentTask,
  );

  nextTask.main_status = mergeProjectDetectionTaskStatus(nextTask.main_status, message.main_status);
  nextTask.main_task_uuid = message.main_task_uuid;

  if (message.node_code !== EMPTY_STRING) {
    assignDetectionNode(nextTask, {
      node_code: message.node_code,
      sub_status: message.sub_status,
      sub_task_uuid: message.sub_task_uuid,
    });
  }
  inferPreviousDetectionNodes(nextTask);

  return {
    ...currentTasks,
    [message.image_uuid]: nextTask,
  };
}

function isStaleDetectionTaskMessage(
  currentTask: ProjectDetectionStatus | undefined,
  message: ProjectDetectionTaskStatusChangedMessage,
): boolean {
  if (!currentTask || currentTask.main_task_uuid === EMPTY_STRING || message.main_task_uuid === EMPTY_STRING) {
    return false;
  }
  if (currentTask.main_task_uuid === message.main_task_uuid) {
    return false;
  }
  return !isDetectionTaskTerminal(currentTask);
}

function shouldReplaceDetectionTask(
  currentTask: ProjectDetectionStatus | undefined,
  message: ProjectDetectionTaskStatusChangedMessage,
): boolean {
  if (!currentTask || currentTask.main_task_uuid === EMPTY_STRING || message.main_task_uuid === EMPTY_STRING) {
    return false;
  }
  if (currentTask.main_task_uuid === message.main_task_uuid) {
    return false;
  }
  return (
    isDetectionTaskTerminal(currentTask) &&
    message.main_status !== ProjectDetectionTaskStatus_Value.Failed &&
    message.main_status !== ProjectDetectionTaskStatus_Value.Succeeded
  );
}

export function detectionTaskCompletedSteps(task: ProjectDetectionStatus | undefined): number {
  if (!task) {
    return EMPTY_DETECTION_STEPS;
  }
  if (isTaskSucceeded(task) || hasSummaryStatus(task, ProjectDetectionSubTaskStatus_Value.Succeeded)) {
    return DETECTION_SUMMARY_COMPLETED_STEPS;
  }
  if (isTaskFailed(task) || hasSummaryStatus(task, ProjectDetectionSubTaskStatus_Value.Failed)) {
    return DETECTION_SUMMARY_COMPLETED_STEPS;
  }
  if (isTaskSummarizing(task) || hasSummaryStatus(task, ProjectDetectionSubTaskStatus_Value.Pending)) {
    return DETECTION_REASONING_COMPLETED_STEPS;
  }

  const reasoningNodes = taskReasoningNodes(task);
  if (areReasoningNodesSucceeded(reasoningNodes)) {
    return DETECTION_REASONING_COMPLETED_STEPS;
  }

  const reasoningSucceededSteps = reasoningNodes.filter(
    (node) => node.sub_status === ProjectDetectionSubTaskStatus_Value.Succeeded,
  ).length;
  if (reasoningSucceededSteps > EMPTY_DETECTION_STEPS) {
    return DETECTION_CLASSIFICATION_COMPLETED_STEPS + reasoningSucceededSteps;
  }
  if (
    isTaskDetecting(task) ||
    hasNodeStatus(task, PROJECT_DETECTION_NODE_CLASSIFICATION, ProjectDetectionSubTaskStatus_Value.Succeeded)
  ) {
    return DETECTION_CLASSIFICATION_COMPLETED_STEPS;
  }
  return EMPTY_DETECTION_STEPS;
}

export function detectionTaskProgressPercent(task: ProjectDetectionStatus | undefined): number {
  return Math.round((detectionTaskCompletedSteps(task) / DETECTION_PROGRESS_TOTAL_STEPS) * PERCENT_FULL);
}

export function isDetectionTaskTerminal(task: ProjectDetectionStatus | undefined): boolean {
  return (
    task?.main_status === ProjectDetectionTaskStatus_Value.Succeeded ||
    task?.main_status === ProjectDetectionTaskStatus_Value.Failed
  );
}

export function isDetectionTaskSucceeded(task: ProjectDetectionStatus | undefined): boolean {
  return task?.main_status === ProjectDetectionTaskStatus_Value.Succeeded;
}

export function isDetectionTaskRunning(task: ProjectDetectionStatus | undefined): boolean {
  return task !== undefined && !isDetectionTaskTerminal(task);
}

export function hasDetectionTasks(taskMap: ProjectDetectionStatusMap): boolean {
  return Object.keys(taskMap).length > EMPTY_DETECTION_STEPS;
}

export function hasFailedDetectionTask(taskMap: ProjectDetectionStatusMap): boolean {
  return Object.values(taskMap).some((task) => task?.main_status === ProjectDetectionTaskStatus_Value.Failed);
}

export function allDetectionTasksSucceeded(taskMap: ProjectDetectionStatusMap, imageCount: number): boolean {
  const tasks = Object.values(taskMap);
  return (
    imageCount > EMPTY_DETECTION_STEPS &&
    tasks.length >= imageCount &&
    tasks.every((task) => task?.main_status === ProjectDetectionTaskStatus_Value.Succeeded)
  );
}

export function detectionTaskStatusStats(taskMap: ProjectDetectionStatusMap): { failed: number; running: number } {
  return Object.values(taskMap).reduce(
    (stats, task) => ({
      failed:
        stats.failed +
        (task?.main_status === ProjectDetectionTaskStatus_Value.Failed
          ? DETECTION_TASK_COUNT_INCREMENT
          : EMPTY_DETECTION_STEPS),
      running: stats.running + (isDetectionTaskRunning(task) ? DETECTION_TASK_COUNT_INCREMENT : EMPTY_DETECTION_STEPS),
    }),
    {
      failed: EMPTY_DETECTION_STEPS,
      running: EMPTY_DETECTION_STEPS,
    },
  );
}

export function detectionActionState({
  detectionCompleted,
  detectionFailed,
  detectionStarted,
  pageLoading,
  readOnly,
  uploadedImageCount,
}: DetectionActionStateParams): DetectionActionState {
  const showRetry = detectionStarted && !detectionCompleted;

  return {
    canComplete: !readOnly && !pageLoading && detectionCompleted,
    canRetry: !readOnly && !pageLoading && showRetry && detectionFailed,
    canStart: !readOnly && !pageLoading && !detectionStarted && uploadedImageCount > EMPTY_DETECTION_STEPS,
    showComplete: detectionCompleted,
    showRetry,
    showStart: !detectionStarted,
  };
}

export function parseProjectDetectionTaskStatusChangedMessage(
  data: unknown,
): ProjectDetectionTaskStatusChangedMessage | null {
  if (typeof data !== 'string') {
    return null;
  }

  let value: unknown;
  try {
    value = JSON.parse(data) as unknown;
  } catch {
    return null;
  }

  if (!isRecord(value) || value.type !== PROJECT_EVENT_TYPE_DETECTION_TASK_STATUS_CHANGED) {
    return null;
  }
  if (!isProjectDetectionTaskStatusChangedRecord(value)) {
    return null;
  }

  return {
    image_uuid: value.image_uuid,
    main_status: value.main_status,
    main_task_uuid: value.main_task_uuid,
    node_code: value.node_code,
    occurred_at: value.occurred_at,
    project_id: value.project_id,
    sub_status: value.sub_status,
    sub_task_uuid: value.sub_task_uuid,
    type: PROJECT_EVENT_TYPE_DETECTION_TASK_STATUS_CHANGED,
  };
}

function hasNodeStatus(
  task: ProjectDetectionStatus,
  nodeCode: string,
  subStatus: ProjectDetectionSubTaskStatus_Value,
): boolean {
  if (nodeCode === PROJECT_DETECTION_NODE_CLASSIFICATION) {
    return task.classification_status?.sub_status === subStatus;
  }
  if (nodeCode === PROJECT_DETECTION_NODE_SUMMARY) {
    return task.summary_status?.sub_status === subStatus;
  }
  return taskReasoningNodes(task).some((node) => node.node_code === nodeCode && node.sub_status === subStatus);
}

function hasSummaryStatus(task: ProjectDetectionStatus, subStatus: ProjectDetectionSubTaskStatus_Value): boolean {
  return task.summary_status?.sub_status === subStatus;
}

function isTaskSucceeded(task: ProjectDetectionStatus): boolean {
  return task.main_status === ProjectDetectionTaskStatus_Value.Succeeded;
}

function isTaskFailed(task: ProjectDetectionStatus): boolean {
  return task.main_status === ProjectDetectionTaskStatus_Value.Failed;
}

function isTaskSummarizing(task: ProjectDetectionStatus): boolean {
  return task.main_status === ProjectDetectionTaskStatus_Value.Summarizing;
}

function isTaskDetecting(task: ProjectDetectionStatus): boolean {
  return task.main_status === ProjectDetectionTaskStatus_Value.Detecting;
}

function taskReasoningNodes(task: ProjectDetectionStatus): ProjectDetectionNodeStatus[] {
  return task.detection_status.filter((node) => node.node_code.startsWith(PROJECT_DETECTION_NODE_REASONING_PREFIX));
}

function areReasoningNodesSucceeded(nodes: ProjectDetectionNodeStatus[]): boolean {
  return (
    nodes.length > EMPTY_DETECTION_STEPS &&
    nodes.every((node) => node.sub_status === ProjectDetectionSubTaskStatus_Value.Succeeded)
  );
}

function normalizeDetectionTask(task: ProjectDetectionStatus): ProjectDetectionStatus {
  return {
    ...task,
    detection_status: task.detection_status,
  };
}

function assignDetectionNode(task: ProjectDetectionStatus, nextNode: ProjectDetectionNodeStatus): void {
  if (nextNode.node_code === PROJECT_DETECTION_NODE_CLASSIFICATION) {
    task.classification_status = mergeProjectDetectionNode(task.classification_status, nextNode);
    return;
  }
  if (nextNode.node_code === PROJECT_DETECTION_NODE_SUMMARY) {
    task.summary_status = mergeProjectDetectionNode(task.summary_status, nextNode);
    return;
  }
  task.detection_status = upsertDetectionNode(task.detection_status, nextNode);
}

function inferPreviousDetectionNodes(task: ProjectDetectionStatus): void {
  if (
    task.main_status === ProjectDetectionTaskStatus_Value.Detecting ||
    task.main_status === ProjectDetectionTaskStatus_Value.Summarizing ||
    task.main_status === ProjectDetectionTaskStatus_Value.Succeeded ||
    (task.main_status === ProjectDetectionTaskStatus_Value.Failed &&
      task.classification_status?.sub_status !== ProjectDetectionSubTaskStatus_Value.Failed)
  ) {
    task.classification_status = {
      node_code: PROJECT_DETECTION_NODE_CLASSIFICATION,
      sub_status: ProjectDetectionSubTaskStatus_Value.Succeeded,
      sub_task_uuid: task.classification_status?.sub_task_uuid ?? EMPTY_STRING,
    };
  }
}

function upsertDetectionNode(
  nodes: ProjectDetectionNodeStatus[] | undefined,
  nextNode: ProjectDetectionNodeStatus,
): ProjectDetectionNodeStatus[] {
  const currentNodes = nodes ?? [];
  const nodeIndex = currentNodes.findIndex((node) => node.node_code === nextNode.node_code);
  if (nodeIndex < EMPTY_DETECTION_STEPS) {
    return [...currentNodes, nextNode];
  }
  return currentNodes.map((node, index) => {
    if (index !== nodeIndex) {
      return node;
    }
    return mergeProjectDetectionNode(node, nextNode);
  });
}

function mergeProjectDetectionNode(
  currentNode: ProjectDetectionNodeStatus | undefined,
  nextNode: ProjectDetectionNodeStatus,
): ProjectDetectionNodeStatus {
  if (!currentNode) {
    return nextNode;
  }
  if (
    currentNode.sub_status === ProjectDetectionSubTaskStatus_Value.Succeeded &&
    nextNode.sub_status === ProjectDetectionSubTaskStatus_Value.Failed
  ) {
    return currentNode;
  }
  if (projectDetectionSubStatusRank(nextNode.sub_status) >= projectDetectionSubStatusRank(currentNode.sub_status)) {
    return nextNode;
  }
  return currentNode;
}

function mergeProjectDetectionTaskStatus(
  currentStatus: ProjectDetectionTaskStatus_Value,
  nextStatus: ProjectDetectionTaskStatus_Value,
): ProjectDetectionTaskStatus_Value {
  if (nextStatus === ProjectDetectionTaskStatus_Value.Succeeded) {
    return nextStatus;
  }
  if (currentStatus === ProjectDetectionTaskStatus_Value.Succeeded) {
    return currentStatus;
  }
  const currentRank = projectDetectionMainStatusRank(currentStatus);
  const nextRank = projectDetectionMainStatusRank(nextStatus);
  if (nextRank >= currentRank) {
    return nextStatus;
  }
  return currentStatus;
}

function projectDetectionMainStatusRank(status: ProjectDetectionTaskStatus_Value): number {
  switch (status) {
    case ProjectDetectionTaskStatus_Value.Pending:
      return STATUS_RANK_PENDING;
    case ProjectDetectionTaskStatus_Value.Classifying:
      return STATUS_RANK_CLASSIFYING;
    case ProjectDetectionTaskStatus_Value.Detecting:
      return STATUS_RANK_DETECTING;
    case ProjectDetectionTaskStatus_Value.Summarizing:
      return STATUS_RANK_SUMMARIZING;
    case ProjectDetectionTaskStatus_Value.Succeeded:
    case ProjectDetectionTaskStatus_Value.Failed:
      return STATUS_RANK_TERMINAL;
    case ProjectDetectionTaskStatus_Value.Unknown:
    case ProjectDetectionTaskStatus_Value.UNRECOGNIZED:
    default:
      return STATUS_RANK_INITIAL;
  }
}

function projectDetectionSubStatusRank(status: ProjectDetectionSubTaskStatus_Value): number {
  switch (status) {
    case ProjectDetectionSubTaskStatus_Value.Pending:
      return STATUS_RANK_PENDING;
    case ProjectDetectionSubTaskStatus_Value.Succeeded:
    case ProjectDetectionSubTaskStatus_Value.Failed:
      return STATUS_RANK_TERMINAL;
    case ProjectDetectionSubTaskStatus_Value.Unknown:
    case ProjectDetectionSubTaskStatus_Value.UNRECOGNIZED:
    default:
      return STATUS_RANK_INITIAL;
  }
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null;
}

function isProjectDetectionTaskStatusChangedRecord(
  value: Record<string, unknown>,
): value is ProjectDetectionTaskStatusChangedRecord {
  return (
    typeof value.project_id === 'string' &&
    typeof value.image_uuid === 'string' &&
    typeof value.node_code === 'string' &&
    typeof value.main_task_uuid === 'string' &&
    isProjectDetectionTaskStatusValue(value.main_status) &&
    typeof value.sub_task_uuid === 'string' &&
    isProjectDetectionSubTaskStatusValue(value.sub_status) &&
    typeof value.occurred_at === 'string'
  );
}

function isProjectDetectionTaskStatusValue(value: unknown): value is ProjectDetectionTaskStatus_Value {
  return (
    value === ProjectDetectionTaskStatus_Value.Pending ||
    value === ProjectDetectionTaskStatus_Value.Classifying ||
    value === ProjectDetectionTaskStatus_Value.Detecting ||
    value === ProjectDetectionTaskStatus_Value.Summarizing ||
    value === ProjectDetectionTaskStatus_Value.Succeeded ||
    value === ProjectDetectionTaskStatus_Value.Failed
  );
}

function isProjectDetectionSubTaskStatusValue(value: unknown): value is ProjectDetectionSubTaskStatus_Value {
  return (
    value === ProjectDetectionSubTaskStatus_Value.Pending ||
    value === ProjectDetectionSubTaskStatus_Value.Succeeded ||
    value === ProjectDetectionSubTaskStatus_Value.Failed
  );
}
