import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import type { RouterTopic } from '../../../shared/api/types';
import { RouteBuilder } from '../RouteBuilder';
import { horizonSlice, useRouteCompute } from '../useRouteCompute';

vi.mock('../api', () => ({
  computeRoute: vi.fn(async (_learnerId: string, goalTopicId: string) => {
    if (goalTopicId === 'unreachable') {
      throw new Error('topic_not_found');
    }
    return [
      { topic_id: 'math-5-1', order: 0 },
      { topic_id: 'math-5-11', order: 1, link_type: 'hasStrictPrerequisite' },
      { topic_id: goalTopicId, order: 2 },
    ];
  }),
}));

describe('<RouteBuilder>', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });
  afterEach(() => {
    vi.restoreAllMocks();
  });

  function Harness() {
    const route = useRouteCompute('l1');
    return <RouteBuilder route={route} />;
  }

  it('computes a route and renders the timeline', async () => {
    render(<Harness />);
    fireEvent.change(screen.getByLabelText('Goal module id'), { target: { value: 'math-5-12' } });
    fireEvent.click(screen.getByRole('button', { name: 'Compute route' }));
    await waitFor(() => expect(screen.getByTestId('route-timeline')).toBeInTheDocument());
    expect(screen.getByTestId('route-timeline').textContent).toContain('math-5-11');
    expect(screen.getByTestId('route-timeline').textContent).toContain('math-5-12');
  });

  it('switches between far/mid/near horizons', async () => {
    render(<Harness />);
    fireEvent.change(screen.getByLabelText('Goal module id'), { target: { value: 'math-5-12' } });
    fireEvent.click(screen.getByRole('button', { name: 'Compute route' }));
    await waitFor(() => expect(screen.getByTestId('route-timeline')).toBeInTheDocument());

    fireEvent.click(screen.getByRole('button', { name: 'Near (current module)' }));
    expect(screen.getByTestId('route-timeline').children).toHaveLength(1);
  });

  it('shows an error when the goal is unreachable', async () => {
    render(<Harness />);
    fireEvent.change(screen.getByLabelText('Goal module id'), { target: { value: 'unreachable' } });
    fireEvent.click(screen.getByRole('button', { name: 'Compute route' }));
    await waitFor(() => expect(screen.getByRole('alert')).toHaveTextContent('topic_not_found'));
  });

  it('slices the route by horizon', () => {
    const route: RouterTopic[] = Array.from({ length: 10 }, (_, i) => ({
      topic_id: `m${i}`,
      order: i,
    }));
    expect(horizonSlice(route, 'near')).toHaveLength(1);
    expect(horizonSlice(route, 'mid')).toHaveLength(5);
    expect(horizonSlice(route, 'far')).toHaveLength(10);
  });
});
