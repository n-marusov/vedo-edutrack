import { Badge } from '../shared/components/Badge';
import { Card } from '../shared/components/Card';
import { useAuthStore } from '../store/authStore';

const stats = [
  { label: 'Active Plans', value: '2' },
  { label: 'Completed Modules', value: '14' },
  { label: 'Coverage %', value: '68%' },
];

/** Dashboard — role-aware home for authenticated users (stub data). */
export function DashboardPage() {
  const user = useAuthStore((s) => s.user);
  const roles = useAuthStore((s) => s.roles);
  const primaryRole = roles[0] ?? 'learner';

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <h1 className="text-2xl font-semibold text-gray-900">Welcome, {user?.userId}</h1>
        <Badge color="info">{primaryRole}</Badge>
      </div>

      <div className="grid gap-6 md:grid-cols-3">
        {stats.map((s) => (
          <Card key={s.label} title={s.label}>
            <p className="text-3xl font-semibold text-gray-900">{s.value}</p>
          </Card>
        ))}
      </div>

      <Card title={`${primaryRole} dashboard (stub)`}>
        <p className="text-sm text-gray-600">
          {primaryRole === 'learner' && 'Your route, plan, progress and gaps will appear here.'}
          {primaryRole === 'parent' && "Your children's progress overview will appear here."}
          {primaryRole === 'teacher' && 'Class / group management will appear here.'}
          {primaryRole === 'methodologist' && 'FGOS coverage will appear here.'}
          {primaryRole === 'admin' && 'User management will appear here.'}
        </p>
      </Card>
    </div>
  );
}
