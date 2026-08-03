import type { ReactNode } from 'react';

type Color = 'success' | 'warning' | 'danger' | 'info';

const colorClasses: Record<Color, string> = {
  success: 'bg-green-100 text-green-800',
  warning: 'bg-amber-100 text-amber-800',
  danger: 'bg-red-100 text-red-800',
  info: 'bg-blue-100 text-blue-800',
};

/** Badge — small status label with color variants. */
export function Badge({ color = 'info', children }: { color?: Color; children: ReactNode }) {
  return (
    <span
      className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${colorClasses[color]}`}
    >
      {children}
    </span>
  );
}
