interface AvatarProps {
  name?: string;
  /** Optional image URL; initials fallback when absent. */
  src?: string;
  size?: 'sm' | 'md' | 'lg';
}

const sizeClasses = {
  sm: 'h-8 w-8 text-xs',
  md: 'h-10 w-10 text-sm',
  lg: 'h-12 w-12 text-base',
};

function initials(name?: string): string {
  if (!name) return '?';
  return name
    .split(/[\s@.]+/)
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join('');
}

/** Avatar — image with initials fallback. */
export function Avatar({ name, src, size = 'md' }: AvatarProps) {
  if (src) {
    return (
      <img
        src={src}
        alt={name ?? 'avatar'}
        className={`${sizeClasses[size]} rounded-full object-cover`}
      />
    );
  }
  return (
    <span
      className={`${sizeClasses[size]} inline-flex items-center justify-center rounded-full bg-blue-100 font-semibold text-blue-700`}
    >
      {initials(name)}
    </span>
  );
}
