import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { request } from '@umijs/max';
import React, { useRef } from 'react';

type Customer = {
  id: number;
  username: string;
  email: string;
  phone: string;
  status: number;
  orderCount: number;
  createdAt: string;
};

const columns: ProColumns<Customer>[] = [
  { title: '用户名', dataIndex: 'username' },
  { title: '邮箱', dataIndex: 'email', search: false, copyable: true },
  { title: '手机号', dataIndex: 'phone', search: false },
  { title: '订单数', dataIndex: 'orderCount', search: false },
  {
    title: '状态',
    dataIndex: 'status',
    valueEnum: {
      0: { text: '禁用', status: 'Error' },
      1: { text: '正常', status: 'Success' },
    },
  },
  { title: '注册时间', dataIndex: 'createdAt', search: false, valueType: 'date' },
  {
    title: '操作',
    valueType: 'option',
    render: () => [
      <a key="detail" href="#">详情</a>,
      <a key="toggle" href="#">禁用</a>,
    ],
  },
];

const CustomerList: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);

  return (
    <PageContainer>
      <ProTable<Customer>
        headerTitle="用户列表"
        actionRef={actionRef}
        rowKey="id"
        search={{ labelWidth: 80 }}
        request={async (params) => {
          return request('/api/customer/list', { params });
        }}
        columns={columns}
      />
    </PageContainer>
  );
};

export default CustomerList;
