import { Alert, Button, Card, Form, Input } from 'antd';
import { useState } from 'react';
import { getCaptcha, login, registerUrl, safeRedirect, type Captcha, type LoginCredentials } from '../service';

export default function LoginPage() {
  const [form] = Form.useForm<LoginCredentials>();
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [captcha, setCaptcha] = useState<Captcha>();
  const [failedAttempts, setFailedAttempts] = useState(0);
  const submit = async (values: LoginCredentials) => {
    setLoading(true); setError('');
    try {
      await login({ ...values, ...(captcha ? { captcha_id: captcha.id } : {}) });
      setFailedAttempts(0);
      window.location.assign(safeRedirect(new URLSearchParams(window.location.search).get('redirect')));
    } catch (err) {
      const text = err instanceof Error ? err.message : '登录失败';
      setError(text);
      const nextAttempt = failedAttempts + 1;
      setFailedAttempts(nextAttempt);
      if (captcha) form.setFields([{ name: 'captcha_code', errors: ['登录失败，请重新输入验证码'] }]);
      if (nextAttempt >= 2) {
        void getCaptcha().then(setCaptcha).catch(() => setError('验证码加载失败，请重试'));
      }
    } finally { setLoading(false); }
  };
  return <main className="account-page"><Card className="account-card" bordered={false}>
    <a className="brand" href="/">MALL<span>.</span></a><h1>欢迎回来</h1><p>登录后继续你的购物旅程</p>
    {error && <Alert type="error" showIcon message={error} />}
    <Form form={form} layout="vertical" onFinish={submit}>
      <Form.Item name="identifier" label="用户名或邮箱" rules={[{ required: true }]}><Input autoComplete="username" size="large" /></Form.Item>
      <Form.Item name="password" label="密码" rules={[{ required: true }]}><Input.Password autoComplete="current-password" size="large" /></Form.Item>
      {captcha && <Form.Item name="captcha_code" label="验证码" rules={[{ required: true }]}><Input size="large" addonAfter={<button className="captcha-refresh" type="button" aria-label="刷新验证码" onClick={() => void getCaptcha().then(setCaptcha)}><img src={captcha.image} alt="图形验证码" width="100" height="34" /></button>} /></Form.Item>}
      <Button type="primary" htmlType="submit" size="large" block loading={loading}>登录</Button>
    </Form>
    <div className="account-switch">还没有账户？ <a href={registerUrl(new URLSearchParams(window.location.search).get('redirect') ?? '/')}>创建账户</a></div>
  </Card></main>;
}
