import { Badge } from '../../shared/components/Badge';
import { Card } from '../../shared/components/Card';
import { LoadingSpinner } from '../../shared/components/LoadingSpinner';
import { useCoverage } from '../gap-coverage/useCoverage';

export type PlanStatus = 'on-track' | 'at-risk' | 'off-track';

export interface PlanModuleProgress {
  module_id: string;
  planned_end: string;
  mastered: boolean;
  /** Deviation in days; positive = late. */
  deviationDays: number;
}

interface ProgressDashboardProps {
  learnerId: string;
  modules: PlanModuleProgress[];
  /** Override forecast (defaults to on-track). */
  forecast?: PlanStatus;
}

function forecastStatus(modules: PlanModuleProgress[]): PlanStatus {
  const late = modules.filter((module) => module.mastered && module.deviationDays > 0);
  if (late.length === 0) return 'on-track';
  const ratio = late.length / modules.length;
  if (ratio > 0.3) return 'off-track';
  return 'at-risk';
}

const statusBadge: Record<PlanStatus, 'success' | 'warning' | 'danger'> = {
  'on-track': 'success',
  'at-risk': 'warning',
  'off-track': 'danger',
};

/**
 * ProgressDashboard — plan-vs-actual comparison with deviation flags and a
 * binary readiness forecast (F2).
 */
export function ProgressDashboard({ learnerId, modules, forecast }: ProgressDashboardProps) {
  const coverage = useCoverage(learnerId);
  const status = forecast ?? forecastStatus(modules);

  return (
    <div className="space-y-4">
      <Card title="Plan vs actual" header={<Badge color={statusBadge[status]}>{status}</Badge>}>
        {coverage.loading && <LoadingSpinner label="Loading progress…" />}
        {coverage.error && (
          <p role="alert" className="text-sm text-red-700">
            {coverage.error}
          </p>
        )}
        <ul data-testid="progress-list" className="space-y-2">
          {modules.map((module) => (
            <li key={module.module_id} className="flex items-center gap-3 text-sm">
              <span className="flex-1 text-gray-800">{module.module_id}</span>
              {module.mastered ? (
                <Badge color="success">mastered</Badge>
              ) : (
                <Badge color="warning">pending</Badge>
              )}
              {module.mastered && module.deviationDays > 0 && (
                <Badge color="danger">+{module.deviationDays}d late</Badge>
              )}
            </li>
          ))}
        </ul>
      </Card>
    </div>
  );
}
