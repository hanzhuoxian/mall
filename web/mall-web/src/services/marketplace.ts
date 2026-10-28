import { api } from '../pages/user/service';

export type Page<T> = { items: T[]; total: number };
export type Category = { id: string; name: string };
export type Product = {
  id: string; shopId: string; categoryId: string; name: string; description: string;
  image: string; status: string; price: number; stock: number;
};
export type CartItem = { id: string; productId: string; quantity: number; product?: Product };
export type Address = { id: string; recipient: string; phone: string; detail: string };
export type OrderItem = { id: string; name: string; price: number; quantity: number };
export type Order = {
  id: string; shopId: string; status: string; amount: number; address: string;
  items: OrderItem[]; expiresAt: string; refundReason?: string;
};
export type MarketplaceContext = { userId: string; sandboxPayment: boolean };

const query = (params: Record<string, string | number | undefined>) => {
  const search = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => { if (value !== undefined && value !== '') search.set(key, String(value)); });
  return search.toString();
};

async function allPages<T>(path: string): Promise<Page<T>> {
  const pageSize = 100;
  const separator = path.includes('?') ? '&' : '?';
  const load = (current: number) => api<Page<T>>(`${path}${separator}${query({ current, pageSize })}`);
  const first = await load(1);
  const pages: Page<T>[] = [first];
  for (let start = 2; start <= Math.ceil(first.total / pageSize); start += 5) {
    const batch = await Promise.all(Array.from({ length: Math.min(5, Math.ceil(first.total / pageSize) - start + 1) }, (_, index) => load(start + index)));
    pages.push(...batch);
  }
  return { items: pages.flatMap((page) => page.items), total: first.total };
}

export const marketplace = {
  catalog: (params: { name?: string; categoryId?: string; current?: number; pageSize?: number } = {}) =>
    api<Page<Product>>(`/v1/market/catalog?${query(params)}`, { auth: false }),
  categories: (params: { current?: number; pageSize?: number } = {}) =>
    api<Page<Category>>(`/v1/market/categories?${query(params)}`, { auth: false }),
  product: (id: string) => api<Product>(`/v1/market/catalog/${encodeURIComponent(id)}`, { auth: false }),
  putCart: (productId: string, quantity: number) => api<CartItem>('/v1/market/cart', { method: 'PUT', body: JSON.stringify({ productId, quantity }) }),
  cart: () => allPages<CartItem>('/v1/market/cart'),
  updateCart: (item: CartItem) => marketplace.putCart(item.productId, item.quantity),
  deleteCart: (id: string) => api<null>(`/v1/market/cart/${encodeURIComponent(id)}`, { method: 'DELETE' }),
  addresses: () => allPages<Address>('/v1/market/addresses'),
  createAddress: (input: Omit<Address, 'id'>) => api<Address>('/v1/market/addresses', { method: 'POST', body: JSON.stringify(input) }),
  deleteAddress: (id: string) => api<null>(`/v1/market/addresses/${encodeURIComponent(id)}`, { method: 'DELETE' }),
  checkout: (input: { addressId: string; cartIds: string[]; requestKey: string }) => api<Order[]>('/v1/market/checkout', { method: 'POST', body: JSON.stringify(input) }),
  orders: (status?: string) => allPages<Order>(`/v1/market/orders?scope=buyer${status ? `&status=${encodeURIComponent(status)}` : ''}`),
  context: () => api<MarketplaceContext>('/v1/market/context'),
  orderAction: (id: string, action: string, reason?: string) => api<null>(`/v1/market/orders/${encodeURIComponent(id)}/actions`, { method: 'POST', body: JSON.stringify({ action, reason }) }),
};

export function checkoutRequestKey(signature: string): string {
  const storageKey = 'mall.checkout.pending';
  const stored = sessionStorage.getItem(storageKey);
  if (stored) {
    try {
      const previous = JSON.parse(stored) as { signature: string; key: string };
      if (previous.signature === signature) return previous.key;
    } catch { /* discard malformed pending checkout */ }
  }
  const key = globalThis.crypto?.randomUUID?.() ?? `${Date.now()}-${Math.random().toString(36).slice(2)}`;
  sessionStorage.setItem(storageKey, JSON.stringify({ signature, key }));
  return key;
}

export function checkoutSignature(userId: string, addressId: string, cartIds: string[]): string {
  return `${userId}|${addressId}|${[...cartIds].sort().join(',')}`;
}

export type PendingCheckout = {
  userId: string;
  signature: string;
  key: string;
  addressId: string;
  cartIds: string[];
  address: Address;
  items: CartItem[];
};

export function savePendingCheckout(input: Omit<PendingCheckout, 'signature' | 'key'>): PendingCheckout {
  const cartIds = [...input.cartIds].sort();
  const signature = checkoutSignature(input.userId, input.addressId, cartIds);
  const attempt = { ...input, cartIds, signature, key: checkoutRequestKey(signature) };
  sessionStorage.setItem('mall.checkout.pending', JSON.stringify(attempt));
  return attempt;
}

export function readPendingCheckout(): PendingCheckout | undefined {
  try {
    const value = sessionStorage.getItem('mall.checkout.pending');
    if (!value) return undefined;
    const attempt = JSON.parse(value) as PendingCheckout;
    if (!attempt.userId || !attempt.addressId || !attempt.key || !Array.isArray(attempt.cartIds) || !Array.isArray(attempt.items)) return undefined;
    return attempt;
  } catch { return undefined; }
}
