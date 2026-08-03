/**
 * Visualization feature — barrel exports.
 */
export { KnowledgeMap } from './KnowledgeMap';
export type { KnowledgeMapProps } from './KnowledgeMap';
export { LearnerDashboard } from './LearnerDashboard';
export type { LearnerDashboardProps, LearnerDashboardData, Widget } from './LearnerDashboard';
export { ParentDashboard } from './ParentDashboard';
export type { ParentDashboardProps, ParentDashboardData, ChildOverview } from './ParentDashboard';
export { MethodologistDashboard } from './MethodologistDashboard';
export type {
  MethodologistDashboardProps,
  MethodologistDashboardData,
  ClassCoverage,
  LaggingTopic,
} from './MethodologistDashboard';
export { GapMap } from './GapMap';
export type { GapMapProps, RootGap } from './GapMap';
export { GroupPanel } from './GroupPanel';
export type { GroupPanelProps, LearnerCard } from './GroupPanel';
export * from './types';
export * from './api';
