import {
  CaretDownOutlined,
  CloseCircleFilled,
  DeleteOutlined,
  DragOutlined,
  EyeOutlined,
  FileImageOutlined,
  LoadingOutlined,
  UploadOutlined,
} from '@ant-design/icons';
import { Button, Checkbox, Input, Tooltip } from 'antd';
import type { ChangeEvent, DragEvent, ReactElement } from 'react';
import { useRef } from 'react';

import { ProjectAssessGroupHeaderActions } from '@/components/project/assess/ProjectAssessGroupHeaderActions';
import type { ProjectGroup } from '@/gen/core/api/common';
import { type ProjectImage, ProjectImageStatus_Value } from '@/gen/core/common';
import { canMoveProjectImage, EMPTY_ITEMS_COUNT, NEXT_INDEX_OFFSET, projectImageStats } from '@/utils/assetsStage';

interface ProjectAssetImageCardProps {
  batchMode: boolean;
  deleting: boolean;
  image: ProjectImage;
  onDeleteImage: (imageUuid: string) => void;
  onDragEnd: () => void;
  onDragStart: (event: DragEvent<HTMLDivElement>, imageUuid: string, imageStatus: ProjectImageStatus_Value) => void;
  onOpenImageViewer: (imageUuid: string) => void;
  onToggleSelectedImage: (imageUuid: string, checked?: boolean) => void;
  readOnly: boolean;
  selected: boolean;
}

interface ProjectImageUploadTileProps {
  accept: string;
  onUploadFiles: (files: File[]) => void;
  uploading: boolean;
}

interface ProjectImageActionOverlayProps {
  deleting: boolean;
  imageUuid: string;
  imageUploaded: boolean;
  onDeleteImage: (imageUuid: string) => void;
  onOpenImageViewer: (imageUuid: string) => void;
  readOnly: boolean;
}

interface ProjectImageCardState {
  actionVisible: boolean;
  batchSelectable: boolean;
  movable: boolean;
  ready: boolean;
  uploaded: boolean;
}

interface ProjectGroupDragHandleProps {
  readOnly: boolean;
}

interface ProjectAssessGroupProps {
  accept: string;
  batchMode: boolean;
  collapsed: boolean;
  deletingGroup: boolean;
  deletingImageUuids: Set<string>;
  dragging: boolean;
  editing: boolean;
  editingGroupName: string;
  group: ProjectGroup;
  groupCount: number;
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
    imageStatus: ProjectImageStatus_Value,
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
  savingGroup: boolean;
  selectedImageUuids: Set<string>;
  uploadingImages: boolean;
}

function ImageUnavailableIcon(): ReactElement {
  return (
    <span className="relative inline-flex size-8 items-center justify-center text-slate-400">
      <FileImageOutlined className="text-3xl" />
      <CloseCircleFilled className="absolute -right-1 -bottom-1 rounded-full bg-white text-sm text-red-500" />
    </span>
  );
}

function imageStatusNode(status: ProjectImageStatus_Value): ReactElement {
  switch (status) {
    case ProjectImageStatus_Value.Pending:
      return <LoadingOutlined className="text-2xl text-slate-400" />;
    case ProjectImageStatus_Value.Failed:
      return (
        <Tooltip title="图像已失效">
          <ImageUnavailableIcon />
        </Tooltip>
      );
    case ProjectImageStatus_Value.Uploaded:
      return (
        <Tooltip title="图像已失效">
          <ImageUnavailableIcon />
        </Tooltip>
      );
    case ProjectImageStatus_Value.Unknown:
    case ProjectImageStatus_Value.UNRECOGNIZED:
    default:
      return (
        <Tooltip title="图像状态异常">
          <ImageUnavailableIcon />
        </Tooltip>
      );
  }
}

function projectImageCardState(image: ProjectImage, batchMode: boolean, readOnly: boolean): ProjectImageCardState {
  const uploaded = image.status === ProjectImageStatus_Value.Uploaded;
  const failed = image.status === ProjectImageStatus_Value.Failed;

  return {
    actionVisible: !batchMode && (uploaded || (!readOnly && failed)),
    batchSelectable: batchMode && image.status !== ProjectImageStatus_Value.Pending,
    movable: !readOnly && !batchMode && canMoveProjectImage(image.status),
    ready: uploaded && image.thumbnail_url !== '',
    uploaded,
  };
}

