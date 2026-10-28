import { PageContainer } from '@ant-design/pro-components';
import { useLocation } from '@umijs/max';
import {
  Alert,
  Button,
  Card,
  Col,
  Descriptions,
  Empty,
  Form,
  Input,
  Modal,
  message,
  Popconfirm,
  Result,
  Row,
  Select,
  Space,
  Spin,
  Statistic,
  Table,
  Tag,
  Typography,
} from 'antd';
import type { TableColumnsType, TableProps } from 'antd';
import React, { useCallback, useEffect, useRef, useState } from 'react';
import {
  type Address,
  api,
  type Cart,
  type Category,
  type Context,
  labels,
  money,
  type Order,
  type Page,
  type Product,
  type Shop,
} from './api';
import './style.less';

type Editor = {
  title: string;
  fields: {
    name: string;
    label: string;
    type?: 'number' | 'text' | 'select' | 'textarea';
    options?: { label: string; value: string }[];
    min?: number;
    max?: number;
    optional?: boolean;
  }[];
  initial?: Record<string, unknown>;
  submit: (values: any) => Promise<unknown>;
};
const titles: Record<string, string> = {
  dashboard: '经营概览',
  shops: '商家与入驻',
  products: '商品管理',
  categories: '商品分类',
  catalog: '逛商城',
  cart: '购物车与结算',
  orders: '订单与售后',
  addresses: '收货地址',
  finance: '结算概览',
  audit: '操作审计',
};
const Status = ({
  value,
  order = false,
}: {
  value: string;
  order?: boolean;
}) => (
  <Tag
    color={
      ['approved', 'completed'].includes(value)
        ? 'green'
        : ['rejected', 'cancelled', 'suspended'].includes(value)
          ? 'red'
          : 'blue'
    }
  >
    {order && value === 'pending' ? '待付款' : labels[value] || value}
  </Tag>
);

