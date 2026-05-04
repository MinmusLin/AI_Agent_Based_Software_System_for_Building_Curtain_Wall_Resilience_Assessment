import type { ProjectImageStatus } from '@/types/common';
import {
  PROJECT_EVENT_TYPE_IMAGE_STATUS_CHANGED,
  PROJECT_IMAGE_STATUS_FAILED,
  PROJECT_IMAGE_STATUS_PENDING,
  PROJECT_IMAGE_STATUS_UPLOADED,
} from '@/types/common';
import type { ProjectGroup, ProjectImage, UploadProjectImageItem } from '@/types/project/assets';
import type { ProjectImageStatusChangedMessage } from '@/types/socket';
import { buildProjectAssetImageBlobs, PROJECT_ASSET_IMAGE_OUTPUT_CONTENT_TYPE } from '@/utils/images';

export const EMPTY_GROUP_ID = '';
export const EMPTY_ITEMS_COUNT = 0;
export const FIRST_INDEX = 0;
export const GROUP_MOVE_ANIMATION_MS = 220;
export const GROUP_MOVE_DELTA_EPSILON = 0.5;
export const KILOBYTE_SIZE_BYTES = 1024;
export const NEXT_INDEX_OFFSET = 1;
export const NOT_FOUND_INDEX = -1;
export const UPLOAD_ACCEPT = 'image/jpeg,image/png,image/webp';
export const WEBSOCKET_RECONNECT_DELAY_MS = 2000;

const JSON_FORMAT_INDENT = 2;

export interface DraggingImage {
  imageUuid: string;
  sourceGroupId: string;
}

export interface ViewerImage {
  groupId: string;
  image: ProjectImage;
}

export interface ImageViewerState {
  imageUuid: string;
  loading: boolean;
  originalUrl: string;
}

export interface PreparedImageUpload {
  height: number;
  metadata: string;
  originalBlob: Blob;
  sourceFile: File;
  thumbnailBlob: Blob;
  width: number;
}

export function sortProjectImages(images: ProjectImage[]): ProjectImage[] {
  return [...images].sort((left, right) => {
    const leftTime = new Date(left.created_at).getTime();
    const rightTime = new Date(right.created_at).getTime();
    return leftTime - rightTime;
  });
}

export function normalizeGroups(groups: ProjectGroup[]): ProjectGroup[] {
  return groups.map((group) => ({
    ...group,
    images: sortProjectImages(group.images),
  }));
}

export function replaceGroup(groups: ProjectGroup[], nextGroup: ProjectGroup): ProjectGroup[] {
  return groups.map((group) => {
    if (group.id !== nextGroup.id) {
      return group;
    }
    return {
      ...group,
      ...nextGroup,
      images: nextGroup.images.length > EMPTY_ITEMS_COUNT ? sortProjectImages(nextGroup.images) : group.images,
    };
  });
}

export function removeGroup(groups: ProjectGroup[], groupId: string): ProjectGroup[] {
  return groups.filter((group) => group.id !== groupId);
}

export function appendImagesToGroup(groups: ProjectGroup[], groupId: string, images: ProjectImage[]): ProjectGroup[] {
  return groups.map((group) => {
    if (group.id !== groupId) {
      return group;
    }
    return {
      ...group,
      images: sortProjectImages([...group.images, ...images]),
    };
  });
}

export function removeImages(groups: ProjectGroup[], imageUuids: string[]): ProjectGroup[] {
  const imageUuidSet = new Set(imageUuids);
  return groups.map((group) => ({
    ...group,
    images: group.images.filter((image) => !imageUuidSet.has(image.uuid)),
  }));
}

export function imageUuidSetFromGroups(groups: ProjectGroup[]): Set<string> {
  return new Set(groups.flatMap((group) => group.images.map((image) => image.uuid)));
}

