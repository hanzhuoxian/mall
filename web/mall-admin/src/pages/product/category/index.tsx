import type { ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { request } from '@umijs/max';
import { Button, Tag } from 'antd';
import React from 'react';

type Category = {
  id: number;
  name: string;
  parentId: number;
  sort: number;
  status: number;
  productCount: number;
};

const columns: ProColumns<Category>[] = [
  { title: '分类名称', dataIndex: 'name' },
  {
    title: '上级分类',
    dataIndex: 'parentId',
    renderText: (val: number) => (val === 0 ? '顶级分类' : `分类 #${val}`),
  },
  { title: '排序', dataIndex: 'sort', search: false },
  { title: '商品数', dataIndex: 'productCount', search: false },
  {
    title: '状态',
    dataIndex: 'status',
    render: (_, record) =>
      record.status === 1 ? <Tag color="success">启用</Tag> : <Tag color="default">停用</Tag>,
  },
  {
    title: '操作',
    valueType: 'option',
    render: () => [
      <a key="edit" href="#">编辑</a>,
      <a key="delete" href="#" style={{ color: '#ff4d4f' }}>删除</a>,
    ],
  },
];

const ProductCategory: React.FC = () => {
  return (
    <PageContainer>
      <ProTable<Category>
        headerTitle="商品分类"
        rowKey="id"
        search={{ labelWidth: 80 }}
        toolBarRender={() => [
          <Button key="add" type="primary">新增分类</Button>,
        ]}
        request={async () => {
          return request('/api/product/category');
        }}
        columns={columns}
        pagination={false}
      />
    </PageContainer>
  );
};

export default ProductCategory;
