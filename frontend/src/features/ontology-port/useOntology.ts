import { useCallback, useEffect, useState } from 'react';
import type { Concept } from '../../shared/api/types';
import { getConcept } from './api';

/** Graph node — a concept expanded from the ontology. */
export interface GraphNode {
  id: string;
  title: string;
  description?: string;
  /** Link type used to reach this node from its parent (undefined for the root). */
  linkType?: string;
}

/** Ontology graph state for the browser. */
export interface OntologyGraphState {
  nodes: GraphNode[];
  /** edges: [source, target, linkType] */
  edges: [string, string, string][];
  selectedId: string | null;
  loading: boolean;
  error: string | null;
}

export interface OntologyBrowserState extends OntologyGraphState {
  selectNode: (id: string) => void;
  expandNode: (id: string) => Promise<void>;
  reload: (rootId: string) => Promise<void>;
}

/**
 * useOntology — fetches the ontology around a root concept and supports
 * incremental expansion via linked topics (F0/F4 graph browser).
 *
 * A concept response carries its outgoing links (prerequisite -> consequent),
 * so expanding a node adds its linked topics as child nodes in one step.
 */
export function useOntology(initialRootId = 'math-5-1'): OntologyBrowserState {
  const [state, setState] = useState<OntologyGraphState>({
    nodes: [],
    edges: [],
    selectedId: null,
    loading: false,
    error: null,
  });

  const applyConcept = useCallback((concept: Concept, parentLinkType?: string) => {
    setState((prev) => {
      const nodes = new Map(prev.nodes.map((node) => [node.id, node]));
      const edges = new Set(prev.edges.map((edge) => edge.join('|')));

      // The fetched concept itself (root or re-fetch of an existing node).
      if (!nodes.has(concept.id)) {
        nodes.set(concept.id, {
          id: concept.id,
          title: concept.title,
          description: concept.description,
          linkType: parentLinkType,
        });
      } else if (concept.title) {
        const existing = nodes.get(concept.id) as GraphNode;
        nodes.set(concept.id, {
          ...existing,
          title: existing.title || concept.title,
          description: existing.description ?? concept.description,
        });
      }

      // Outgoing links add child topics as new nodes.
      for (const link of concept.links ?? []) {
        if (!nodes.has(link.topic_id)) {
          nodes.set(link.topic_id, {
            id: link.topic_id,
            title: link.topic_id,
            linkType: link.link_type,
          });
        }
        edges.add([concept.id, link.topic_id, link.link_type].join('|'));
      }

      return {
        nodes: [...nodes.values()],
        edges: [...edges].map((edge) => edge.split('|') as [string, string, string]),
        selectedId: prev.selectedId ?? concept.id,
        loading: false,
        error: null,
      };
    });
  }, []);

  const expandNode = useCallback(
    async (id: string) => {
      setState((prev) => ({ ...prev, loading: true, selectedId: id, error: null }));
      try {
        const concept = await getConcept(id);
        applyConcept(concept);
      } catch (err) {
        const message = err instanceof Error ? err.message : String(err);
        setState((prev) => ({ ...prev, loading: false, error: message }));
      }
    },
    [applyConcept],
  );

  const reload = useCallback(
    async (rootId: string) => {
      setState({ nodes: [], edges: [], selectedId: null, loading: true, error: null });
      try {
        const concept = await getConcept(rootId);
        applyConcept(concept);
      } catch (err) {
        const message = err instanceof Error ? err.message : String(err);
        setState({ nodes: [], edges: [], selectedId: null, loading: false, error: message });
      }
    },
    [applyConcept],
  );

  const selectNode = useCallback((id: string) => {
    setState((prev) => ({ ...prev, selectedId: id }));
  }, []);

  useEffect(() => {
    void reload(initialRootId);
  }, [initialRootId, reload]);

  return { ...state, selectNode, expandNode, reload };
}
