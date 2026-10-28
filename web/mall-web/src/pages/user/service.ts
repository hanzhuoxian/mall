export const ACCESS_TOKEN = 'mall.access_token';
export const REFRESH_TOKEN = 'mall.refresh_token';

type Envelope<T> = { code: number; message?: string; data: T };
export type ApiOptions = RequestInit & { auth?: boolean; retry?: boolean };

export class ApiError extends Error {
  constructor(message: string, readonly status: number, readonly code?: number) {
    super(message);
    this.name = 'ApiError';
  }
}

export function clearSession() {
  localStorage.removeItem(ACCESS_TOKEN);
  localStorage.removeItem(REFRESH_TOKEN);
}

function redirectToLogin() {
  if (typeof window === 'undefined' || window.location.pathname.startsWith('/user/')) return;
  const redirect = window.location.pathname + window.location.search + window.location.hash;
  window.location.assign(`/user/login?redirect=${encodeURIComponent(redirect)}`);
}

async function readEnvelope<T>(response: Response): Promise<T> {
  let body: Envelope<T>;
  try {
    body = await response.json() as Envelope<T>;
  } catch {
    throw new ApiError('服务器响应格式错误', response.status);
  }
  if (!response.ok || body.code !== 0) {
    throw new ApiError(body.message || '请求失败，请稍后重试', response.status, body.code);
  }
  return body.data;
}

let refreshInFlight: Promise<boolean> | undefined;

async function performRefresh(): Promise<boolean> {
  const refreshToken = localStorage.getItem(REFRESH_TOKEN);
  if (!refreshToken) return false;
  try {
    const response = await fetch('/refresh', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ refresh_token: refreshToken }),
    });
    if (!response.ok) return false;
    const token = await readEnvelope<TokenPair>(response);
    storeTokens(token);
    return true;
  } catch {
    return false;
  }
}

export function refreshSession(): Promise<boolean> {
  if (!refreshInFlight) {
    refreshInFlight = performRefresh().finally(() => { refreshInFlight = undefined; });
  }
  return refreshInFlight;
}

export async function api<T>(path: string, options: ApiOptions = {}): Promise<T> {
  const { auth = true, retry = true, headers: inputHeaders, ...init } = options;
  const headers = new Headers(inputHeaders);
  if (init.body && !headers.has('Content-Type')) headers.set('Content-Type', 'application/json');
  if (auth) {
    const token = localStorage.getItem(ACCESS_TOKEN);
    if (token) headers.set('Authorization', `Bearer ${token}`);
  }
  const response = await fetch(path, { ...init, headers });
  if (response.status === 401 && auth) {
    if (retry && await refreshSession()) return api<T>(path, { ...options, retry: false });
    clearSession();
    redirectToLogin();
  }
  return readEnvelope<T>(response);
}

export type TokenPair = {
  access_token: string;
  refresh_token?: string;
  token_type?: string;
  expires_at?: number;
  created_at?: number;
};

function storeTokens(token: TokenPair) {
  if (token.access_token) localStorage.setItem(ACCESS_TOKEN, token.access_token);
  if (token.refresh_token) localStorage.setItem(REFRESH_TOKEN, token.refresh_token);
}

export type LoginCredentials = {
  identifier: string;
  password: string;
  captcha_id?: string;
  captcha_code?: string;
};

export async function login(credentials: LoginCredentials): Promise<TokenPair> {
  const token = await api<TokenPair>('/login', {
    method: 'POST', auth: false, body: JSON.stringify(credentials),
  });
  storeTokens(token);
  return token;
}

export type RegistrationInput = {
  username: string;
  email: string;
  password: string;
  captchaId: string;
  captchaCode: string;
};

export async function register(input: RegistrationInput): Promise<{ registered: boolean }> {
  return api('/register', { method: 'POST', auth: false, body: JSON.stringify(input) });
}

export async function logout(): Promise<void> {
  try {
    await api('/logout', {
      method: 'POST', retry: false,
      body: JSON.stringify({ refresh_token: localStorage.getItem(REFRESH_TOKEN) ?? '' }),
    });
  } finally {
    clearSession();
  }
}

export function safeRedirect(value: string | null | undefined): string {
  if (!value?.startsWith('/') || value.startsWith('//') || value.includes('\\')) return '/';
  try {
    const parsed = new URL(value, window.location.origin);
    return parsed.origin === window.location.origin ? parsed.pathname + parsed.search + parsed.hash : '/';
  } catch {
    return '/';
  }
}

export function loginUrl(redirect: string): string {
  return `/user/login?redirect=${encodeURIComponent(safeRedirect(redirect))}`;
}

export function registerUrl(redirect: string): string {
  return `/user/register?redirect=${encodeURIComponent(safeRedirect(redirect))}`;
}

export type Captcha = { id: string; image: string };
export const getCaptcha = () => api<Captcha>('/captcha', { auth: false });
