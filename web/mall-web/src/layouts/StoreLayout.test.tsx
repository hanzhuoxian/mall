import { afterEach, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import StoreLayout from './StoreLayout';
import { ACCESS_TOKEN, REFRESH_TOKEN } from '../pages/user/service';

afterEach(() => vi.unstubAllGlobals());

it('exposes address management and sign out for an authenticated shopper', async () => {
  localStorage.setItem(ACCESS_TOKEN, 'access');
  localStorage.setItem(REFRESH_TOKEN, 'refresh');
  const fetcher = vi.fn().mockResolvedValue(new Response(JSON.stringify({ code: 0, data: null }), { status: 200 }));
  vi.stubGlobal('fetch', fetcher);
  render(
    <MemoryRouter>
      <Routes>
        <Route element={<StoreLayout />}>
          <Route index element={<div>商品目录</div>} />
        </Route>
      </Routes>
    </MemoryRouter>,
  );
  expect(screen.getByText('商品目录')).toBeTruthy();
  expect(screen.getByLabelText('全站搜索')).toBeTruthy();
  expect(screen.getByRole('link', { name: '账户' }).getAttribute('href')).toBe('/account/addresses');
  fireEvent.click(screen.getByRole('button', { name: '退出登录' }));
  await waitFor(() => expect(fetcher).toHaveBeenCalledTimes(1));
  expect(JSON.parse(fetcher.mock.calls[0][1].body)).toEqual({ refresh_token: 'refresh' });
  expect(await screen.findByRole('link', { name: '账户' })).toHaveProperty('href', expect.stringContaining('/user/login'));
});
