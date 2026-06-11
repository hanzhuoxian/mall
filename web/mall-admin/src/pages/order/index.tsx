import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { request } from '@umijs/max';
import React, { useRef } from 'react';

type Order = {
  id: number;
  orderNo: string;
  customer: string;
  amount: number;
  status: number;
  itemCount: number;
  createdAt: string;
};

const columns: ProColumns<Order>[] = [
  { title: '订单编号', dataIndex: 'orderNo', copyable: true },
  { title: '客户', dataIndex: 'customer', search: false },
  {
    title: '金额',
    dataIndex: 'amount',
    search: false,
    renderText: (val: number) => `¥ ${val.toFixed(2)}`,
  },
  { title: '商品数', dataIndex: 'itemCount', search: false },
  {
    title: '状态',
    dataIndex: 'status',
    valueEnum: {
      0: { text: '待付款', status: 'Default' },
      1: { text: '待发货', status: 'Warning' },
      2: { text: '已发货', status: 'Processing' },
      3: { text: '已完成', status: 'Success' },
      4: { text: '已取消', status: 'Error' },
    },
  },
  { title: '下单时间', dataIndex: 'createdAt', search: false, valueType: 'date' },
  {
    title: '操作',
    valueType: 'option',
    render: () => [
      <a key="detail" href="#">详情</a>,
    ],
  },
];

const OrderList: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);

  return (
    <PageContainer>
      <ProTable<Order>
        headerTitle="订单列表"
        actionRef={actionRef}
        rowKey="id"
        search={{ labelWidth: 80 }}
        request={async (params) => {
          return request('/api/order/list', { params });
        }}
        columns={columns}
      />
    </PageContainer>
  );
};

export default OrderList;
