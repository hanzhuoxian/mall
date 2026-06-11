import { defineMock } from 'umi';

const statuses = [0, 1, 2, 3, 4];
const customers = ['张三', '李四', '王五', '赵六', '孙七', '周八', '吴九', '郑十'];

const orders = Array.from({ length: 40 }, (_, i) => ({
  id: i + 1,
  orderNo: `ORD-2024${String(i + 1).padStart(5, '0')}`,
  customer: customers[i % customers.length],
  amount: Number(((i + 1) * 88.8).toFixed(2)),
  status: statuses[i % statuses.length],
  itemCount: (i % 5) + 1,
  createdAt: `2024-0${(i % 9) + 1}-${String((i % 28) + 1).padStart(2, '0')}`,
}));

export default defineMock({
  'GET /api/order/list': (req, res) => {
    const { current = 1, pageSize = 10, orderNo, status } = req.query;
    let list = [...orders];
    if (orderNo) list = list.filter((o) => o.orderNo.includes(String(orderNo)));
    if (status !== undefined && status !== '') list = list.filter((o) => o.status === Number(status));
    const start = (Number(current) - 1) * Number(pageSize);
    res.json({ data: list.slice(start, start + Number(pageSize)), total: list.length, success: true });
  },
});
