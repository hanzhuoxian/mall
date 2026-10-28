import { afterEach, expect, it, vi } from 'vitest';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import OrdersPage from './index';

afterEach(() => vi.unstubAllGlobals());

const order = { id: 'order-12345678', shopId: 's1', status: 'pending', amount: 2500, address: '上海市', items: [{ id: 'i1', name: '陶杯', price: 2500, quantity: 1 }] };

it('hides sandbox payment when disabled and keeps clear disclosure when enabled and used', async () => {
  const fetcher = vi.fn()
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { userId: 'u1', sandboxPayment: false } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [order], total: 1 } }), { status: 200 }));
  vi.stubGlobal('fetch', fetcher);
  const { unmount } = render(<OrdersPage />);
  expect(await screen.findByText(/陶杯/)).toBeTruthy();
  expect(screen.queryByRole('button', { name: '开发沙箱付款' })).toBeNull();
  unmount();

  const sandboxFetch = vi.fn()
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { userId: 'u1', sandboxPayment: true } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [order], total: 1 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: null }), { status: 200 }));
  vi.stubGlobal('fetch', sandboxFetch);
  const sandboxView = render(<OrdersPage />);
  const pay = await screen.findByRole('button', { name: '开发沙箱付款' });
  expect(screen.getAllByText(/开发沙箱付款仅供开发演示，不产生真实扣款/).length).toBeGreaterThan(0);
  fireEvent.click(pay);
  fireEvent.click(await screen.findByRole('button', { name: '确认开发沙箱付款' }));
  await waitFor(() => expect(sandboxFetch).toHaveBeenCalledTimes(3));
  expect(JSON.parse(sandboxFetch.mock.calls[2][1].body)).toMatchObject({ action: 'pay' });
  expect(screen.getAllByText(/开发沙箱付款仅供开发演示，不产生真实扣款/).length).toBeGreaterThan(0);
  sandboxView.unmount();
});

it('offers only buyer actions allowed by each order state', async () => {
  for (const [status, visible, hidden] of [
    ['pending', ['取消订单'], ['申请售后', '确认收货']],
    ['paid', ['申请售后'], ['取消订单', '确认收货']],
    ['shipped', ['申请售后', '确认收货'], ['取消订单']],
    ['completed', ['申请售后'], ['取消订单', '确认收货']],
    ['refunding', [], ['申请售后', '取消订单', '确认收货']],
  ] as const) {
    const fetcher = vi.fn()
      .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { userId: 'u1', sandboxPayment: false } }), { status: 200 }))
      .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [{ ...order, status }], total: 1 } }), { status: 200 }));
    vi.stubGlobal('fetch', fetcher);
    const view = render(<OrdersPage />);
    await screen.findByText(/陶杯/);
    for (const label of visible) expect(screen.getByRole('button', { name: label })).toBeTruthy();
    for (const label of hidden) expect(screen.queryByRole('button', { name: label }), `${status}: ${label}`).toBeNull();
    view.unmount(); vi.unstubAllGlobals();
  }
});

it('does not offer sandbox payment after the order expires', async () => {
  const expiredOrder = { ...order, expiresAt: '2020-01-01T00:00:00Z' };
  const fetcher = vi.fn()
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { userId: 'u1', sandboxPayment: true } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [expiredOrder], total: 1 } }), { status: 200 }));
  vi.stubGlobal('fetch', fetcher);
  render(<OrdersPage />);
  await screen.findByText(/陶杯/);
  expect(screen.queryByRole('button', { name: '开发沙箱付款' })).toBeNull();
});

it('keeps the latest order filter when earlier requests finish later', async () => {
  let finishInitial!: (response: Response) => void;
  const initialOrders = new Promise<Response>((resolve) => { finishInitial = resolve; });
  const response = (items: unknown[]) => new Response(JSON.stringify({ code: 0, data: { items, total: items.length } }), { status: 200 });
  const fetcher = vi.fn((input: RequestInfo | URL) => {
    const url = String(input);
    if (url.includes('/context')) return Promise.resolve(new Response(JSON.stringify({ code: 0, data: { userId: 'u1', sandboxPayment: false } }), { status: 200 }));
    if (url.includes('status=pending')) return Promise.resolve(response([{ ...order, id: 'new-order', items: [{ id: 'n', name: '新订单', price: 100, quantity: 1 }] }]));
    return initialOrders;
  });
  vi.stubGlobal('fetch', fetcher);
  render(<OrdersPage />);
  await waitFor(() => expect(fetcher).toHaveBeenCalledTimes(2));
  fireEvent.click(screen.getByRole('button', { name: '筛选待付款' }));
  expect(await screen.findByText(/新订单/)).toBeTruthy();
  finishInitial(response([{ ...order, id: 'old-order', items: [{ id: 'o', name: '旧订单', price: 100, quantity: 1 }] }]));
  await waitFor(() => expect(screen.queryByText(/旧订单/)).toBeNull());
  expect(screen.getByText(/新订单/)).toBeTruthy();
});

it('removes an order from the active filter after its state changes', async () => {
  const fetcher = vi.fn()
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { userId: 'u1', sandboxPayment: false } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [order], total: 1 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { userId: 'u1', sandboxPayment: false } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { items: [order], total: 1 } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: null }), { status: 200 }));
  vi.stubGlobal('fetch', fetcher);
  render(<OrdersPage />);
  await screen.findByText(/陶杯/);
  fireEvent.click(screen.getByRole('button', { name: '筛选待付款' }));
  await waitFor(() => expect(fetcher).toHaveBeenCalledTimes(4));
  fireEvent.click(screen.getByRole('button', { name: '取消订单' }));
  fireEvent.click(await screen.findByRole('button', { name: '确认取消订单' }));
  await waitFor(() => expect(fetcher).toHaveBeenCalledTimes(5));
  expect(await screen.findByText('还没有相关订单')).toBeTruthy();
});
