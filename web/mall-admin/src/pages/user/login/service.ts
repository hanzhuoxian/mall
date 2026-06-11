import { request } from '@umijs/max';

export type LoginParams = {
  identifier: string;
  password: string;
};

export type LoginToken = {
  access_token: string;
  token_type: string;
  refresh_token?: string;
  expires_at: number;
  created_at: number;
};

export type UserInfo = {
  instanceId: string;
  name: string;
  email: string;
  phone: string;
  username: string;
  nickname: string;
  status: number;
};

type ApiResponse<T> = {
  success: boolean;
  data: T;
  errorCode?: number;
  errorMessage?: string;
};

export async function login(params: LoginParams) {
  return request<ApiResponse<LoginToken>>('/login', {
    method: 'POST',
    data: params,
  });
}

export async function getCurrentUser() {
  return request<ApiResponse<{ user: UserInfo }>>('/v1/users/me', {
    method: 'GET',
    skipErrorHandler: true,
  });
}
