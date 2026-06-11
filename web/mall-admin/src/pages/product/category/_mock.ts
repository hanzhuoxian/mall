import { defineMock } from 'umi';

const categories = [
  { id: 1, name: '服装', parentId: 0, sort: 1, status: 1, productCount: 120 },
  { id: 2, name: '男装', parentId: 1, sort: 1, status: 1, productCount: 60 },
  { id: 3, name: '女装', parentId: 1, sort: 2, status: 1, productCount: 60 },
  { id: 4, name: '鞋类', parentId: 0, sort: 2, status: 1, productCount: 80 },
  { id: 5, name: '运动鞋', parentId: 4, sort: 1, status: 1, productCount: 40 },
  { id: 6, name: '皮鞋', parentId: 4, sort: 2, status: 1, productCount: 40 },
  { id: 7, name: '配件', parentId: 0, sort: 3, status: 1, productCount: 50 },
  { id: 8, name: '包袋', parentId: 7, sort: 1, status: 0, productCount: 25 },
];

export default defineMock({
  'GET /api/product/category': (_req, res) => {
    res.json({ data: categories, success: true });
  },
});
