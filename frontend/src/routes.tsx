import { Suspense, lazy } from 'react';
import { createBrowserRouter } from 'react-router';
import { NotFoundPage } from './pages';
import { LoadingSpinner } from './shared/components';
import { ProtectedRoute } from './shared/guards';
import { AuthLayout, LandingLayout, MainLayout } from './shared/layouts';

const LandingPage = lazy(() => import('./pages/Landing').then((m) => ({ default: m.LandingPage })));
const LoginPage = lazy(() => import('./pages/Login').then((m) => ({ default: m.LoginPage })));
const DashboardPage = lazy(() =>
  import('./pages/Dashboard').then((m) => ({ default: m.DashboardPage })),
);
const RouteView = lazy(() => import('./pages/RouteView').then((m) => ({ default: m.RouteView })));
const PlanView = lazy(() => import('./pages/PlanView').then((m) => ({ default: m.PlanView })));
const ProgressView = lazy(() =>
  import('./pages/ProgressView').then((m) => ({ default: m.ProgressView })),
);
const OntologyView = lazy(() =>
  import('./pages/OntologyView').then((m) => ({ default: m.OntologyView })),
);
const DemoPage = lazy(() => import('./pages/DemoPage').then((m) => ({ default: m.DemoPage })));
const LearnerDashboardPage = lazy(() =>
  import('./pages/LearnerDashboardPage').then((m) => ({ default: m.LearnerDashboardPage })),
);
const ParentDashboardPage = lazy(() =>
  import('./pages/ParentDashboardPage').then((m) => ({ default: m.ParentDashboardPage })),
);
const MethodologistDashboardPage = lazy(() =>
  import('./pages/MethodologistDashboardPage').then((m) => ({
    default: m.MethodologistDashboardPage,
  })),
);
const GapMapPage = lazy(() =>
  import('./pages/GapMapPage').then((m) => ({ default: m.GapMapPage })),
);
const GroupPanelPage = lazy(() =>
  import('./pages/GroupPanelPage').then((m) => ({ default: m.GroupPanelPage })),
);
const PracticePage = lazy(() =>
  import('./pages/PracticePage').then((m) => ({ default: m.PracticePage })),
);

function withSuspense(node: React.ReactNode) {
  return <Suspense fallback={<LoadingSpinner />}>{node}</Suspense>;
}

/**
 * SPA route tree (react-router v8):
 *
 *   /                   → Landing (public)
 *   /login              → Login (public)
 *   /dashboard          → DashboardLayout (protected, role-aware)
 *   /dashboard/route    → RouteView (protected)
 *   /dashboard/plan     → PlanView (protected)
 *   /dashboard/progress → ProgressView (protected)
 *   /*                  → NotFound
 */
export const router = createBrowserRouter([
  {
    path: '/',
    element: <LandingLayout />,
    children: [{ index: true, element: withSuspense(<LandingPage />) }],
  },
  {
    path: '/login',
    element: <AuthLayout />,
    children: [{ index: true, element: withSuspense(<LoginPage />) }],
  },
  {
    path: '/demo',
    element: <LandingLayout />,
    children: [{ index: true, element: withSuspense(<DemoPage />) }],
  },
  {
    path: '/dashboard',
    element: <ProtectedRoute />,
    children: [
      {
        element: <MainLayout />,
        children: [
          { index: true, element: withSuspense(<DashboardPage />) },
          { path: 'route', element: withSuspense(<RouteView />) },
          { path: 'plan', element: withSuspense(<PlanView />) },
          { path: 'progress', element: withSuspense(<ProgressView />) },
          { path: 'ontology', element: withSuspense(<OntologyView />) },
          { path: 'learner', element: withSuspense(<LearnerDashboardPage />) },
          { path: 'parent', element: withSuspense(<ParentDashboardPage />) },
          { path: 'methodologist', element: withSuspense(<MethodologistDashboardPage />) },
          { path: 'gaps', element: withSuspense(<GapMapPage />) },
          { path: 'group', element: withSuspense(<GroupPanelPage />) },
          { path: 'practice', element: withSuspense(<PracticePage />) },
        ],
      },
    ],
  },
  { path: '*', element: <NotFoundPage /> },
]);
