/**
 * DemoPage — guided product validation flow for M2 (Family Education).
 *
 * Walks through the family education journey end-to-end without custom
 * customer setup: pick a learner → build a route → learner dashboard →
 * knowledge map → FGOS coverage & gaps → stories & projects.
 *
 * Uses the already-built M2 components with demo fixture data.
 */
import { useMemo, useState } from 'react';
import {
  type ProjectCardData,
  RecommendationPanel,
  type StoryCardData,
} from '../features/practice-life';
import { RouteBuilder, useRouteCompute } from '../features/route-planning';
import {
  GapMap,
  type GraphEdge,
  type GraphNode,
  KnowledgeMap,
  LearnerDashboard,
  type LearnerDashboardData,
  type ModuleStatus,
  type RootGap,
} from '../features/visualization';
import { Badge, Button, Card } from '../shared/components';

interface DemoChild {
  id: string;
  name: string;
  grade: string;
}

const demoChildren: DemoChild[] = [
  { id: 'demo-misha', name: 'Миша', grade: '7 класс' },
  { id: 'demo-katya', name: 'Катя', grade: '5 класс' },
];

const demoDashboard: LearnerDashboardData = {
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

const demoGraph: GraphNode[] = [
  { id: 'percent', title: 'Проценты', subject: 'Математика' },
  { id: 'solutions', title: 'Растворы', subject: 'Химия' },
  { id: 'chemistry', title: 'Химия', subject: 'Химия' },
  { id: 'photosynthesis', title: 'Фотосинтез', subject: 'Биология' },
  { id: 'demography', title: 'Демография', subject: 'География' },
  { id: 'inflation', title: 'Инфляция', subject: 'Обществознание' },
];

const demoEdges: GraphEdge[] = [
  { id: 'e1', source: 'percent', target: 'solutions', type: 'hasStrictPrerequisite' },
  { id: 'e2', source: 'solutions', target: 'chemistry', type: 'hasStrictPrerequisite' },
  { id: 'e3', source: 'percent', target: 'demography', type: 'appliesTo' },
  { id: 'e4', source: 'percent', target: 'inflation', type: 'appliesTo' },
  { id: 'e5', source: 'chemistry', target: 'photosynthesis', type: 'enriches' },
];

const demoProgress: Record<string, ModuleStatus> = {
  percent: 'in-progress',
  solutions: 'available',
  chemistry: 'blocked',
  photosynthesis: 'available',
  demography: 'available',
  inflation: 'available',
};

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
    id: 'story-soc-inflation',
    title: 'Инфляция: почему дорожают товары',
    linkedModules: ['inflation', 'percent'],
    subjects: ['Обществознание', 'Математика'],
    realWorld: 'Центральные банки регулируют инфляцию процентными ставками.',
    readingMinutes: 4,
  },
];

const demoProjects: ProjectCardData[] = [
  {
    id: 'proj-flagship-budget',
    title: 'Семейный бюджет: математика и обществознание',
    modules: ['percent', 'taxes', 'inflation', 'economics'],
    difficultyLevel: 'basic',
    expectedOutcome: 'Составить и проанализировать семейный бюджет с учётом инфляции.',
  },
];

type DemoStep = 'select' | 'route' | 'dashboard' | 'map' | 'gaps' | 'recommendations';

const stepLabels: Record<DemoStep, string> = {
  select: '1. Выбор ученика',
  route: '2. Построение маршрута',
  dashboard: '3. Дашборд ученика',
  map: '4. Карта знаний',
  gaps: '5. Лакуны и покрытие',
  recommendations: '6. Истории и проекты',
};

const stepOrder: DemoStep[] = ['select', 'route', 'dashboard', 'map', 'gaps', 'recommendations'];

export function DemoPage() {
  const [step, setStep] = useState<DemoStep>('select');
  const [selectedChild, setSelectedChild] = useState<DemoChild | null>(null);
  const route = useRouteCompute(selectedChild?.id ?? 'demo-misha');

  const stepIndex = useMemo(() => stepOrder.indexOf(step), [step]);

  const next = () => {
    const nextStep = stepOrder[stepIndex + 1];
    if (nextStep) {
      setStep(nextStep);
    }
  };
  const prev = () => {
    const prevStep = stepOrder[stepIndex - 1];
    if (prevStep) {
      setStep(prevStep);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-semibold text-neutral-900">
            Демо: семейное образование — «Дай пять»
          </h1>
          <p className="text-sm text-neutral-500 mt-1">
            Сквозной сценарий валидации: маршрут за ≤5 минут, карта знаний, покрытие ФГОС, лакуны и
            мотивационный контент.
          </p>
        </div>
        <Badge color="info">{stepLabels[step]}</Badge>
      </div>

      {/* Step progress */}
      <div className="flex flex-wrap gap-2">
        {stepOrder.map((s) => (
          <button
            key={s}
            type="button"
            onClick={() => setStep(s)}
            className={`px-3 py-1 text-xs rounded-full transition-colors ${
              s === step
                ? 'bg-primary-500 text-white'
                : stepIndex > stepOrder.indexOf(s)
                  ? 'bg-green-100 text-green-700'
                  : 'bg-neutral-100 text-neutral-500'
            }`}
          >
            {stepLabels[s]}
          </button>
        ))}
      </div>

      {/* Step content */}
      <Card>
        {step === 'select' && (
          <div className="space-y-4">
            <h2 className="text-lg font-semibold">Выберите ученика</h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
              {demoChildren.map((child) => (
                <button
                  key={child.id}
                  type="button"
                  onClick={() => {
                    setSelectedChild(child);
                    next();
                  }}
                  className="rounded-lg border border-neutral-200 p-4 text-left hover:border-primary-500 hover:shadow transition-all"
                >
                  <div className="font-medium text-neutral-800">{child.name}</div>
                  <div className="text-xs text-neutral-500">{child.grade}</div>
                </button>
              ))}
            </div>
          </div>
        )}

        {step === 'route' && (
          <div className="space-y-4">
            <h2 className="text-lg font-semibold">Постройте маршрут за ≤5 минут</h2>
            <RouteBuilder route={route} />
          </div>
        )}

        {step === 'dashboard' && (
          <div className="space-y-4">
            <h2 className="text-lg font-semibold">Дашборд ученика</h2>
            <LearnerDashboard data={demoDashboard} />
          </div>
        )}

        {step === 'map' && (
          <div className="space-y-4">
            <h2 className="text-lg font-semibold">Карта знаний с прогрессом</h2>
            <KnowledgeMap nodes={demoGraph} edges={demoEdges} progress={demoProgress} />
          </div>
        )}

        {step === 'gaps' && (
          <div className="space-y-4">
            <h2 className="text-lg font-semibold">Лакуны и покрытие ФГОС</h2>
            <GapMap gaps={demoGaps} />
          </div>
        )}

        {step === 'recommendations' && (
          <div className="space-y-4">
            <h2 className="text-lg font-semibold">Мотивационный контент</h2>
            <RecommendationPanel
              moduleName="Проценты"
              stories={demoStories}
              projects={demoProjects}
            />
          </div>
        )}
      </Card>

      {/* Navigation */}
      <div className="flex justify-between">
        <Button variant="secondary" onClick={prev} disabled={stepIndex === 0}>
          ← Назад
        </Button>
        {stepIndex < stepOrder.length - 1 ? (
          <Button variant="primary" onClick={next}>
            Далее →
          </Button>
        ) : (
          <Button variant="primary" onClick={() => setStep('select')}>
            Заново
          </Button>
        )}
      </div>
    </div>
  );
}
