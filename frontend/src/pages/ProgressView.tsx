import { ProgressDashboard } from '../features/execution-progress';
import type { PlanModuleProgress } from '../features/execution-progress';
import { GapCoverage } from '../features/gap-coverage';
import { ResourceCatalog } from '../features/resources';
import { RoleGate } from '../shared/guards/RoleGate';

const demoModules: PlanModuleProgress[] = [
  { module_id: 'math-5-1', planned_end: '2026-09-01', mastered: true, deviationDays: 0 },
  { module_id: 'math-5-2', planned_end: '2026-09-15', mastered: true, deviationDays: 3 },
  { module_id: 'math-5-3', planned_end: '2026-10-01', mastered: false, deviationDays: 0 },
];

/** ProgressView — plan-vs-actual, gap diagnosis, FGOS coverage, resources. */
export function ProgressView() {
  return (
    <RoleGate requiredRole={['learner', 'parent', 'teacher', 'methodologist', 'admin']}>
      <div className="space-y-6">
        <h1 className="text-2xl font-semibold text-gray-900">Progress & coverage</h1>
        <ProgressDashboard learnerId="cli-user" modules={demoModules} />
        <GapCoverage learnerId="cli-user" lagModuleId="chemistry" />
        <ResourceCatalog />
      </div>
    </RoleGate>
  );
}
