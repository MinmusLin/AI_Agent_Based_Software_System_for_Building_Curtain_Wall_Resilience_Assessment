import { useEffect, useState } from 'react';

const ONE_SECOND_MS = 1000;

interface EmailCodeCountdown {
  buttonText: string;
  isCounting: boolean;
  startCountdown: (expiresIn: number) => void;
}

export function useEmailCodeCountdown(): EmailCodeCountdown {
  const [secondsLeft, setSecondsLeft] = useState(0);

  useEffect(() => {
    if (secondsLeft <= 0) {
      return;
    }
    const timer = window.setTimeout(() => {
      setSecondsLeft((current) => Math.max(current - 1, 0));
    }, ONE_SECOND_MS);
    return () => {
      window.clearTimeout(timer);
    };
  }, [secondsLeft]);

  const startCountdown = (expiresIn: number): void => {
    setSecondsLeft(Math.max(Math.ceil(expiresIn), 0));
  };

  return {
    buttonText: secondsLeft > 0 ? `${String(secondsLeft)}s 后重试` : '发送验证码',
    isCounting: secondsLeft > 0,
    startCountdown,
  };
}
