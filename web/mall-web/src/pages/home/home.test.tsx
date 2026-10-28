import { afterEach, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import HomePage from './index';

afterEach(() => vi.unstubAllGlobals());

it('loads public categories and products and filters by name', async () => {
  const fetcher = vi.fn()
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ id: 'cat-1', name: '家居' }], total: 1 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ id: 'p-1', name: '木质台灯', price: 12900, stock: 4, categoryId: 'cat-1', image: '' }], total: 25 } }), { status: 200 }))
    .mockResolvedValue(new Response(JSON.stringify({ code: 0, data: { items: [], total: 25 } }), { status: 200 }));
  vi.stubGlobal('fetch', fetcher);
  render(<HomePage />);
  expect(await screen.findByText('木质台灯')).toBeTruthy();
  expect(screen.getByText('家居')).toBeTruthy();
  fireEvent.change(screen.getByLabelText('搜索商品'), { target: { value: '不存在' } });
  fireEvent.click(screen.getByRole('button', { name: /搜\s*索/ }));
  await waitFor(() => expect(fetcher).toHaveBeenCalledTimes(3));
  expect(fetcher.mock.calls[2][0]).toContain('name=%E4%B8%8D%E5%AD%98%E5%9C%A8');
  fireEvent.click(screen.getByText('家居'));
  await waitFor(() => expect(fetcher).toHaveBeenCalledTimes(4));
  expect(fetcher.mock.calls[3][0]).toContain('categoryId=cat-1');
  fireEvent.click(screen.getByTitle('2'));
  await waitFor(() => expect(fetcher).toHaveBeenCalledTimes(5));
  expect(fetcher.mock.calls[4][0]).toContain('current=2');
});
