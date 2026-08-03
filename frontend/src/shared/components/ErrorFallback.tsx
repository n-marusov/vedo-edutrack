import { useEffect } from 'react';

interface ErrorFallbackProps {
  error?: Error;
  resetErrorBoundary?: () => void;
}

/** ErrorFallback — "Something went wrong" UI with a retry button. */
export function ErrorFallback({ error, resetErrorBoundary }: ErrorFallbackProps) {
  useEffect(() => {
    if (error) {
      console.error('[app] render error', error);
    }
  }, [error]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50">
      <div className="max-w-md rounded-lg bg-white p-8 text-center shadow">
        <h1 className="text-xl font-semibold text-gray-900">Something went wrong</h1>
        <p className="mt-2 text-sm text-gray-600">
          An unexpected error occurred. Please try again.
        </p>
        {error && (
          <p className="mt-3 truncate rounded bg-gray-100 px-3 py-2 font-mono text-xs text-gray-500">
            {error.message}
          </p>
        )}
        {resetErrorBoundary && (
          <button
            type="button"
            onClick={resetErrorBoundary}
            className="mt-6 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
          >
            Try again
          </button>
        )}
      </div>
    </div>
  );
}
