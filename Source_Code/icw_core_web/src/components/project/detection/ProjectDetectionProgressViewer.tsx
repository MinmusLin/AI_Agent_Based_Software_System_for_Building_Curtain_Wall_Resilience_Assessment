import {
  CheckCircleFilled,
  ClockCircleOutlined,
  CloseCircleFilled,
  LoadingOutlined,
  MinusCircleOutlined,
} from '@ant-design/icons';
import { Modal, Progress, Spin, Tooltip } from 'antd';
import type { ReactElement } from 'react';

import {
  DETECTION_PROGRESS_NODE_STATUS_FAILED,
  DETECTION_PROGRESS_NODE_STATUS_PENDING,
  DETECTION_PROGRESS_NODE_STATUS_RUNNING,
  DETECTION_PROGRESS_NODE_STATUS_SKIPPED,
  DETECTION_PROGRESS_NODE_STATUS_SUCCEEDED,
  type DetectionProgressNodeStatus,
  reasoningDetectionNodeCode,
} from '@/constants/common';
import type { ProjectDetectionNodeStatus, ProjectDetectionStatus } from '@/gen/core/common';
import {
  DetectionTaskCode_Value,
  ProjectDetectionSubTaskStatus_Value,
  ProjectDetectionTaskStatus_Value,
} from '@/gen/core/common';
import { detectionTaskProgressPercent } from '@/utils/detectionStage';

interface ProjectDetectionProgressViewerProps {
  onClose: () => void;
  open: boolean;
  task: ProjectDetectionStatus | undefined;
}

interface ReasoningNodeMeta {
  label: string;
  nodeCode: string;
}

const EMPTY_NODE_COUNT = 0;
const EMPTY_TASK_UUID = '';

const REASONING_NODES: ReasoningNodeMeta[] = [
  { label: '金属锈蚀', nodeCode: reasoningDetectionNodeCode(DetectionTaskCode_Value.Corrosion) },
  { label: '石材裂缝', nodeCode: reasoningDetectionNodeCode(DetectionTaskCode_Value.Crack) },
  { label: '石材污渍', nodeCode: reasoningDetectionNodeCode(DetectionTaskCode_Value.Stain) },
  { label: '玻璃平整度', nodeCode: reasoningDetectionNodeCode(DetectionTaskCode_Value.Flatness) },
  { label: '玻璃爆裂', nodeCode: reasoningDetectionNodeCode(DetectionTaskCode_Value.Spalling) },
];

function statusText(status: DetectionProgressNodeStatus): string {
  switch (status) {
    case DETECTION_PROGRESS_NODE_STATUS_FAILED:
      return '失败';
    case DETECTION_PROGRESS_NODE_STATUS_RUNNING:
      return '执行中';
    case DETECTION_PROGRESS_NODE_STATUS_SKIPPED:
      return '不执行';
    case DETECTION_PROGRESS_NODE_STATUS_SUCCEEDED:
      return '完成';
    case DETECTION_PROGRESS_NODE_STATUS_PENDING:
    default:
      return '等待中';
  }
}

function statusIcon(status: DetectionProgressNodeStatus): ReactElement {
  switch (status) {
    case DETECTION_PROGRESS_NODE_STATUS_FAILED:
      return <CloseCircleFilled className="text-red-500" />;
    case DETECTION_PROGRESS_NODE_STATUS_RUNNING:
      return <Spin indicator={<LoadingOutlined className="text-[#1677FF]" spin />} />;
    case DETECTION_PROGRESS_NODE_STATUS_SKIPPED:
      return <MinusCircleOutlined className="text-slate-400" />;
    case DETECTION_PROGRESS_NODE_STATUS_SUCCEEDED:
      return <CheckCircleFilled className="text-emerald-500" />;
    case DETECTION_PROGRESS_NODE_STATUS_PENDING:
    default:
      return <ClockCircleOutlined className="text-slate-400" />;
  }
}

function statusClassName(status: DetectionProgressNodeStatus): string {
  switch (status) {
    case DETECTION_PROGRESS_NODE_STATUS_FAILED:
      return 'border-red-200 bg-red-50 text-red-700';
    case DETECTION_PROGRESS_NODE_STATUS_RUNNING:
      return 'border-blue-200 bg-blue-50 text-[#1677FF]';
    case DETECTION_PROGRESS_NODE_STATUS_SKIPPED:
      return 'border-dashed border-slate-200 bg-slate-50 text-slate-400';
    case DETECTION_PROGRESS_NODE_STATUS_SUCCEEDED:
      return 'border-emerald-200 bg-emerald-50 text-emerald-700';
    case DETECTION_PROGRESS_NODE_STATUS_PENDING:
    default:
      return 'border-slate-200 bg-slate-50 text-slate-500';
  }
}

