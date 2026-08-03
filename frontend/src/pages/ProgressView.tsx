import { Card } from '../shared/components/Card';
import { RoleGate } from '../shared/guards/RoleGate';

/** ProgressView — "Progress" placeholder (protected, learner+). */
export function ProgressView() {
  return (
    <RoleGate requiredRole={['learner', 'parent', 'teacher', 'methodologist', 'admin']}>
      <Card title="Progress">
        <p className="text-sm text-gray-600">
          Plan vs actual progress tracking is coming soon (M1).
        </p>
      </Card>
    </RoleGate>
  );
}
