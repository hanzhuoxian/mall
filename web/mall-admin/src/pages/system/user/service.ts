import { request } from '@umijs/max';

// 后端用 protojson 序列化 proto 消息：字段为 camelCase，时间戳为 RFC3339 字符串。
export type User = {
  instanceId: string;
  name: string;
  email: string;
  phone: string;
  username: string;
  nickname: string;
  status?: number;
  lastLoginAt?: string;
  createdAt?: string;
  updatedAt?: string;
};

type ApiResponse<T> = {
  success: boolean;
  data: T;
  errorCode?: number;
  errorMessage?: string;
};

export type ListUsersParams = {
  page?: number;
  pageSize?: number;
};

export type ListUsersResult = {
  success: boolean;
  users: User[];
  total: number;
};

export type CreateUserBody = {
  name: string;
  email: string;
  phone?: string;
  username: string;
  nickname: string;
  password: string;
};

// UpdateUserBody 对应 proto UpdateUserRequest，仅这些字段可更新。
export type UpdateUserBody = {
  email?: string;
  phone?: string;
  nickname?: string;
  password?: string;
  status?: number;
};

// GET /v1/users
export async function listUsers(
  params: ListUsersParams,
): Promise<ListUsersResult> {
  const res = await request<ApiResponse<{ users?: User[]; total?: number }>>(
    '/v1/users',
    { method: 'GET', params },
  );
  return {
    success: res.success,
    users: res.data?.users ?? [],
    total: Number(res.data?.total ?? 0),
  };
}

// POST /v1/users
export async function createUser(body: CreateUserBody) {
  const res = await request<ApiResponse<{ user: User }>>('/v1/users', {
    method: 'POST',
    data: body,
  });
  return {
    success: res.success,
    user: res.data?.user,
  };
}

// PUT /v1/users/:id
export async function updateUser(id: string, body: UpdateUserBody) {
  const res = await request<ApiResponse<{ user: User }>>(`/v1/users/${id}`, {
    method: 'PUT',
    data: body,
  });
  return {
    success: res.success,
    user: res.data?.user,
  };
}

// DELETE /v1/users/:id
export async function deleteUser(id: string) {
  return request<ApiResponse<unknown>>(`/v1/users/${id}`, {
    method: 'DELETE',
  });
}

// DELETE /v1/users  批量删除
export async function deleteUsers(instanceIds: string[]) {
  return request<ApiResponse<unknown>>('/v1/users', {
    method: 'DELETE',
    data: { instanceIds },
  });
}
