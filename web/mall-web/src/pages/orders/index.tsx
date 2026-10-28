import { Alert, Button, Card, Empty, Form, Input, Modal, Spin, message } from 'antd';
import { useCallback, useEffect, useRef, useState } from 'react';
import { marketplace, type MarketplaceContext, type Order } from '../../services/marketplace';

const money = (cents: number) => `¥${(cents / 100).toFixed(2)}`;
const statuses: Record<string, string> = { pending: '待付款', paid: '已付款', shipped: '待收货', completed: '已完成', cancelled: '已取消', refunding: '售后处理中', refunded: '已退款' };
type Action = 'cancel' | 'pay' | 'receive' | 'refund';

export default function OrdersPage() {
  const [context, setContext] = useState<MarketplaceContext>();
  const [orders, setOrders] = useState<Order[]>([]);
  const [status, setStatus] = useState('');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [active, setActive] = useState<{ order: Order; action: Action }>();
  const [busy, setBusy] = useState(false);
  const latestLoad = useRef(0);
  const [form] = Form.useForm<{ reason: string }>();

  const load = useCallback(async (filter: string) => {
    const requestId = ++latestLoad.current;
    setLoading(true); setError('');
    try {
      const [identity, page] = await Promise.all([marketplace.context(), marketplace.orders(filter || undefined)]);
      if (requestId === latestLoad.current) { setContext(identity); setOrders(page.items); }
    } catch (err) { if (requestId === latestLoad.current) setError(err instanceof Error ? err.message : '订单加载失败'); }
    finally { if (requestId === latestLoad.current) setLoading(false); }
  }, []);
  useEffect(() => { void load(''); }, [load]);

  const openAction = (order: Order, action: Action) => { setActive({ order, action }); form.resetFields(); };
  const confirmAction = async (values: { reason?: string }) => {
    if (!active) return;
    setBusy(true); setError('');
    try {
      await marketplace.orderAction(active.order.id, active.action, values.reason);
      const next: Record<Action, string> = { pay: 'paid', cancel: 'cancelled', receive: 'completed', refund: 'refunding' };
      setOrders((current) => current.flatMap((order) => {
        if (order.id !== active.order.id) return [order];
        if (status && status !== next[active.action]) return [];
        return [{ ...order, status: next[active.action] }];
      }));
      if (active.action === 'pay') message.success('开发沙箱付款已完成；没有产生真实扣款。');
      setActive(undefined);
    } catch (err) { setError(err instanceof Error ? err.message : '订单操作失败'); }
    finally { setBusy(false); }
  };
  const actionLabel = active?.action === 'pay' ? '确认开发沙箱付款' : active?.action === 'cancel' ? '确认取消订单' : active?.action === 'receive' ? '确认收货' : '提交售后申请';

  return <main className="page-wrap"><div className="eyebrow">Your orders</div><h1 className="page-title">我的订单</h1>
    {context?.sandboxPayment && <Alert className="sandbox-notice" type="warning" showIcon title="开发沙箱付款仅供开发演示，不产生真实扣款。" />}
    {error && <Alert className="checkout-error" type="error" showIcon title={error} />}
    <nav className="order-filter" aria-label="订单状态筛选">{[
      { value: '', label: '全部订单' }, { value: 'pending', label: '待付款' }, { value: 'paid', label: '已付款' }, { value: 'shipped', label: '待收货' }, { value: 'completed', label: '已完成' },
    ].map((filter) => <Button key={filter.value || 'all'} type={status === filter.value ? 'primary' : 'default'} aria-label={`筛选${filter.label}`} onClick={() => { setStatus(filter.value); void load(filter.value); }}>{filter.label}</Button>)}</nav>
    {loading ? <div className="loading"><Spin /></div> : orders.length === 0 ? <Empty description="还没有相关订单" /> : <div className="orders-list">{orders.map((order) => <Card key={order.id} className="order-card" variant="borderless">
      <header><span>订单号　{order.id}</span><strong>{statuses[order.status] ?? order.status}</strong></header>
      <div className="order-items">{order.items.map((item) => <div key={item.id}><span>{item.name} × {item.quantity}</span><span>{money(item.price * item.quantity)}</span></div>)}</div>
      <div className="order-address">收货信息：{order.address}</div>
      <footer><span>订单金额</span><strong>{money(order.amount)}</strong><div className="order-actions">
        {order.status === 'pending' && <>
          <Button onClick={() => openAction(order, 'cancel')}>取消订单</Button>
          {context?.sandboxPayment && (!order.expiresAt || Date.parse(order.expiresAt) > Date.now()) && <Button type="primary" onClick={() => openAction(order, 'pay')}>开发沙箱付款</Button>}
        </>}
        {order.status === 'shipped' && <Button type="primary" onClick={() => openAction(order, 'receive')}>确认收货</Button>}
        {['paid', 'shipped', 'completed'].includes(order.status) && <Button onClick={() => openAction(order, 'refund')}>申请售后</Button>}
      </div></footer>
    </Card>)}</div>}
    <Modal open={Boolean(active)} title={active?.action === 'pay' ? '开发沙箱付款' : active?.action === 'cancel' ? '取消订单' : active?.action === 'receive' ? '确认收货' : '申请售后'} okText={actionLabel} cancelText="再想想" confirmLoading={busy} onCancel={() => setActive(undefined)} onOk={() => form.submit()}>
      {active?.action === 'pay' && <Alert type="warning" showIcon title="开发沙箱付款仅供开发演示，不产生真实扣款。" />}
      {active?.action === 'refund' ? <Form form={form} onFinish={confirmAction}><Form.Item name="reason" label="售后原因" rules={[{ required: true, message: '请填写售后原因' }]}><Input.TextArea rows={3} maxLength={500} /></Form.Item></Form> : <Form form={form} onFinish={confirmAction} />}
    </Modal>
  </main>;
}
