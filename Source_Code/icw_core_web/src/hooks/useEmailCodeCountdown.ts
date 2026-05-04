import { useEffect, useState } from 'react';

const ONE_SECOND_MS = 1000;
const INITIAL_SECONDS_LEFT = 0;
const MIN_SECONDS_LEFT = 0;
const COUNTDOWN_STEP_SECONDS = 1;

interface EmailCodeCountdown {
  buttonText: string;
  isCounting: boolean;
  startCountdown: (expiresIn: number) => void;
}

export function useEmailCodeCountdown(): EmailCodeCountdown {
  const [secondsLeft, setSecondsLeft] = useState(INITIAL_SECONDS_LEFT);

  useEffect(() => {
    if (secondsLeft <= MIN_SECONDS_LEFT) {
      return;
    }
    const timer = window.setTimeout(() => {
      setSecondsLeft((current) => Math.max(current - COUNTDOWN_STEP_SECONDS, MIN_SECONDS_LEFT));
    }, ONE_SECOND_MS);
    return () => {
      window.clearTimeout(timer);
    };
  }, [secondsLeft]);

  const startCountdown = (expiresIn: number): void => {
    setSecondsLeft(Math.max(Math.ceil(expiresIn), MIN_SECONDS_LEFT));
  };

  return {
    buttonText: secondsLeft > MIN_SECONDS_LEFT ? `${String(secondsLeft)}s 后重试` : '发送验证码',
    isCounting: secondsLeft > MIN_SECONDS_LEFT,
    startCountdown,
  };
}
