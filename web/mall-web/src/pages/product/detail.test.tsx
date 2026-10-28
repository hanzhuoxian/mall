import { afterEach, expect, it, vi } from 'vitest';
import { fireEvent, render, screen } from '@testing-library/react';
import ProductPage from './detail';

afterEach(() => vi.unstubAllGlobals());

it('loads a direct product URL and adds it to cart for a signed-in shopper', async () => {
  const fetcher = vi.fn()
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { id: 'p-2', name: '陶杯', description: '手作陶杯', price: 6800, stock: 5, image: '' } }), { status: 200 }))
    .mockResolvedValueOnce(new Response(JSON.stringify({ code: 0, data: { id: 'cart-1' } }), { status: 200 }));
  vi.stubGlobal('fetch', fetcher);
  vi.stubGlobal('localStorage', { getItem: () => 'token', setItem: () => {}, removeItem: () => {} });
  window.history.replaceState({}, '', '/product/p-2');
  render(<ProductPage />);
  expect(await screen.findByText('手作陶杯')).toBeTruthy();
  fireEvent.click(screen.getByRole('button', { name: '加入购物车' }));
  expect(await screen.findByText(/已加入购物车/)).toBeTruthy();
});

it('shows an unavailable state when a deep-linked product is hidden', async () => {
  vi.stubGlobal('fetch', vi.fn().mockResolvedValue(new Response(JSON.stringify({ code: 404, message: 'record not found' }), { status: 404 })));
  window.history.replaceState({}, '', '/product/hidden-product');
  render(<ProductPage />);
  expect(await screen.findByText(/暂时无法查看/)).toBeTruthy();
});

it('builds a safe login return URL for the product page', async () => {
  const { loginUrl } = await import('../user/service');
  expect(loginUrl('/product/p-2?variant=blue')).toBe('/user/login?redirect=%2Fproduct%2Fp-2%3Fvariant%3Dblue');
});
