import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  PlayCircleOutlined,
  ReloadOutlined,
  RetweetOutlined,
  StepForwardOutlined,
} from '@ant-design/icons';
import { Button } from 'antd';
import type { ReactElement } from 'react';

import { EMPTY_ITEMS_COUNT } from '@/utils/assetsStage';

interface ProjectDetectionToolbarProps {
  actionsHidden: boolean;
  advancing: boolean;
  canComplete: boolean;
  canRetry: boolean;
  canStart: boolean;
  failedTaskCount: number;
  groupCount: number;
  loading: boolean;
  onCollapseAllGroups: () => void;
  onComplete: () => void;
  onExpandAllGroups: () => void;
  onRefresh: () => void;
  onRetry: () => void;
  onStart: () => void;
  readOnly: boolean;
  retrying: boolean;
  runningTaskCount: number;
  showComplete: boolean;
  showRetry: boolean;
  showStart: boolean;
  starting: boolean;
  totalImageCount: number;
}

interface ProjectDetectionStatusTagsProps {
  failedTaskCount: number;
  groupCount: number;
  loading: boolean;
  runningTaskCount: number;
  totalImageCount: number;
}

interface ProjectDetectionPrimaryActionProps {
  actionBusy: boolean;
  advancing: boolean;
  canComplete: boolean;
  canRetry: boolean;
  canStart: boolean;
  onComplete: () => void;
  onRetry: () => void;
  onStart: () => void;
  retrying: boolean;
  showComplete: boolean;
  showRetry: boolean;
  showStart: boolean;
  starting: boolean;
}

export function ProjectDetectionToolbar({
  actionsHidden,
  advancing,
  canComplete,
  canRetry,
  canStart,
  failedTaskCount,
  groupCount,
  loading,
  onCollapseAllGroups,
  onComplete,
  onExpandAllGroups,
  onRefresh,
  onRetry,
  onStart,
  readOnly,
  retrying,
  runningTaskCount,
  showComplete,
  showRetry,
  showStart,
  starting,
  totalImageCount,
}: ProjectDetectionToolbarProps): ReactElement {
  const actionBusy = loading || starting || retrying || advancing;

  return (
    <div className="mb-4 flex items-center justify-between gap-4">
      <div>
        <div className="flex items-center gap-2">
          <h2 className="text-base font-semibold text-slate-900">Agent 智能检测</h2>
          <ProjectDetectionStatusTags
            failedTaskCount={failedTaskCount}
            groupCount={groupCount}
            loading={loading}
            runningTaskCount={runningTaskCount}
            totalImageCount={totalImageCount}
          />
        </div>
        <p className="mt-1 text-sm text-slate-500">
          由 Agent 自主按需调度基于计算机视觉算法和深度学习模型的原子检测能力，汇总生成图像级 LLM 总结
        </p>
      </div>
      {actionsHidden ? null : (
        <div className="flex shrink-0 items-center gap-3">
          <Button icon={<MenuFoldOutlined />} onClick={onCollapseAllGroups}>
            全部收起
          </Button>
          <Button icon={<MenuUnfoldOutlined />} onClick={onExpandAllGroups}>
            全部展开
          </Button>
          {readOnly ? null : (
            <>
              <Button disabled={actionBusy} icon={<ReloadOutlined />} onClick={onRefresh}>
                刷新
              </Button>
              <ProjectDetectionPrimaryAction
                actionBusy={actionBusy}
                advancing={advancing}
                canComplete={canComplete}
                canRetry={canRetry}
                canStart={canStart}
                onComplete={onComplete}
                onRetry={onRetry}
                onStart={onStart}
                retrying={retrying}
                showComplete={showComplete}
                showRetry={showRetry}
                showStart={showStart}
                starting={starting}
              />
            </>
          )}
        </div>
      )}
    </div>
  );
}

function ProjectDetectionPrimaryAction({
  actionBusy,
  advancing,
  canComplete,
  canRetry,
  canStart,
  onComplete,
  onRetry,
  onStart,
  retrying,
  showComplete,
  showRetry,
  showStart,
  starting,
}: ProjectDetectionPrimaryActionProps): ReactElement | null {
  if (showRetry) {
    return (
      <Button danger disabled={!canRetry || actionBusy} icon={<RetweetOutlined />} loading={retrying} onClick={onRetry}>
        失败重试
      </Button>
    );
  }

  if (showStart) {
    return (
      <Button
        disabled={!canStart || actionBusy}
        icon={<PlayCircleOutlined />}
        loading={starting}
        onClick={onStart}
        type="primary"
      >
        启动 Agent 智能检测
      </Button>
    );
  }

  if (showComplete) {
    return (
      <Button
        disabled={!canComplete || actionBusy}
        icon={<StepForwardOutlined />}
        loading={advancing}
        onClick={onComplete}
        type="primary"
      >
        完成并进入下一步
      </Button>
    );
  }

  return null;
}

function ProjectDetectionStatusTags({
  failedTaskCount,
  groupCount,
  loading,
  runningTaskCount,
  totalImageCount,
}: ProjectDetectionStatusTagsProps): ReactElement | null {
  if (loading) {
    return null;
  }

  return (
    <>
      <span className="rounded bg-blue-50 px-2 py-0.5 text-xs text-[#1677FF]">
        <span>
          {groupCount} 组 {totalImageCount} 张图像
        </span>
      </span>
      {failedTaskCount > EMPTY_ITEMS_COUNT ? (
        <span className="rounded bg-red-50 px-2 py-0.5 text-xs text-red-500">{failedTaskCount} 张失败</span>
      ) : null}
      {runningTaskCount > EMPTY_ITEMS_COUNT ? (
        <span className="rounded bg-emerald-50 px-2 py-0.5 text-xs text-emerald-600">{runningTaskCount} 张检测中</span>
      ) : null}
    </>
  );
}
