import type { ReactElement } from 'react';

interface ModalTitleProps {
  text: string;
}

export function ModalTitle({ text }: ModalTitleProps): ReactElement {
  return (
    <div className="block max-w-[calc(100%-16px)] truncate pr-2" title={text}>
      {text}
    </div>
  );
}
