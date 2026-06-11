import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { request } from '@umijs/max';
import { Button, Popconfirm, Tag, message } from 'antd';
import React, { useRef } from 'react';

type Product = {
  id: number;
  name: string;
  category: string;
  price: number;
  stock: number;
  status: number;
  createdAt: string;
};

const columns: ProColumns<Product>[] = [
  { title: '商品名称', dataIndex: 'name', ellipsis: true },
  { title: '分类', dataIndex: 'category', search: false },
  {
    title: '价格',
    dataIndex: 'price',
    search: false,
    renderText: (val: number) => `¥ ${val.toFixed(2)}`,
  },
  {
    title: '库存',
    dataIndex: 'stock',
    search: false,
    render: (_, record) =>
      record.stock === 0 ? <Tag color="error">缺货</Tag> : record.stock,
  },
  {
    title: '状态',
    dataIndex: 'status',
    valueEnum: {
      0: { text: '下架', status: 'Default' },
      1: { text: '上架', status: 'Success' },
    },
  },
  { title: '创建时间', dataIndex: 'createdAt', search: false, valueType: 'date' },
  {
    title: '操作',
    valueType: 'option',
    render: (_, record, __, action) => [
      <a key="edit" href="#">编辑</a>,
      <Popconfirm
        key="delete"
        title="确认删除该商品？"
        onConfirm={async () => {
          await request(`/api/product/${record.id}`, { method: 'DELETE' });
          message.success('删除成功');
          action?.reload();
        }}
      >
        <a href="#" style={{ color: '#ff4d4f' }}>删除</a>
      </Popconfirm>,
    ],
  },
];

const ProductList: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);

  return (
    <PageContainer>
      <ProTable<Product>
        headerTitle="商品列表"
        actionRef={actionRef}
        rowKey="id"
        search={{ labelWidth: 80 }}
        toolBarRender={() => [
          <Button key="add" type="primary">新增商品</Button>,
        ]}
        request={async (params) => {
          return request('/api/product/list', { params });
        }}
        columns={columns}
      />
    </PageContainer>
  );
};

export default ProductList;
