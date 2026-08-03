import {
  type ChildOverview,
  ParentDashboard,
  type ParentDashboardData,
} from '../features/visualization';
import type { LearnerDashboardData } from '../features/visualization';
import { RoleGate } from '../shared/guards/RoleGate';

const demoChildren: ChildOverview[] = [
  {
    id: 'demo-misha',
    name: 'Миша',
    fgosCoverage: 70,
    forecastStatus: 'at-risk',
    lagPercent: 15,
    lagModule: 'percent',
    recommendedCount: 4,
  },
  {
    id: 'demo-katya',
    name: 'Катя',
    fgosCoverage: 90,
    forecastStatus: 'on-track',
    lagPercent: 2,
    lagModule: '',
    recommendedCount: 2,
  },
];

const demoChildDetail: Record<string, LearnerDashboardData> = {
  'demo-misha': {
    learnerName: 'Миша',
    currentModule: 'percent',
    currentStatus: 'in-progress',
    horizons: { far: 'chemistry', mid: 'solutions', near: 'percent' },
    planVsActual: { onTime: 2, delayed: 1, ahead: 0 },
    subjectProgress: [
      { subject: 'Математика', percent: 80 },
      { subject: 'Химия', percent: 40 },
    ],
    fgosCoverage: 70,
    recommendedCount: 4,
    lastUpdated: '2026-08-03T12:00:00Z',
  },
  'demo-katya': {
    learnerName: 'Катя',
    currentModule: 'fractions',
    currentStatus: 'mastered',
    horizons: { far: 'algebra', mid: 'equations', near: 'fractions' },
    planVsActual: { onTime: 3, delayed: 0, ahead: 1 },
    subjectProgress: [
      { subject: 'Математика', percent: 90 },
      { subject: 'Литература', percent: 85 },
    ],
    fgosCoverage: 90,
    recommendedCount: 2,
    lastUpdated: '2026-08-03T12:00:00Z',
  },
};

const demoData: ParentDashboardData = {
  parentName: 'Марина',
  children: demoChildren,
  childDetail: demoChildDetail,
};

/** ParentDashboardPage — role-aware parent dashboard with 2+ children (M2 F4.4). */
export function ParentDashboardPage() {
  return (
    <RoleGate requiredRole={['parent', 'admin']}>
      <ParentDashboard data={demoData} />
    </RoleGate>
  );
}
