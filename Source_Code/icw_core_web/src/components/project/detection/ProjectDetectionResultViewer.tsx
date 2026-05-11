import { LeftOutlined, RightOutlined } from '@ant-design/icons';
import { Button, Modal, Spin, Tabs } from 'antd';
import type { ReactElement } from 'react';
import { useMemo, useState } from 'react';

import { ModalTitle } from '@/components/ModalTitle';
import type { GetImageDetectionResultResponse } from '@/gen/core/api/project_detection';
import type { ViewerImage } from '@/utils/assetsStage';
import { FIRST_INDEX, NEXT_INDEX_OFFSET } from '@/utils/assetsStage';

import { detectionResultTabs, viewerTitle } from './ProjectDetectionResultTabs';
import {
  DEFAULT_ACTIVE_TAB,
  DetectionPreview,
  DetectionResultTab,
  EMPTY_ITEMS_COUNT,
  OriginalInfoTab,
  type ProjectDetectionReviewProps,
  ReviewPanel,
  SUMMARY_TAB,
} from './ProjectDetectionResultViewerParts';
import { SummaryResultTab } from './ProjectDetectionSummaryResultTab';

interface ProjectDetectionResultViewerProps {
  loading: boolean;
  onClose: () => void;
  onOpenImage: (imageUuid: string) => void;
  open: boolean;
  review?: ProjectDetectionReviewProps;
  result: GetImageDetectionResultResponse | null;
  uploadedImages: ViewerImage[];
  viewerIndex: number;
}

export function ProjectDetectionResultViewer({
  loading,
  onClose,
  onOpenImage,
  open,
  review,
  result,
  uploadedImages,
  viewerIndex,
}: ProjectDetectionResultViewerProps): ReactElement {
  const [activeTab, setActiveTab] = useState(DEFAULT_ACTIVE_TAB);
  const [groupIndex, setGroupIndex] = useState(FIRST_INDEX);
  const image = result?.image;
  const title = viewerTitle(image, uploadedImages, viewerIndex);
  const reasoningTabs = useMemo(() => (loading ? [] : detectionResultTabs(result)), [loading, result]);
  const activeDetectionResult = loading ? undefined : reasoningTabs.find((item) => item.key === activeTab)?.result;
  const tabs = useMemo(
    () => [
      {
        children: (
          <OriginalInfoTab image={image} mainTaskUuid={result?.status?.main_task_uuid} status={result?.status} />
        ),
        key: DEFAULT_ACTIVE_TAB,
        label: '图像数据',
      },
      ...reasoningTabs.map((item) => ({
        children: (
          <DetectionResultTab
            groupIndex={groupIndex}
            onGroupIndexChange={setGroupIndex}
            result={item.result}
            taskKey={item.key}
          />
        ),
        key: item.key,
        label: item.label,
      })),
      ...(reasoningTabs.length > EMPTY_ITEMS_COUNT && result?.summary_result
        ? [
            {
              children: <SummaryResultTab result={result.summary_result} />,
              key: SUMMARY_TAB,
              label: 'Agent 总结',
            },
          ]
        : []),
    ],
    [groupIndex, image, reasoningTabs, result],
  );
  const currentActiveTab = !loading && tabs.some((tab) => tab.key === activeTab) ? activeTab : DEFAULT_ACTIVE_TAB;
  if (loading) {
    return (
      <Modal centered footer={null} onCancel={onClose} open={open} title={<ModalTitle text={title} />} width={1040}>
        <div className="flex h-[520px] flex-col items-center justify-center gap-3 text-sm text-slate-500">
          <Spin />
          <div>正在加载检测结果</div>
        </div>
      </Modal>
    );
  }

  return (
    <Modal centered footer={null} onCancel={onClose} open={open} title={<ModalTitle text={title} />} width={1040}>
      <div>
        <div className="grid h-[520px] grid-cols-[minmax(0,1fr)_320px] gap-5">
          <div className="relative flex min-h-0 items-center justify-center overflow-hidden rounded-lg bg-slate-100">
            <DetectionPreview
              activeTab={currentActiveTab}
              groupIndex={groupIndex}
              image={image}
              loading={loading}
              originalUrl={result?.original_url}
              result={activeDetectionResult}
              taskKey={currentActiveTab}
            />
          </div>
          <div className="flex min-h-0 flex-col">
            <Tabs
              activeKey={currentActiveTab}
              className="min-h-0 flex-1 [&_.ant-tabs-content-holder]:min-h-0 [&_.ant-tabs-content]:h-full [&_.ant-tabs-tabpane]:h-full"
              items={tabs}
              onChange={(key) => {
                setActiveTab(key);
                setGroupIndex(FIRST_INDEX);
              }}
              size="small"
            />
            <div className="flex justify-between gap-3 pt-4">
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
        {review?.enabled ? <ReviewPanel review={review} /> : null}
      </div>
    </Modal>
  );
}
