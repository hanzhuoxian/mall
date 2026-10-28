import React from 'react';
import {
  cleanup,
  fireEvent,
  render,
  screen,
  waitFor,
} from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
const mocks = vi.hoisted(() => ({ path: '/market/orders', api: vi.fn() }));
vi.mock('@umijs/max', () => ({
  useLocation: () => ({ pathname: mocks.path }),
}));
vi.mock('@ant-design/pro-components', () => ({
  PageContainer: ({ children }: any) => <div>{children}</div>,
}));
vi.mock('./api', async (importOriginal) => ({
  ...(await importOriginal<object>()),
  api: mocks.api,
}));
import Marketplace from './index';
const context = {
  userId: 'buyer',
  admin: false,
  shops: [],
  sandboxPayment: false,
};
beforeEach(() => {
  mocks.path = '/market/orders';
  mocks.api.mockReset();
  mocks.api.mockImplementation(async (path: string) => {
    if (path === '/context') return context;
    if (path === '/categories') return { items: [], total: 0 };
    return {
      items: [
        {
          id: 'order-one',
          status: 'pending',
          amount: 1990,
          createdAt: '2026-09-16T00:00:00Z',
          items: [],
        },
      ],
      total: 1,
    };
  });
});
afterEach(cleanup);
describe('Marketplace workflows', () => {
  it('disables payment when sandbox is off and shows server money', async () => {
    render(<Marketplace />);
    expect(await screen.findByText('¥19.90')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '沙箱付款' })).toBeDisabled();
  });
  it('submits a confirmed cancellation with the order id', async () => {
    render(<Marketplace />);
    fireEvent.click(await screen.findByRole('button', { name: /取\s*消/ }));
    fireEvent.click(screen.getByRole('button', { name: /OK|确 定|确定/ }));
    await waitFor(() =>
      expect(mocks.api).toHaveBeenCalledWith(
        '/orders/order-one/actions',
        'POST',
        { action: 'cancel' },
      ),
    );
  });
  it('does not display fake data after a failed request', async () => {
    mocks.api.mockRejectedValue(new Error('offline'));
    render(<Marketplace />);
    expect(await screen.findByText('数据加载失败')).toBeInTheDocument();
    expect(screen.queryByText('¥19.90')).not.toBeInTheDocument();
  });
  it('approves a shop with the entered commission rate', async () => {
    mocks.path = '/market/shops';
    mocks.api.mockImplementation(async (path: string) => {
      if (path === '/context')
        return { ...context, userId: 'admin', admin: true };
      if (path === '/shops')
        return {
          items: [
            {
              id: 'shop-one',
              name: '测试店铺',
              status: 'pending',
              reason: '',
              contact: '张三 13800000000',
              license: 'LICENSE',
              commissionBps: 500,
            },
          ],
          total: 1,
        };
      return { items: [], total: 0 };
    });
    render(<Marketplace />);
    fireEvent.click(await screen.findByRole('button', { name: /通\s*过/ }));
    fireEvent.change(screen.getByRole('spinbutton'), {
      target: { value: '600' },
    });
    fireEvent.click(screen.getByRole('button', { name: /OK|确 定|确定/ }));
    await waitFor(() =>
      expect(mocks.api).toHaveBeenCalledWith('/shops/shop-one/review', 'POST', {
        status: 'approved',
        commissionBps: 600,
      }),
    );
  });
  it('does not offer merchant approval to a buyer', async () => {
    mocks.path = '/market/shops';
    mocks.api.mockImplementation(async (path: string) =>
      path === '/context' ? context : { items: [], total: 0 },
    );
    render(<Marketplace />);
    await screen.findByText('申请入驻／重新提交');
    expect(
      screen.queryByRole('button', { name: '通过' }),
    ).not.toBeInTheDocument();
  });

});
