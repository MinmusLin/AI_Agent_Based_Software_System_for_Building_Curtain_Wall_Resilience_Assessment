import { LockOutlined, MailOutlined, SafetyCertificateOutlined } from '@ant-design/icons';
import { App, Button, Form, Input, Segmented } from 'antd';
import type { ReactElement } from 'react';
import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';

import { sendEmailCode } from '@/api/auth';
import { getErrorMessage } from '@/api/http';
import { AuthShell } from '@/components/AuthShell';
import { useAuth } from '@/contexts/AuthContext';
import { EmailCodeScene_Value, LoginScene_Value } from '@/gen/core/common';
import { useEmailCodeCountdown } from '@/hooks/useEmailCodeCountdown';
import {
  EMAIL_MAX_LENGTH,
  normalizeEmailAddress,
  normalizeEmailCode,
  normalizeNonWhitespaceAscii,
} from '@/utils/validation';

interface LoginFormValues {
  email: string;
  password?: string;
  email_code?: string;
}

export default function LoginPage(): ReactElement {
  const { message } = App.useApp();
  const navigate = useNavigate();
  const { login } = useAuth();
  const [mode, setMode] = useState<LoginScene_Value>(LoginScene_Value.Password);
  const [loading, setLoading] = useState(false);
  const [sending, setSending] = useState(false);
  const { buttonText, isCounting, startCountdown } = useEmailCodeCountdown();
  const [form] = Form.useForm<LoginFormValues>();

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
        scene: EmailCodeScene_Value.Login,
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

  const submit = (values: LoginFormValues): void => {
    void handleSubmit(values);
  };

  async function handleSubmit(values: LoginFormValues): Promise<void> {
    setLoading(true);
    try {
      const email = normalizeEmailAddress(values.email);
      await login({
        email,
        scene: mode,
        code: mode === LoginScene_Value.Password ? (values.password ?? '') : (values.email_code ?? ''),
      });
      void navigate('/dashboard', { replace: true });
    } catch (error: unknown) {
      message.error(getErrorMessage(error));
    } finally {
      setLoading(false);
    }
  }

  return (
    <AuthShell subtitle="欢迎使用建筑幕墙韧性评估软件系统" title="登录账号">
      <div className="mb-4">
        <Segmented
          block
          onChange={(value) => {
            setMode(value === LoginScene_Value.Email ? LoginScene_Value.Email : LoginScene_Value.Password);
          }}
          options={[
            { label: '密码登录', value: LoginScene_Value.Password },
            { label: '邮箱验证码登录', value: LoginScene_Value.Email },
          ]}
          value={mode}
        />
      </div>
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
        {mode === LoginScene_Value.Password ? (
          <Form.Item
            label="密码"
            name="password"
            normalize={normalizeNonWhitespaceAscii}
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password placeholder="请输入密码" prefix={<LockOutlined />} size="large" />
          </Form.Item>
        ) : (
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
        )}
        <div className="mb-5 flex items-center justify-between text-sm">
          <Link to="/register">注册账号</Link>
          <Link to="/forget-password">忘记密码</Link>
        </div>
        <Button block htmlType="submit" loading={loading} size="large" type="primary">
          登录
        </Button>
      </Form>
    </AuthShell>
  );
}
