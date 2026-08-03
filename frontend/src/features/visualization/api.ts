/**
 * Visualization feature API calls.
 *
 * Contracts mirror the OpenAPI spec (backend/api/openapi/v1.yaml):
 * - GET /learners/{learner_id}/progress → ProgressResponse
 * - GET /ontology/concepts → ConceptResponse (graph building blocks)
 */
import { api } from '../../shared/api/client';
import { getConcept } from '../ontology-port/api';
import type { GraphEdge, GraphNode, ModuleStatus } from './types';

/** Raw API contract for GET /learners/{learner_id}/progress. */
export interface ProgressResponse {
  learner_id: string;
  plan_id?: string;
  generated_at?: string;
  modules: ModuleProgressItem[];
}

/** Per-module plan-vs-actual item (mirrors OpenAPI ModuleProgressItem). */
export interface ModuleProgressItem {
  module_id: string;
  status: 'not_started' | 'in_progress' | 'mastered' | 'skipped';
  planned_date?: string;
  actual_date?: string;
  deviation_days?: number;
  deviation_cause?: string;
}

const statusToModuleStatus: Record<ModuleProgressItem['status'], ModuleStatus> = {
  not_started: 'available',
  in_progress: 'in-progress',
  mastered: 'mastered',
  skipped: 'available',
};

/** Convert a modules array to a { module_id → ModuleStatus } map for color-coding. */
export function toStatusMap(modules: ModuleProgressItem[]): Record<string, ModuleStatus> {
  const map: Record<string, ModuleStatus> = {};
  for (const m of modules) {
    map[m.module_id] = statusToModuleStatus[m.status] ?? 'available';
  }
  return map;
}

/** Fetch learner progress for color-coding (returns the raw API contract). */
export async function fetchProgress(learnerId: string): Promise<ProgressResponse> {
  return api.get<ProgressResponse>(`/api/v1/learners/${learnerId}/progress`);
}

/**
 * Build the knowledge graph for a set of module IDs from the ontology.
 * The ontology-port exposes single-concept reads; a graph is assembled by
 * fetching each module and expanding its linked topics.
 */
export async function buildGraph(
  moduleIds: string[],
): Promise<{ nodes: GraphNode[]; edges: GraphEdge[] }> {
  const nodes = new Map<string, GraphNode>();
  const edges = new Map<string, GraphEdge>();

  for (const moduleId of moduleIds) {
    const concept = await getConcept(moduleId);
    if (!nodes.has(concept.id)) {
      nodes.set(concept.id, {
        id: concept.id,
        title: concept.title,
        subject: concept.description ?? '',
      });
    }
    for (const link of concept.links ?? []) {
      if (!nodes.has(link.topic_id)) {
        // Linked topic node: id known, title resolved lazily by the caller.
        nodes.set(link.topic_id, { id: link.topic_id, title: link.topic_id, subject: '' });
      }
      const edgeId = `${concept.id}->${link.topic_id}`;
      if (!edges.has(edgeId)) {
        edges.set(edgeId, {
          id: edgeId,
          source: concept.id,
          target: link.topic_id,
          type: link.link_type,
        });
      }
    }
  }

  return { nodes: [...nodes.values()], edges: [...edges.values()] };
}
