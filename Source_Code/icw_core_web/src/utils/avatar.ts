const AVATAR_SIZE = 128;

const ALLOWED_SOURCE_TYPES = new Set(['image/jpeg', 'image/png', 'image/webp']);

export const AVATAR_OUTPUT_CONTENT_TYPE = 'image/png';
export const AVATAR_MAX_SOURCE_SIZE = 2 * 1024 * 1024;

export function isAllowedAvatarFile(file: File): boolean {
  return ALLOWED_SOURCE_TYPES.has(file.type);
}

export async function resizeAvatarToPng(file: File): Promise<Blob> {
  const image = await createImageBitmap(file);
  try {
    const canvas = document.createElement('canvas');
    canvas.width = AVATAR_SIZE;
    canvas.height = AVATAR_SIZE;

    const context = canvas.getContext('2d');
    if (!context) {
      throw new Error('canvas context is unavailable');
    }

    const sourceSize = Math.min(image.width, image.height);
    const sourceX = Math.floor((image.width - sourceSize) / 2);
    const sourceY = Math.floor((image.height - sourceSize) / 2);
    context.clearRect(0, 0, AVATAR_SIZE, AVATAR_SIZE);
    context.drawImage(image, sourceX, sourceY, sourceSize, sourceSize, 0, 0, AVATAR_SIZE, AVATAR_SIZE);

    return await canvasToPngBlob(canvas);
  } finally {
    image.close();
  }
}

function canvasToPngBlob(canvas: HTMLCanvasElement): Promise<Blob> {
  return new Promise((resolve, reject) => {
    canvas.toBlob((blob) => {
      if (!blob) {
        reject(new Error('failed to encode avatar'));
        return;
      }
      resolve(blob);
    }, AVATAR_OUTPUT_CONTENT_TYPE);
  });
}
