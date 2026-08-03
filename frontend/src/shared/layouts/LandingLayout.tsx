import { Outlet } from 'react-router';

/** LandingLayout — full-width, no sidebar (landing page). */
export function LandingLayout() {
  return (
    <div className="min-h-screen bg-white">
      <Outlet />
    </div>
  );
}
