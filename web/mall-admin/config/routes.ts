export default [
  {
    path: '/user',
    layout: false,
    routes: [
      { path: '/user/register', component: './user/register' },
      {
        name: 'login',
        path: '/user/login',
        component: './user/login',
      },
    ],
  },
  {
    path: '/dashboard',
    name: 'dashboard',
    icon: 'dashboard',
    component: './marketplace',
  },
  {
    path: '/market/catalog',
    name: 'catalog',
    icon: 'shopping',
    component: './marketplace',
  },
  {
    path: '/market/shops',
    name: 'shops',
    icon: 'shop',
    component: './marketplace',
  },
  {
    path: '/market/products',
    name: 'products',
    icon: 'appstore',
    component: './marketplace',
  },
  {
    path: '/market/categories',
    name: 'categories',
    icon: 'tags',
    component: './marketplace',
  },
  {
    path: '/market/cart',
    name: 'cart',
    icon: 'shoppingCart',
    component: './marketplace',
  },
  {
    path: '/market/orders',
    name: 'orders',
    icon: 'fileText',
    component: './marketplace',
  },
  {
    path: '/market/addresses',
    name: 'addresses',
    icon: 'environment',
    component: './marketplace',
  },
  {
    path: '/market/finance',
    name: 'finance',
    icon: 'accountBook',
    component: './marketplace',
  },
  {
    path: '/market/audit',
    name: 'audit',
    icon: 'audit',
    component: './marketplace',
  },
  {
    path: '/system',
    access: 'canAdmin',
    name: 'system',
    icon: 'setting',
    routes: [
      {
        path: '/system',
        redirect: '/system/user',
      },
      {
        path: '/system/user',
        name: 'user',
        component: './system/user',
      },
    ],
  },
  {
    path: '/',
    redirect: '/dashboard',
  },
  {
    component: './exception/404',
    layout: false,
    path: './*',
  },
];
