import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { OntologyBrowser } from '../OntologyBrowser';
import { useOntology } from '../useOntology';

// Mock the API layer so tests never hit the network.
vi.mock('../api', () => ({
  getConcept: vi.fn(async (topicId: string) => {
    if (topicId === 'math-5-2') {
      throw new Error('concept not found');
    }
    return {
      id: topicId,
      title: topicId === 'math-5-1' ? 'Натуральные числа' : `Topic ${topicId}`,
      description: 'Test module',
      links: [
        { topic_id: 'math-5-2', link_type: 'hasStrictPrerequisite' },
        { topic_id: 'math-5-3', link_type: 'hasSoftPrerequisite' },
      ],
    };
  }),
}));

describe('<OntologyBrowser>', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });
  afterEach(() => {
    vi.restoreAllMocks();
  });

  function Harness() {
    const ontology = useOntology('math-5-1');
    return <OntologyBrowser ontology={ontology} />;
  }

  it('renders the ontology graph after loading the root concept', async () => {
    render(<Harness />);
    await waitFor(() => expect(screen.getByTestId('ontology-graph')).toBeInTheDocument());
    await waitFor(() => expect(screen.getByTestId('node-math-5-1')).toBeInTheDocument());
    // Linked topics appear as nodes.
    expect(screen.getByTestId('node-math-5-2')).toBeInTheDocument();
    expect(screen.getByTestId('node-math-5-3')).toBeInTheDocument();
    // Color-coded edges for strict + soft prerequisites.
    expect(screen.getByTestId('edge-math-5-1-math-5-2')).toBeInTheDocument();
    expect(screen.getByTestId('edge-math-5-1-math-5-3')).toBeInTheDocument();
    // Detail panel shows the selected (root) module.
    expect(screen.getByTestId('module-detail')).toHaveTextContent('Натуральные числа');
  });

  it('selects a node and shows its detail', async () => {
    render(<Harness />);
    await waitFor(() => expect(screen.getByTestId('node-math-5-2')).toBeInTheDocument());
    fireEvent.click(screen.getByTestId('node-math-5-2'));
    await waitFor(() => expect(screen.getByTestId('module-detail')).toHaveTextContent('math-5-2'));
    expect(screen.getByTestId('module-detail')).toHaveTextContent('strict');
  });

  it('shows the error message when a concept cannot be loaded', async () => {
    render(<Harness />);
    await waitFor(() => expect(screen.getByTestId('node-math-5-1')).toBeInTheDocument());
    fireEvent.click(screen.getByTestId('node-math-5-2'));
    fireEvent.click(screen.getByRole('button', { name: 'Expand module' }));
    await waitFor(() => expect(screen.getByRole('alert')).toHaveTextContent('concept not found'));
  });
});
