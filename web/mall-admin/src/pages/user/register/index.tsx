import { history, request } from '@umijs/max';
import { Button, Card, Form, Input, Space, Typography, message } from 'antd';
import React, { useEffect, useState } from 'react';
import { getCaptcha, type CaptchaData } from '../login/service';

export default function Register() {
  const [captcha, setCaptcha] = useState<CaptchaData>();
  const [busy, setBusy] = useState(false);
  const refresh = async () => {
    try {
      const result = await getCaptcha();
      setCaptcha(result.data);
    } catch {
      message.error('验证码加载失败，请重试');
    }
  };
  useEffect(() => {
    void refresh();
  }, []);
  return (
    <div
      style={{
        minHeight: '100vh',
        background: '#f4f8f7',
        display: 'grid',
        placeItems: 'center',
        padding: 24,
      }}
    >
      <Card style={{ width: 440, maxWidth: '100%' }}>
        <Typography.Title level={2}>加入 Mall</Typography.Title>
        <Typography.Paragraph type="secondary">
          创建账号即可购物，也可以在登录后申请商家入驻。
        </Typography.Paragraph>
        <Form
          layout="vertical"
          onFinish={async (v) => {
            setBusy(true);
            try {
              await request('/register', {
                method: 'POST',
                data: { ...v, captchaId: captcha?.id },
              });
              message.success('注册成功，请登录');
              history.push('/user/login');
            } catch {
              await refresh();
            } finally {
              setBusy(false);
            }
          }}
        >
          <Form.Item
            name="username"
            label="用户名"
            rules={[
              {
                required: true,
                pattern: /^[a-zA-Z][a-zA-Z0-9_]{2,31}$/,
                message: '3–32 位字母、数字、下划线，以字母开头',
              },
            ]}
          >
            <Input autoComplete="username" />
          </Form.Item>
          <Form.Item
            name="email"
            label="邮箱"
            rules={[{ required: true, type: 'email' }]}
          >
            <Input autoComplete="email" />
          </Form.Item>
          <Form.Item
            name="password"
            label="密码"
            rules={[
              {
                required: true,
                min: 12,
                max: 72,
                message: '请设置 12–72 位密码',
              },
            ]}
          >
            <Input.Password autoComplete="new-password" />
          </Form.Item>
          <Form.Item
            name="confirm"
            label="确认密码"
            dependencies={['password']}
            rules={[
              { required: true },
              ({ getFieldValue }) => ({
                validator: (_, v) =>
                  getFieldValue('password') === v
                    ? Promise.resolve()
                    : Promise.reject(new Error('两次输入的密码不一致')),
              }),
            ]}
          >
            <Input.Password autoComplete="new-password" />
          </Form.Item>
          <Form.Item
            name="captchaCode"
            label="图形验证码"
            rules={[{ required: true }]}
          >
            <Input maxLength={6} />
          </Form.Item>
          <Space style={{ marginBottom: 20 }}>
            {captcha && (
              <img src={captcha.image} alt="注册验证码" width={160} />
            )}
            <Button onClick={() => void refresh()}>刷新验证码</Button>
          </Space>
          <Button
            type="primary"
            htmlType="submit"
            block
            loading={busy}
            disabled={!captcha}
          >
            注册账号
          </Button>
        </Form>
        <div style={{ marginTop: 20, textAlign: 'center' }}>
          <a href="/user/login">已有账号，返回登录</a>
        </div>
      </Card>
    </div>
  );
}
