const ALLOWED_SOURCE_IMAGE_TYPES = new Set(['image/jpeg', 'image/png', 'image/webp']);

const AVATAR_SIZE = 128;
const PROJECT_THUMBNAIL_SIZE = 512;
const PROJECT_ASSET_THUMBNAIL_SIZE = 256;
const CENTER_OFFSET_DIVISOR = 2;
const CANVAS_ORIGIN = 0;

const PNG_CONTENT_TYPE = 'image/png';

export const AVATAR_OUTPUT_CONTENT_TYPE = PNG_CONTENT_TYPE;
export const PROJECT_THUMBNAIL_OUTPUT_CONTENT_TYPE = PNG_CONTENT_TYPE;
export const PROJECT_ASSET_IMAGE_OUTPUT_CONTENT_TYPE = PNG_CONTENT_TYPE;
export const PROJECT_ASSET_THUMBNAIL_OUTPUT_CONTENT_TYPE = PNG_CONTENT_TYPE;

export interface ProjectAssetImageBlobs {
  height: number;
  originalBlob: Blob;
  thumbnailBlob: Blob;
  width: number;
}

function isAllowedSourceImageFile(file: File): boolean {
  return ALLOWED_SOURCE_IMAGE_TYPES.has(file.type);
}

export function isAllowedAvatarFile(file: File): boolean {
  return isAllowedSourceImageFile(file);
}

export function isAllowedProjectThumbnailFile(file: File): boolean {
  return isAllowedSourceImageFile(file);
}

export function isAllowedProjectAssetImageFile(file: File): boolean {
  return isAllowedSourceImageFile(file);
}

export async function resizeAvatarToPng(file: File): Promise<Blob> {
  return resizeCenteredSquareImageToPng(file, AVATAR_SIZE);
}

export async function resizeProjectThumbnailToPng(file: File): Promise<Blob> {
  return resizeCenteredSquareImageToPng(file, PROJECT_THUMBNAIL_SIZE);
}

export async function buildProjectAssetImageBlobs(file: File): Promise<ProjectAssetImageBlobs> {
  const image = await createImageBitmap(file);
  try {
    const originalCanvas = document.createElement('canvas');
    originalCanvas.width = image.width;
    originalCanvas.height = image.height;

    const originalContext = originalCanvas.getContext('2d');
    if (!originalContext) {
      throw new Error('canvas context is unavailable');
    }
    originalContext.clearRect(CANVAS_ORIGIN, CANVAS_ORIGIN, image.width, image.height);
    originalContext.drawImage(image, CANVAS_ORIGIN, CANVAS_ORIGIN);

    const thumbnailCanvas = document.createElement('canvas');
    thumbnailCanvas.width = PROJECT_ASSET_THUMBNAIL_SIZE;
    thumbnailCanvas.height = PROJECT_ASSET_THUMBNAIL_SIZE;

    const thumbnailContext = thumbnailCanvas.getContext('2d');
    if (!thumbnailContext) {
      throw new Error('canvas context is unavailable');
    }
    drawContainedImage(thumbnailContext, image, PROJECT_ASSET_THUMBNAIL_SIZE);

    return {
      height: image.height,
      originalBlob: await canvasToPngBlob(originalCanvas),
      thumbnailBlob: await canvasToPngBlob(thumbnailCanvas),
      width: image.width,
    };
  } finally {
    image.close();
  }
}

async function resizeCenteredSquareImageToPng(file: File, targetSize: number): Promise<Blob> {
  const image = await createImageBitmap(file);
  try {
    const canvas = document.createElement('canvas');
    canvas.width = targetSize;
    canvas.height = targetSize;

    const context = canvas.getContext('2d');
    if (!context) {
      throw new Error('canvas context is unavailable');
    }

    drawCenteredSquareImage(context, image, targetSize);

    return await canvasToPngBlob(canvas);
  } finally {
    image.close();
  }
}

function drawCenteredSquareImage(context: CanvasRenderingContext2D, image: ImageBitmap, targetSize: number): void {
  const sourceSize = Math.min(image.width, image.height);
  const sourceX = Math.floor((image.width - sourceSize) / CENTER_OFFSET_DIVISOR);
  const sourceY = Math.floor((image.height - sourceSize) / CENTER_OFFSET_DIVISOR);
  context.clearRect(CANVAS_ORIGIN, CANVAS_ORIGIN, targetSize, targetSize);
  context.drawImage(
    image,
    sourceX,
    sourceY,
    sourceSize,
    sourceSize,
    CANVAS_ORIGIN,
    CANVAS_ORIGIN,
    targetSize,
    targetSize,
  );
}

function drawContainedImage(context: CanvasRenderingContext2D, image: ImageBitmap, targetSize: number): void {
  const scale = Math.min(targetSize / image.width, targetSize / image.height);
  const targetWidth = Math.round(image.width * scale);
  const targetHeight = Math.round(image.height * scale);
  const targetX = Math.floor((targetSize - targetWidth) / CENTER_OFFSET_DIVISOR);
  const targetY = Math.floor((targetSize - targetHeight) / CENTER_OFFSET_DIVISOR);

  context.clearRect(CANVAS_ORIGIN, CANVAS_ORIGIN, targetSize, targetSize);
  context.drawImage(image, targetX, targetY, targetWidth, targetHeight);
}

function canvasToPngBlob(canvas: HTMLCanvasElement): Promise<Blob> {
  return new Promise((resolve, reject) => {
    canvas.toBlob((blob) => {
      if (!blob) {
        reject(new Error('failed to encode image'));
        return;
      }
      resolve(blob);
    }, PNG_CONTENT_TYPE);
  });
}
