import { Button, Result } from 'antd';
import type { ReactElement } from 'react';
import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';

const REDIRECT_SECONDS = 3;
const COUNTDOWN_STEP_SECONDS = 1;
const MIN_REDIRECT_SECONDS = 0;
const ONE_SECOND_MS = 1000;

export default function NotFoundPage(): ReactElement {
  const navigate = useNavigate();
  const [seconds, setSeconds] = useState(REDIRECT_SECONDS);

  useEffect(() => {
    const redirectTimer = window.setTimeout(() => {
      void navigate('/dashboard', { replace: true });
    }, REDIRECT_SECONDS * ONE_SECOND_MS);

    const countdownTimer = window.setInterval(() => {
      setSeconds((value) => Math.max(value - COUNTDOWN_STEP_SECONDS, MIN_REDIRECT_SECONDS));
    }, ONE_SECOND_MS);

    return () => {
      window.clearTimeout(redirectTimer);
      window.clearInterval(countdownTimer);
    };
  }, [navigate]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-slate-50 px-6">
      <Result
        extra={
          <Button
            onClick={() => {
              void navigate('/dashboard', { replace: true });
            }}
            type="primary"
          >
            返回工作台
          </Button>
        }
        status="404"
        subTitle={`${String(seconds)} 秒后自动跳转到工作台首页`}
        title="页面不存在"
      />
    </div>
  );
}
