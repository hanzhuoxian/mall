import { afterEach, beforeEach, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import CheckoutPage from './index';
import { checkoutRequestKey, checkoutSignature } from '../../services/marketplace';

afterEach(() => vi.unstubAllGlobals());
beforeEach(() => sessionStorage.clear());

it('reuses checkout key only while address and cart ID set are unchanged', () => {
  const values = new Map<string, string>();
  vi.stubGlobal('sessionStorage', { getItem: (key: string) => values.get(key) ?? null, setItem: (key: string, value: string) => values.set(key, value) });
  const first = checkoutRequestKey(checkoutSignature('u1', 'a1', ['c1', 'c2']));
  expect(checkoutRequestKey(checkoutSignature('u1', 'a1', ['c1', 'c2']))).toBe(first);
  expect(checkoutRequestKey(checkoutSignature('u1', 'a2', ['c1', 'c2']))).not.toBe(first);
  expect(checkoutRequestKey(checkoutSignature('u1', 'a2', ['c2', 'c1']))).toBe(checkoutRequestKey(checkoutSignature('u1', 'a2', ['c1', 'c2'])));
});

it('submits selected cart IDs and displays the server returned order total', async () => {
  sessionStorage.setItem('mall.checkout.cartIds', '["c1"]');
  const fetcher = vi.fn()
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ id: 'c1', productId: 'p1', quantity: 2, product: { id: 'p1', name: '台灯', price: 1200, stock: 4 } }], total: 1 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ id: 'a1', recipient: '林女士', phone: '13800000000', detail: '上海市' }], total: 1 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { userId: 'u1', sandboxPayment: true } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: [{ id: 'o1', amount: 2400, status: 'pending' }] }), { status: 200 }));
  vi.stubGlobal('fetch', fetcher);
  render(<CheckoutPage />);
  expect(await screen.findByRole('link', { name: '管理地址' })).toBeTruthy();
  fireEvent.click(await screen.findByRole('button', { name: /确认并提交订单/ }));
  await waitFor(() => expect(fetcher).toHaveBeenCalledTimes(4));
  expect(JSON.parse(fetcher.mock.calls[3][1].body)).toMatchObject({ addressId: 'a1', cartIds: ['c1'] });
  expect(await screen.findByText('¥24.00')).toBeTruthy();
});

it('keeps selected cart IDs available after checkout validation fails', async () => {
  sessionStorage.setItem('mall.checkout.cartIds', '["c2"]');
  const fetcher = vi.fn()
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ id: 'c2', productId: 'p2', quantity: 1, product: { id: 'p2', name: '杯子', price: 900, stock: 2 } }], total: 1 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ id: 'a2', recipient: '林女士', phone: '138', detail: '上海' }], total: 1 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { userId: 'u1', sandboxPayment: false } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 400, message: '商品库存不足' }), { status: 400 }));
  vi.stubGlobal('fetch', fetcher);
  render(<CheckoutPage />);
  fireEvent.click(await screen.findByRole('button', { name: /确认并提交订单/ }));
  expect(await screen.findByText('商品库存不足')).toBeTruthy();
  expect(sessionStorage.getItem('mall.checkout.cartIds')).toBe('["c2"]');
  expect(screen.getByText(/杯子/)).toBeTruthy();
});

it('retries a committed checkout after reload with the exact original payload', async () => {
  sessionStorage.setItem('mall.checkout.cartIds', '["c1"]');
  const fetcher = vi.fn()
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ id: 'c1', productId: 'p1', quantity: 1, product: { id: 'p1', name: '台灯', price: 1200, stock: 3 } }], total: 1 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ id: 'a1', recipient: '林女士', phone: '138', detail: '上海' }, { id: 'a2', recipient: '王先生', phone: '139', detail: '杭州' }], total: 2 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { userId: 'u1', sandboxPayment: false } }), { status: 200 }))
    .mockRejectedValueOnce(new TypeError('network timeout'))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [], total: 0 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ id: 'a1', recipient: '林女士', phone: '138', detail: '上海' }, { id: 'a2', recipient: '王先生', phone: '139', detail: '杭州' }], total: 2 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { userId: 'u1', sandboxPayment: false } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: [{ id: 'o1', amount: 1200, status: 'pending' }] }), { status: 200 }));
  vi.stubGlobal('fetch', fetcher);
  const firstView = render(<CheckoutPage />);
  fireEvent.click(await screen.findByRole('radio', { name: /王先生/ }));
  fireEvent.click(await screen.findByRole('button', { name: /确认并提交订单/ }));
  await waitFor(() => expect(fetcher).toHaveBeenCalledTimes(4));
  const firstPayload = fetcher.mock.calls[3][1].body;
  firstView.unmount();

  render(<CheckoutPage />);
  expect(await screen.findByText(/台灯/)).toBeTruthy();
  fireEvent.click(screen.getByRole('button', { name: /确认并提交订单/ }));
  await waitFor(() => expect(fetcher).toHaveBeenCalledTimes(8));
  expect(fetcher.mock.calls[7][1].body).toBe(firstPayload);
  expect(JSON.parse(fetcher.mock.calls[7][1].body).addressId).toBe('a2');
  expect(await screen.findByText('订单已提交')).toBeTruthy();
});
