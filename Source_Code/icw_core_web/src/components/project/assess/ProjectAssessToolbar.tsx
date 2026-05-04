import {
  CloseOutlined,
  CopyOutlined,
  DeleteOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  PlusOutlined,
  ReloadOutlined,
  StepForwardOutlined,
  SwapOutlined,
} from '@ant-design/icons';
import type { MenuProps } from 'antd';
import { Button, Dropdown } from 'antd';
import type { ReactElement } from 'react';

import { EMPTY_ITEMS_COUNT } from '@/utils/assetsStage';

interface ProjectAssessToolbarProps {
  advancing: boolean;
  assetsLoading: boolean;
  batchDeleting: boolean;
  batchMode: boolean;
  batchMoveDisabled: boolean;
  batchMoveMenuItems: NonNullable<MenuProps['items']>;
  batchMoving: boolean;
  canComplete: boolean;
  creatingGroup: boolean;
  failedImageCount: number;
  groupCount: number;
  hasSelectedImages: boolean;
  onBatchDeleteImages: () => void;
  onBatchModeToggle: () => void;
  onBatchMoveImages: (targetGroupId: string) => void;
  onCollapseAllGroups: () => void;
  onComplete: () => void;
  onCreateGroup: () => void;
  onExpandAllGroups: () => void;
  onRefresh: () => void;
  pendingImageCount: number;
  readOnly: boolean;
  totalImageCount: number;
}

interface ProjectAssessStatusTagsProps {
  assetsLoading: boolean;
  failedImageCount: number;
  groupCount: number;
  pendingImageCount: number;
  totalImageCount: number;
}

interface ProjectAssessBatchActionsProps {
  assetsLoading: boolean;
  batchDeleting: boolean;
  batchMode: boolean;
  batchMoveDisabled: boolean;
  batchMoveMenuItems: NonNullable<MenuProps['items']>;
  batchMoving: boolean;
  hasSelectedImages: boolean;
  onBatchDeleteImages: () => void;
  onBatchMoveImages: (targetGroupId: string) => void;
  onCollapseAllGroups: () => void;
  onExpandAllGroups: () => void;
  readOnly: boolean;
}

export function ProjectAssessToolbar({
  advancing,
  assetsLoading,
  batchDeleting,
  batchMode,
  batchMoveDisabled,
  batchMoveMenuItems,
  batchMoving,
  canComplete,
  creatingGroup,
  failedImageCount,
  groupCount,
  hasSelectedImages,
  onBatchDeleteImages,
  onBatchModeToggle,
  onBatchMoveImages,
  onCollapseAllGroups,
  onComplete,
  onCreateGroup,
  onExpandAllGroups,
  onRefresh,
  pendingImageCount,
  readOnly,
  totalImageCount,
}: ProjectAssessToolbarProps): ReactElement {
  return (
    <div className="mb-4 flex items-center justify-between gap-4">
      <div>
        <div className="flex items-center gap-2">
          <h2 className="text-base font-semibold text-slate-900">图像资产构建</h2>
          <ProjectAssessStatusTags
            assetsLoading={assetsLoading}
            failedImageCount={failedImageCount}
            groupCount={groupCount}
            pendingImageCount={pendingImageCount}
            totalImageCount={totalImageCount}
          />
        </div>
        <p className="mt-1 text-sm text-slate-500">按建筑立面或区域组织幕墙图像，上传完成后执行 Agent 智能检测</p>
      </div>
      <div className="flex shrink-0 items-center gap-3">
        <ProjectAssessBatchActions
          assetsLoading={assetsLoading}
          batchDeleting={batchDeleting}
          batchMode={batchMode}
          batchMoveDisabled={batchMoveDisabled}
          batchMoveMenuItems={batchMoveMenuItems}
          batchMoving={batchMoving}
          hasSelectedImages={hasSelectedImages}
          onBatchDeleteImages={onBatchDeleteImages}
          onBatchMoveImages={onBatchMoveImages}
          onCollapseAllGroups={onCollapseAllGroups}
          onExpandAllGroups={onExpandAllGroups}
          readOnly={readOnly}
        />
        {!readOnly && batchMode ? (
          <Button disabled={assetsLoading} icon={<CloseOutlined />} onClick={onBatchModeToggle}>
            退出批量
          </Button>
        ) : null}
        {!readOnly && !batchMode ? (
          <Button disabled={assetsLoading} icon={<CopyOutlined />} onClick={onBatchModeToggle}>
            批量操作
          </Button>
        ) : null}
        {!readOnly ? (
          <Button disabled={assetsLoading} icon={<ReloadOutlined />} onClick={onRefresh}>
            刷新
          </Button>
        ) : null}
        {!readOnly ? (
          <Button
            disabled={assetsLoading}
            icon={<PlusOutlined />}
            loading={creatingGroup}
            onClick={onCreateGroup}
            type="primary"
          >
            新建图像组
          </Button>
        ) : null}
        {canComplete ? (
          <Button
            disabled={assetsLoading}
            icon={<StepForwardOutlined />}
            loading={advancing}
            onClick={onComplete}
            type="primary"
          >
            完成并进入下一步
          </Button>
        ) : null}
      </div>
    </div>
  );
}

