import { LockOutlined, MailOutlined, SafetyCertificateOutlined, UserOutlined } from '@ant-design/icons';
import { App, Button, Form, Input } from 'antd';
import type { ReactElement } from 'react';
import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { register, sendEmailCode } from '@/api/auth';
import { getErrorMessage } from '@/api/http';
import { AuthShell } from '@/components/AuthShell';
import type { RegisterRequest } from '@/gen/core/api/auth';
import { EmailCodeScene_Value } from '@/gen/core/common';
import { useEmailCodeCountdown } from '@/hooks/useEmailCodeCountdown';
import {
  EMAIL_MAX_LENGTH,
  normalizeEmailAddress,
  normalizeEmailCode,
  normalizeNonWhitespaceAscii,
  passwordRules,
} from '@/utils/validation';

interface RegisterFormValues extends RegisterRequest {
  confirm_password?: string;
}

export default function RegisterPage(): ReactElement {
  const { message } = App.useApp();
  const navigate = useNavigate();
  const [form] = Form.useForm<RegisterFormValues>();
  const [loading, setLoading] = useState(false);
  const [sending, setSending] = useState(false);
  const { buttonText, isCounting, startCountdown } = useEmailCodeCountdown();

  async function handleSendCode(): Promise<void> {
    const email = normalizeEmailAddress(String(form.getFieldValue('email') ?? ''));
    if (!email) {
      message.warning('请输入邮箱');
      return;
    }
    setSending(true);
    try {
      const result = await sendEmailCode({
        email,
        scene: EmailCodeScene_Value.Register,
      });
      startCountdown(result.expires_in);
      message.success('已发送邮箱验证码');
    } catch (error: unknown) {
      message.error(getErrorMessage(error));
    } finally {
      setSending(false);
    }
  }

  const sendCode = (): void => {
    void handleSendCode();
  };

  const submit = (values: RegisterFormValues): void => {
    void handleSubmit(values);
  };

  async function handleSubmit(values: RegisterFormValues): Promise<void> {
    setLoading(true);
    try {
      const email = normalizeEmailAddress(values.email);
      await register({
        email,
        email_code: values.email_code,
        password: values.password,
        name: values.name,
      });
      message.success('账号注册成功，请登录');
      void navigate('/login', { replace: true });
    } catch (error: unknown) {
      message.error(getErrorMessage(error));
    } finally {
      setLoading(false);
    }
  }

  return (
    <AuthShell subtitle="欢迎使用建筑幕墙韧性评估软件系统" title="注册账号">
      <Form form={form} layout="vertical" onFinish={submit} requiredMark={false}>
        <Form.Item
          label="邮箱"
          name="email"
          normalize={normalizeEmailAddress}
          rules={[
            { required: true, message: '请输入邮箱' },
            { type: 'email', message: '邮箱格式错误' },
          ]}
        >
          <Input maxLength={EMAIL_MAX_LENGTH} placeholder="user@example.com" prefix={<MailOutlined />} size="large" />
        </Form.Item>
        <Form.Item label="验证码" required>
          <div className="flex gap-2">
            <Form.Item
              name="email_code"
              normalize={normalizeEmailCode}
              noStyle
              rules={[
                { required: true, message: '请输入 6 位数字验证码' },
                { pattern: /^\d+$/, message: '请输入 6 位数字验证码' },
                { len: 6, message: '请输入 6 位数字验证码' },
              ]}
            >
              <Input
                autoComplete="one-time-code"
                inputMode="numeric"
                maxLength={6}
                placeholder="邮箱验证码"
                prefix={<SafetyCertificateOutlined />}
                size="large"
              />
            </Form.Item>
            <Button disabled={isCounting} loading={sending} onClick={sendCode} size="large">
              {buttonText}
            </Button>
          </div>
        </Form.Item>
        <Form.Item
          label="用户名称"
          name="name"
          rules={[
            { required: true, message: '请输入用户名称' },
            { max: 8, message: '用户名称不能超过 8 个字符' },
          ]}
        >
          <Input maxLength={8} placeholder="不能超过 8 个字符" prefix={<UserOutlined />} size="large" />
        </Form.Item>
        <Form.Item
          label="密码"
          name="password"
          normalize={normalizeNonWhitespaceAscii}
          rules={passwordRules('请输入密码')}
        >
          <Input.Password maxLength={24} placeholder="请输入密码" prefix={<LockOutlined />} size="large" />
        </Form.Item>
        <Form.Item
          dependencies={['password']}
          label="确认密码"
          name="confirm_password"
          normalize={normalizeNonWhitespaceAscii}
          rules={[
            { required: true, message: '请再次输入密码' },
            ({ getFieldValue }) => ({
              validator(_, value) {
                if (!value || getFieldValue('password') === value) {
                  return Promise.resolve();
                }
                return Promise.reject(new Error('两次输入的密码不一致'));
              },
            }),
          ]}
        >
          <Input.Password placeholder="请再次输入密码" prefix={<LockOutlined />} size="large" />
        </Form.Item>
        <div className="mb-5 text-sm">
          已有账号？<Link to="/login">前往登录</Link>
        </div>
        <Button block htmlType="submit" loading={loading} size="large" type="primary">
          注册
        </Button>
      </Form>
    </AuthShell>
  );
}
