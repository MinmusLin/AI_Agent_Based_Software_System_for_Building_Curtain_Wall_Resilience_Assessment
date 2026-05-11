import { DownloadOutlined, LeftOutlined, RightOutlined } from '@ant-design/icons';
import { Button, Modal, Spin } from 'antd';
import type { ReactElement } from 'react';
import { useState } from 'react';

import { ModalTitle } from '@/components/ModalTitle';
import type { ProjectImage } from '@/gen/core/common';
import type { ImageViewerState, ViewerImage } from '@/utils/assetsStage';
import { FIRST_INDEX, formatProjectImageMetadata, KILOBYTE_SIZE_BYTES, NEXT_INDEX_OFFSET } from '@/utils/assetsStage';
import { formatDateTime } from '@/utils/datetime';

const ORIGINAL_IMAGE_DOWNLOAD_EXTENSION = '.png';

interface ProjectAssessViewerProps {
  onClose: () => void;
  onOpenImage: (imageUuid: string) => void;
  uploadedImages: ViewerImage[];
  viewer: ImageViewerState | null;
  viewerImage: ViewerImage | null;
  viewerIndex: number;
}

interface ProjectImageViewerDetailsProps {
  image: ProjectImage | undefined;
}

function imageDimensionText(image?: ProjectImage): string {
  if (!image) {
    return '-';
  }
  return `${String(image.width)} × ${String(image.height)}`;
}

function imageSizeText(image?: ProjectImage): string {
  if (!image) {
    return '-';
  }
  return `${String(Math.round(image.size_bytes / KILOBYTE_SIZE_BYTES))} KB`;
}

function imageDateText(value?: string): string {
  if (!value) {
    return '-';
  }
  return formatDateTime(value, true);
}

async function downloadPresignedObject(downloadUrl: string, fileName: string): Promise<void> {
  const response = await fetch(downloadUrl);
  if (!response.ok) {
    throw new Error('download image failed');
  }

  const blob = await response.blob();
  const objectUrl = URL.createObjectURL(blob);
  try {
    const anchor = document.createElement('a');
    anchor.download = fileName;
    anchor.href = objectUrl;
    anchor.rel = 'noopener';
    document.body.append(anchor);
    anchor.click();
    anchor.remove();
  } finally {
    URL.revokeObjectURL(objectUrl);
  }
}

function ProjectImageViewerDetails({ image }: ProjectImageViewerDetailsProps): ReactElement {
  return (
    <>
      <div className="space-y-4">
        <div>
          <div className="mb-1 font-medium text-slate-900">图像 ID</div>
          <div className="break-all">{image?.uuid ?? '-'}</div>
        </div>
        <div>
          <div className="mb-1 font-medium text-slate-900">文件名称</div>
          <div className="break-all">{image?.file_name ?? '-'}</div>
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <div className="mb-1 font-medium text-slate-900">图像尺寸（宽 × 高）</div>
            <div>{imageDimensionText(image)}</div>
          </div>
          <div>
            <div className="mb-1 font-medium text-slate-900">文件大小</div>
            <div>{imageSizeText(image)}</div>
          </div>
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <div className="mb-1 font-medium text-slate-900">上传开始时间</div>
            <div>{imageDateText(image?.created_at)}</div>
          </div>
          <div>
            <div className="mb-1 font-medium text-slate-900">上传完成时间</div>
            <div>{imageDateText(image?.uploaded_at)}</div>
          </div>
        </div>
      </div>
      <div className="mt-4 flex min-h-0 flex-1 flex-col">
        <div className="mb-1 font-medium text-slate-900">元数据</div>
        <pre className="min-h-0 flex-1 overflow-auto overscroll-contain whitespace-pre-wrap break-all rounded bg-slate-50 p-3 text-xs leading-5">
          {formatProjectImageMetadata(image?.metadata ?? '{}')}
        </pre>
      </div>
    </>
  );
}

export function ProjectAssessViewer({
  onClose,
  onOpenImage,
  uploadedImages,
  viewer,
  viewerImage,
  viewerIndex,
}: ProjectAssessViewerProps): ReactElement {
  const image = viewerImage?.image;
  const title = viewerImage?.image.file_name ?? '图像详情';
  const [downloading, setDownloading] = useState(false);

  const handleDownloadOriginalImage = async (): Promise<void> => {
    if (!image || !viewer || viewer.loading || viewer.originalUrl === '') {
      return;
    }

    setDownloading(true);
    try {
      await downloadPresignedObject(viewer.originalUrl, `${image.uuid}${ORIGINAL_IMAGE_DOWNLOAD_EXTENSION}`);
    } finally {
      setDownloading(false);
    }
  };

  return (
    <Modal
      centered
      footer={null}
      onCancel={onClose}
      open={viewer !== null}
      title={<ModalTitle text={title} />}
      width={1040}
    >
      {viewer ? (
        <div className="grid h-[520px] grid-cols-[minmax(0,1fr)_320px] gap-5">
          <div className="relative flex min-h-0 items-center justify-center overflow-hidden rounded-lg bg-slate-100">
            {viewer.loading ? (
              <Spin />
            ) : (
              <img
                alt={viewerImage?.image.file_name ?? '图像原图'}
                className="max-h-[520px] max-w-full object-contain"
                src={viewer.originalUrl}
              />
            )}
          </div>
          <div className="flex min-h-0 flex-col text-sm text-slate-600">
            <ProjectImageViewerDetails image={image} />
            <div className="flex justify-between pt-4">
              <Button
                disabled={!image || viewer.loading || viewer.originalUrl === ''}
                icon={<DownloadOutlined />}
                loading={downloading}
                onClick={() => {
                  void handleDownloadOriginalImage();
                }}
              >
                下载原图
              </Button>
              <Button
                disabled={viewerIndex <= FIRST_INDEX}
                icon={<LeftOutlined />}
                onClick={() => {
                  if (viewerIndex > FIRST_INDEX) {
                    onOpenImage(uploadedImages[viewerIndex - NEXT_INDEX_OFFSET].image.uuid);
                  }
                }}
              >
                上一张
              </Button>
              <Button
                disabled={viewerIndex < FIRST_INDEX || viewerIndex >= uploadedImages.length - NEXT_INDEX_OFFSET}
                onClick={() => {
                  if (viewerIndex >= FIRST_INDEX && viewerIndex < uploadedImages.length - NEXT_INDEX_OFFSET) {
                    onOpenImage(uploadedImages[viewerIndex + NEXT_INDEX_OFFSET].image.uuid);
                  }
                }}
              >
                下一张
                <RightOutlined />
              </Button>
            </div>
          </div>
        </div>
      ) : null}
    </Modal>
  );
}
