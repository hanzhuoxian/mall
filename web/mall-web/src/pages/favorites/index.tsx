import { ArrowLeftOutlined, HeartFilled } from '@ant-design/icons';
import { Button, Empty, message, Spin } from 'antd';
import { useEffect, useState } from 'react';
import { marketplace, type Favorite } from '../../services/marketplace';
import { ACCESS_TOKEN, loginUrl } from '../user/service';

const money = (cents: number) => `¥${(cents / 100).toFixed(2)}`;

export default function FavoritesPage() {
  const [items, setItems] = useState<Favorite[]>([]);
  const [loading, setLoading] = useState(true);
  const [removing, setRemoving] = useState<string>();
  const load = () => marketplace.favorites().then(({ items: favorites }) => setItems(favorites)).catch((error) => message.error(error instanceof Error ? error.message : '收藏列表加载失败')).finally(() => setLoading(false));
  useEffect(() => {
    if (!localStorage.getItem(ACCESS_TOKEN)) {
      window.location.assign(loginUrl('/favorites'));
      setLoading(false);
      return;
    }
    void load();
  }, []);
  const remove = async (productId: string) => {
    setRemoving(productId);
    try { await marketplace.deleteFavorite(productId); setItems((current) => current.filter((item) => item.productId !== productId)); message.success('已取消收藏'); }
    catch (error) { message.error(error instanceof Error ? error.message : '取消收藏失败'); }
    finally { setRemoving(undefined); }
  };
  return <main className="page-wrap"><a href="/" className="back-link"><ArrowLeftOutlined /> 返回好物精选</a><section className="section-heading"><div><div className="eyebrow">YOUR SHORTLIST</div><h1>我的收藏</h1></div></section>
    {loading ? <div className="loading"><Spin /></div> : items.length === 0 ? <Empty description="还没有收藏商品，遇到喜欢的好物就收进来吧" /> : <div className="product-grid">{items.map(({ product, productId }) => product && <article className="product-card" key={productId}><a href={`/product/${product.id}`} className="product-card-link">{product.image ? <img src={product.image} alt={product.name} /> : <div className="image-placeholder" role="img" aria-label={`${product.name}商品图片`} />}<div className="product-card-copy"><h2>{product.name}</h2><p>{money(product.price)}</p></div></a><div className="favorite-actions"><Button icon={<HeartFilled />} loading={removing === productId} onClick={() => void remove(productId)}>取消收藏</Button><Button type="primary" href={`/product/${product.id}`}>查看商品</Button></div></article>)}</div>}
  </main>;
}
