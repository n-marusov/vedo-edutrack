/**
 * Practice Life feature — barrel exports.
 */
export { StoryCard, ProjectCard, RecommendationPanel } from './PracticeComponents';
export type {
  StoryCardData,
  ProjectCardData,
  RecommendationPanelProps,
} from './PracticeComponents';

/** API functions for practice-life data. */
import { api } from '../../shared/api/client';
import type { ProjectCardData, StoryCardData } from './PracticeComponents';

export interface PracticeDataResponse {
  stories: StoryCardData[];
  projects: ProjectCardData[];
}

export async function fetchStoriesForModule(moduleId: string): Promise<StoryCardData[]> {
  return api.get<StoryCardData[]>(`/api/v1/modules/${moduleId}/stories`);
}

export async function fetchProjectsForModule(moduleId: string): Promise<ProjectCardData[]> {
  return api.get<ProjectCardData[]>(`/api/v1/modules/${moduleId}/projects`);
}

export async function fetchRecommendations(learnerId: string): Promise<PracticeDataResponse> {
  return api.get<PracticeDataResponse>(`/api/v1/learners/${learnerId}/recommended-stories`);
}
