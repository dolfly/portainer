import { PropsWithChildren } from 'react';

export function SortableListCard({ children }: PropsWithChildren<unknown>) {
  return (
    <div className="flex flex-col rounded-lg shadow-[0_1px_2px_rgba(16,24,40,.06),0_6px_16px_rgba(16,24,40,.08)]">
      {children}
    </div>
  );
}
