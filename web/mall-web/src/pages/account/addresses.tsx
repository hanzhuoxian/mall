import { Button, Card, Empty, Form, Input, Spin } from 'antd';
import { useCallback, useEffect, useState } from 'react';
import { marketplace, type Address } from '../../services/marketplace';

export default function AddressesPage() {
  const [form] = Form.useForm<Omit<Address, 'id'>>();
  const [items, setItems] = useState<Address[]>([]);
  const [loading, setLoading] = useState(true);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const load = useCallback(() => marketplace.addresses().then((page) => setItems(page.items)).catch((err) => setError(err instanceof Error ? err.message : '地址加载失败')).finally(() => setLoading(false)), []);
  useEffect(() => { void load(); }, [load]);
  const save = async (values: Omit<Address, 'id'>) => {
    setBusy(true); setError('');
    try { const saved = await marketplace.createAddress(values); setItems((current) => [...current, saved]); form.resetFields(); }
    catch (err) { setError(err instanceof Error ? err.message : '地址保存失败'); }
    finally { setBusy(false); }
  };
  const remove = async (id: string) => {
    try { await marketplace.deleteAddress(id); setItems((current) => current.filter((item) => item.id !== id)); }
    catch (err) { setError(err instanceof Error ? err.message : '地址删除失败'); }
  };
  return <main className="page-wrap"><div className="eyebrow">Your account</div><h1 className="page-title">收货地址</h1>
    {error && <p className="inline-error" role="alert">{error}</p>}
    <div className="address-layout"><Card title="新增地址" className="address-form-card">
      <Form form={form} layout="vertical" onFinish={save}>
        <Form.Item name="recipient" label="收件人" rules={[{ required: true }]}><Input /></Form.Item>
        <Form.Item name="phone" label="联系电话" rules={[{ required: true }]}><Input /></Form.Item>
        <Form.Item name="detail" label="详细地址" rules={[{ required: true }]}><Input.TextArea rows={3} /></Form.Item>
        <Button type="primary" htmlType="submit" loading={busy}>保存地址</Button>
      </Form>
    </Card><section className="address-list">{loading ? <Spin /> : items.length === 0 ? <Empty description="还没有收货地址" /> : items.map((item) => <Card key={item.id} className="address-card">
      <div><strong>{item.recipient}</strong><span>{item.phone}</span></div><p>{item.detail}</p>
      <Button type="link" danger aria-label="删除地址" onClick={() => void remove(item.id)}>删除</Button>
    </Card>)}</section></div>
  </main>;
}
