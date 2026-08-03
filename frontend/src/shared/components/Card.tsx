import type { ReactNode } from 'react';

interface CardProps {
  title?: string;
  header?: ReactNode;
  footer?: ReactNode;
  className?: string;
  children: ReactNode;
}

/** Card — padded container with shadow, rounded corners, optional header/footer. */
export function Card({ title, header, footer, className = '', children }: CardProps) {
  return (
    <div className={`rounded-lg border border-gray-200 bg-white shadow-sm ${className}`}>
      {(title || header) && (
        <div className="border-b border-gray-100 px-5 py-4">
          {header ?? <h3 className="text-sm font-semibold text-gray-900">{title}</h3>}
        </div>
      )}
      <div className="px-5 py-4">{children}</div>
      {footer && <div className="border-t border-gray-100 px-5 py-3">{footer}</div>}
    </div>
  );
}
