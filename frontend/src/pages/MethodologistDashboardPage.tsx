import { MethodologistDashboard, type MethodologistDashboardData } from '../features/visualization';
import { RoleGate } from '../shared/guards/RoleGate';

const demoData: MethodologistDashboardData = {
  schoolName: 'Школа «Ковчег»',
  classes: [
    { className: '7А', learnerCount: 24, fgosCoverage: 72 },
    { className: '7Б', learnerCount: 26, fgosCoverage: 68 },
    { className: '8А', learnerCount: 30, fgosCoverage: 81 },
  ],
  laggingTopics: [
    { moduleId: 'percent', title: 'Проценты', learnersWithGap: 12, subject: 'Математика' },
    { moduleId: 'solutions', title: 'Растворы', learnersWithGap: 9, subject: 'Химия' },
  ],
  ontologyContributions: 14,
};

/** MethodologistDashboardPage — role-aware methodologist dashboard (M2 F4.5). */
export function MethodologistDashboardPage() {
  return (
    <RoleGate requiredRole={['methodologist', 'admin']}>
      <MethodologistDashboard data={demoData} />
    </RoleGate>
  );
}
