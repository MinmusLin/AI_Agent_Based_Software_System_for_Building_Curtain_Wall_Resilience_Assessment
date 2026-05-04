import { Empty, Spin } from 'antd';
import type { DragEvent, ReactElement } from 'react';

import { ProjectAssessGroup } from '@/components/project/assess/ProjectAssessGroup';
import type { ProjectImageStatus } from '@/types/common';
import type { ProjectGroup } from '@/types/project/assets';
import { EMPTY_ITEMS_COUNT } from '@/utils/assetsStage';

interface ProjectAssessContentProps {
  accept: string;
  assetsLoading: boolean;
  batchMode: boolean;
  collapsedGroupIds: Set<string>;
  deletingGroupIds: Set<string>;
  deletingImageUuids: Set<string>;
  draggingGroupId: string;
  editingGroupId: string;
  editingGroupName: string;
  groups: ProjectGroup[];
  onDeleteGroup: (groupId: string) => void;
  onDeleteImage: (imageUuid: string) => void;
  onDeselectGroupImages: (imageUuids: string[]) => void;
  onEditingGroupNameChange: (name: string) => void;
  onGroupDragEnd: () => void;
  onGroupDragOver: (event: DragEvent<HTMLElement>, groupId: string) => void;
  onGroupDragStart: (event: DragEvent<HTMLDivElement>, group: ProjectGroup) => void;
  onGroupDrop: (event: DragEvent<HTMLElement>, groupId: string) => void;
  onImageDragEnd: () => void;
  onImageDragStart: (
    event: DragEvent<HTMLDivElement>,
    imageUuid: string,
    sourceGroupId: string,
    imageStatus: ProjectImageStatus,
  ) => void;
  onOpenImageViewer: (imageUuid: string) => void;
  onSaveEditGroup: () => void;
  onSelectGroupImages: (imageUuids: string[]) => void;
  onStartEditGroup: (group: ProjectGroup) => void;
  onToggleCollapsed: (groupId: string) => void;
  onToggleSelectedImage: (imageUuid: string, checked?: boolean) => void;
  onUploadFiles: (groupId: string, files: File[]) => void;
  readOnly: boolean;
  registerGroupRef: (groupId: string, node: HTMLElement | null) => void;
  savingGroupId: string;
  selectedImageUuids: Set<string>;
  uploadingImages: boolean;
}

export function ProjectAssessContent({
  accept,
  assetsLoading,
  batchMode,
  collapsedGroupIds,
  deletingGroupIds,
  deletingImageUuids,
  draggingGroupId,
  editingGroupId,
  editingGroupName,
  groups,
  onDeleteGroup,
  onDeleteImage,
  onDeselectGroupImages,
  onEditingGroupNameChange,
  onGroupDragEnd,
  onGroupDragOver,
  onGroupDragStart,
  onGroupDrop,
  onImageDragEnd,
  onImageDragStart,
  onOpenImageViewer,
  onSaveEditGroup,
  onSelectGroupImages,
  onStartEditGroup,
  onToggleCollapsed,
  onToggleSelectedImage,
  onUploadFiles,
  readOnly,
  registerGroupRef,
  savingGroupId,
  selectedImageUuids,
  uploadingImages,
}: ProjectAssessContentProps): ReactElement {
  if (assetsLoading) {
    return (
      <div className="flex h-full items-center justify-center">
        <Spin description="正在加载图像资产" />
      </div>
    );
  }

  if (groups.length === EMPTY_ITEMS_COUNT) {
    return (
      <div className="flex h-full items-center justify-center rounded-lg border border-dashed border-slate-200">
        <Empty description="暂无图像组，请新建图像组后上传图像" />
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {groups.map((group) => (
        <ProjectAssessGroup
          accept={accept}
          batchMode={batchMode}
          collapsed={collapsedGroupIds.has(group.id)}
          deletingGroup={deletingGroupIds.has(group.id)}
          deletingImageUuids={deletingImageUuids}
          dragging={draggingGroupId === group.id}
          editing={editingGroupId === group.id}
          editingGroupName={editingGroupName}
          group={group}
          groupCount={groups.length}
          key={group.id}
          onDeleteGroup={onDeleteGroup}
          onDeleteImage={onDeleteImage}
          onDeselectGroupImages={onDeselectGroupImages}
          onEditingGroupNameChange={onEditingGroupNameChange}
          onGroupDragEnd={onGroupDragEnd}
          onGroupDragOver={onGroupDragOver}
          onGroupDragStart={onGroupDragStart}
          onGroupDrop={onGroupDrop}
          onImageDragEnd={onImageDragEnd}
          onImageDragStart={onImageDragStart}
          onOpenImageViewer={onOpenImageViewer}
          onSaveEditGroup={onSaveEditGroup}
          onSelectGroupImages={onSelectGroupImages}
          onStartEditGroup={onStartEditGroup}
          onToggleCollapsed={onToggleCollapsed}
          onToggleSelectedImage={onToggleSelectedImage}
          onUploadFiles={onUploadFiles}
          readOnly={readOnly}
          registerGroupRef={registerGroupRef}
          savingGroup={savingGroupId === group.id}
          selectedImageUuids={selectedImageUuids}
          uploadingImages={uploadingImages}
        />
      ))}
    </div>
  );
}
