import { GapMap, type RootGap } from '../features/visualization';
import { RoleGate } from '../shared/guards/RoleGate';

const demoGaps: RootGap[] = [
  {
    moduleId: 'percent',
    title: 'Проценты',
    blockedModules: 4,
    blockedSubjects: ['Химия', 'Биология', 'Обществознание'],
    rank: 1,
  },
  {
    moduleId: 'fractions',
    title: 'Дроби',
    blockedModules: 1,
    blockedSubjects: ['Математика'],
    rank: 2,
  },
];

/** GapMapPage — role-aware root-gap diagnostic map (M2 F4.6). */
export function GapMapPage() {
  return (
    <RoleGate requiredRole={['learner', 'parent', 'teacher', 'methodologist', 'admin']}>
      <GapMap gaps={demoGaps} />
    </RoleGate>
  );
}
