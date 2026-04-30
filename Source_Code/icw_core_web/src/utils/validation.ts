import type { Rule } from 'antd/es/form';

export function normalizeEmailCode(value: unknown): string {
  if (typeof value !== 'string' && typeof value !== 'number') {
    return '';
  }
  return String(value).replace(/\D/g, '');
}

export function passwordRules(requiredMessage: string): Rule[] {
  return [
    {
      required: true,
      message: requiredMessage,
    },
    {
      min: 8,
      message: '密码长度必须不小于 8 位',
    },
    {
      pattern: /^(?=.*\p{Ll})(?=.*\p{Lu})(?=.*\p{Nd})(?=.*[\p{P}\p{S}]).+$/u,
      message: '密码必须同时包含大小写字母、数字和符号',
    },
  ];
}
