import { Button, Empty } from 'antd';
import type { ReactElement } from 'react';

export default function ProjectsPage(): ReactElement {
  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold text-slate-900">项目管理</h1>
          <p className="mt-1 text-sm text-slate-500">项目创建和图像上传会在下一阶段接入。</p>
        </div>
        <Button disabled type="primary">
          新建项目
        </Button>
      </div>
      <div className="rounded border border-slate-200 bg-white py-16">
        <Empty description="暂无项目" />
      </div>
    </div>
  );
}
