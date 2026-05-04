import { CheckSquareOutlined, CloseSquareOutlined, DeleteOutlined } from '@ant-design/icons';
import { Button } from 'antd';
import type { ReactElement } from 'react';

import { EMPTY_ITEMS_COUNT } from '@/utils/assetsStage';

interface ProjectAssessGroupHeaderActionsProps {
  allSelectableImagesSelected: boolean;
  batchMode: boolean;
  canDeleteGroup: boolean;
  deletingGroup: boolean;
  groupId: string;
  groupImageCount: number;
  hasSelectableImages: boolean;
  hasSelectedGroupImages: boolean;
  onDeleteGroup: (groupId: string) => void;
  onDeselectGroupImages: (imageUuids: string[]) => void;
  onSelectGroupImages: (imageUuids: string[]) => void;
  selectableImageUuids: string[];
}

export function ProjectAssessGroupHeaderActions({
  allSelectableImagesSelected,
  batchMode,
  canDeleteGroup,
  deletingGroup,
  groupId,
  groupImageCount,
  hasSelectableImages,
  hasSelectedGroupImages,
  onDeleteGroup,
  onDeselectGroupImages,
  onSelectGroupImages,
  selectableImageUuids,
}: ProjectAssessGroupHeaderActionsProps): ReactElement {
  return (
    <div className="ml-auto flex items-center gap-2 opacity-0 transition-opacity duration-200 group-hover:opacity-100">
      {batchMode && groupImageCount > EMPTY_ITEMS_COUNT ? (
        <>
          <Button
            disabled={!hasSelectableImages || allSelectableImagesSelected}
            icon={<CheckSquareOutlined />}
            onClick={() => {
              onSelectGroupImages(selectableImageUuids);
            }}
            size="small"
            type="text"
          >
            全选本组
          </Button>
          <Button
            disabled={!hasSelectedGroupImages}
            icon={<CloseSquareOutlined />}
            onClick={() => {
              onDeselectGroupImages(selectableImageUuids);
            }}
            size="small"
            type="text"
          >
            取消全选本组
          </Button>
        </>
      ) : null}
      {canDeleteGroup ? (
        <Button
          aria-label="删除图像组"
          danger
          icon={<DeleteOutlined />}
          loading={deletingGroup}
          onClick={() => {
            onDeleteGroup(groupId);
          }}
          shape="circle"
          size="small"
          type="text"
        />
      ) : null}
    </div>
  );
}