function nodeStatus(node: ProjectDetectionNodeStatus | undefined): DetectionProgressNodeStatus {
  if (!node) {
    return DETECTION_PROGRESS_NODE_STATUS_SKIPPED;
  }
  switch (node.sub_status) {
    case ProjectDetectionSubTaskStatus_Value.Failed:
      return DETECTION_PROGRESS_NODE_STATUS_FAILED;
    case ProjectDetectionSubTaskStatus_Value.Pending:
      return DETECTION_PROGRESS_NODE_STATUS_RUNNING;
    case ProjectDetectionSubTaskStatus_Value.Succeeded:
      return DETECTION_PROGRESS_NODE_STATUS_SUCCEEDED;
    case ProjectDetectionSubTaskStatus_Value.Unknown:
    case ProjectDetectionSubTaskStatus_Value.UNRECOGNIZED:
    default:
      return DETECTION_PROGRESS_NODE_STATUS_PENDING;
  }
}

function classificationStatus(task: ProjectDetectionStatus): DetectionProgressNodeStatus {
  return task.classification_status ? nodeStatus(task.classification_status) : DETECTION_PROGRESS_NODE_STATUS_PENDING;
}

function reasoningNodeStatus(
  node: ProjectDetectionNodeStatus | undefined,
  currentClassificationStatus: DetectionProgressNodeStatus,
): DetectionProgressNodeStatus {
  if (currentClassificationStatus === DETECTION_PROGRESS_NODE_STATUS_FAILED) {
    return DETECTION_PROGRESS_NODE_STATUS_SKIPPED;
  }
  if (node) {
    return nodeStatus(node);
  }
  if (
    currentClassificationStatus === DETECTION_PROGRESS_NODE_STATUS_SUCCEEDED ||
    currentClassificationStatus === DETECTION_PROGRESS_NODE_STATUS_SKIPPED
  ) {
    return DETECTION_PROGRESS_NODE_STATUS_SKIPPED;
  }
  return DETECTION_PROGRESS_NODE_STATUS_PENDING;
}

function reasoningAggregateStatus(
  task: ProjectDetectionStatus,
  currentClassificationStatus: DetectionProgressNodeStatus,
): DetectionProgressNodeStatus {
  if (currentClassificationStatus === DETECTION_PROGRESS_NODE_STATUS_FAILED) {
    return DETECTION_PROGRESS_NODE_STATUS_SKIPPED;
  }

  const nodes = task.detection_status;
  if (nodes.length === EMPTY_NODE_COUNT) {
    if (
      currentClassificationStatus === DETECTION_PROGRESS_NODE_STATUS_SUCCEEDED ||
      task.main_status === ProjectDetectionTaskStatus_Value.Succeeded
    ) {
      return DETECTION_PROGRESS_NODE_STATUS_SKIPPED;
    }
    return DETECTION_PROGRESS_NODE_STATUS_PENDING;
  }

  const statuses = nodes.map((node) => nodeStatus(node));
  if (statuses.includes(DETECTION_PROGRESS_NODE_STATUS_FAILED)) {
    return DETECTION_PROGRESS_NODE_STATUS_FAILED;
  }
  if (statuses.includes(DETECTION_PROGRESS_NODE_STATUS_RUNNING)) {
    return DETECTION_PROGRESS_NODE_STATUS_RUNNING;
  }
  if (statuses.every((status) => status === DETECTION_PROGRESS_NODE_STATUS_SUCCEEDED)) {
    return DETECTION_PROGRESS_NODE_STATUS_SUCCEEDED;
  }
  return DETECTION_PROGRESS_NODE_STATUS_PENDING;
}

function endStatus(task: ProjectDetectionStatus): DetectionProgressNodeStatus {
  if (task.main_status === ProjectDetectionTaskStatus_Value.Failed) {
    return DETECTION_PROGRESS_NODE_STATUS_FAILED;
  }
  if (task.main_status === ProjectDetectionTaskStatus_Value.Succeeded) {
    return DETECTION_PROGRESS_NODE_STATUS_SUCCEEDED;
  }
  return DETECTION_PROGRESS_NODE_STATUS_PENDING;
}

