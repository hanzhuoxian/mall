import { Button, message } from 'antd';
import { ShopOutlined, ShoppingCartOutlined, UserOutlined } from '@ant-design/icons';
import { Outlet } from 'react-router-dom';
import { useEffect, useState } from 'react';
import { ACCESS_TOKEN, logout } from '../pages/user/service';

export default function StoreLayout() {
  const [signedIn, setSignedIn] = useState(false);
  useEffect(() => { setSignedIn(Boolean(localStorage.getItem(ACCESS_TOKEN))); }, []);
  const authPage = typeof window !== 'undefined' && window.location.pathname.startsWith('/user/');
  if (authPage) return <Outlet />;
  const signOut = async () => {
    try { await logout(); message.success('已退出登录'); }
    catch { message.error('已清除本地登录状态，服务器退出确认失败'); }
    finally { setSignedIn(false); }
  };
  return <div className="store-shell">
    <div className="announcement">为日常挑选好物 · 安心逛，慢慢选</div>
    <header className="store-header"><a className="brand" href="/">MALL<span>.</span></a>
      <form className="store-search" action="/" method="get"><input aria-label="全站搜索" name="name" placeholder="搜索你喜欢的好物" /></form>
      <nav><a href="/">发现好物</a><a href="/orders">我的订单</a></nav>
      <div className="store-actions"><a aria-label="购物车" href="/cart"><ShoppingCartOutlined /></a><a aria-label="账户" href={signedIn ? '/account/addresses' : '/user/login'}><UserOutlined /></a>{signedIn && <Button type="text" size="small" aria-label="退出登录" onClick={() => void signOut()}>退出登录</Button>}</div>
    </header>
    <Outlet />
    <footer className="store-footer"><ShopOutlined /> MALL · 为日常带来一点好物</footer>
  </div>;
}
