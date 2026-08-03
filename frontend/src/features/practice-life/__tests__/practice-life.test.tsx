import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import {
  ProjectCard,
  type ProjectCardData,
  RecommendationPanel,
  StoryCard,
  type StoryCardData,
} from '../PracticeComponents';

const story: StoryCardData = {
  id: 's1',
  title: 'Проценты в жизни',
  linkedModules: ['percent', 'math-5-11'],
  subjects: ['Математика'],
  realWorld: 'Проценты используются в банковских вкладах и скидках.',
  readingMinutes: 3,
};

const project: ProjectCardData = {
  id: 'p1',
  title: 'Биохимическая лаборатория дома',
  modules: ['solutions', 'chemistry'],
  difficultyLevel: 'medium',
  expectedOutcome: 'Провести серию опытов и описать результаты.',
};

describe('<StoryCard>', () => {
  it('renders story title, subjects, and real-world section', () => {
    render(<StoryCard story={story} />);
    expect(screen.getByText('Проценты в жизни')).toBeInTheDocument();
    expect(screen.getByText('Математика')).toBeInTheDocument();
    expect(screen.getByText(/Проценты используются/)).toBeInTheDocument();
  });

  it('shows reading time', () => {
    render(<StoryCard story={story} />);
    expect(screen.getByText('3 min read')).toBeInTheDocument();
  });

  it('shows linked modules', () => {
    render(<StoryCard story={story} />);
    expect(screen.getByText(/Linked modules: percent, math-5-11/)).toBeInTheDocument();
  });
});

describe('<ProjectCard>', () => {
  it('renders project title, outcome, and difficulty', () => {
    render(<ProjectCard project={project} />);
    expect(screen.getByText('Биохимическая лаборатория дома')).toBeInTheDocument();
    expect(screen.getByText('medium')).toBeInTheDocument();
    expect(screen.getByText(/Провести серию опытов/)).toBeInTheDocument();
  });

  it('shows required modules', () => {
    render(<ProjectCard project={project} />);
    expect(screen.getByText(/Requires: solutions, chemistry/)).toBeInTheDocument();
  });
});

describe('<RecommendationPanel>', () => {
  it('shows recommendations when present', () => {
    render(<RecommendationPanel moduleName="Проценты" stories={[story]} projects={[project]} />);
    expect(screen.getByText(/You've mastered/)).toBeInTheDocument();
    expect(screen.getByText('Проценты')).toBeInTheDocument();
    expect(screen.getByText(/📖 Related Stories/)).toBeInTheDocument();
    expect(screen.getByText(/🔬 Project Ideas/)).toBeInTheDocument();
  });

  it('shows an empty state when no recommendations exist', () => {
    render(<RecommendationPanel moduleName="Проценты" stories={[]} projects={[]} />);
    expect(screen.getByText(/No recommendations available/)).toBeInTheDocument();
  });
});