function summaryStatus(
  task: ProjectDetectionStatus,
  currentClassificationStatus: DetectionProgressNodeStatus,
  currentReasoningStatus: DetectionProgressNodeStatus,
): DetectionProgressNodeStatus {
  if (
    currentClassificationStatus === DETECTION_PROGRESS_NODE_STATUS_FAILED ||
    currentReasoningStatus === DETECTION_PROGRESS_NODE_STATUS_FAILED ||
    currentReasoningStatus === DETECTION_PROGRESS_NODE_STATUS_SKIPPED
  ) {
    return DETECTION_PROGRESS_NODE_STATUS_SKIPPED;
  }
  if (task.summary_status) {
    return nodeStatus(task.summary_status);
  }
  if (task.main_status === ProjectDetectionTaskStatus_Value.Succeeded) {
    return DETECTION_PROGRESS_NODE_STATUS_SKIPPED;
  }
  return DETECTION_PROGRESS_NODE_STATUS_PENDING;
}

function taskUuidTooltip(taskUuid: string | undefined): string {
  return taskUuid && taskUuid !== EMPTY_TASK_UUID ? `任务 ID：${taskUuid}` : '暂无任务 ID';
}

function FlowNode({
  label,
  status,
  taskUuid,
}: {
  label: string;
  status: DetectionProgressNodeStatus;
  taskUuid?: string;
}): ReactElement {
  const node = (
    <div
      className={`flex size-24 flex-col items-center justify-center rounded-lg border px-2 py-3 text-center ${statusClassName(status)}`}
    >
      <div className="mb-2 flex h-5 items-center justify-center">{statusIcon(status)}</div>
      <div className="text-sm font-medium">{label}</div>
      <div className="mt-1 text-xs">{statusText(status)}</div>
    </div>
  );
  if (taskUuid === undefined) {
    return node;
  }
  return <Tooltip title={taskUuidTooltip(taskUuid)}>{node}</Tooltip>;
}

function Connector(): ReactElement {
  return <div className="h-px min-w-4 flex-1 bg-slate-200" />;
}

export function ProjectDetectionProgressViewer({
  onClose,
  open,
  task,
}: ProjectDetectionProgressViewerProps): ReactElement {
  const reasoningStatusMap = new Map((task?.detection_status ?? []).map((node) => [node.node_code, node]));
  const progressPercent = detectionTaskProgressPercent(task);
  const currentClassificationStatus = task ? classificationStatus(task) : DETECTION_PROGRESS_NODE_STATUS_PENDING;
  const currentReasoningStatus = task
    ? reasoningAggregateStatus(task, currentClassificationStatus)
    : DETECTION_PROGRESS_NODE_STATUS_PENDING;
  const currentSummaryStatus = task
    ? summaryStatus(task, currentClassificationStatus, currentReasoningStatus)
    : DETECTION_PROGRESS_NODE_STATUS_PENDING;

  return (
    <Modal centered footer={null} onCancel={onClose} open={open} title="检测进度" width={980}>
      {task ? (
        <div className="space-y-5">
          <div className="mb-8">
            <Progress
              percent={progressPercent}
              status={task.main_status === ProjectDetectionTaskStatus_Value.Failed ? 'exception' : 'active'}
            />
          </div>
          <div className="flex items-center gap-3">
            <FlowNode label="开始" status={DETECTION_PROGRESS_NODE_STATUS_SUCCEEDED} />
            <Connector />
            <FlowNode
              label="材质分类"
              status={currentClassificationStatus}
              taskUuid={task.classification_status?.sub_task_uuid ?? task.main_task_uuid}
            />
            <Connector />
            <FlowNode label="原子检测" status={currentReasoningStatus} />
            <Connector />
            <FlowNode label="结果总结" status={currentSummaryStatus} taskUuid={task.summary_status?.sub_task_uuid} />
            <Connector />
            <FlowNode label="结束" status={endStatus(task)} />
          </div>
          <div className="flex justify-center">
            <div className="w-[660px] rounded-lg border border-slate-200 bg-slate-50 p-3">
              <div className="mb-3 text-center text-sm font-medium text-slate-900">原子能力检测器</div>
              <div className="grid grid-cols-5 justify-items-center gap-3">
                {REASONING_NODES.map((node) => (
                  <FlowNode
                    key={node.nodeCode}
                    label={node.label}
                    status={reasoningNodeStatus(reasoningStatusMap.get(node.nodeCode), currentClassificationStatus)}
                    taskUuid={reasoningStatusMap.get(node.nodeCode)?.sub_task_uuid}
                  />
                ))}
              </div>
            </div>
          </div>
        </div>
      ) : null}
    </Modal>
  );
}
