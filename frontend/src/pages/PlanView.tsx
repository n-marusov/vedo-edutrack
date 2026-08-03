import { Card } from '../shared/components/Card';
import { RoleGate } from '../shared/guards/RoleGate';

/** PlanView — "My Plan" placeholder (protected, learner+). */
export function PlanView() {
  return (
    <RoleGate requiredRole={['learner', 'parent', 'admin']}>
      <Card title="My Plan">
        <p className="text-sm text-gray-600">
          Your fixed learning plan with timeline is coming soon (M1).
        </p>
      </Card>
    </RoleGate>
  );
}
