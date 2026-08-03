import { useState } from 'react';
import { Link, NavLink, Outlet, useNavigate } from 'react-router';
import { useAuthStore } from '../../store/authStore';
import { Avatar } from '../components/Avatar';

const navItems = [
  { to: '/dashboard', label: 'Dashboard', roles: [] },
  {
    to: '/dashboard/route',
    label: 'Route Builder',
    roles: ['learner', 'parent', 'teacher', 'admin'],
  },
  { to: '/dashboard/plan', label: 'My Plan', roles: ['learner', 'parent', 'admin'] },
  {
    to: '/dashboard/progress',
    label: 'Progress',
    roles: ['learner', 'parent', 'teacher', 'methodologist', 'admin'],
  },
];

/** MainLayout — sidebar + header + content (Outlet). */
export function MainLayout() {
  const [collapsed, setCollapsed] = useState(false);
  const user = useAuthStore((s) => s.user);
  const roles = useAuthStore((s) => s.roles);
  const logout = useAuthStore((s) => s.logout);
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  const visible = navItems.filter(
    (item) => item.roles.length === 0 || item.roles.some((r) => roles.includes(r)),
  );

  return (
    <div className="flex min-h-screen bg-gray-50">
      <aside
        className={`flex flex-col border-r border-gray-200 bg-white transition-all ${collapsed ? 'w-16' : 'w-60'}`}
      >
        <div className="flex h-16 items-center gap-2 border-b border-gray-100 px-4">
          <span className="text-xl font-bold text-blue-600">
            {collapsed ? 'VE' : 'VEDO EduTrack'}
          </span>
        </div>
        <nav className="flex-1 space-y-1 p-3">
          {visible.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.to === '/dashboard'}
              className={({ isActive }) =>
                `block rounded-md px-3 py-2 text-sm font-medium ${
                  isActive ? 'bg-blue-50 text-blue-700' : 'text-gray-700 hover:bg-gray-100'
                } ${collapsed ? 'text-center' : ''}`
              }
            >
              {collapsed ? item.label.charAt(0) : item.label}
            </NavLink>
          ))}
        </nav>
        <div className="border-t border-gray-100 p-3">
          <button
            type="button"
            onClick={() => setCollapsed((c) => !c)}
            className="w-full rounded-md px-3 py-2 text-left text-xs text-gray-500 hover:bg-gray-100"
          >
            {collapsed ? '»' : '« Collapse'}
          </button>
        </div>
      </aside>

      <div className="flex min-w-0 flex-1 flex-col">
        <header className="flex h-16 items-center justify-between border-b border-gray-200 bg-white px-6">
          <h1 className="text-sm font-medium text-gray-900">VEDO EduTrack</h1>
          <div className="flex items-center gap-3">
            <Link to="/dashboard" className="flex items-center gap-2 text-sm text-gray-700">
              <Avatar name={user?.userId} size="sm" />
              <span>{user?.userId ?? 'guest'}</span>
            </Link>
            <button
              type="button"
              onClick={handleLogout}
              className="rounded-md px-3 py-1.5 text-sm text-gray-600 hover:bg-gray-100"
            >
              Sign out
            </button>
          </div>
        </header>
        <main className="flex-1 overflow-y-auto p-6">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
