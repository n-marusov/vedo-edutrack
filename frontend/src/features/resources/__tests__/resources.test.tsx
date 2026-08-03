import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { ResourceCatalog } from '../ResourceCatalog';

vi.mock('../../execution-progress/api', () => ({
  listResources: vi.fn(async (params: { format?: string } = {}) => {
    const all = [
      { id: 'r1', title: 'Percent video', type: 'content', format: 'video' },
      { id: 'r2', title: 'Percent text', type: 'content', format: 'text' },
      { id: 'r3', title: 'Lab access', type: 'enabling', format: 'lab' },
    ];
    const filtered = params.format
      ? all.filter((resource) => resource.format === params.format)
      : all;
    return { items: filtered, total: filtered.length };
  }),
}));

describe('<resources>', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('renders the resource catalog', async () => {
    render(<ResourceCatalog />);
    await waitFor(() => expect(screen.getByTestId('resource-list')).toBeInTheDocument());
    expect(screen.getByTestId('resource-list').textContent).toContain('Percent video');
  });

  it('filters by format', async () => {
    render(<ResourceCatalog />);
    await waitFor(() => expect(screen.getByTestId('resource-list')).toBeInTheDocument());
    fireEvent.change(screen.getByLabelText('Format filter'), { target: { value: 'video' } });
    await waitFor(() =>
      expect(screen.getByTestId('resource-list').textContent).toContain('Percent video'),
    );
    expect(screen.getByTestId('resource-list').textContent).not.toContain('Percent text');
  });

  it('marks enabling resources with a badge', async () => {
    render(<ResourceCatalog />);
    await waitFor(() => expect(screen.getByTestId('resource-list')).toBeInTheDocument());
    expect(screen.getByText('enabling')).toBeInTheDocument();
  });
});
