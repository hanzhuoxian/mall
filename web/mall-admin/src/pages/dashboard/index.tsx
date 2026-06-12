import { PageContainer } from '@ant-design/pro-components';
import { Card, Col, Row, Statistic, Table, Tag } from 'antd';
import React from 'react';

const recentOrders = [
  { key: '1', orderNo: 'ORD-20240001', customer: '张三', amount: 299.0, status: 1 },
  { key: '2', orderNo: 'ORD-20240002', customer: '李四', amount: 1580.0, status: 2 },
  { key: '3', orderNo: 'ORD-20240003', customer: '王五', amount: 88.0, status: 0 },
  { key: '4', orderNo: 'ORD-20240004', customer: '赵六', amount: 450.0, status: 3 },
  { key: '5', orderNo: 'ORD-20240005', customer: '孙七', amount: 720.0, status: 2 },
];

const statusMap: Record<number, { text: string; color: string }> = {
  0: { text: '待付款', color: 'default' },
  1: { text: '待发货', color: 'processing' },
  2: { text: '已发货', color: 'success' },
  3: { text: '已完成', color: 'cyan' },
};

const orderColumns = [
  { title: '订单编号', dataIndex: 'orderNo' },
  { title: '客户', dataIndex: 'customer' },
  {
    title: '金额',
    dataIndex: 'amount',
    render: (val: number) => `¥ ${val.toFixed(2)}`,
  },
  {
    title: '状态',
    dataIndex: 'status',
    render: (val: number) => (
      <Tag color={statusMap[val]?.color}>{statusMap[val]?.text}</Tag>
    ),
  },
];

const Dashboard: React.FC = () => {
  return (
    <PageContainer title="仪表盘">
      <Row gutter={16} style={{ marginBottom: 24 }}>
        <Col span={6}>
          <Card>
            <Statistic title="今日销售额" value={12680} prefix="¥" precision={2} styles={{ content: { color: '#1677ff' } }} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="今日订单数" value={128} styles={{ content: { color: '#52c41a' } }} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="商品总数" value={864} styles={{ content: { color: '#faad14' } }} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="用户总数" value={3256} styles={{ content: { color: '#722ed1' } }} />
          </Card>
        </Col>
      </Row>
      <Card title="最近订单">
        <Table
          dataSource={recentOrders}
          columns={orderColumns}
          pagination={false}
          size="small"
        />
      </Card>
    </PageContainer>
  );
};

export default Dashboard;
