import { Alert, Button, Empty, Radio, Spin } from 'antd';
import { useEffect, useMemo, useState } from 'react';
import { marketplace, readPendingCheckout, savePendingCheckout, type Address, type CartItem, type Order } from '../../services/marketplace';

const money = (cents: number) => `¥${(cents / 100).toFixed(2)}`;

export default function CheckoutPage() {
  const cartIds = useMemo<string[]>(() => {
    try { const value = JSON.parse(sessionStorage.getItem('mall.checkout.cartIds') ?? '[]'); return Array.isArray(value) ? value.filter((id): id is string => typeof id === 'string') : []; }
    catch { return []; }
  }, []);
  const [items, setItems] = useState<CartItem[]>([]);
  const [addresses, setAddresses] = useState<Address[]>([]);
  const [pendingAddress, setPendingAddress] = useState<Address>();
  const [recovering, setRecovering] = useState(false);
  const [userId, setUserId] = useState('');
  const [addressId, setAddressId] = useState('');
  const [orders, setOrders] = useState<Order[]>([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  useEffect(() => {
    let active = true;
    Promise.all([marketplace.cart(), marketplace.addresses(), marketplace.context()])
      .then(([cart, addressPage, context]) => {
        if (!active) return;
        const pending = readPendingCheckout();
        const sameCart = pending && pending.userId === context.userId && pending.cartIds.length === cartIds.length && pending.cartIds.every((id, index) => [...cartIds].sort()[index] === id);
        const liveItems = cart.items.filter((item) => cartIds.includes(item.id));
        setRecovering(Boolean(sameCart));
        setItems(liveItems.length ? liveItems : sameCart ? pending.items : []);
        setAddresses(addressPage.items); setUserId(context.userId);
        const selectedAddress = sameCart ? addressPage.items.find((item) => item.id === pending.addressId) ?? pending.address : addressPage.items[0];
        setAddressId(selectedAddress?.id ?? '');
        setPendingAddress(sameCart ? selectedAddress : undefined);
      })
      .catch((err) => { if (active) setError(err instanceof Error ? err.message : '结算信息加载失败'); })
      .finally(() => { if (active) setLoading(false); });
    return () => { active = false; };
  }, [cartIds]);
  const submit = async () => {
    const address = addresses.find((item) => item.id === addressId) ?? pendingAddress;
    if (cartIds.length === 0 || cartIds.length > 100 || !addressId || !address || !items.length) {
      setError(cartIds.length > 100 ? '一次最多结算 100 件商品' : '请选择商品和收货地址'); return;
    }
    setSubmitting(true); setError('');
    const attempt = savePendingCheckout({ userId, addressId, cartIds, address, items });
    try {
      const created = await marketplace.checkout({ addressId, cartIds: attempt.cartIds, requestKey: attempt.key });
      setOrders(created);
      sessionStorage.removeItem('mall.checkout.pending');
      sessionStorage.removeItem('mall.checkout.cartIds');
    } catch (err) { setError(err instanceof Error ? err.message : '提交失败，请检查信息后重试'); }
    finally { setSubmitting(false); }
  };
  if (loading) return <div className="page-wrap loading"><Spin /></div>;
  if (orders.length) return <main className="page-wrap"><div className="success-mark">✓</div><div className="eyebrow">Order placed</div><h1>订单已提交</h1><Alert type="success" title="订单金额由服务端确认" showIcon />
    <section className="placed-orders">{orders.map((order) => <article key={order.id}><span>订单 {order.id.slice(0, 8)}</span><strong>{money(order.amount)}</strong></article>)}</section>
    <Button type="primary" href="/orders">查看我的订单</Button>
  </main>;
  return <main className="page-wrap"><div className="eyebrow">Almost yours</div><h1 className="page-title">确认订单</h1>
    {error && <Alert className="checkout-error" type="error" message={error} showIcon />}
    {!cartIds.length || (!items.length && !recovering) ? <Empty description="没有选中的商品，请返回购物车重新选择"><Button href="/cart">返回购物车</Button></Empty> : <div className="checkout-layout">
      <section><div className="checkout-section-title"><h2>收货地址</h2><Button type="link" href="/account/addresses">管理地址</Button></div>{addresses.length || pendingAddress ? <Radio.Group className="address-options" value={addressId} onChange={(event) => setAddressId(event.target.value)}>
        {pendingAddress && !addresses.some((item) => item.id === pendingAddress.id) && <Radio value={pendingAddress.id}><strong>{pendingAddress.recipient}</strong>　{pendingAddress.phone}<div>{pendingAddress.detail}（原结算地址）</div></Radio>}
        {addresses.map((address) => <Radio key={address.id} value={address.id}><strong>{address.recipient}</strong>　{address.phone}<div>{address.detail}</div></Radio>)}
      </Radio.Group> : <Alert type="warning" title="请先添加收货地址" action={<Button href="/account/addresses">管理地址</Button>} />}
      <h2>商品清单</h2>{items.map((item) => <div className="checkout-item" key={item.id}><span>{item.product?.name} × {item.quantity}</span><strong>{money((item.product?.price ?? 0) * item.quantity)}</strong></div>)}
      </section><aside className="checkout-summary"><span>商品金额（提交后服务端核价）</span><strong>{money(items.reduce((sum, item) => sum + (item.product?.price ?? 0) * item.quantity, 0))}</strong><Button type="primary" size="large" block loading={submitting} disabled={!addressId || (!addresses.length && !recovering)} onClick={() => void submit()}>确认并提交订单</Button><small>订单将按店铺分别创建，最终金额以服务端返回为准。</small></aside>
    </div>}
  </main>;
}
