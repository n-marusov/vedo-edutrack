import { render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { GapCoverage } from '../GapCoverage';

vi.mock('../useCoverage', () => ({
  useCoverage: () => ({
    coverage: {
      covered: 1,
      total: 2,
      percent: 50,
      ready: false,
      deficits: [{ requirement_id: 'r2', blocking_module_id: 'percent' }],
    },
    gaps: {
      status: 'root-causes-found',
      root_causes: [{ module_id: 'percent', mastery: 0.7, blocked_modules: 3 }],
    },
    loading: false,
    error: null,
  }),
}));

describe('<gap-coverage>', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('renders live FGOS coverage with deficits', async () => {
    render(<GapCoverage learnerId="l1" lagModuleId="chemistry" />);
    await waitFor(() => expect(screen.getByTestId('coverage-report')).toBeInTheDocument());
    expect(screen.getByTestId('coverage-report')).toHaveTextContent('50.0% covered');
    expect(screen.getByTestId('coverage-report')).toHaveTextContent('r2');
  });

  it('renders ranked root causes', async () => {
    render(<GapCoverage learnerId="l1" lagModuleId="chemistry" />);
    await waitFor(() => expect(screen.getByTestId('gap-diagnosis')).toBeInTheDocument());
    expect(screen.getByTestId('gap-diagnosis').textContent).toContain('percent');
    expect(screen.getByTestId('gap-diagnosis').textContent).toContain('blocks 3');
  });

  it('prompts for a lag module when none is provided', () => {
    render(<GapCoverage learnerId="l1" />);
    expect(screen.getByText('Pass a lag module id to run diagnosis.')).toBeInTheDocument();
  });
});