export function selectedGroupIdsFromGroups(groups: ProjectGroup[], selectedImageUuids: Set<string>): Set<string> {
  const selectedGroupIds = new Set<string>();
  groups.forEach((group) => {
    if (group.images.some((image) => selectedImageUuids.has(image.uuid))) {
      selectedGroupIds.add(group.id);
    }
  });
  return selectedGroupIds;
}

export function pruneSelectedImageUuids(selectedImageUuids: Set<string>, groups: ProjectGroup[]): Set<string> {
  const existingImageUuids = imageUuidSetFromGroups(groups);
  const nextSelectedImageUuids = new Set(
    [...selectedImageUuids].filter((imageUuid) => existingImageUuids.has(imageUuid)),
  );
  if (nextSelectedImageUuids.size === selectedImageUuids.size) {
    return selectedImageUuids;
  }
  return nextSelectedImageUuids;
}

export function replaceImage(groups: ProjectGroup[], nextImage: ProjectImage): ProjectGroup[] {
  const targetGroupIndex = groups.findIndex((group) => group.images.some((image) => image.uuid === nextImage.uuid));
  if (targetGroupIndex < FIRST_INDEX) {
    return groups;
  }

  return groups.map((group, index) => {
    if (index !== targetGroupIndex) {
      return group;
    }
    return {
      ...group,
      images: sortProjectImages(
        group.images.map((image) => {
          if (image.uuid !== nextImage.uuid) {
            return image;
          }
          return nextImage;
        }),
      ),
    };
  });
}

export function moveImagesToGroup(
  groups: ProjectGroup[],
  targetGroupId: string,
  movedImages: ProjectImage[],
): ProjectGroup[] {
  const movedImageUuidSet = new Set(movedImages.map((image) => image.uuid));
  return groups.map((group) => {
    const remainingImages = group.images.filter((image) => !movedImageUuidSet.has(image.uuid));
    if (group.id !== targetGroupId) {
      return {
        ...group,
        images: remainingImages,
      };
    }
    return {
      ...group,
      images: sortProjectImages([...remainingImages, ...movedImages]),
    };
  });
}

export function flattenUploadedImages(groups: ProjectGroup[]): ViewerImage[] {
  return groups.flatMap((group) =>
    group.images
      .filter((image) => image.status === PROJECT_IMAGE_STATUS_UPLOADED)
      .map((image) => ({
        groupId: group.id,
        image,
      })),
  );
}

export function buildUploadMetadata(file: File): string {
  return JSON.stringify({
    source_content_type: file.type,
    source_file_name: file.name,
    source_size_bytes: file.size,
  });
}

export function formatProjectImageMetadata(metadata: string): string {
  const rawMetadata = metadata.trim();
  if (rawMetadata === '') {
    return '{}';
  }

  try {
    const parsedMetadata = JSON.parse(rawMetadata) as unknown;
    return JSON.stringify(parsedMetadata, null, JSON_FORMAT_INDENT);
  } catch {
    return metadata;
  }
}

export function groupMovePosition(
  groups: ProjectGroup[],
  sourceGroupId: string,
  targetGroupId: string,
): ProjectGroup[] {
  const sourceIndex = groups.findIndex((group) => group.id === sourceGroupId);
  const targetIndex = groups.findIndex((group) => group.id === targetGroupId);
  if (sourceIndex < FIRST_INDEX || targetIndex < FIRST_INDEX || sourceIndex === targetIndex) {
    return groups;
  }

  const nextGroups = [...groups];
  const [sourceGroup] = nextGroups.splice(sourceIndex, NEXT_INDEX_OFFSET);
  nextGroups.splice(targetIndex, EMPTY_ITEMS_COUNT, sourceGroup);
  return nextGroups;
}

export function isSameGroupOrder(leftGroups: ProjectGroup[], rightGroups: ProjectGroup[]): boolean {
  if (leftGroups.length !== rightGroups.length) {
    return false;
  }
  return leftGroups.every((group, index) => group.id === rightGroups[index]?.id);
}

export function groupIdsSet(groups: ProjectGroup[]): Set<string> {
  return new Set(groups.map((group) => group.id));
}

