import { RouteBuilder, useRouteCompute } from '../features/route-planning';
import { Card } from '../shared/components/Card';
import { RoleGate } from '../shared/guards/RoleGate';

/** RouteView — "Route Builder" (protected, learner+). */
export function RouteView() {
  const route = useRouteCompute('cli-user');
  return (
    <RoleGate requiredRole={['learner', 'parent', 'teacher', 'admin']}>
      <div className="space-y-4">
        <h1 className="text-2xl font-semibold text-gray-900">Route Builder</h1>
        <RouteBuilder route={route} />
        <Card title="Plan fixation">
          <p className="text-sm text-gray-600">
            Review the computed route, then fix it as an immutable learning plan (F1 plan fixation).
          </p>
        </Card>
      </div>
    </RoleGate>
  );
}
