import { defineMock } from 'umi';

const products = Array.from({ length: 30 }, (_, i) => ({
  id: i + 1,
  name: ['运动T恤', '休闲牛仔裤', '夏季连衣裙', '商务衬衫', '运动鞋', '皮带', '帆布包', '太阳镜'][i % 8] + ` ${i + 1}号`,
  category: ['服装', '鞋类', '配件'][i % 3],
  price: Number(((i + 1) * 29.9).toFixed(2)),
  stock: (i % 4 === 0) ? 0 : Math.floor(Math.random() * 200) + 10,
  status: i % 5 === 0 ? 0 : 1,
  createdAt: `2024-0${(i % 9) + 1}-${String((i % 28) + 1).padStart(2, '0')}`,
}));

export default defineMock({
  'GET /api/product/list': (req, res) => {
    const { current = 1, pageSize = 10, name, status } = req.query;
    let list = [...products];
    if (name) list = list.filter((p) => p.name.includes(String(name)));
    if (status !== undefined && status !== '') list = list.filter((p) => p.status === Number(status));
    const start = (Number(current) - 1) * Number(pageSize);
    res.json({ data: list.slice(start, start + Number(pageSize)), total: list.length, success: true });
  },
});