export function groupMovePayload(
  projectId: string,
  groupId: string,
  nextGroups: ProjectGroup[],
): {
  move_to_first: boolean;
  move_to_last: boolean;
  next_group_id: string;
  previous_group_id: string;
  project_id: string;
  group_id: string;
} {
  const nextIndex = nextGroups.findIndex((group) => group.id === groupId);
  const previousGroupId = nextIndex > FIRST_INDEX ? nextGroups[nextIndex - NEXT_INDEX_OFFSET].id : EMPTY_GROUP_ID;
  const nextGroupId =
    nextIndex >= FIRST_INDEX && nextIndex < nextGroups.length - NEXT_INDEX_OFFSET
      ? nextGroups[nextIndex + NEXT_INDEX_OFFSET].id
      : EMPTY_GROUP_ID;
  return {
    group_id: groupId,
    move_to_first: nextIndex === FIRST_INDEX,
    move_to_last: nextIndex === nextGroups.length - NEXT_INDEX_OFFSET,
    next_group_id: nextGroupId,
    previous_group_id: previousGroupId,
    project_id: projectId,
  };
}

export function parseProjectImageStatusChangedMessage(data: unknown): ProjectImageStatusChangedMessage | null {
  if (typeof data !== 'string') {
    return null;
  }

  let value: unknown;
  try {
    value = JSON.parse(data) as unknown;
  } catch {
    return null;
  }

  if (
    !isRecord(value) ||
    value.type !== PROJECT_EVENT_TYPE_IMAGE_STATUS_CHANGED ||
    typeof value.project_id !== 'string'
  ) {
    return null;
  }
  const { image } = value;
  if (!isProjectImage(image)) {
    return null;
  }
  return {
    image,
    project_id: value.project_id,
    type: PROJECT_EVENT_TYPE_IMAGE_STATUS_CHANGED,
  };
}

export async function putPresignedObject(uploadUrl: string, blob: Blob, contentType: string): Promise<void> {
  if (uploadUrl.trim() === '') {
    throw new Error('presigned upload url is empty');
  }

  const headers = new Headers();
  headers.set('Content-Type', contentType);
  const response = await fetch(uploadUrl, {
    body: blob,
    headers,
    method: 'PUT',
  });
  if (!response.ok) {
    throw new Error('presigned upload failed');
  }
}

export async function prepareUploadImage(file: File): Promise<PreparedImageUpload> {
  const blobs = await buildProjectAssetImageBlobs(file);
  return {
    height: blobs.height,
    metadata: buildUploadMetadata(file),
    originalBlob: blobs.originalBlob,
    sourceFile: file,
    thumbnailBlob: blobs.thumbnailBlob,
    width: blobs.width,
  };
}

export function uploadItemFromPreparedImage(image: PreparedImageUpload): UploadProjectImageItem {
  return {
    content_type: PROJECT_ASSET_IMAGE_OUTPUT_CONTENT_TYPE,
    file_name: image.sourceFile.name,
    height: image.height,
    metadata: image.metadata,
    size_bytes: image.originalBlob.size,
    width: image.width,
  };
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null;
}

function isProjectImageStatus(value: unknown): value is ProjectImageStatus {
  return (
    value === PROJECT_IMAGE_STATUS_PENDING ||
    value === PROJECT_IMAGE_STATUS_UPLOADED ||
    value === PROJECT_IMAGE_STATUS_FAILED
  );
}

function isProjectImage(value: unknown): value is ProjectImage {
  if (!isRecord(value)) {
    return false;
  }
  return (
    typeof value.uuid === 'string' &&
    typeof value.file_name === 'string' &&
    typeof value.content_type === 'string' &&
    typeof value.size_bytes === 'number' &&
    typeof value.width === 'number' &&
    typeof value.height === 'number' &&
    typeof value.metadata === 'string' &&
    isProjectImageStatus(value.status) &&
    typeof value.thumbnail_url === 'string' &&
    typeof value.uploaded_at === 'string' &&
    typeof value.created_at === 'string'
  );
}
