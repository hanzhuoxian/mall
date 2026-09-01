export default [
  {
    path: '/user',
    layout: false,
    routes: [
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
    component: './dashboard',
  },
  {
    path: '/product',
    name: 'product',
    icon: 'shopping',
    routes: [
      {
        path: '/product',
        redirect: '/product/list',
      },
      {
        path: '/product/list',
        name: 'list',
        component: './product/list',
      },
      {
        path: '/product/category',
        name: 'category',
        component: './product/category',
      },
    ],
  },
  {
    path: '/order',
    name: 'order',
    icon: 'fileText',
    component: './order',
  },
  {
    path: '/customer',
    name: 'customer',
    icon: 'team',
    component: './customer',
  },
  {
    path: '/system',
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
