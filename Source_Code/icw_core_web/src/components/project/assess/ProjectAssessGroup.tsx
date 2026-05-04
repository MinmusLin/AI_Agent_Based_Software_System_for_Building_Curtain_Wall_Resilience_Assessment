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

import type { ProjectImageStatus } from '@/types/common';
import {
  PROJECT_IMAGE_STATUS_FAILED,
  PROJECT_IMAGE_STATUS_PENDING,
  PROJECT_IMAGE_STATUS_UPLOADED,
} from '@/types/common';
import type { ProjectGroup, ProjectImage } from '@/types/project/assets';
import { NEXT_INDEX_OFFSET } from '@/utils/assetsStage';

interface ProjectAssetImageCardProps {
  batchMode: boolean;
  deleting: boolean;
  image: ProjectImage;
  onDeleteImage: (imageUuid: string) => void;
  onDragEnd: () => void;
  onDragStart: (event: DragEvent<HTMLDivElement>, imageUuid: string) => void;
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
  onCancelEditGroup: () => void;
  onDeleteGroup: (groupId: string) => void;
  onDeleteImage: (imageUuid: string) => void;
  onEditingGroupNameChange: (name: string) => void;
  onGroupDragEnd: () => void;
  onGroupDragOver: (event: DragEvent<HTMLElement>, groupId: string) => void;
  onGroupDragStart: (event: DragEvent<HTMLDivElement>, group: ProjectGroup) => void;
  onGroupDrop: (event: DragEvent<HTMLElement>, groupId: string) => void;
  onImageDragEnd: () => void;
  onImageDragStart: (event: DragEvent<HTMLDivElement>, imageUuid: string, sourceGroupId: string) => void;
  onOpenImageViewer: (imageUuid: string) => void;
  onSaveEditGroup: () => void;
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

function imageStatusNode(status: ProjectImageStatus): ReactElement {
  switch (status) {
    case PROJECT_IMAGE_STATUS_PENDING:
      return <LoadingOutlined className="text-2xl text-slate-400" />;
    case PROJECT_IMAGE_STATUS_FAILED:
      return (
        <Tooltip title="图像已失效">
          <ImageUnavailableIcon />
        </Tooltip>
      );
    case PROJECT_IMAGE_STATUS_UPLOADED:
      return (
        <Tooltip title="图像已失效">
          <ImageUnavailableIcon />
        </Tooltip>
      );
    default:
      return (
        <Tooltip title="图像状态异常">
          <ImageUnavailableIcon />
        </Tooltip>
      );
  }
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
  const imageUploaded = image.status === PROJECT_IMAGE_STATUS_UPLOADED;
  const imageFailed = image.status === PROJECT_IMAGE_STATUS_FAILED;
  const imageBatchSelectable = batchMode && image.status !== PROJECT_IMAGE_STATUS_PENDING;
  const imageReady = imageUploaded && image.thumbnail_url !== '';
  const imageActionVisible = !batchMode && (imageUploaded || (!readOnly && imageFailed));

  return (
    <div
      className={`group/image relative aspect-square overflow-hidden rounded-lg border bg-white ${
        selected ? 'border-[#1677FF]' : 'border-slate-200'
      } ${batchMode ? 'cursor-pointer' : ''}`}
      draggable={!readOnly && !batchMode}
      onClick={() => {
        if (imageBatchSelectable) {
          onToggleSelectedImage(image.uuid);
        }
      }}
      onDragEnd={onDragEnd}
      onDragStart={(event: DragEvent<HTMLDivElement>) => {
        onDragStart(event, image.uuid);
      }}
    >
      {imageBatchSelectable ? (
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
      {imageReady ? (
        <img alt={image.file_name} className="size-full object-cover" src={image.thumbnail_url} />
      ) : (
        <div className="flex size-full flex-col items-center justify-center gap-2 px-3 text-center text-xs text-slate-500">
          {imageStatusNode(image.status)}
        </div>
      )}
      {imageActionVisible ? (
        <ProjectImageActionOverlay
          deleting={deleting}
          imageUploaded={imageUploaded}
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
  onCancelEditGroup,
  onDeleteGroup,
  onDeleteImage,
  onEditingGroupNameChange,
  onGroupDragEnd,
  onGroupDragOver,
  onGroupDragStart,
  onGroupDrop,
  onImageDragEnd,
  onImageDragStart,
  onOpenImageViewer,
  onSaveEditGroup,
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

  return (
    <section
      className={`overflow-hidden rounded-lg border bg-slate-50 transition-[border-color,box-shadow,background-color] duration-200 ${
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
        className="group flex items-center gap-3 rounded-t-lg border-b border-slate-200 bg-white px-4 py-3"
        draggable={!readOnly}
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
        <DragOutlined className="cursor-grab text-slate-400" />
        {editing ? (
          <Input
            autoFocus
            className="h-7 max-w-72 text-sm"
            disabled={savingGroup}
            maxLength={32}
            onBlur={() => {
              onSaveEditGroup();
            }}
            onChange={(event: ChangeEvent<HTMLInputElement>) => {
              onEditingGroupNameChange(event.target.value);
            }}
            onKeyDown={(event) => {
              if (event.key === 'Escape') {
                onCancelEditGroup();
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
        <span className="rounded bg-slate-100 px-2 py-0.5 text-xs text-slate-500">{group.images.length} 张图像</span>
        {canDeleteGroup ? (
          <Button
            aria-label="删除图像组"
            className="ml-auto opacity-0 transition-opacity duration-200 group-hover:opacity-100"
            danger
            icon={<DeleteOutlined />}
            loading={deletingGroup}
            onClick={() => {
              onDeleteGroup(group.id);
            }}
            shape="circle"
            size="small"
            type="text"
          />
        ) : null}
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
                onDragStart={(event, imageUuid) => {
                  onImageDragStart(event, imageUuid, group.id);
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
