import { Link, Outlet } from 'react-router';

/** AuthLayout — centered card layout for login/register pages. */
export function AuthLayout() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4">
      <div className="w-full max-w-md">
        <div className="mb-8 text-center">
          <Link to="/" className="text-2xl font-bold text-blue-600">
            VEDO EduTrack
          </Link>
          <p className="mt-1 text-sm text-gray-500">Персональная образовательная траектория</p>
        </div>
        <div className="rounded-lg bg-white p-8 shadow">
          <Outlet />
        </div>
      </div>
    </div>
  );
}