function ProjectAssetImageCard({
  batchMode,
  deleting,
  image,
  onDeleteImage,
  onDragEnd,
  onDragStart,
  onOpenImageViewer,
  onToggleSelectedImage,
  readOnly,
  selected,
}: ProjectAssetImageCardProps): ReactElement {
  const imageState = projectImageCardState(image, batchMode, readOnly);

  return (
    <div
      className={`group/image relative aspect-square overflow-hidden rounded-lg border bg-white ${
        selected ? 'border-[#1677FF]' : 'border-slate-200'
      } ${batchMode ? 'cursor-pointer' : ''}`}
      draggable={imageState.movable}
      onClick={() => {
        if (imageState.batchSelectable) {
          onToggleSelectedImage(image.uuid);
        }
      }}
      onDragEnd={onDragEnd}
      onDragStart={(event: DragEvent<HTMLDivElement>) => {
        if (!imageState.movable) {
          event.preventDefault();
          event.stopPropagation();
          return;
        }
        onDragStart(event, image.uuid, image.status);
      }}
    >
      {imageState.batchSelectable ? (
        <Checkbox
          checked={selected}
          className="absolute left-2 top-2 z-20"
          onChange={(event) => {
            onToggleSelectedImage(image.uuid, event.target.checked);
          }}
          onClick={(event) => {
            event.stopPropagation();
          }}
        />
      ) : null}
      {imageState.ready ? (
        <img alt={image.file_name} className="size-full object-cover" draggable={false} src={image.thumbnail_url} />
      ) : (
        <div className="flex size-full flex-col items-center justify-center gap-2 px-3 text-center text-xs text-slate-500">
          {imageStatusNode(image.status)}
        </div>
      )}
      {imageState.actionVisible ? (
        <ProjectImageActionOverlay
          deleting={deleting}
          imageUploaded={imageState.uploaded}
          imageUuid={image.uuid}
          onDeleteImage={onDeleteImage}
          onOpenImageViewer={onOpenImageViewer}
          readOnly={readOnly}
        />
      ) : null}
    </div>
  );
}

function ProjectImageActionOverlay({
  deleting,
  imageUploaded,
  imageUuid,
  onDeleteImage,
  onOpenImageViewer,
  readOnly,
}: ProjectImageActionOverlayProps): ReactElement {
  return (
    <div className="absolute inset-0 flex items-center justify-center gap-2 bg-slate-950/0 opacity-0 transition duration-200 group-hover/image:bg-slate-950/35 group-hover/image:opacity-100">
      {imageUploaded ? (
        <Tooltip title="图像详情">
          <Button
            aria-label="图像详情"
            icon={<EyeOutlined />}
            onClick={() => {
              onOpenImageViewer(imageUuid);
            }}
            shape="circle"
            size="small"
          />
        </Tooltip>
      ) : null}
      {!readOnly ? (
        <Tooltip title="删除图像">
          <Button
            aria-label="删除图像"
            danger
            icon={<DeleteOutlined />}
            loading={deleting}
            onClick={() => {
              onDeleteImage(imageUuid);
            }}
            shape="circle"
            size="small"
          />
        </Tooltip>
      ) : null}
    </div>
  );
}

function ProjectImageUploadTile({ accept, onUploadFiles, uploading }: ProjectImageUploadTileProps): ReactElement {
  const inputRef = useRef<HTMLInputElement | null>(null);

  return (
    <button
      className={`flex aspect-square flex-col items-center justify-center gap-2 rounded-lg border border-dashed bg-white text-sm transition duration-200 ${
        uploading
          ? 'cursor-not-allowed border-slate-200 text-slate-400'
          : 'border-slate-300 text-slate-500 hover:border-[#1677FF] hover:text-[#1677FF]'
      }`}
      disabled={uploading}
      onClick={() => {
        if (uploading) {
          return;
        }
        inputRef.current?.click();
      }}
      type="button"
    >
      {uploading ? <LoadingOutlined className="text-2xl" /> : <UploadOutlined className="text-2xl" />}
      <input
        accept={accept}
        className="hidden"
        disabled={uploading}
        multiple
        onChange={(event: ChangeEvent<HTMLInputElement>) => {
          if (uploading) {
            event.target.value = '';
            return;
          }
          const selectedFiles = Array.from(event.target.files ?? []);
          event.target.value = '';
          onUploadFiles(selectedFiles);
        }}
        ref={inputRef}
        type="file"
      />
    </button>
  );
}

function ProjectGroupDragHandle({ readOnly }: ProjectGroupDragHandleProps): ReactElement | null {
  if (readOnly) {
    return null;
  }
  return <DragOutlined className="cursor-grab text-slate-400" />;
}

