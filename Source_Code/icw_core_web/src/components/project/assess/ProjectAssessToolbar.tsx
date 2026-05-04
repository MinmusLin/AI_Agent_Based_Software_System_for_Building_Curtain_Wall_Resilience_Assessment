import {
  CloseOutlined,
  CopyOutlined,
  DeleteOutlined,
  PlusOutlined,
  ReloadOutlined,
  StepForwardOutlined,
  SwapOutlined,
} from '@ant-design/icons';
import type { MenuProps } from 'antd';
import { Button, Dropdown } from 'antd';
import type { ReactElement } from 'react';

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
  groupCount: number;
  hasSelectedImages: boolean;
  onBatchDeleteImages: () => void;
  onBatchModeToggle: () => void;
  onBatchMoveImages: (targetGroupId: string) => void;
  onComplete: () => void;
  onCreateGroup: () => void;
  onRefresh: () => void;
  readOnly: boolean;
  totalImageCount: number;
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
  groupCount,
  hasSelectedImages,
  onBatchDeleteImages,
  onBatchModeToggle,
  onBatchMoveImages,
  onComplete,
  onCreateGroup,
  onRefresh,
  readOnly,
  totalImageCount,
}: ProjectAssessToolbarProps): ReactElement {
  const readonlyOrLoading = readOnly || assetsLoading;

  return (
    <div className="mb-4 flex items-center justify-between gap-4">
      <div>
        <div className="flex items-center gap-2">
          <h2 className="text-base font-semibold text-slate-900">图像资产构建</h2>
          {assetsLoading ? null : (
            <span className="flex items-center gap-2 rounded bg-slate-100 px-2 py-0.5 text-xs text-slate-500">
              <span>{groupCount} 个图像组</span>
              <span>{totalImageCount} 张图像</span>
            </span>
          )}
        </div>
        <p className="mt-1 text-sm text-slate-500">按建筑立面或区域组织幕墙图像，上传完成后执行 Agent 智能检测</p>
      </div>
      <div className="flex shrink-0 items-center gap-3">
        {batchMode && hasSelectedImages ? (
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
            >
              <Button disabled={assetsLoading || batchMoveDisabled} icon={<SwapOutlined />} loading={batchMoving}>
                批量移动
              </Button>
            </Dropdown>
          </>
        ) : null}
        {batchMode ? (
          <Button disabled={readonlyOrLoading} icon={<CloseOutlined />} onClick={onBatchModeToggle}>
            退出批量
          </Button>
        ) : (
          <Button disabled={readonlyOrLoading} icon={<CopyOutlined />} onClick={onBatchModeToggle}>
            批量操作
          </Button>
        )}
        <Button disabled={assetsLoading} icon={<ReloadOutlined />} onClick={onRefresh}>
          刷新
        </Button>
        <Button
          disabled={readonlyOrLoading}
          icon={<PlusOutlined />}
          loading={creatingGroup}
          onClick={onCreateGroup}
          type="primary"
        >
          新建图像组
        </Button>
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
