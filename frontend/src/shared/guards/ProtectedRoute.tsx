import { Navigate, Outlet, useLocation } from 'react-router';
import { useAuthStore } from '../../store/authStore';

/** ProtectedRoute — redirects to /login when unauthenticated. */
export function ProtectedRoute() {
  const isAuthenticated = useAuthStore((s) => s.isAuthenticated);
  const location = useLocation();

  if (!isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }
  return <Outlet />;
}
