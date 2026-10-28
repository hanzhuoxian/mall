import { afterEach, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import AddressesPage from './addresses';

afterEach(() => vi.unstubAllGlobals());

it('creates and deletes a shipping address', async () => {
  const fetcher = vi.fn()
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [], total: 0 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { id: 'a1', recipient: '林女士', phone: '13800000000', detail: '上海市' } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: null }), { status: 200 }));
  vi.stubGlobal('fetch', fetcher);
  render(<AddressesPage />);
  fireEvent.change(await screen.findByLabelText('收件人'), { target: { value: '林女士' } });
  fireEvent.change(screen.getByLabelText('联系电话'), { target: { value: '13800000000' } });
  fireEvent.change(screen.getByLabelText('详细地址'), { target: { value: '上海市' } });
  fireEvent.click(screen.getByRole('button', { name: '保存地址' }));
  expect(await screen.findByText('林女士')).toBeTruthy();
  fireEvent.click(screen.getByRole('button', { name: '删除地址' }));
  await waitFor(() => expect(fetcher).toHaveBeenCalledTimes(3));
});
