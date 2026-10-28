import { afterEach, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import CartPage from './index';
import { marketplace } from '../../services/marketplace';

afterEach(() => vi.unstubAllGlobals());

it('updates quantity, removes items, and stores selected cart IDs for checkout', async () => {
  const fetcher = vi.fn()
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ id: 'c1', productId: 'p1', quantity: 1, product: { id: 'p1', name: '台灯', price: 1200, stock: 8 } }], total: 1 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: {} }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: {} }), { status: 200 }));
  vi.stubGlobal('fetch', fetcher);
  render(<CartPage />);
  expect(await screen.findByText('台灯')).toBeTruthy();
  fireEvent.change(screen.getByRole('spinbutton', { name: '商品数量' }), { target: { value: '2' } });
  await waitFor(() => expect(fetcher).toHaveBeenCalledTimes(2));
  expect(JSON.parse(fetcher.mock.calls[1][1].body)).toMatchObject({ productId: 'p1', quantity: 2 });
  fireEvent.click(screen.getByRole('checkbox'));
  fireEvent.click(screen.getByRole('button', { name: /去结算/ }));
  expect(sessionStorage.getItem('mall.checkout.cartIds')).toBe('["c1"]');
});

it('loads every cart page instead of silently dropping items past the API default', async () => {
  const fetcher = vi.fn()
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ id: 'c1', productId: 'p1', quantity: 1 }], total: 101 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ id: 'c101', productId: 'p101', quantity: 1 }], total: 101 } }), { status: 200 }));
  vi.stubGlobal('fetch', fetcher);
  await expect(marketplace.cart()).resolves.toMatchObject({ total: 101, items: [{ id: 'c1' }, { id: 'c101' }] });
  expect(fetcher).toHaveBeenCalledTimes(2);
});

it('loads later address and order pages as well', async () => {
  const fetcher = vi.fn()
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ id: 'a1' }], total: 101 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ id: 'a101' }], total: 101 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ id: 'o1' }], total: 101 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ id: 'o101' }], total: 101 } }), { status: 200 }));
  vi.stubGlobal('fetch', fetcher);
  await expect(marketplace.addresses()).resolves.toMatchObject({ total: 101, items: [{ id: 'a1' }, { id: 'a101' }] });
  await expect(marketplace.orders()).resolves.toMatchObject({ total: 101, items: [{ id: 'o1' }, { id: 'o101' }] });
  expect(fetcher.mock.calls.map(([url]) => String(url))).toEqual(expect.arrayContaining([
    expect.stringContaining('/addresses?current=1&pageSize=100'),
    expect.stringContaining('/addresses?current=2&pageSize=100'),
    expect.stringContaining('/orders?scope=buyer&current=1&pageSize=100'),
    expect.stringContaining('/orders?scope=buyer&current=2&pageSize=100'),
  ]));
});
