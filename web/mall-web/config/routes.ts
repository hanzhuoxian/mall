export default [{
  path: '/', component: '@/layouts/StoreLayout', routes: [
    { path: '/', component: '@/pages/home' },
    { path: '/product/:id', component: '@/pages/product/detail' },
    { path: '/cart', component: '@/pages/cart' },
    { path: '/checkout', component: '@/pages/checkout' },
    { path: '/orders', component: '@/pages/orders' },
    { path: '/account/addresses', component: '@/pages/account/addresses' },
    { path: '/user/login', component: '@/pages/user/login' },
    { path: '/user/register', component: '@/pages/user/register' },
  ],
}];
