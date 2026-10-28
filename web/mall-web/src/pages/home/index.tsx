import { Button, Empty, Input, Pagination, Spin, Tag } from 'antd';
import { useEffect, useState } from 'react';
import { marketplace, type Category, type Product } from '../../services/marketplace';

const money = (cents: number) => `¥${(cents / 100).toFixed(2)}`;
const initialSearch = () => typeof window === 'undefined' ? '' : new URLSearchParams(window.location.search).get('name') ?? '';

export default function HomePage() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [total, setTotal] = useState(0);
  const [name, setName] = useState(initialSearch);
  const [draft, setDraft] = useState(initialSearch);
  const [categoryId, setCategoryId] = useState('');
  const [current, setCurrent] = useState(1);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [retry, setRetry] = useState(0);

  useEffect(() => { marketplace.categories({ pageSize: 100 }).then((page) => setCategories(page.items)).catch(() => setError('分类加载失败')); }, []);
  useEffect(() => {
    void retry;
    let active = true;
    setLoading(true); setError('');
    marketplace.catalog({ name, categoryId, current, pageSize: 12 })
      .then((page) => { if (active) { setProducts(page.items); setTotal(page.total); } })
      .catch(() => { if (active) setError('商品加载失败，请稍后重试'); })
      .finally(() => { if (active) setLoading(false); });
    return () => { active = false; };
  }, [name, categoryId, current, retry]);

  return <main className="page-wrap">
    <section className="hero"><div><div className="eyebrow">Thoughtfully picked</div><h1>把喜欢的日常，<br />过得更好一点。</h1><p>发现值得带回家的设计与好物</p></div><div className="hero-mark" aria-hidden="true">好<br />物</div></section>
    <section className="catalog-tools" aria-label="商品筛选">
      <Input aria-label="搜索商品" placeholder="搜索商品名称" value={draft} onChange={(event) => setDraft(event.target.value)} onPressEnter={() => { setCurrent(1); setName(draft.trim()); }} />
      <Button type="primary" onClick={() => { setCurrent(1); setName(draft.trim()); }}>搜索</Button>
    </section>
    <div className="category-row"><Tag color={!categoryId ? 'green' : undefined} onClick={() => { setCategoryId(''); setCurrent(1); }}>全部好物</Tag>
      {categories.map((category) => <Tag key={category.id} color={categoryId === category.id ? 'green' : undefined} onClick={() => { setCategoryId(category.id); setCurrent(1); }}>{category.name}</Tag>)}
    </div>
    <div className="section-heading"><div><div className="eyebrow">Shop the collection</div><h2>当季精选</h2></div><span>{total} 件好物</span></div>
    {error ? <Empty description={error}><Button onClick={() => setRetry((value) => value + 1)}>重试</Button></Empty> : loading ? <div className="loading"><Spin /></div> : products.length === 0 ? <Empty description="暂时没有找到商品，换个关键词试试" /> : <div className="product-grid">
      {products.map((product) => <a className="product-card" key={product.id} href={`/product/${product.id}`}>
        {product.image ? <img className="product-image" src={product.image} alt={product.name} /> : <div className="product-image image-placeholder" role="img" aria-label={`${product.name}商品图片`} />}
        <div className="product-info"><h3>{product.name}</h3><div className="product-meta">{product.stock > 0 ? `库存 ${product.stock}` : '暂时售罄'}</div><div className="price">{money(product.price)}</div></div>
      </a>)}
    </div>}
    {total > 12 && <Pagination className="pagination" current={current} pageSize={12} total={total} onChange={setCurrent} />}
  </main>;
}