function ProjectAssessStatusTags({
  assetsLoading,
  failedImageCount,
  groupCount,
  pendingImageCount,
  totalImageCount,
}: ProjectAssessStatusTagsProps): ReactElement | null {
  if (assetsLoading) {
    return null;
  }

  return (
    <>
      <span className="rounded bg-blue-50 px-2 py-0.5 text-xs text-[#1677FF]">
        <span>
          {groupCount} 组 {totalImageCount} 张图像
        </span>
      </span>
      {failedImageCount > EMPTY_ITEMS_COUNT ? (
        <span className="rounded bg-red-50 px-2 py-0.5 text-xs text-red-500">{failedImageCount} 张失败</span>
      ) : null}
      {pendingImageCount > EMPTY_ITEMS_COUNT ? (
        <span className="rounded bg-emerald-50 px-2 py-0.5 text-xs text-emerald-600">{pendingImageCount} 张上传中</span>
      ) : null}
    </>
  );
}

function ProjectAssessBatchActions({
  assetsLoading,
  batchDeleting,
  batchMode,
  batchMoveDisabled,
  batchMoveMenuItems,
  batchMoving,
  hasSelectedImages,
  onBatchDeleteImages,
  onBatchMoveImages,
  onCollapseAllGroups,
  onExpandAllGroups,
  readOnly,
}: ProjectAssessBatchActionsProps): ReactElement | null {
  if (readOnly || !batchMode || !hasSelectedImages) {
    return (
      <>
        <Button disabled={assetsLoading} icon={<MenuFoldOutlined />} onClick={onCollapseAllGroups}>
          全部收起
        </Button>
        <Button disabled={assetsLoading} icon={<MenuUnfoldOutlined />} onClick={onExpandAllGroups}>
          全部展开
        </Button>
      </>
    );
  }
  return (
    <>
      <Button
        danger
        disabled={assetsLoading}
        icon={<DeleteOutlined />}
        loading={batchDeleting}
        onClick={onBatchDeleteImages}
      >
        批量删除
      </Button>
      <Dropdown
        disabled={assetsLoading || batchMoveDisabled}
        menu={{
          items: batchMoveMenuItems,
          onClick: ({ key }) => {
            onBatchMoveImages(key);
          },
        }}
        popupRender={(menus) => (
          <div className="max-h-64 overflow-y-auto rounded-md border border-slate-200 bg-white shadow-sm">{menus}</div>
        )}
      >
        <Button disabled={assetsLoading || batchMoveDisabled} icon={<SwapOutlined />} loading={batchMoving}>
          批量移动
        </Button>
      </Dropdown>
    </>
  );
}