export default function Marketplace() {
  const location = useLocation();
  const section = location.pathname.split('/').pop() || 'dashboard';
  const [ctx, setCtx] = useState<Context>();
  const [rows, setRows] = useState<Record<string, unknown>[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [query, setQuery] = useState('');
  const [scope, setScope] = useState('buyer');
  const [categories, setCategories] = useState<Category[]>([]);
  const [addresses, setAddresses] = useState<Address[]>([]);
  const [addressId, setAddressId] = useState<string>();
  const [selected, setSelected] = useState<React.Key[]>([]);
  const [summary, setSummary] = useState({
    orders: 0,
    gmv: 0,
    commission: 0,
    payable: 0,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);
  const [busy, setBusy] = useState(false);
  const [editor, setEditor] = useState<Editor>();
  const [detail, setDetail] = useState<Order>();
  const [form] = Form.useForm();
  const requestKey = useRef<{ signature: string; key: string } | undefined>(
    undefined,
  );
  const sequence = useRef(0);
  const load = useCallback(async () => {
    const seq = ++sequence.current;
    setLoading(true);
    setError(false);
    try {
      const [context, cats] = await Promise.all([
        api<Context>('/context'),
        api<Page<Category>>('/categories', 'GET', undefined, { pageSize: 100 }),
      ]);
      if (seq !== sequence.current) return;
      setCtx(context);
      setCategories(cats.items);
      if (['dashboard', 'finance'].includes(section)) {
        const data = await api<typeof summary>('/dashboard', 'GET', undefined, {
          scope,
        });
        if (seq === sequence.current) setSummary(data);
      } else {
        const endpoint =
          section === 'categories' ? '/categories' : `/${section}`;
        const data = await api<Page<any>>(endpoint, 'GET', undefined, {
          current: page,
          pageSize: 20,
          name: query,
          scope,
        });
        if (seq !== sequence.current) return;
        setRows(data.items);
        setTotal(data.total);
        if (section === 'cart') {
          const result = await api<Page<Address>>(
            '/addresses',
            'GET',
            undefined,
            { pageSize: 100 },
          );
          if (seq === sequence.current) {
            setAddresses(result.items);
            setSelected([]);
          }
        }
      }
    } catch {
      if (seq === sequence.current) setError(true);
    } finally {
      if (seq === sequence.current) setLoading(false);
    }
  }, [section, page, query, scope]);
  useEffect(() => {
    void load();
    return () => {
      sequence.current++;
    };
  }, [load]);
  useEffect(() => {
    setPage(1);
    setQuery('');
    setSelected([]);
    setDetail(undefined);
  }, [section]);
  useEffect(() => {
    if (editor) form.setFieldsValue(editor.initial || {});
  }, [editor, form]);
  const run = async (fn: () => Promise<unknown>) => {
    setBusy(true);
    try {
      await fn();
      message.success('操作成功');
      await load();
    } catch {
      /* Global request handler displays failures. */
    } finally {
      setBusy(false);
    }
  };
  const edit = (config: Editor) => {
    form.resetFields();
    setEditor(config);
  };
  const action = (o: Order, name: string, fields: Editor['fields'] = []) => {
    edit({
      title: (
        {
          pay: '确认沙箱付款（不会扣款）',
          cancel: '取消订单并释放库存',
          ship: '填写发货信息',
          receive: '确认收货',
          refund: '申请售后',
          'approve-refund': '同意沙箱退款（无真实资金）',
          'reject-refund': '驳回售后申请',
        } as Record<string, string>
      )[name],
      fields,
      submit: (values) =>
        api(`/orders/${o.id}/actions`, 'POST', { action: name, ...values }),
    });
  };
  const productEditor = (p?: Product) =>
    edit({
      title: p ? '编辑商品（将重新审核）' : '发布商品',
      initial: p
        ? { ...p, price: p.price / 100 }
        : {
            shopId: ctx?.shops.find((s) => s.status === 'approved')?.id,
            stock: 0,
          },
      fields: [
        {
          name: 'shopId',
          label: '所属店铺',
          type: 'select',
          options: ctx?.shops
            .filter((s) => s.status === 'approved')
            .map((s) => ({ value: s.id, label: s.name })),
        },
        { name: 'name', label: '商品名称' },
        {
          name: 'categoryId',
          label: '商品分类',
          type: 'select',
          options: categories.map((c) => ({ value: c.id, label: c.name })),
        },
        {
          name: 'price',
          label: '售价（元，最多两位小数）',
          type: 'number',
          min: 0.01,
          max: 1000000,
        },
        ...(!p
          ? [
              {
                name: 'stock',
                label: '初始库存',
                type: 'number' as const,
                min: 0,
                max: 100000000,
              },
            ]
          : []),
        { name: 'image', label: '图片 HTTPS 地址', optional: true },
        {
          name: 'description',
          label: '商品介绍',
          type: 'textarea',
          optional: true,
        },
      ],
      submit: (v) =>
        api(p ? `/products/${p.id}` : '/products', p ? 'PUT' : 'POST', {
          ...v,
          price: Math.round(v.price * 100),
        }),
    });
  const review = (
    type: 'shops' | 'products',
    record: Shop | Product,
    status: string,
  ) =>
    edit({
      title: `${labels[status]} · ${record.name}`,
      initial: {
        commissionBps: 'commissionBps' in record ? record.commissionBps : 500,
      },
      fields: [
        {
          name: 'reason',
          label: '审核意见',
          type: 'textarea',
          optional: status === 'approved' || status === 'offline',
        },
        ...(type === 'shops'
          ? [
              {
                name: 'commissionBps',
                label: '平台佣金（基点，500 = 5%）',
                type: 'number' as const,
                min: 0,
                max: 10000,
              },
            ]
          : []),
      ],
      submit: (v) =>
        api(`/${type}/${record.id}/review`, 'POST', { ...v, status }),
    });
  const scopePicker = (
    <Select
      aria-label="数据范围"
      value={scope}
      style={{ width: 160 }}
      onChange={(v) => {
        setScope(v);
        setPage(1);
      }}
      options={[
        { label: '我的消费订单', value: 'buyer' },
        ...(ctx?.shops.length
          ? [{ label: '我的店铺订单', value: 'merchant' }]
          : []),
        ...(ctx?.admin ? [{ label: '全平台订单', value: 'platform' }] : []),
      ]}
    />
  );
  const table = <T extends { id?: string }>(
    columns: TableColumnsType<T>,
    extra?: Pick<TableProps<T>, 'rowSelection' | 'onRow'>,
  ) => (
    <Table
      rowKey="id"
      loading={loading}
      dataSource={rows as T[]}
      columns={columns}
      scroll={{ x: 800 }}
      pagination={{
        current: page,
        pageSize: 20,
        total,
        showSizeChanger: false,
        onChange: setPage,
      }}
      {...extra}
    />
  );
  let content: React.ReactNode;
  if (!ctx && loading) return <Spin fullscreen />;
  if (['dashboard', 'finance'].includes(section))
    content = (
      <>
        <div className="market-toolbar">{scopePicker}</div>
        <Row gutter={[16, 16]}>
          {[
            ['订单总数', summary.orders],
            ['有效成交金额', money(summary.gmv)],
            ['预计平台佣金', money(summary.commission)],
            ['已完成订单商家应收', money(summary.payable)],
          ].map(([label, value]) => (
            <Col xs={24} sm={12} xl={6} key={label}>
              <Card>
                <Statistic title={label} value={value} />
              </Card>
            </Col>
          ))}
        </Row>
        <Card title="运营工作台" style={{ marginTop: 24 }}>
          <Typography.Paragraph>
            入驻审核 → 发布商品 → 平台审核 → 消费者下单 → 支付 → 商家发货 →
            确认收货 → 售后与结算。
          </Typography.Paragraph>
          <Alert
            type="info"
            showIcon
            title="金额口径"
            description="成交金额包含待发货、待收货、已完成订单，不含取消和售后处理中订单。商家应收仅为已完成订单扣除佣金后的统计，不代表已打款；自动分账、提现及银行对账尚未接入。"
          />
        </Card>
      </>
    );
  else if (section === 'shops')
    content = (
      <>
        <div className="market-toolbar">
          <Button
            type="primary"
            disabled={!!ctx?.shops.length && ctx.shops[0].status !== 'rejected'}
            onClick={() =>
              edit({
                title: '申请商家入驻',
                fields: [
                  { name: 'name', label: '店铺名称' },
                  { name: 'contact', label: '联系人与联系电话' },
                  { name: 'license', label: '营业执照号' },
                ],
                submit: (v) => api('/shops', 'POST', v),
              })
            }
          >
            申请入驻／重新提交
          </Button>
          <Typography.Text type="secondary">
            每个账号可入驻一家店铺；平台审核通过后可以发布商品。
          </Typography.Text>
        </div>
        {table([
          { title: '店铺', dataIndex: 'name' },
          { title: '联系人', dataIndex: 'contact' },
          { title: '营业执照号', dataIndex: 'license' },
          {
            title: '状态',
            dataIndex: 'status',
            render: (v: string) => <Status value={v} />,
          },
          {
            title: '佣金比例',
            dataIndex: 'commissionBps',
            render: (v: number) => `${v / 100}%`,
          },
          { title: '审核意见', dataIndex: 'reason' },
          ...(ctx?.admin
            ? [
                {
                  title: '审核',
                  render: (_: unknown, s: Shop) => (
                    <Space wrap>
                      <Button
                        size="small"
                        onClick={() => review('shops', s, 'approved')}
                      >
                        通过
                      </Button>
                      <Button
                        size="small"
                        onClick={() => review('shops', s, 'rejected')}
                      >
                        驳回
                      </Button>
                      <Button
                        danger
                        size="small"
                        onClick={() => review('shops', s, 'suspended')}
                      >
                        冻结
                      </Button>
                    </Space>
                  ),
                },
              ]
            : []),
        ])}
      </>
    );
  else if (section === 'categories')
    content = (
      <>
        {ctx?.admin && (
          <div className="market-toolbar">
            <Button
              type="primary"
              onClick={() =>
                edit({
                  title: '新增分类',
                  fields: [{ name: 'name', label: '分类名称' }],
                  submit: (v) => api('/categories', 'POST', v),
                })
              }
            >
              新增分类
            </Button>
          </div>
        )}
        {table([
          { title: '分类名称', dataIndex: 'name' },
          { title: '分类编号', dataIndex: 'id' },
        ])}
      </>
    );
  else if (section === 'products')
    content = (
      <>
        <div className="market-toolbar">
          <Button
            type="primary"
            disabled={!ctx?.shops.some((s) => s.status === 'approved')}
            onClick={() => productEditor()}
          >
            发布商品
          </Button>
          <Input.Search
            placeholder="搜索商品名称"
            allowClear
            onSearch={(v) => {
              setQuery(v);
              setPage(1);
            }}
            style={{ maxWidth: 300 }}
          />
        </div>
        {table([
          { title: '商品', dataIndex: 'name' },
          { title: '售价', dataIndex: 'price', render: money },
          { title: '可售库存', dataIndex: 'stock' },
          {
            title: '状态',
            dataIndex: 'status',
            render: (v: string) => <Status value={v} />,
          },
          { title: '审核意见', dataIndex: 'reason' },
          {
            title: '操作',
            render: (_: unknown, p: Product) => (
              <Space wrap>
                {ctx?.shops.some((s) => s.id === p.shopId) && (
                  <Button size="small" onClick={() => productEditor(p)}>
                    编辑
                  </Button>
                )}
                <Button
                  size="small"
                  onClick={() =>
                    edit({
                      title: '调整可售库存',
                      fields: [
                        {
                          name: 'delta',
                          label: '变化数量（正数入库，负数出库）',
                          type: 'number',
                          min: -1000000,
                          max: 1000000,
                        },
                      ],
                      submit: (v) => api(`/products/${p.id}/stock`, 'POST', v),
                    })
                  }
                >
                  库存
                </Button>
                <Button
                  size="small"
                  onClick={() => review('products', p, 'offline')}
                >
                  下架
                </Button>
                {ctx?.admin && p.status === 'pending' && (
                  <>
                    <Button
                      size="small"
                      onClick={() => review('products', p, 'approved')}
                    >
                      审核通过
                    </Button>
                    <Button
                      size="small"
                      onClick={() => review('products', p, 'rejected')}
                    >
                      驳回
                    </Button>
                  </>
                )}
              </Space>
            ),
          },
        ])}
      </>
    );
  else if (section === 'catalog')
    content = (
      <>
        <div className="market-toolbar">
          <Input.Search
            placeholder="搜索在售商品"
            allowClear
            onSearch={(v) => {
              setQuery(v);
              setPage(1);
            }}
            style={{ maxWidth: 360 }}
          />
          <Tag>仅展示通过审核且店铺正常的商品</Tag>
        </div>
        {table([
          {
            title: '商品',
            render: (_: unknown, p: Product) => (
              <Space>
                {p.image && (
                  <img
                    alt={p.name}
                    src={p.image}
                    width={60}
                    height={60}
                    style={{ objectFit: 'cover', borderRadius: 8 }}
                  />
                )}
                <div>
                  <Typography.Text strong>{p.name}</Typography.Text>
                  <div>
                    <Typography.Text type="secondary">
                      {p.description}
                    </Typography.Text>
                  </div>
                </div>
              </Space>
            ),
          },
          { title: '价格', dataIndex: 'price', render: money },
          { title: '库存', dataIndex: 'stock' },
          {
            title: '购买',
            render: (_: unknown, p: Product) => (
              <Button
                disabled={p.stock < 1}
                type="primary"
                onClick={() =>
                  edit({
                    title: `加入购物车 · ${p.name}`,
                    initial: { quantity: 1 },
                    fields: [
                      {
                        name: 'quantity',
                        label: '购物车数量（覆盖原数量）',
                        type: 'number',
                        min: 1,
                        max: Math.min(p.stock, 999),
                      },
                    ],
                    submit: (v) =>
                      api('/cart', 'PUT', { ...v, productId: p.id }),
                  })
                }
              >
                加入购物车
              </Button>
            ),
          },
        ])}
      </>
    );
  else if (section === 'addresses')
    content = (
      <>
        <div className="market-toolbar">
          <Button
            type="primary"
            onClick={() =>
              edit({
                title: '新增收货地址',
                fields: [
                  { name: 'recipient', label: '收货人' },
                  { name: 'phone', label: '联系电话' },
                  { name: 'detail', label: '完整收货地址', type: 'textarea' },
                ],
                submit: (v) => api('/addresses', 'POST', v),
              })
            }
          >
            新增地址
          </Button>
        </div>
        {table([
          { title: '收货人', dataIndex: 'recipient' },
          { title: '电话', dataIndex: 'phone' },
          { title: '地址', dataIndex: 'detail' },
          {
            title: '操作',
            render: (_: unknown, a: Address) => (
              <Popconfirm
                title="删除此地址？已有订单地址不受影响。"
                onConfirm={() => run(() => api(`/addresses/${a.id}`, 'DELETE'))}
              >
                <Button danger>删除</Button>
              </Popconfirm>
            ),
          },
        ])}
      </>
    );
  else if (section === 'cart')
    content = (
      <>
        <Alert
          showIcon
          type="info"
          title="结算按店铺自动拆单，商品价格和库存以下单时服务器校验为准。"
          style={{ marginBottom: 16 }}
        />
        {table(
          [
            {
              title: '商品',
              render: (_: unknown, c: Cart) =>
                c.product?.name || '商品已不可用',
            },
            {
              title: '单价',
              render: (_: unknown, c: Cart) =>
                c.product ? money(c.product.price) : '—',
            },
            { title: '数量', dataIndex: 'quantity' },
            {
              title: '小计',
              render: (_: unknown, c: Cart) =>
                c.product ? money(c.product.price * c.quantity) : '—',
            },
            {
              title: '操作',
              render: (_: unknown, c: Cart) => (
                <Space>
                  <Button
                    onClick={() =>
                      edit({
                        title: '修改购物车数量',
                        initial: { quantity: c.quantity },
                        fields: [
                          {
                            name: 'quantity',
                            label: '数量',
                            type: 'number',
                            min: 1,
                            max: 999,
                          },
                        ],
                        submit: (v) =>
                          api('/cart', 'PUT', { ...v, productId: c.productId }),
                      })
                    }
                  >
                    修改
                  </Button>
                  <Popconfirm
                    title="移除此商品？"
                    onConfirm={() => run(() => api(`/cart/${c.id}`, 'DELETE'))}
                  >
                    <Button danger>移除</Button>
                  </Popconfirm>
                </Space>
              ),
            },
          ],
          {
            rowSelection: {
              selectedRowKeys: selected,
              onChange: setSelected,
              preserveSelectedRowKeys: false,
            },
          },
        )}
        <div className="market-toolbar">
          <Select
            placeholder="选择收货地址（请先在收货地址页添加）"
            value={addressId}
            onChange={setAddressId}
            style={{ minWidth: 340, maxWidth: '100%' }}
            options={addresses.map((a) => ({
              value: a.id,
              label: `${a.recipient} · ${a.detail}`,
            }))}
          />
          <Popconfirm
            title={`确认购买选中的 ${selected.length} 项商品？`}
            disabled={!addressId || !selected.length || busy}
            onConfirm={() =>
              run(async () => {
                const signature = JSON.stringify({
                  addressId,
                  selected: [...selected].sort(),
                });
                if (requestKey.current?.signature !== signature)
                  requestKey.current = { signature, key: crypto.randomUUID() };
                await api('/checkout', 'POST', {
                  addressId,
                  cartIds: selected,
                  requestKey: requestKey.current.key,
                });
                requestKey.current = undefined;
              })
            }
          >
            <Button
              type="primary"
              loading={busy}
              disabled={!addressId || !selected.length}
            >
              提交订单
            </Button>
          </Popconfirm>
        </div>
      </>
    );
  else if (section === 'orders')
    content = (
      <>
        <div className="market-toolbar">
          {scopePicker}
          <Typography.Text type="secondary">
            未付款订单 30 分钟后自动关闭并返还库存。
          </Typography.Text>
        </div>
        {table([
          { title: '订单编号', dataIndex: 'id', ellipsis: true },
          { title: '订单金额', dataIndex: 'amount', render: money },
          {
            title: '状态',
            dataIndex: 'status',
            render: (v: string) => <Status value={v} order />,
          },
          {
            title: '创建时间',
            dataIndex: 'createdAt',
            render: (v: string) => new Date(v).toLocaleString(),
          },
          {
            title: '操作',
            render: (_: unknown, o: Order) => (
              <Space wrap>
                <Button size="small" onClick={() => setDetail(o)}>
                  详情
                </Button>
                {scope === 'buyer' && o.status === 'pending' && (
                  <>
                    <Button
                      size="small"
                      disabled={!ctx?.sandboxPayment}
                      onClick={() => action(o, 'pay')}
                    >
                      沙箱付款
                    </Button>
                    <Button size="small" onClick={() => action(o, 'cancel')}>
                      取消
                    </Button>
                  </>
                )}
                {scope === 'merchant' && o.status === 'paid' && (
                  <Button
                    size="small"
                    onClick={() =>
                      action(o, 'ship', [
                        { name: 'carrier', label: '物流公司' },
                        { name: 'tracking', label: '运单号' },
                      ])
                    }
                  >
                    发货
                  </Button>
                )}
                {scope === 'buyer' && o.status === 'shipped' && (
                  <Button size="small" onClick={() => action(o, 'receive')}>
                    确认收货
                  </Button>
                )}
                {scope === 'buyer' &&
                  ['paid', 'shipped', 'completed'].includes(o.status) && (
                    <Button
                      size="small"
                      onClick={() =>
                        action(o, 'refund', [
                          {
                            name: 'reason',
                            label: '售后原因',
                            type: 'textarea',
                          },
                        ])
                      }
                    >
                      申请售后
                    </Button>
                  )}
                {scope !== 'buyer' && o.status === 'refunding' && (
                  <>
                    <Button
                      size="small"
                      disabled={!ctx?.sandboxPayment}
                      onClick={() => action(o, 'approve-refund')}
                    >
                      同意退款
                    </Button>
                    <Button
                      size="small"
                      onClick={() =>
                        action(o, 'reject-refund', [
                          {
                            name: 'reason',
                            label: '驳回原因',
                            type: 'textarea',
                          },
                        ])
                      }
                    >
                      驳回售后
                    </Button>
                  </>
                )}
              </Space>
            ),
          },
        ])}
      </>
    );
  else if (section === 'audit')
    content = ctx?.admin ? (
      table([
        { title: '时间', dataIndex: 'createdAt' },
        { title: '操作人', dataIndex: 'actorId' },
        { title: '动作', dataIndex: 'action' },
        { title: '对象', dataIndex: 'targetId' },
        { title: '说明', dataIndex: 'detail' },
      ])
    ) : (
      <Result status="403" title="仅平台管理员可查看审计记录" />
    );
  else content = <Empty />;
  return (
    <PageContainer
      title={titles[section] || '多商家商城'}
      subTitle="连接消费者、商家与平台"
      extra={<Button onClick={() => void load()}>刷新数据</Button>}
    >
      <div className="marketplace">
        <div className="market-banner">
          <div>
            <Typography.Title level={3} style={{ margin: 0 }}>
              让每一家好店，被看见。
            </Typography.Title>
            <Typography.Text type="secondary">
              MALL · 多商家经营平台
            </Typography.Text>
          </div>
          <Space>
            <Tag>
              {ctx?.admin
                ? '平台管理员'
                : ctx?.shops.some((s) => s.status === 'approved')
                  ? '入驻商家'
                  : '消费者'}
            </Tag>
            <Tag color={ctx?.sandboxPayment ? 'orange' : 'default'}>
              {ctx?.sandboxPayment
                ? '沙箱环境 · 不发生真实扣款'
                : '真实支付尚未接入'}
            </Tag>
          </Space>
        </div>
        {error ? (
          <Result
            status="error"
            title="数据加载失败"
            subTitle="请检查登录状态、访问权限及后端服务。页面不会显示模拟经营数据。"
            extra={<Button onClick={() => void load()}>重试</Button>}
          />
        ) : (
          content
        )}
      </div>
      <Modal
        title={editor?.title}
        open={!!editor}
        forceRender
        confirmLoading={busy}
        onCancel={() => !busy && setEditor(undefined)}
        onOk={async () => {
          try {
            const values = await form.validateFields();
            editor?.fields.forEach((field) => {
              if (field.type === 'number' && values[field.name] !== undefined) {
                values[field.name] = Number(values[field.name]);
              }
            });
            setBusy(true);
            await editor?.submit(values);
            setEditor(undefined);
            message.success('操作成功');
            await load();
          } catch {
            /* Keep the form open for corrections. */
          } finally {
            setBusy(false);
          }
        }}
      >
        <Form form={form} layout="vertical">
          {editor?.fields.map((f) => (
            <Form.Item
              key={f.name}
              name={f.name}
              label={f.label}
              initialValue={editor?.initial?.[f.name]}
              rules={[
                {
                  required: !f.optional,
                  // async-validator switches to string validation whenever
                  // `whitespace` is present, which fails numeric values.
                  ...(f.type !== 'number' && f.type !== 'select'
                    ? { whitespace: true }
                    : {}),
                  message: `请填写${f.label}`,
                },
              ]}
            >
              {f.type === 'number' ? (
                <Input
                  type="number"
                  min={f.min}
                  max={f.max}
                  step={f.name === 'price' ? '0.01' : '1'}
                />
              ) : f.type === 'select' ? (
                <Select options={f.options} />
              ) : f.type === 'textarea' ? (
                <Input.TextArea
                  rows={3}
                  maxLength={f.name === 'description' ? 10000 : 500}
                />
              ) : (
                <Input maxLength={f.name === 'image' ? 1000 : 200} />
              )}
            </Form.Item>
          ))}
        </Form>
        {editor && !editor.fields.length && (
          <Typography.Paragraph>请确认执行此操作。</Typography.Paragraph>
        )}
      </Modal>
      <Modal
        title="订单详情"
        open={!!detail}
        onCancel={() => setDetail(undefined)}
        footer={null}
        width={720}
      >
        {detail && (
          <>
            <Descriptions
              column={1}
              items={[
                { key: 'id', label: '订单编号', children: detail.id },
                { key: 'shop', label: '店铺编号', children: detail.shopId },
                {
                  key: 'address',
                  label: '收货信息',
                  children: `${detail.recipient} ${detail.phone} ${detail.address}`,
                },
                {
                  key: 'tracking',
                  label: '物流',
                  children: detail.tracking
                    ? `${detail.carrier} ${detail.tracking}`
                    : '尚未发货',
                },
                {
                  key: 'refund',
                  label: '售后原因',
                  children: detail.refundReason || '无',
                },
              ]}
            />
            <Table
              rowKey="id"
              dataSource={detail.items}
              pagination={false}
              columns={[
                { title: '商品', dataIndex: 'name' },
                { title: '单价', dataIndex: 'price', render: money },
                { title: '数量', dataIndex: 'quantity' },
              ]}
            />
            <Typography.Title level={4}>
              合计 {money(detail.amount)}
            </Typography.Title>
          </>
        )}
      </Modal>
    </PageContainer>
  );
}
