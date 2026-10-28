import { ArrowLeftOutlined, CheckCircleFilled } from '@ant-design/icons';
import { Alert, Button, Empty, InputNumber, Spin, message } from 'antd';
import { useEffect, useState } from 'react';
import { marketplace, type Product } from '../../services/marketplace';
import { ACCESS_TOKEN, loginUrl } from '../user/service';

const money = (cents: number) => `¥${(cents / 100).toFixed(2)}`;

export default function ProductDetailPage() {
  const id = typeof window === 'undefined' ? '' : window.location.pathname.split('/').filter(Boolean).at(-1) ?? '';
  const [product, setProduct] = useState<Product>();
  const [failed, setFailed] = useState(false);
  const [quantity, setQuantity] = useState(1);
  const [loading, setLoading] = useState(true);
  const [adding, setAdding] = useState(false);
  useEffect(() => {
    let active = true;
    marketplace.product(id).then((item) => { if (active) setProduct(item); }).catch(() => { if (active) setFailed(true); }).finally(() => { if (active) setLoading(false); });
    return () => { active = false; };
  }, [id]);
  const add = async () => {
    if (!localStorage.getItem(ACCESS_TOKEN)) {
      window.location.assign(loginUrl(window.location.pathname + window.location.search));
      return;
    }
    if (!product) return;
    setAdding(true);
    try { await marketplace.putCart(product.id, quantity); message.success('已加入购物车'); }
    catch (error) { message.error(error instanceof Error ? error.message : '加入购物车失败'); }
    finally { setAdding(false); }
  };
  if (loading) return <div className="page-wrap loading"><Spin /></div>;
  if (failed || !product) return <div className="page-wrap"><a href="/"><ArrowLeftOutlined /> 返回商城</a><Empty className="unavailable" description="这件商品暂时无法查看，可能已下架或不再销售" /></div>;
  return <main className="page-wrap"><a href="/" className="back-link"><ArrowLeftOutlined /> 返回好物精选</a>
    <section className="product-layout">
      {product.image ? <img className="detail-image" src={product.image} alt={product.name} /> : <div className="detail-image image-placeholder" role="img" aria-label={`${product.name}商品图片`} />}
      <div className="detail-copy"><div className="eyebrow">A little something for you</div><h1>{product.name}</h1><div className="detail-price">{money(product.price)}</div>
        <p className="detail-description">{product.description || '为日常精心挑选的好物。'}</p>
        <div className="stock-note"><CheckCircleFilled /> {product.stock > 0 ? `现货 · 库存 ${product.stock} 件` : '暂时售罄'}</div>
        <div className="buy-row"><InputNumber aria-label="购买数量" min={1} max={Math.max(product.stock, 1)} value={quantity} onChange={(value) => setQuantity(value ?? 1)} disabled={product.stock < 1} /><Button type="primary" size="large" onClick={() => void add()} loading={adding} disabled={product.stock < 1}>加入购物车</Button></div>
        {product.stock < 1 && <Alert type="warning" message="商品目前没有库存" showIcon />}
        <p className="detail-footnote">库存与价格以提交订单时服务端核验为准。</p>
      </div>
    </section>
  </main>;
}
