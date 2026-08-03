import { render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { ProgressDashboard } from '../ProgressDashboard';
import type { PlanModuleProgress } from '../ProgressDashboard';

vi.mock('../api', () => ({
  getFgosCoverage: vi.fn(async () => ({
    covered: 2,
    total: 3,
    percent: 66.7,
    ready: false,
    deficits: [{ requirement_id: 'r3', blocking_module_id: 'm3' }],
  })),
  diagnoseGaps: vi.fn(async () => ({ status: 'root-causes-found', root_causes: [] })),
  listResources: vi.fn(async () => ({ items: [], total: 0 })),
}));

describe('<execution-progress>', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });
  afterEach(() => {
    vi.restoreAllMocks();
  });

  const modules: PlanModuleProgress[] = [
    { module_id: 'm1', planned_end: '2026-09-01', mastered: true, deviationDays: 0 },
    { module_id: 'm2', planned_end: '2026-09-15', mastered: true, deviationDays: 5 },
  ];

  it('renders plan-vs-actual comparison with status badges', () => {
    render(<ProgressDashboard learnerId="l1" modules={modules} />);
    const list = screen.getByTestId('progress-list');
    expect(list.textContent).toContain('m1');
    expect(list.textContent).toContain('m2');
    expect(list.textContent).toContain('mastered');
  });

  it('flags late modules with deviation', () => {
    render(<ProgressDashboard learnerId="l1" modules={modules} />);
    expect(screen.getByTestId('progress-list').textContent).toContain('+5d late');
  });

  it('shows at-risk forecast when a large share of modules is late', () => {
    const late = [
      { module_id: 'a', planned_end: '2026-09-01', mastered: true, deviationDays: 4 },
      { module_id: 'b', planned_end: '2026-09-01', mastered: true, deviationDays: 6 },
      { module_id: 'c', planned_end: '2026-09-01', mastered: true, deviationDays: 2 },
    ];
    render(<ProgressDashboard learnerId="l1" modules={late} />);
    expect(screen.getByText('off-track')).toBeInTheDocument();
  });

  it('loads FGOS coverage through the coverage hook without errors', async () => {
    render(<ProgressDashboard learnerId="l1" modules={modules} />);
    await waitFor(() => expect(screen.queryByRole('alert')).not.toBeInTheDocument());
  });
});
