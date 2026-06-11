import { defineMock } from 'umi';

const names = ['张三', '李四', '王五', '赵六', '孙七', '周八', '吴九', '郑十'];

const customers = Array.from({ length: 35 }, (_, i) => ({
  id: i + 1,
  username: names[i % names.length] + (i + 1),
  email: `user${i + 1}@example.com`,
  phone: `138${String(i + 1).padStart(8, '0')}`,
  status: i % 7 === 0 ? 0 : 1,
  orderCount: Math.floor(Math.random() * 20),
  createdAt: `2024-0${(i % 9) + 1}-${String((i % 28) + 1).padStart(2, '0')}`,
}));

export default defineMock({
  'GET /api/customer/list': (req, res) => {
    const { current = 1, pageSize = 10, username, status } = req.query;
    let list = [...customers];
    if (username) list = list.filter((c) => c.username.includes(String(username)));
    if (status !== undefined && status !== '') list = list.filter((c) => c.status === Number(status));
    const start = (Number(current) - 1) * Number(pageSize);
    res.json({ data: list.slice(start, start + Number(pageSize)), total: list.length, success: true });
  },
});
