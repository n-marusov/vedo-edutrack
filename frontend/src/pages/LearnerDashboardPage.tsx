import { LearnerDashboard, type LearnerDashboardData } from '../features/visualization';
import { RoleGate } from '../shared/guards/RoleGate';

const demoData: LearnerDashboardData = {
  learnerName: 'Миша',
  currentModule: 'percent',
  currentStatus: 'in-progress',
  horizons: { far: 'chemistry', mid: 'solutions', near: 'percent' },
  planVsActual: { onTime: 2, delayed: 1, ahead: 0 },
  subjectProgress: [
    { subject: 'Математика', percent: 80 },
    { subject: 'Химия', percent: 40 },
    { subject: 'Биология', percent: 55 },
  ],
  fgosCoverage: 70,
  recommendedCount: 4,
  lastUpdated: '2026-08-03T12:00:00Z',
};

/** LearnerDashboardPage — role-aware learner dashboard (M2 F4.2). */
export function LearnerDashboardPage() {
  return (
    <RoleGate requiredRole={['learner', 'parent', 'teacher', 'methodologist', 'admin']}>
      <LearnerDashboard data={demoData} />
    </RoleGate>
  );
}
