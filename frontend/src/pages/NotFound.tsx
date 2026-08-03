import { Link } from 'react-router';

/** NotFound — 404 page with a link back. */
export function NotFoundPage() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50">
      <div className="text-center">
        <p className="text-6xl font-bold text-gray-300">404</p>
        <h1 className="mt-4 text-xl font-semibold text-gray-900">Page not found</h1>
        <p className="mt-2 text-sm text-gray-600">The page you're looking for doesn't exist.</p>
        <div className="mt-6 flex items-center justify-center gap-4">
          <Link to="/" className="text-sm font-medium text-blue-600 hover:text-blue-700">
            Go home
          </Link>
          <Link to="/dashboard" className="text-sm font-medium text-blue-600 hover:text-blue-700">
            Go to dashboard
          </Link>
        </div>
      </div>
    </div>
  );
}
