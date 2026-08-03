import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { GapMap, type RootGap } from '../GapMap';
import { GroupPanel, type LearnerCard } from '../GroupPanel';
import { LearnerDashboard, type LearnerDashboardData } from '../LearnerDashboard';

const dashboardData: LearnerDashboardData = {
  learnerName: 'Misha',
  currentModule: 'percent',
  currentStatus: 'mastered',
  horizons: { far: 'chemistry', mid: 'solutions', near: 'percent' },
  planVsActual: { onTime: 2, delayed: 1, ahead: 0 },
  fgosCoverage: 70,
  subjectProgress: [
    { subject: 'Math', percent: 80 },
    { subject: 'Chemistry', percent: 40 },
  ],
  recommendedCount: 3,
  lastUpdated: '2026-08-03T12:00:00Z',
};

describe('<LearnerDashboard>', () => {
  it('renders all 6 mandatory widgets', () => {
    render(<LearnerDashboard data={dashboardData} />);
    expect(screen.getByText("Misha's Dashboard")).toBeInTheDocument();
    expect(screen.getByText('Current Position')).toBeInTheDocument();
    expect(screen.getByText('Three Horizons')).toBeInTheDocument();
    expect(screen.getByText('Plan vs Actual')).toBeInTheDocument();
    expect(screen.getByText('FGOS Coverage')).toBeInTheDocument();
    expect(screen.getByText('Subject Progress')).toBeInTheDocument();
    expect(screen.getByText('Recommendations')).toBeInTheDocument();
  });

  it('shows FGOS coverage with percentage', () => {
    render(<LearnerDashboard data={dashboardData} />);
    expect(screen.getByText('70%')).toBeInTheDocument();
    expect(screen.getByText('30% remaining')).toBeInTheDocument();
  });

  it('shows plan-vs-actual summary', () => {
    render(<LearnerDashboard data={dashboardData} />);
    expect(screen.getByText('2 on time')).toBeInTheDocument();
    expect(screen.getByText('1 delayed, 0 ahead')).toBeInTheDocument();
  });
});

const learners: LearnerCard[] = [
  {
    id: 'l1',
    name: 'Misha',
    currentModule: 'percent',
    fgosCoverage: 70,
    forecastStatus: 'at-risk',
    attentionFlag: true,
  },
  {
    id: 'l2',
    name: 'Katya',
    currentModule: 'fractions',
    fgosCoverage: 90,
    forecastStatus: 'on-track',
    attentionFlag: false,
  },
];

describe('<GroupPanel>', () => {
  it('shows X of Y at-risk summary', () => {
    render(<GroupPanel title="My Children" learners={learners} />);
    expect(screen.getByText('1 of 2 at risk')).toBeInTheDocument();
  });

  it('renders learner cards with required fields', () => {
    render(<GroupPanel title="My Children" learners={learners} />);
    expect(screen.getByText('Misha')).toBeInTheDocument();
    expect(screen.getByText('Katya')).toBeInTheDocument();
    expect(screen.getByText('Module: percent')).toBeInTheDocument();
    expect(screen.getByText(/FGOS: 70%/)).toBeInTheDocument();
    expect(screen.getByText('Forecast: at-risk')).toBeInTheDocument();
  });

  it('flags attention learners', () => {
    render(<GroupPanel title="My Children" learners={learners} />);
    expect(screen.getByText('⚠ Attention')).toBeInTheDocument();
  });

  it('invokes onSelect when a card is clicked', () => {
    const onSelect = vi.fn();
    render(<GroupPanel title="My Children" learners={learners} onSelect={onSelect} />);
    fireEvent.click(screen.getByText('Misha'));
    expect(onSelect).toHaveBeenCalledWith('l1');
  });
});

const gaps: RootGap[] = [
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

describe('<GapMap>', () => {
  it('shows ranked root gaps', () => {
    render(<GapMap gaps={gaps} />);
    expect(screen.getByText('Gap Diagnostic Map')).toBeInTheDocument();
    // Проценты appears in the top-gap summary and in the gap list.
    expect(screen.getAllByText('Проценты').length).toBeGreaterThan(0);
    expect(screen.getByText('Дроби')).toBeInTheDocument();
  });

  it('shows the top gap summary', () => {
    render(<GapMap gaps={gaps} />);
    expect(screen.getByText(/2 root gaps/)).toBeInTheDocument();
  });

  it('shows cascade impact per gap', () => {
    render(<GapMap gaps={gaps} />);
    expect(screen.getByText(/Blocks 4 modules in 3 subjects/)).toBeInTheDocument();
  });

  it('shows an empty state when no gaps exist', () => {
    render(<GapMap gaps={[]} />);
    expect(screen.getByText(/No root gaps detected/)).toBeInTheDocument();
  });
});
