import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { api, login, logout, register, ACCESS_TOKEN, REFRESH_TOKEN } from './service';
import LoginPage from './login';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';

const storage = new Map<string, string>();
beforeEach(() => {
  storage.clear();
  vi.stubGlobal('localStorage', {
    getItem: (key: string) => storage.get(key) ?? null,
    setItem: (key: string, value: string) => storage.set(key, value),
    removeItem: (key: string) => storage.delete(key),
    clear: () => storage.clear(),
  });
});
afterEach(() => { vi.unstubAllGlobals(); vi.restoreAllMocks(); });

describe('consumer account session', () => {
  it('parses the raw API envelope', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(JSON.stringify({ code: 0, message: 'ok', data: { value: 7 } }), { status: 200 })));
    await expect(api<{ value: number }>('/v1/example')).resolves.toEqual({ value: 7 });
  });

  it('sends the stored bearer access token', async () => {
    localStorage.setItem(ACCESS_TOKEN, 'access-123');
    const fetcher = vi.fn().mockResolvedValue(new Response(JSON.stringify({ code: 0, data: true }), { status: 200 }));
    vi.stubGlobal('fetch', fetcher);
    await api('/v1/private');
    expect(fetcher.mock.calls[0][1].headers.get('Authorization')).toBe('Bearer access-123');
  });

  it('refreshes once then retries the original request', async () => {
    localStorage.setItem(ACCESS_TOKEN, 'expired'); localStorage.setItem(REFRESH_TOKEN, 'refresh-1');
    const fetcher = vi.fn()
      .mockResolvedValueOnce(new Response('{}', { status: 401 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { access_token: 'fresh', refresh_token: 'refresh-2' } }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: 'ok' }), { status: 200 }));
    vi.stubGlobal('fetch', fetcher);
    await expect(api('/v1/private')).resolves.toBe('ok');
    expect(fetcher).toHaveBeenCalledTimes(3);
    expect(fetcher.mock.calls[1][1].body).toBe(JSON.stringify({ refresh_token: 'refresh-1' }));
    expect(localStorage.getItem(ACCESS_TOKEN)).toBe('fresh');
  });

  it('clears both tokens when refresh fails', async () => {
    localStorage.setItem(ACCESS_TOKEN, 'expired'); localStorage.setItem(REFRESH_TOKEN, 'bad');
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response('{}', { status: 401 })));
    await expect(api('/v1/private')).rejects.toThrow();
    expect(localStorage.getItem(ACCESS_TOKEN)).toBeNull();
    expect(localStorage.getItem(REFRESH_TOKEN)).toBeNull();
  });

  it('clears tokens after the retried request is still unauthorized', async () => {
    localStorage.setItem(ACCESS_TOKEN, 'expired'); localStorage.setItem(REFRESH_TOKEN, 'refresh-1');
    const fetcher = vi.fn()
      .mockResolvedValueOnce(new Response('{}', { status: 401 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { access_token: 'fresh', refresh_token: 'refresh-2' } }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ code: 100200, message: 'Token invalid' }), { status: 401 }));
    vi.stubGlobal('fetch', fetcher);
    await expect(api('/v1/private')).rejects.toThrow();
    expect(localStorage.getItem(ACCESS_TOKEN)).toBeNull();
    expect(localStorage.getItem(REFRESH_TOKEN)).toBeNull();
  });

  it('clears tokens if the refresh request cannot be reached', async () => {
    localStorage.setItem(ACCESS_TOKEN, 'expired'); localStorage.setItem(REFRESH_TOKEN, 'refresh-1');
    const fetcher = vi.fn().mockResolvedValueOnce(new Response('{}', { status: 401 })).mockRejectedValueOnce(new TypeError('network unavailable'));
    vi.stubGlobal('fetch', fetcher);
    await expect(api('/v1/private')).rejects.toThrow();
    expect(localStorage.getItem(ACCESS_TOKEN)).toBeNull();
    expect(localStorage.getItem(REFRESH_TOKEN)).toBeNull();
  });

  it('shares one rotating refresh request between concurrent API calls', async () => {
    localStorage.setItem(ACCESS_TOKEN, 'expired'); localStorage.setItem(REFRESH_TOKEN, 'refresh-1');
    let finishRefresh!: (response: Response) => void;
    const refreshResponse = new Promise<Response>((resolve) => { finishRefresh = resolve; });
    let refreshCalls = 0;
    const fetcher = vi.fn((input: RequestInfo | URL, init?: RequestInit) => {
      const url = String(input);
      if (url === '/refresh') { refreshCalls += 1; return refreshResponse; }
      if (new Headers(init?.headers).get('Authorization') === 'Bearer fresh') return Promise.resolve(new Response(JSON.stringify({ code: 0, data: url }), { status: 200 }));
      return Promise.resolve(new Response('{}', { status: 401 }));
    });
    vi.stubGlobal('fetch', fetcher);
    const pending = Promise.all([api('/v1/one'), api('/v1/two')]);
    await waitFor(() => expect(refreshCalls).toBe(1));
    finishRefresh(new Response(JSON.stringify({ code: 0, data: { access_token: 'fresh', refresh_token: 'refresh-2' } }), { status: 200 }));
    await expect(pending).resolves.toEqual(['/v1/one', '/v1/two']);
    expect(refreshCalls).toBe(1);
  });

  it('logs in and stores token pair', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(JSON.stringify({ code: 0, data: { access_token: 'a', refresh_token: 'r' } }), { status: 200 })));
    await login({ identifier: 'buyer', password: 'secret' });
    expect(localStorage.getItem(ACCESS_TOKEN)).toBe('a');
    expect(localStorage.getItem(REFRESH_TOKEN)).toBe('r');
  });

  it('registers without creating an authenticated session', async () => {
    vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(JSON.stringify({ code: 0, data: { registered: true } }), { status: 200 })));
    await register({ username: 'buyer', email: 'buyer@example.com', password: '123456789012', captchaId: 'id', captchaCode: 'code' });
    expect(localStorage.getItem(ACCESS_TOKEN)).toBeNull();
    expect(localStorage.getItem(REFRESH_TOKEN)).toBeNull();
  });

  it('sends refresh token when logging out and clears both tokens', async () => {
    localStorage.setItem(ACCESS_TOKEN, 'access'); localStorage.setItem(REFRESH_TOKEN, 'refresh');
    const fetcher = vi.fn().mockResolvedValue(new Response(JSON.stringify({ code: 0, data: null }), { status: 200 }));
    vi.stubGlobal('fetch', fetcher);
    await logout();
    expect(JSON.parse(fetcher.mock.calls[0][1].body)).toEqual({ refresh_token: 'refresh' });
    expect(localStorage.getItem(ACCESS_TOKEN)).toBeNull();
    expect(localStorage.getItem(REFRESH_TOKEN)).toBeNull();
  });

  it('rejects unsafe login return paths', async () => {
    const { safeRedirect } = await import('./service');
    expect(safeRedirect('//evil.example')).toBe('/');
    expect(safeRedirect('/cart?x=1')).toBe('/cart?x=1');
  });

  it('keeps the intended destination when continuing from registration to login', async () => {
    const { registerUrl, loginUrl } = await import('./service');
    expect(registerUrl('/cart?from=detail')).toBe('/user/register?redirect=%2Fcart%3Ffrom%3Ddetail');
    expect(loginUrl('/cart?from=detail')).toBe('/user/login?redirect=%2Fcart%3Ffrom%3Ddetail');
  });

  it('shows captcha after two failed passwords when the server uses a generic login error', async () => {
    const fetcher = vi.fn()
      .mockResolvedValueOnce(new Response(JSON.stringify({ code: 100200, message: 'Token invalid' }), { status: 401 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ code: 100200, message: 'Token invalid' }), { status: 401 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { id: 'captcha-id', image: 'data:image/png;base64,AA==' } }), { status: 200 }));
    vi.stubGlobal('fetch', fetcher);
    render(<LoginPage />);
    fireEvent.change(screen.getByLabelText('用户名或邮箱'), { target: { value: 'buyer' } });
    fireEvent.change(screen.getByLabelText('密码'), { target: { value: 'correct-horse' } });
    fireEvent.click(screen.getByRole('button', { name: /登\s*录/ }));
    await waitFor(() => expect(fetcher).toHaveBeenCalledTimes(1));
    fireEvent.click(screen.getByRole('button', { name: /登\s*录/ }));
    expect(await screen.findByLabelText('验证码')).toBeTruthy();
    await waitFor(() => expect(fetcher).toHaveBeenCalledTimes(3));
  });
});