export function ProjectAssessGroup({
  accept,
  batchMode,
  collapsed,
  deletingGroup,
  deletingImageUuids,
  dragging,
  editing,
  editingGroupName,
  group,
  groupCount,
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
  savingGroup,
  selectedImageUuids,
  uploadingImages,
}: ProjectAssessGroupProps): ReactElement {
  const canDeleteGroup = !readOnly && groupCount > NEXT_INDEX_OFFSET;
  const imageStats = projectImageStats(group.images);
  const selectableImageUuids = group.images
    .filter((image) => image.status !== ProjectImageStatus_Value.Pending)
    .map((image) => image.uuid);
  const hasSelectableImages = selectableImageUuids.length > EMPTY_ITEMS_COUNT;
  const allSelectableImagesSelected =
    hasSelectableImages && selectableImageUuids.every((imageUuid) => selectedImageUuids.has(imageUuid));
  const hasSelectedGroupImages = selectableImageUuids.some((imageUuid) => selectedImageUuids.has(imageUuid));

  return (
    <section
      className={`group overflow-hidden rounded-lg border bg-slate-50 transition-[border-color,box-shadow,background-color] duration-200 ${
        dragging ? 'border-slate-300 shadow-sm' : 'border-slate-200'
      }`}
      onDragOver={(event: DragEvent<HTMLElement>) => {
        onGroupDragOver(event, group.id);
      }}
      onDrop={(event: DragEvent<HTMLElement>) => {
        onGroupDrop(event, group.id);
      }}
      ref={(node) => {
        registerGroupRef(group.id, node);
      }}
    >
      <div
        className="flex items-center gap-3 rounded-t-lg border-b border-slate-200 bg-white px-4 py-3"
        draggable={!readOnly && !editing}
        onDragEnd={onGroupDragEnd}
        onDragStart={(event: DragEvent<HTMLDivElement>) => {
          onGroupDragStart(event, group);
        }}
      >
        <Button
          aria-label={collapsed ? '展开图像组' : '折叠图像组'}
          icon={
            <CaretDownOutlined
              className={`transition-transform duration-300 ${collapsed ? '-rotate-90' : 'rotate-0'}`}
            />
          }
          onClick={() => {
            onToggleCollapsed(group.id);
          }}
          shape="circle"
          size="small"
          type="text"
        />
        <ProjectGroupDragHandle readOnly={readOnly} />
        {editing ? (
          <Input
            autoFocus
            className="h-7 max-w-116 text-sm"
            disabled={savingGroup}
            draggable={false}
            maxLength={32}
            onBlur={() => {
              onSaveEditGroup();
            }}
            onChange={(event: ChangeEvent<HTMLInputElement>) => {
              onEditingGroupNameChange(event.target.value);
            }}
            onDragStart={(event) => {
              event.stopPropagation();
            }}
            onKeyDown={(event) => {
              if (event.key === 'Escape') {
                event.preventDefault();
                onSaveEditGroup();
              }
            }}
            onPressEnter={() => {
              onSaveEditGroup();
            }}
            size="small"
            value={editingGroupName}
          />
        ) : (
          <button
            className="h-7 min-w-0 truncate rounded px-2 text-left text-sm font-semibold leading-7 text-slate-900 hover:bg-slate-100 disabled:cursor-default disabled:hover:bg-transparent"
            disabled={readOnly}
            onClick={() => {
              onStartEditGroup(group);
            }}
            type="button"
          >
            {group.name}
          </button>
        )}
        <span className="rounded bg-blue-50 px-2 py-0.5 text-xs text-[#1677FF]">{group.images.length} 张图像</span>
        {imageStats.failed > EMPTY_ITEMS_COUNT ? (
          <span className="rounded bg-red-50 px-2 py-0.5 text-xs text-red-500">{imageStats.failed} 张失败</span>
        ) : null}
        {imageStats.pending > EMPTY_ITEMS_COUNT ? (
          <span className="rounded bg-emerald-50 px-2 py-0.5 text-xs text-emerald-600">
            {imageStats.pending} 张上传中
          </span>
        ) : null}
        <ProjectAssessGroupHeaderActions
          allSelectableImagesSelected={allSelectableImagesSelected}
          batchMode={batchMode}
          canDeleteGroup={canDeleteGroup}
          deletingGroup={deletingGroup}
          groupId={group.id}
          groupImageCount={group.images.length}
          hasSelectableImages={hasSelectableImages}
          hasSelectedGroupImages={hasSelectedGroupImages}
          onDeleteGroup={onDeleteGroup}
          onDeselectGroupImages={onDeselectGroupImages}
          onSelectGroupImages={onSelectGroupImages}
          selectableImageUuids={selectableImageUuids}
        />
      </div>
      <div
        className={`grid transition-[grid-template-rows,opacity] duration-300 ease-in-out ${
          collapsed ? 'grid-rows-[0fr] opacity-0' : 'grid-rows-[1fr] opacity-100'
        }`}
      >
        <div className="min-h-0 overflow-hidden">
          <div className="grid grid-cols-[repeat(auto-fill,minmax(88px,1fr))] gap-3 p-4">
            {group.images.map((image) => (
              <ProjectAssetImageCard
                batchMode={batchMode}
                deleting={deletingImageUuids.has(image.uuid)}
                image={image}
                key={image.uuid}
                onDeleteImage={onDeleteImage}
                onDragEnd={onImageDragEnd}
                onDragStart={(event, imageUuid, imageStatus) => {
                  onImageDragStart(event, imageUuid, group.id, imageStatus);
                }}
                onOpenImageViewer={onOpenImageViewer}
                onToggleSelectedImage={onToggleSelectedImage}
                readOnly={readOnly}
                selected={selectedImageUuids.has(image.uuid)}
              />
            ))}
            {!readOnly ? (
              <ProjectImageUploadTile
                accept={accept}
                onUploadFiles={(files) => {
                  onUploadFiles(group.id, files);
                }}
                uploading={uploadingImages}
              />
            ) : null}
          </div>
        </div>
      </div>
    </section>
  );
}
