import { request } from '@umijs/max';
export type Shop = {
  id: string;
  name: string;
  status: string;
  reason: string;
  contact: string;
  license: string;
  commissionBps: number;
};
export type Context = {
  userId: string;
  admin: boolean;
  shops: Shop[];
  sandboxPayment: boolean;
};
export type Product = {
  id: string;
  shopId: string;
  categoryId: string;
  name: string;
  description: string;
  image: string;
  price: number;
  stock: number;
  status: string;
  reason: string;
};
export type Address = {
  id: string;
  recipient: string;
  phone: string;
  detail: string;
};
export type Cart = {
  id: string;
  productId: string;
  quantity: number;
  product?: Product;
};
export type Order = {
  id: string;
  shopId: string;
  userId: string;
  status: string;
  amount: number;
  commission: number;
  recipient: string;
  phone: string;
  address: string;
  carrier: string;
  tracking: string;
  refundReason: string;
  createdAt: string;
  items: { id: string; name: string; price: number; quantity: number }[];
};
export type Category = { id: string; name: string };
export type Page<T> = { items: T[]; total: number };
export async function api<T>(
  path: string,
  method = 'GET',
  data?: unknown,
  params?: Record<string, unknown>,
): Promise<T> {
  const result = await request<{
    success: boolean;
    data: T;
    errorMessage?: string;
  }>(`/v1/market${path}`, { method, data, params });
  if (!result.success) throw new Error(result.errorMessage || '操作失败');
  return result.data;
}
export const money = (cents: number) => `¥${(cents / 100).toFixed(2)}`;
export const labels: Record<string, string> = {
  pending: '待审核',
  approved: '已通过',
  rejected: '已驳回',
  suspended: '已冻结',
  offline: '已下架',
  paid: '待发货',
  shipped: '待收货',
  completed: '已完成',
  cancelled: '已取消',
  refunding: '售后处理中',
  refunded: '已退款',
};
