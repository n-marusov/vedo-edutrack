/**
 * Visualization feature API calls.
 */
import { api } from '../../shared/api/client';
import type { GraphEdge, GraphNode, ModuleStatus } from './types';

export interface GraphDataResponse {
  nodes: GraphNode[];
  edges: GraphEdge[];
}

export interface ProgressResponse {
  progress: Record<string, ModuleStatus>;
}

/** Fetch the knowledge graph for visualization. */
export async function fetchGraphData(): Promise<GraphDataResponse> {
  return api.get<GraphDataResponse>('/api/v1/visualization/graph');
}

/** Fetch learner progress for color-coding. */
export async function fetchProgress(learnerId: string): Promise<ProgressResponse> {
  return api.get<ProgressResponse>(`/api/v1/learners/${learnerId}/progress`);
}
