const ALLOWED_SOURCE_IMAGE_TYPES = new Set(['image/jpeg', 'image/png', 'image/webp']);

const AVATAR_SIZE = 128;
const PROJECT_THUMBNAIL_SIZE = 512;

const PNG_CONTENT_TYPE = 'image/png';

export const AVATAR_OUTPUT_CONTENT_TYPE = PNG_CONTENT_TYPE;
export const PROJECT_THUMBNAIL_OUTPUT_CONTENT_TYPE = PNG_CONTENT_TYPE;

function isAllowedSourceImageFile(file: File): boolean {
  return ALLOWED_SOURCE_IMAGE_TYPES.has(file.type);
}

export function isAllowedAvatarFile(file: File): boolean {
  return isAllowedSourceImageFile(file);
}

export function isAllowedProjectThumbnailFile(file: File): boolean {
  return isAllowedSourceImageFile(file);
}

export async function resizeAvatarToPng(file: File): Promise<Blob> {
  return resizeCenteredSquareImageToPng(file, AVATAR_SIZE);
}

export async function resizeProjectThumbnailToPng(file: File): Promise<Blob> {
  return resizeCenteredSquareImageToPng(file, PROJECT_THUMBNAIL_SIZE);
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

    const sourceSize = Math.min(image.width, image.height);
    const sourceX = Math.floor((image.width - sourceSize) / 2);
    const sourceY = Math.floor((image.height - sourceSize) / 2);
    context.clearRect(0, 0, targetSize, targetSize);
    context.drawImage(image, sourceX, sourceY, sourceSize, sourceSize, 0, 0, targetSize, targetSize);

    return await canvasToPngBlob(canvas);
  } finally {
    image.close();
  }
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
