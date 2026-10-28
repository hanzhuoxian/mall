import { Alert, Button, Card, Form, Input } from 'antd';
import { useCallback, useEffect, useState } from 'react';
import { getCaptcha, loginUrl, register, safeRedirect, type Captcha } from '../service';

export default function RegisterPage() {
  const [captcha, setCaptcha] = useState<Captcha>();
  const [error, setError] = useState('');
  const [done, setDone] = useState(false);
  const [loading, setLoading] = useState(false);
  const returnTo = safeRedirect(new URLSearchParams(window.location.search).get('redirect'));
  const loadCaptcha = useCallback(() => getCaptcha().then(setCaptcha).catch(() => setError('验证码加载失败，请刷新页面重试')), []);
  useEffect(() => { void loadCaptcha(); }, [loadCaptcha]);
  const submit = async (values: { username: string; email: string; password: string; captchaCode: string }) => {
    if (!captcha) return;
    setLoading(true); setError('');
    try {
      await register({ ...values, captchaId: captcha.id });
      setDone(true);
    } catch (err) { setError(err instanceof Error ? err.message : '注册失败'); void loadCaptcha(); }
    finally { setLoading(false); }
  };
  return <main className="account-page"><Card className="account-card" bordered={false}>
    <a className="brand" href="/">MALL<span>.</span></a><h1>创建账户</h1><p>注册成功后，请使用新账户登录</p>
    {done ? <><Alert type="success" showIcon message="注册成功" description="账户已创建，请登录后开始购物。" /><Button type="primary" block size="large" href={loginUrl(returnTo)} style={{ marginTop: 20 }}>前往登录</Button></> : <>
      {error && <Alert type="error" showIcon message={error} />}
      <Form layout="vertical" onFinish={submit}>
        <Form.Item name="username" label="用户名" rules={[{ required: true }, { min: 3, max: 32 }]}><Input autoComplete="username" size="large" /></Form.Item>
        <Form.Item name="email" label="邮箱" rules={[{ required: true, type: 'email' }]}><Input autoComplete="email" size="large" /></Form.Item>
        <Form.Item name="password" label="密码（至少 12 位）" rules={[{ required: true }, { min: 12 }]}><Input.Password autoComplete="new-password" size="large" /></Form.Item>
        <Form.Item name="captchaCode" label="图形验证码" rules={[{ required: true }]}>
          <Input size="large" addonAfter={captcha && <button className="captcha-refresh" type="button" aria-label="刷新验证码" onClick={() => void loadCaptcha()}><img src={captcha.image} alt="图形验证码" width="110" height="36" /></button>} />
        </Form.Item>
        <Button type="primary" htmlType="submit" size="large" block loading={loading} disabled={!captcha}>创建账户</Button>
      </Form>
      <div className="account-switch">已有账户？ <a href={loginUrl(returnTo)}>立即登录</a></div>
    </>}
  </Card></main>;
}
