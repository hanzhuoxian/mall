import { request } from '@umijs/max';

export type LoginParams = {
  identifier: string;
  password: string;
  captcha_id?: string;
  captcha_code?: string;
};

export type CaptchaData = {
  id: string;
  image: string;
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
    skipErrorHandler: true,
  });
}

export async function getCaptcha() {
  return request<ApiResponse<CaptchaData>>('/captcha', {
    method: 'GET',
    skipErrorHandler: true,
  });
}

export async function getCurrentUser() {
  return request<ApiResponse<{ user: UserInfo }>>('/v1/users/me', {
    method: 'GET',
    skipErrorHandler: true,
  });
}
