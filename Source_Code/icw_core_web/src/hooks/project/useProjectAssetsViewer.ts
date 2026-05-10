import { useCallback, useMemo, useState } from 'react';

import { getErrorMessage } from '@/api/http';
import { getProjectImageOriginal } from '@/api/project/assets';
import type { ImageViewerState, ViewerImage } from '@/utils/assetsStage';
import { FIRST_INDEX, NOT_FOUND_INDEX } from '@/utils/assetsStage';

interface UseProjectAssetsViewerParams {
  onError: (message: string) => void;
  projectId: string;
  uploadedImages: ViewerImage[];
}

interface UseProjectAssetsViewerResult {
  openImageViewer: (imageUuid: string) => Promise<void>;
  setViewer: (viewer: ImageViewerState | null) => void;
  viewer: ImageViewerState | null;
  viewerImage: ViewerImage | null;
  viewerIndex: number;
}

export function useProjectAssetsViewer({
  onError,
  projectId,
  uploadedImages,
}: UseProjectAssetsViewerParams): UseProjectAssetsViewerResult {
  const [viewer, setViewer] = useState<ImageViewerState | null>(null);

  const openImageViewer = useCallback(
    async (imageUuid: string): Promise<void> => {
      setViewer({
        imageUuid,
        loading: true,
        originalUrl: '',
      });
      try {
        const data = await getProjectImageOriginal({
          image_uuid: imageUuid,
          project_id: projectId,
        });
        setViewer({
          imageUuid,
          loading: false,
          originalUrl: data.original_url,
        });
      } catch (error: unknown) {
        setViewer(null);
        onError(getErrorMessage(error));
      }
    },
    [onError, projectId],
  );

  const viewerIndex = useMemo(() => {
    if (!viewer) {
      return NOT_FOUND_INDEX;
    }
    return uploadedImages.findIndex((item) => item.image.uuid === viewer.imageUuid);
  }, [uploadedImages, viewer]);
  const viewerImage = viewerIndex >= FIRST_INDEX ? uploadedImages[viewerIndex] : null;

  return {
    openImageViewer,
    setViewer,
    viewer,
    viewerImage,
    viewerIndex,
  };
}
