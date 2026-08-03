import { Card } from '../shared/components/Card';
import { RoleGate } from '../shared/guards/RoleGate';

/** RouteView — "Route Builder" placeholder (protected, learner+). */
export function RouteView() {
  return (
    <RoleGate requiredRole={['learner', 'parent', 'teacher', 'admin']}>
      <Card title="Route Builder">
        <p className="text-sm text-gray-600">
          Visual route construction over the knowledge graph is coming soon (M1).
        </p>
      </Card>
    </RoleGate>
  );
}
