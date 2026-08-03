import {
  type ProjectCardData,
  RecommendationPanel,
  type StoryCardData,
} from '../features/practice-life';
import { RoleGate } from '../shared/guards/RoleGate';

const demoStories: StoryCardData[] = [
  {
    id: 'story-math-percent',
    title: 'Проценты в жизни',
    linkedModules: ['percent', 'math-5-11'],
    subjects: ['Математика'],
    realWorld: 'Аналитики считают процентные изменения курсов валют и инфляции.',
    readingMinutes: 3,
  },
  {
    id: 'story-chem-solutions',
    title: 'Растворы вокруг нас',
    linkedModules: ['solutions', 'chemistry'],
    subjects: ['Химия'],
    realWorld: 'Физраствор, уксус и морская вода — растворы разной концентрации.',
    readingMinutes: 4,
  },
  {
    id: 'story-bio-photosynthesis',
    title: 'Фотосинтез — химия в каждом листе',
    linkedModules: ['photosynthesis', 'chemistry'],
    subjects: ['Биология', 'Химия'],
    realWorld: 'Понимание фотосинтеза ведёт к созданию искусственного фотосинтеза для топлива.',
    readingMinutes: 5,
  },
];

const demoProjects: ProjectCardData[] = [
  {
    id: 'proj-math-bio-lab',
    title: 'Биохимическая лаборатория дома',
    modules: ['solutions', 'chemistry', 'cells'],
    difficultyLevel: 'medium',
    expectedOutcome: 'Провести серию опытов по концентрации растворов и описать результаты.',
  },
  {
    id: 'proj-flagship-budget',
    title: 'Семейный бюджет: математика и обществознание',
    modules: ['percent', 'taxes', 'inflation', 'economics'],
    difficultyLevel: 'basic',
    expectedOutcome: 'Составить и проанализировать семейный бюджет с учётом инфляции.',
  },
];

/** PracticePage — browse stories and project ideas (M2 F5.1-F5.3). */
export function PracticePage() {
  return (
    <RoleGate requiredRole={['learner', 'parent', 'teacher', 'methodologist', 'admin']}>
      <RecommendationPanel moduleName="Проценты" stories={demoStories} projects={demoProjects} />
    </RoleGate>
  );
}
