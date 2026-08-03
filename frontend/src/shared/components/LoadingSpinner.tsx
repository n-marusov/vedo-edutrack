/** LoadingSpinner — centered spinner with an optional label. */
export function LoadingSpinner({ label = 'Loading…' }: { label?: string }) {
  return (
    <div aria-live="polite" className="flex items-center justify-center gap-3 py-16">
      <span className="h-8 w-8 animate-spin rounded-full border-4 border-gray-300 border-t-blue-600" />
      <span className="text-sm text-gray-500">{label}</span>
    </div>
  );
}
