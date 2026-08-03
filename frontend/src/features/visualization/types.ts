/**
 * Visualization feature type definitions.
 */

export type ModuleStatus = 'mastered' | 'in-progress' | 'available' | 'blocked' | 'unclosed-prereq';

export type GraphMode = 'critical-path' | 'exploration';

export interface GraphNode {
  id: string;
  title: string;
  subject: string;
  x?: number;
  y?: number;
}

export interface GraphEdge {
  id: string;
  source: string;
  target: string;
  type: string; // hasStrictPrerequisite | hasSoftPrerequisite | enriches | appliesTo | isAnalogousTo
}
