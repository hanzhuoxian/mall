import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import {
  ModalForm,
  PageContainer,
  ProFormSelect,
  ProFormText,
  ProTable,
} from '@ant-design/pro-components';
import { Button, message, Popconfirm, Space } from 'antd';
import dayjs from 'dayjs';
import React, { useRef, useState } from 'react';
import {
  type CreateUserBody,
  createUser,
  deleteUser,
  deleteUsers,
  listUsers,
  type UpdateUserBody,
  type User,
  updateUser,
} from './service';

const statusValueEnum = {
  0: { text: '禁用', status: 'Error' },
  1: { text: '正常', status: 'Success' },
} as const;

// 后端时间戳为 RFC3339 字符串，转为可读时间。
const formatTimestamp = (ts?: string) =>
  ts ? dayjs(ts).format('YYYY-MM-DD HH:mm:ss') : '-';

const UserList: React.FC = () => {
  const actionRef = useRef<ActionType | null>(null);
  // editing 为 undefined 时表示新增；为 User 时表示编辑该用户。
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<User | undefined>(undefined);

  const openCreate = () => {
    setEditing(undefined);
    setModalOpen(true);
  };

  const openEdit = (record: User) => {
    setEditing(record);
    setModalOpen(true);
  };

  const columns: ProColumns<User>[] = [
    { title: '用户名', dataIndex: 'username', copyable: true },
    { title: '姓名', dataIndex: 'name' },
    { title: '昵称', dataIndex: 'nickname' },
    { title: '邮箱', dataIndex: 'email', copyable: true },
    { title: '手机号', dataIndex: 'phone' },
    {
      title: '状态',
      dataIndex: 'status',
      valueEnum: statusValueEnum,
      render: (_, record) =>
        statusValueEnum[(record.status ?? 0) as 0 | 1]?.text ?? '-',
    },
    {
      title: '最近登录',
      dataIndex: 'lastLoginAt',
      render: (_, record) => formatTimestamp(record.lastLoginAt),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      render: (_, record) => formatTimestamp(record.createdAt),
    },
    {
      title: '操作',
      valueType: 'option',
      width: 140,
      render: (_, record) => [
        <a key="edit" onClick={() => openEdit(record)}>
          编辑
        </a>,
        <Popconfirm
          key="delete"
          title="确认删除该用户？"
          onConfirm={async () => {
            const res = await deleteUser(record.instanceId);
            if (res.success) {
              message.success('删除成功');
              actionRef.current?.reload();
            }
          }}
        >
          <a style={{ color: '#ff4d4f' }}>删除</a>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<User>
        headerTitle="用户列表"
        actionRef={actionRef}
        rowKey="instanceId"
        // 后端 ListUsers 仅支持分页，暂不提供搜索表单。
        search={false}
        rowSelection={{}}
        tableAlertOptionRender={({ selectedRowKeys, onCleanSelected }) => (
          <Space>
            <Popconfirm
              title={`确认删除选中的 ${selectedRowKeys.length} 个用户？`}
              onConfirm={async () => {
                const res = await deleteUsers(selectedRowKeys as string[]);
                if (res.success) {
                  message.success('批量删除成功');
                  onCleanSelected();
                  actionRef.current?.reload();
                }
              }}
            >
              <a style={{ color: '#ff4d4f' }}>批量删除</a>
            </Popconfirm>
          </Space>
        )}
        toolBarRender={() => [
          <Button
            key="add"
            type="primary"
            icon={<PlusOutlined />}
            onClick={openCreate}
          >
            新增用户
          </Button>,
        ]}
        request={async (params) => {
          const res = await listUsers({
            page: params.current,
            pageSize: params.pageSize,
          });
          return {
            data: res.users,
            total: res.total,
            success: res.success,
          };
        }}
        columns={columns}
      />

      <ModalForm<CreateUserBody & UpdateUserBody>
        title={editing ? '编辑用户' : '新增用户'}
        open={modalOpen}
        onOpenChange={setModalOpen}
        modalProps={{ destroyOnHidden: true }}
        // key 强制在新增/编辑切换时重建表单，保证 initialValues 生效。
        key={editing?.instanceId ?? 'create'}
        initialValues={editing}
        onFinish={async (values) => {
          if (editing) {
            const res = await updateUser(editing.instanceId, {
              email: values.email,
              phone: values.phone,
              nickname: values.nickname,
              password: values.password || undefined,
              status: values.status,
            });
            if (!res.success) return false;
            message.success('更新成功');
          } else {
            const res = await createUser({
              name: values.name,
              email: values.email,
              phone: values.phone,
              username: values.username,
              nickname: values.nickname,
              password: values.password,
            });
            if (!res.success) return false;
            message.success('创建成功');
          }
          actionRef.current?.reload();
          return true;
        }}
      >
        {!editing && (
          <>
            <ProFormText
              name="username"
              label="用户名"
              rules={[{ required: true, message: '请输入用户名' }]}
            />
            <ProFormText
              name="name"
              label="姓名"
              rules={[{ required: true, message: '请输入姓名' }]}
            />
          </>
        )}
        <ProFormText
          name="nickname"
          label="昵称"
          rules={[{ required: true, message: '请输入昵称' }]}
        />
        <ProFormText
          name="email"
          label="邮箱"
          rules={[
            { required: !editing, message: '请输入邮箱' },
            { type: 'email', message: '邮箱格式不正确' },
          ]}
        />
        <ProFormText name="phone" label="手机号" />
        <ProFormText.Password
          name="password"
          label="密码"
          rules={editing ? [] : [{ required: true, message: '请输入密码' }]}
          fieldProps={{
            placeholder: editing ? '留空则不修改密码' : '请输入密码',
          }}
        />
        {editing && (
          <ProFormSelect
            name="status"
            label="状态"
            valueEnum={{ 0: '禁用', 1: '正常' }}
          />
        )}
      </ModalForm>
    </PageContainer>
  );
};

export default UserList;
