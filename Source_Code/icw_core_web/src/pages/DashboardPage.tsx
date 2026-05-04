import { Card, Col, Row, Statistic } from 'antd';
import type { ReactElement } from 'react';

const DASHBOARD_GRID_GUTTER = 16;

export default function DashboardPage(): ReactElement {
  return (
    <div>
      <div className="mb-6">
        <h1 className="text-xl font-semibold text-slate-900">工作台</h1>
        <p className="mt-1 text-sm text-slate-500">第一阶段已接入登录态，项目与图像处理将在后续迭代实现。</p>
      </div>
      <Row gutter={[DASHBOARD_GRID_GUTTER, DASHBOARD_GRID_GUTTER]}>
        <Col md={8} xs={24}>
          <Card>
            <Statistic title="项目数量" value={0} />
          </Card>
        </Col>
        <Col md={8} xs={24}>
          <Card>
            <Statistic title="待处理图像" value={0} />
          </Card>
        </Col>
        <Col md={8} xs={24}>
          <Card>
            <Statistic title="报告数量" value={0} />
          </Card>
        </Col>
      </Row>
    </div>
  );
}
