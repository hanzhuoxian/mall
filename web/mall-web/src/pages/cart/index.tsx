import { Button, Checkbox, Empty, InputNumber, message, Spin } from 'antd';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { marketplace, type CartItem } from '../../services/marketplace';

const money = (cents: number) => `¥${(cents / 100).toFixed(2)}`;

export default function CartPage() {
  const [items, setItems] = useState<CartItem[]>([]);
  const [selected, setSelected] = useState<string[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [busy, setBusy] = useState(false);
  const load = useCallback(() => marketplace.cart().then((page) => setItems(page.items)).catch((err) => setError(err instanceof Error ? err.message : '购物车加载失败')).finally(() => setLoading(false)), []);
  useEffect(() => { void load(); }, [load]);
  const total = useMemo(() => items.filter((item) => selected.includes(item.id)).reduce((sum, item) => sum + (item.product?.price ?? 0) * item.quantity, 0), [items, selected]);
  const update = async (item: CartItem, quantity: number) => {
    setBusy(true); setError('');
    try { await marketplace.updateCart({ ...item, quantity }); setItems((current) => current.map((value) => value.id === item.id ? { ...value, quantity } : value)); }
    catch (err) { setError(err instanceof Error ? err.message : '数量更新失败'); }
    finally { setBusy(false); }
  };
  const remove = async (id: string) => {
    setBusy(true);
    try { await marketplace.deleteCart(id); setItems((current) => current.filter((item) => item.id !== id)); setSelected((current) => current.filter((value) => value !== id)); }
    catch (err) { setError(err instanceof Error ? err.message : '移除失败'); }
    finally { setBusy(false); }
  };
  const proceed = () => {
    if (!selected.length) { message.info('请先选择要结算的商品'); return; }
    sessionStorage.setItem('mall.checkout.cartIds', JSON.stringify(selected));
    window.location.href = '/checkout';
  };
  if (loading) return <div className="page-wrap loading"><Spin /></div>;
  return <main className="page-wrap"><div className="eyebrow">Your picks</div><h1 className="page-title">购物车</h1>
    {error && <div className="inline-error" role="alert">{error} <Button type="link" onClick={() => void load()}>重试</Button></div>}
    {items.length === 0 ? <Empty description="购物车还是空的，去逛逛喜欢的好物吧"><Button type="primary" href="/">去逛逛</Button></Empty> : <>
      <div className="cart-list">{items.map((item) => <article className="cart-row" key={item.id}>
        <Checkbox checked={selected.includes(item.id)} onChange={(event) => setSelected((current) => event.target.checked ? [...current, item.id] : current.filter((id) => id !== item.id))} aria-label={`选择${item.product?.name ?? '商品'}`} />
        {item.product?.image ? <img src={item.product.image} alt="" /> : <div className="cart-thumb image-placeholder" />}
        <div className="cart-product"><strong>{item.product?.name ?? '商品已不可用'}</strong><span>{item.product ? money(item.product.price) : ''}</span></div>
        <InputNumber aria-label="商品数量" min={1} max={Math.max(item.product?.stock ?? 1, 1)} value={item.quantity} disabled={busy} onChange={(value) => value && value !== item.quantity && void update(item, value)} />
        <Button type="link" danger aria-label="移除商品" disabled={busy} onClick={() => void remove(item.id)}>移除</Button>
      </article>)}</div>
      <section className="cart-summary"><span>已选 {selected.length} 件 · 商品小计</span><strong>{money(total)}</strong><Button type="primary" size="large" onClick={proceed} disabled={busy}>去结算</Button></section>
    </>}
  </main>;
}
